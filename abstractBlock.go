package asciidocgo

import (
	"github.com/VonC/asciidocgo/contentModel"
	"github.com/VonC/asciidocgo/context"
)

/* An abstract class that provides state and methods for managing
a block of AsciiDoc content, which is a node. */
type abstractBlock struct {
	*abstractNode
	cm                contentmodel.ContentModel
	subs              []string
	templateName      string
	blocks            []*abstractBlock
	level             int
	title             string
	style             string
	caption           string
	nextSectionIndex  int
	nextSectionNumber int
	subbedTitle       string
}

func newAbstractBlock(parent Documentable, c context.Context) *abstractBlock {
	templateName := "block_" + c.String()
	level := -1 // there is no 'nil' for an int
	if c == context.Document {
		level = 0
	} else if parent != nil && c != context.Section {
		level = parent.Level()
	}
	abstractBlock := &abstractBlock{newAbstractNode(parent, c), contentmodel.Compound, []string{}, templateName, []*abstractBlock{}, level, "", "", "", 0, 1, ""}
	return abstractBlock
}

/* The types of content that this block can accomodate */
func (ab *abstractBlock) ContentModel() contentmodel.ContentModel {
	return ab.cm
}
func (ab *abstractBlock) SetContentModel(c contentmodel.ContentModel) {
	ab.cm = c
}

/* Substitutions to be applied to content in this block */
func (ab *abstractBlock) Subs() []string {
	return ab.subs
}

/* Get/Set the String name of the render template */
func (ab *abstractBlock) TemplateName() string {
	return ab.templateName
}
func (ab *abstractBlock) SetTemplateName(tn string) {
	ab.templateName = tn
}

/* Array of Asciidoctor::AbstractBlock sub-blocks for this block */
func (ab *abstractBlock) Blocks() []*abstractBlock {
	return ab.blocks
}

/* Get/Set the Integer level of this Section or the Section level
in which this Block resides */
func (ab *abstractBlock) Level() int {
	return ab.level
}
func (ab *abstractBlock) SetLevel(l int) {
	ab.level = l
}

/* Set the String block title. */
func (ab *abstractBlock) setTitle(t string) {
	ab.title = t
}

/* Get/Set the String style (block type qualifier) for this block. */
func (ab *abstractBlock) Style() string {
	return ab.style
}
func (ab *abstractBlock) SetStyle(s string) {
	ab.style = s
}

/* Get/Set the caption for this block. */
func (ab *abstractBlock) Caption() string {
	return ab.caption
}
func (ab *abstractBlock) SetCaption(c string) {
	ab.caption = c
}

/* This method changes the context of this block.
It also updates the template name accordingly. */
func (ab *abstractBlock) SetContext(c context.Context) {
	ab.context = c
	ab.templateName = "block_" + ab.Context().String()
}

/* Get the rendered String content for this Block.
If the block has child blocks, the content method should cause them
to be rendered and returned as content that can be included
in the parent block's template. */
func (ab *abstractBlock) Render() string {
	if ab.Document() != nil {
		ab.Document().PlaybackAttributes(ab.Attributes())
	}
	return ab.Renderer().Render(ab.TemplateName(), ab, []interface{}{})
	// TODO make sure document playback_attributes and renderer hare implemented
}

/* Get an rendered version of the block content, rendering the
children appropriate to content model that this block supports. */
func (ab *abstractBlock) Content() string {
	res := ""
	for _, block := range ab.Blocks() {
		res = res + block.Render() + "\n"
	}
	return res
}

/* A convenience method that checks whether the specified substitution
is enabled for this block.
name - The Symbol substitution name */
func (ab *abstractBlock) HasSub(name string) bool {
	res := false
	for _, sub := range ab.Subs() {
		if sub == name {
			res = true
			break
		}
	}
	return res
}

/* A convenience method that indicates whether the title instance
variable is blank (nil or empty) */
func (ab *abstractBlock) HasTitle() bool {
	return (ab.title != "")
}

/* Get the String title of this Block with title substitions applied
The following substitutions are applied to block and section titles:
:specialcharacters, :quotes, :replacements, :macros, :attributes
and :post_replacements
Examples
   block.title = "Foo 3^ # {two-colons} Bar(1)"
   block.title
   => "Foo 3^ # :: Bar(1)" */
func (ab *abstractBlock) Title() string {
	//if ab.subbedTitle != "" {
	//	return ab.subbedTitle
	//}
	// TODO add substitutor as mixin in Section and Block
	//if ab.Title() != "" {
	// return applyTitleSubs(ab.Title())
	//}
	return ab.title
}

/* Convenience method that returns the interpreted title of the Block
with the caption prepended.
Concatenates the value of this Block's caption instance variable and the
return value of this Block's title method. No space is added between the
two values.
If the Block does not have a caption, the interpreted title is returned.
Returns the String title prefixed with the caption, or just the title if no
caption is set */
func (ab *abstractBlock) CaptionedTitle() string {
	return ab.caption + ab.title
}

/* Determine whether this Block contains block content
Returns A Boolean indicating whether this Block has block content */
func (ab *abstractBlock) HasBlocks() bool {
	return len(ab.Blocks()) > 0
}

/* Append a content block to this block's list of blocks.
   block - The new child block.
   Examples
     block = Block.new(parent, :preamble, :content_model => :compound)
     block << Block.new(block, :paragraph, :source => 'p1')
     block << Block.new(block, :paragraph, :source => 'p2')
     block.blocks?
     # => true
     block.blocks.size
     # => 2
 Returns nothing. */
func (ab *abstractBlock) AppendBlock(block *abstractBlock) {
	// parent assignment pending refactor
	// block.parent = self
	ab.blocks = append(ab.Blocks(), block)
}

/* Get the Array of child Section objects
Only applies to Document and Section instances
Examples
   section = Section.new(parent)
   section << Block.new(section, :paragraph, :source => 'paragraph 1')
   section << Section.new(parent)
   section << Block.new(section, :paragraph, :source => 'paragraph 2')
   section.blocks?
   # => true
   section.blocks.size
   # => 3
   section.sections.size
   # => 1
returns an Array of Section objects
*/
func (ab *abstractBlock) Sections() []*abstractBlock {
	res := []*abstractBlock{}
	for _, block := range ab.Blocks() {
		if block.Context() == context.Section {
			res = append(res, block)
		}
	}
	return res
}

/* Remove a substitution from this block
sub  - The Symbol substitution name
Returns nothing */
func (ab *abstractBlock) RemoveSub(sub string) {
	asub := ""
	i := -1
	// http://stackoverflow.com/a/18203895/6309
	for i, asub = range ab.subs {
		if asub == sub {
			break
		}
	}
	if i >= 0 {
		// http://code.google.com/p/go-wiki/wiki/SliceTricks
		ab.subs = append(ab.subs[:i], ab.subs[i+1:]...)
	}
}

/* Generate a caption and assign it to this block if one
is not already assigned.

If the block has a title and a caption prefix is available
for this block, then build a caption from this information,
assign it a number and store it to the caption attribute on
the block.

If an explicit caption has been specified on this block, then
do nothing.

key         - The prefix of the caption and counter attribute names.
              If not provided, the name of the context for this block
              is used. (default: nil).
returns nothing */
func (ab *abstractBlock) AssignCaption(caption, key string) {
	if !ab.HasTitle() && ab.Caption() != "" {
		return
	}
}
