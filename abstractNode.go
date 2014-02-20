package asciidocgo

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/VonC/asciidocgo/consts/context"
	"github.com/VonC/asciidocgo/consts/regexps"
	"github.com/VonC/asciidocgo/consts/safemode"
)

var testan = ""

type Documentable interface {
	Document() Documentable
	Attributes() map[string]interface{}
	Attr(name string, defaultValue interface{}, inherit bool) interface{}
	HasAttr(name string, expect interface{}, inherit bool) bool
	setAttr(name string, val interface{}, override bool) bool
	HasReftext() bool

	Safe() safemode.SafeMode
	BaseDir() string

	PlaybackAttributes(map[string]interface{})
	Renderer() *Renderer
	CounterIncrement(counterName string, block *abstractNode) string
	Counter(name, seed string) int
	DocType() string
}

/* An abstract base class that provides state and methods for managing
a node of AsciiDoc content.
The state and methods on this class are comment to all content segments
in an AsciiDoc document. */
type abstractNode struct {
	parent     *abstractNode
	id         string
	context    context.Context
	document   Documentable
	attributes map[string]interface{}
	_doc       Documentable
	*substitutors
}

func newAbstractNode(parent *abstractNode, c context.Context) *abstractNode {
	abstractNode := &abstractNode{parent, "", c, nil, make(map[string]interface{}), nil, &substitutors{}}
	if c == context.Document {
		abstractNode.parent = nil
		if parent != nil {
			abstractNode.document = parent._doc
		}
	} else if parent != nil {
		abstractNode.document = parent.Document()
	}
	return abstractNode
}

func (an *abstractNode) MainDocumentable(d Documentable) {
	if an.Context() == context.Document {
		an._doc = d
	}
}

//  Get the element which is the parent of this node
func (an *abstractNode) Parent() *abstractNode {
	return an.parent
}

//  Get the Asciidoctor::Document to which this node belongs
func (an *abstractNode) Document() Documentable {
	return an.document
}

// Get the Id for this node
func (an *abstractNode) Id() string {
	return an.id
}

// Set the Id for this node
func (an *abstractNode) SetId(id string) {
	an.id = id
}

// Get the Symbol context for this node
func (an *abstractNode) Context() context.Context {
	return an.context
}

func (an *abstractNode) Attributes() map[string]interface{} {
	return an.attributes
}

// Associate this Block with a new parent Block
// parent: The Block to set as the parent of this Block
func (an *abstractNode) SetParent(parent *abstractNode) {
	an.parent = parent
	if parent != nil {
		an.document = parent.Document()
	} else {
		an.document = nil
	}
}

/* Get the value for the specified attribute.

First look in the attributes on this node and return the value
of the attribute if found.
Otherwise, if this node is a child of the Document node, look in
the attributes of the Document node and return the value of the attribute
if found.
Otherwise, return the default value, which defaults to nil.

- name: the String or Symbol name of the attribute to lookup
- default_value: the Object value to return if the attribute is not found
(default: nil)
- inherit: a Boolean indicating whether to check for the attribute on
the AsciiDoctor::Document if not found on this node (default: false)

Return the value of the attribute or the default value if the attribute is
not found in the attributes of this node or the document node
*/
func (an *abstractNode) Attr(name string, defaultValue interface{}, inherit bool) interface{} {
	if an._doc == an.document {
		inherit = false
	}
	if an.attributes[name] != nil {
		return an.attributes[name]
	}
	if inherit {
		if an.document != nil && an.document.Attributes()[name] != nil {
			return an.document.Attributes()[name]
		}
	}
	return defaultValue
}

/*
Check if the attribute is defined, optionally performing a comparison of
its value if expected is not nil

Check if the attribute is defined.
First look in the attributes on this node.
If not found, and this node is a child of the Document node,
look in the attributes of the Document node.
If the attribute is found and a comparison value is specified (not nil),
return whether the two values match.
Otherwise, return whether the attribute was found.

name    - the String or Symbol name of the attribute to lookup
expect  - the expected Object value of the attribute (default: nil)
inherit - a Boolean indicating whether to check for the attribute on the
          AsciiDoctor::Document if not found on this node (default: false)

return a Boolean indicating whether the attribute exists and, if a
comparison value is specified, whether the value of the attribute matches
the comparison value
*/
func (an *abstractNode) HasAttr(name string, expect interface{}, inherit bool) bool {
	if an._doc == an.document {
		inherit = false
	}
	if expect == nil {
		if _, hasAttr := an.attributes[name]; hasAttr {
			return true
		}
		if inherit {
			if an.document != nil {
				if _, hasAttr := an.document.Attributes()[name]; hasAttr {
					return true
				}
			}
		}
		return false
	}
	if an.attributes[name] != nil {
		return (expect == an.attributes[name])
	}
	if inherit {
		if an.document != nil && an.document.Attributes()[name] != nil {
			return (expect == an.document.Attributes()[name])
		}
	}
	return false
}

/* Assign the value to the specified key in this block's attributes hash.

- key: The attribute key (or name)
- val: The value to assign to the key

returns a flag indicating whether the assignment was performed
*/
func (an *abstractNode) setAttr(name string, val interface{}, override bool) bool {
	if override {
		an.attributes[name] = val
		return true
	}
	if _, hasName := an.attributes[name]; !hasName {
		an.attributes[name] = val
		return true
	}
	return false
}

/* Enable a specified option attribute on the current node.

This method defines the `%name%-option` attribute on the current node.

- name: the String or Symbol name of the option
*/
func (an *abstractNode) SetOption(option string) {
	val := an.attributes["options"]
	if val == nil {
		valmap := make(map[string]bool)
		an.attributes["options"] = valmap
		val = valmap
	}
	valmap := val.(map[string]bool)
	if _, hasOption := valmap[option]; !hasOption {
		valmap[option] = true
	}
	an.attributes[option+"-option"] = true
}

/*  A convenience method to check if the specified option attribute is enabled
on the current node.

Check if the option is enabled.
This method simply checks to see if the `%name%-option` attribute is defined
on the current node.

- name: the String or Symbol name of the option

return a Boolean indicating whether the option has been specified
*/
func (an *abstractNode) HasOption(option string) bool {
	_, res := an.attributes[option+"-option"]
	return res
}

/* Get the Renderer instance being used for the
Document to which this node belongs */
func (an *abstractNode) Renderer() *Renderer {
	var res *Renderer = nil
	if an.Document() != nil {
		res = an.Document().Renderer()
	}
	return res
}

/* Update the attributes of this node with the new values
in the attributes argument.

If an attribute already exists with the same key,
it's value will be overridden.

- attributes: A Hash of attributes to assign to this node.
*/
func (an *abstractNode) UpdateAttributes(attrs map[string]interface{}) {
	for key, value := range attrs {
		an.attributes[key] = value
	}
}

// A convenience method that checks if the specified role is present
// in the list of roles on this node
func (an *abstractNode) HasRole(role interface{}) bool {
	if role == nil {
		if _, hasRole := an.attributes["role"]; hasRole {
			return true
		}
		if an.Document() != nil {
			if _, hasRole := an.Document().Attributes()["role"]; hasRole {
				return true
			}
		}
		return false
	}
	if anAttr := an.Attr("role", nil, true); anAttr == role {
		return true
	}
	return false
}

// A convenience method that checks if the specified role is present
// in the list of roles on this node
func (an *abstractNode) HasARole(name string) bool {
	if name == "" {
		return false
	}
	// inherit = true: check an.Document() as well
	roles := an.Attr("role", nil, true)
	if roles == nil {
		return false
	}
	rolesString := roles.(string)
	rolesArray := strings.Split(rolesString, " ")
	for _, role := range rolesArray {
		if name == role {
			return true
		}
	}
	return false
}

// A convenience method that returns the value of the role attribute
func (an *abstractNode) Role() interface{} {
	// inherit = true: check an.Document() as well
	return an.Attr("role", nil, true)
}

// A convenience method that returns the role names as an Array
func (an *abstractNode) RoleNames() []string {
	roles := an.Attr("role", nil, true)
	if roles == nil {
		return []string{}
	}
	rolesString := roles.(string)
	return strings.Split(rolesString, " ")
}

// A convenience method that checks if the reftext attribute is specified
func (an *abstractNode) HasReftext() bool {
	reftext := an.Attr("reftext", nil, true)
	return (reftext != nil)
}

// A convenience method that returns the value of the reftext attribute
func (an *abstractNode) Reftext() interface{} {
	// inherit = true: check an.Document() as well
	return an.Attr("reftext", nil, true)
}

// Returns a forward slash if the attribute htmlsyntax has the value "xml".
func (an *abstractNode) ShortTagSlash() *rune {
	if an.Document() == nil {
		return nil
	}
	if an.Document().Attr("htmlsyntax", nil, false) == "xml" {
		r, _ := utf8.DecodeLastRuneInString("/")
		return &r
	} else {
		return nil
	}
}

/* Construct a reference or data URI to an icon image
for the specified icon name.

If the 'icon' attribute is set on this block, the name is ignored
and the value of this attribute is used as the  target image path.
Otherwise, construct a target image path by concatenating the value
of the 'iconsdir' attribute, the icon name and the value of the
'icontype' attribute (defaulting to 'png').

The target image path is then passed through the #image_uri() method.
If the 'data-uri' attribute is set on the document, the image will be
safely converted to a data URI.

The return value of this method can be safely used in an image tag.
Returns A String reference or data URI for an icon image */

func (an *abstractNode) IconUri(name string) string {
	if an.HasAttr("icon", nil, false) {
		return an.ImageUri(an.Attr("icon", nil, false).(string), "")
	} else {
		targetImage := name + "."
		if an.Document() != nil {
			targetImage = targetImage + an.Document().Attr("icontype", "png", false).(string)
		}
		return an.ImageUri(targetImage, "iconsdir")
	}
}

/* Construct a URI reference to the target media.

If the target media is a URI reference, then leave it untouched.

The target media is resolved relative to the directory retrieved from
the specified attribute key, if provided.

The return value can be safely used in a media tag (img, audio, video).

target        - A String reference to the target media
asset_dir_key - The String attribute key used to lookup the directory where
(default: 'imagesdir')

Returns A String reference for the target media
*/
func (an *abstractNode) MediaUri(target string, assetDirKey string) string {
	if assetDirKey == "" {
		assetDirKey = "imagesdir"
	}
	if strings.Contains(target, ":") && regexps.UriSniffRx.MatchString(target) {
		return target
	} else if assetDirKey != "" && an.HasAttr(assetDirKey, nil, true) {
		// normalize_web_path(target, @document.attr(asset_dir_key)) ???
		// How? (BUG?) @document can be nil.
		// At least, ask attr on an, with inherit true.
		return normalizeWebPath(target, an.Attr(assetDirKey, nil, true).(string))
	}
	return normalizeWebPath(target, "")
}

/* Construct a URI reference or data URI to the target image.

If the target image is a URI reference, then leave it untouched.

The target image is resolved relative to the directory retrieved from the
specified attribute key, if provided.

If the 'data-uri' attribute is set on the document, and the safe mode level
is less than SafeMode::SECURE, the image will be safely converted to
a data URI by reading it from the same directory. If neither of these conditions
are satisfied, a relative path (i.e., URL) will be returned.

The return value of this method can be safely used in an image tag.

target_image - A String path to the target image
asset_dir_key - The String attribute key used to lookup the directory where
the image is located (default: 'imagesdir')

Returns A String reference or data URI for the target image */
func (an *abstractNode) ImageUri(targetImage, assetDirKey string) string {
	if assetDirKey == "" {
		assetDirKey = "imagesdir"
	}
	if strings.Contains(targetImage, ":") && regexps.UriSniffRx.MatchString(targetImage) {
		return targetImage
	}
	if an.Document() != nil && an.Document().Safe() < safemode.SECURE && an.Document().HasAttr("data-uri", nil, true) {
		return an.generateDataUri(targetImage, assetDirKey)
	}
	if assetDirKey != "" && an.HasAttr(assetDirKey, nil, true) {
		return normalizeWebPath(targetImage, an.Document().Attr(assetDirKey, nil, false).(string))
	} else {
		return normalizeWebPath(targetImage, "")
	}
}

/* Generate a data URI that can be used to embed an image in the output document

First, and foremost, the target image path is cleaned if the document
safe mode level is set to at least SafeMode::SAFE
(a condition which is true by default) to prevent access to ancestor paths
in the filesystem.
The image data is then read and converted to Base64.
Finally, a data URI is built which can be used in an image tag.

target_image - A String path to the target image
asset_dir_key - The String attribute key used to lookup the directory where
                the image is located (default: nil)

Returns A String data URI containing the content of the target image*/
func (an *abstractNode) generateDataUri(targetImage, assetDirKey string) string {
	ext := filepath.Ext(targetImage)
	if len(ext) > 1 {
		ext = ext[1:]
	}
	mimetype := "image/" + ext
	if ext == "svg" {
		mimetype = mimetype + "+xml"
	}
	//return fmt.Sprintf("ext='%v' for mimetype='%v'", ext, mimetype)
	imagePath := ""
	if assetDirKey != "" && an.Document() != nil && an.Document().Attr(assetDirKey, nil, true) != nil {
		// image_path = normalize_system_path(target_image, @document.attr(asset_dir_key), nil, :target_name => 'image')
		imagePath = an.normalizeSystemPath(targetImage, an.Document().Attr(assetDirKey, nil, true).(string), "", false, "image")
	} else {
		imagePath = an.normalizeSystemPath(targetImage, "", "", false, "")
	}
	if testan == "test_generateDataUri_imagePath" {
		return fmt.Sprintf("imagePath='%v'", imagePath)
	}
	if file, err := os.Open(imagePath); err == nil {
		defer file.Close()
		reader := bufio.NewReader(file)
		if content, err := ioutil.ReadAll(reader); err == nil {
			return string(content)
		}
	}
	fmt.Errorf("asciidocgo: WARNING: image to embed not found or not readable: '%v'", imagePath)
	return "data:" + mimetype + ":base64,"
	// uncomment to return 1 pixel white dot instead
	// return 'data:image/gif;base64,R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw=='
}

/* Read the contents of the file at the specified path.
This method assumes that the path is safe to read. It checks
that the file is readable before attempting to read it.
path            - the String path from which to read the contents
warn_on_failure - a Boolean that controls whether a warning is issued if
                  the file cannot be read
returns the contents of the file at the specified path, or nil
if the file does not exist. */
func ReadAsset(path string, warnOnFailure bool) string {
	if file, err := os.Open(path); err == nil {
		defer file.Close()
		reader := bufio.NewReader(file)
		if content, err := ioutil.ReadAll(reader); err == nil {
			res := string(content)
			// QUESTION should we use strip or rstrip instead of chomp here?
			// Here: uses a more advanced chomp function
			res = strings.TrimRight(res, " \r\n")
			return res
		}
	}
	if warnOnFailure {
		fmt.Errorf("asciidocgo: WARNING: file does not exist or cannot be read: '%v'", path)
	}
	return ""
}

/* Normalize the web page using the PathResolver.
target - the String target path
start  - the String start (i.e, parent) path (optional, default: nil)
returns the resolved String path */
func normalizeWebPath(target, start string) string {
	res := WebPath(target, start)
	return res
}

/* Resolve and normalize a secure path from the target and
start paths using the PathResolver.

See PathResolver::system_path(target, start, jail, opts) for details.

The most important functionality in this method is to prevent resolving a
path outside of the jail (which defaults to the directory of the source
file, stored in the base_dir instance variable on Document) if the document
safe level is set to SafeMode::SAFE or greater (a condition which is true
by default).

target - the String target path
start  - the String start (i.e., parent) path
jail   - the String jail path to confine the resolved path
opts   - an optional Hash of options to control processing (default: {}):
          * :recover is used to control whether the processor should auto-recover
              when an illegal path is encountered
          * :target_name is used in messages to refer to the path being resolved

raises a SecurityError if a jail is specified and the resolved path is
outside the jail.

returns a String path resolved from the start and target paths, with any
parent references resolved and self references removed. If a jail is provided,
this path will be guaranteed to be contained within the jail. */
//def normalize_system_path(target, start = nil, jail = nil, opts = {})
func (an *abstractNode) normalizeSystemPath(target, start, jail string, canrecover bool, targetName string) string {
	if start == "" && an.Document() != nil {
		start = an.Document().BaseDir()
	}
	if jail == "" && an.Document() != nil && (an.Document().Safe() >= safemode.SAFE || testan == "test_normalizeSystemPath_safeDocument") {
		jail = an.Document().BaseDir()
	}
	return NewPathResolver(0, "").SystemPath(target, start, jail, canrecover, targetName)
}

/*Normalize the asset file or directory to a concrete and rinsed path

Delegates to normalize_system_path, with the start path set to the value of
the base_dir instance variable on the Document object. */
func (an *abstractNode) normalizeAssetPath(assetRef, assetName string, autocorrect bool) string {
	if assetName == "" {
		assetName = "path"
	}
	start := ""
	if an.Document() != nil {
		start = an.Document().BaseDir()
	}
	return an.normalizeSystemPath(assetRef, start, "", autocorrect, assetName)
}

/* Calculate the relative path to this absolute filename
from the Document#base_dir */
func (an *abstractNode) relativePath(filename string) string {
	baseDirectory := ""
	if an.Document() != nil {
		baseDirectory = an.Document().BaseDir()
	}
	return NewPathResolver(0, "").RelativePath(filename, baseDirectory)
}

/* an abstract block would have a style */
func (an *abstractNode) Style() string {
	return ""
}

/* Retrieve the list marker keyword for the specified list type.
For use in the HTML type attribute.
list_type - the type of list; default to the @style if not specified
returns the single-character String keyword that represents
the marker for the specified list type */
func (an *abstractNode) listMarkerKeyword(listType string) rune {
	if listType == "" {
		listType = an.Style()
	}
	return regexps.ORDERED_LIST_KEYWORDS[listType]
}
