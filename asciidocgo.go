/*
Asciidocgo implements an AsciiDoc renderer in Go.

Methods for parsing Asciidoc input files and rendering documents using eRuby
templates.

Asciidoc documents comprise a header followed by zero or more sections.
Sections are composed of blocks of content.  For example:

  = Doc Title

  == Section 1

  This is a paragraph block in the first section.

  == Section 2

  This section has a paragraph block and an olist block.

  . Item 1
  . Item 2

Examples:

Use built-in templates:

  lines = File.readlines("your_file.asc")
  doc = Asciidoctor::Document.new(lines)
  html = doc.render
  File.open("your_file.html", "w+") do |file|
    file.puts html
  end

Use custom (Tilt-supported) templates:

  lines = File.readlines("your_file.asc")
  doc = Asciidoctor::Document.new(lines, :template_dir => 'templates')
  html = doc.render
  File.open("your_file.html", "w+") do |file|
    file.puts html
  end

*/
package asciidocgo

import (
	"fmt"
	"io"
	"regexp"

	"github.com/VonC/asciidocgo/utils"
)

// Symbol name for the type of content (e.g., :paragraph).
type safeMode int

const (
	/* A safe mode level that disables any of the security features enforced
	   by Asciidocgo (Go is still subject to its own restrictions). */
	UNSAFE safeMode = iota
	/* A safe mode level that closely parallels safe mode in AsciiDoc.
	   This value prevents access to files which reside outside of the
	   parent directory of the source file and disables any macro other
	   than the include::[] macro. */
	SAFE
	/*A safe mode level that disallows the document from setting attributes
	  that would affect the rendering of the document, in addition to all the
	  security features of SafeMode::SAFE. For instance, this level disallows
	  changing the backend or the source-highlighter using an attribute defined
	  in the source document. This is the most fundamental level of security
	  for server-side deployments (hence the name).*/
	SERVER
	/*A safe mode level that disallows the document from attempting to read
	  files from the file system and including the contents of them into the
	  document, in additional to all the security features of SafeMode::SERVER.
	  For instance, this level disallows use of the include::[] macro and the
	  embedding of binary content (data uri), stylesheets and JavaScripts
	  referenced by the document.(Asciidoctor and trusted extensions may still
	  be allowed to embed trusted content into the document).

	  Since Asciidocgo is aiming for wide adoption, this level is the default
	  and is recommended for server-side deployments.*/
	SECURE
	/*A planned safe mode level that disallows the use of passthrough macros and
	  prevents the document from setting any known attributes,
	  in addition to all the security features of SafeMode::SECURE.

	  Please note that this level is not currently implemented
	  (and therefore not enforced)! */
	PARANOID
)

// Accepts input as a string
func LoadString(input string) *Document {
	return nil
}

// Accepts input as an array of strings
func LoadStrings(inputs ...string) *Document {
	return nil
}

// Accepts input as an IO.
// If the input is a File, information about the file is stored in attributes on
// the Document object.
func Load(input io.Reader) *Document {
	return nil
}

const (
	CC_ALPHA = `a-zA-Z`
	CC_ALNUM = `a-zA-Z0-9`
	CC_BLANK = `[ \t]`
	// non-blank character
	CC_GRAPH = `[\x21-\x7E]`
	CC_EOL   = `(?=\n|$)`
)

var ADMONITION_STYLES utils.Arr = []string{"NOTE", "TIP", "IMPORTANT", "WARNING", "CAUTION"}

/* The following pattern, which appears frequently, captures the contents
between square brackets, ignoring escaped closing brackets
(closing brackets prefixed with a backslash '\' character)

	Pattern:
	(?:\[((?:\\\]|[^\]])*?)\])
	Matches:
	[enclosed text here] or [enclosed [text\] here]
*/
var REGEXP_STRING = map[string]string{
	//:strip_line_wise => /\A(?:\s*\n)?(.*?)\s*\z/m,

	// # NOTE: this is a inline admonition note
	//	:admonition_inline => /^(#{ADMONITION_STYLES.to_a * '|'}):#{CC_BLANK}/,
	":admonition_inline": fmt.Sprintf("^(%v):%v", ADMONITION_STYLES.Mult("|"), CC_BLANK),

	//	http://domain
	//	https://domain
	//	data:info
	//	:uri_sniff        => %r{^[#{CC_ALPHA}][#{CC_ALNUM}.+-]*:/{0,2}},
	":uri_sniff": fmt.Sprintf("^([%v][%v.+-]*:/{0,2}).*", CC_ALPHA, CC_ALNUM),
}

func iniREGEXP(regexps map[string]string) map[string]*regexp.Regexp {
	res := map[string]*regexp.Regexp{}
	for key, regexpString := range regexps {
		regexp, err := regexp.Compile(regexpString)
		if err != nil {
			panic(fmt.Sprintf("iniREGEXP should compile all REGEXP_STRING like %v: %v", regexpString, err))
		}
		res[key] = regexp
	}
	return res
}

/* The following pattern, which appears frequently, captures the contents
between square brackets, ignoring escaped closing brackets
(closing brackets prefixed with a backslash '\' character)

	Pattern:
	(?:\[((?:\\\]|[^\]])*?)\])
	Matches:
	[enclosed text here] or [enclosed [text\] here]
*/
var REGEXP = iniREGEXP(REGEXP_STRING)
