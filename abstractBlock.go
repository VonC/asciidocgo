package asciidocgo

/* An abstract class that provides state and methods for managing
a block of AsciiDoc content, which is a node. */
type abstractBlock struct {
	*abstractNode
	contentModel
	subs              []string
	templateName      string
	blocks            []*abstractBlock
	level             int
	title             string
	style             string
	caption           string
	nextSectionIndex  int
	nextSectionNumber int
}

func newAbstractBlock(parent Documentable, context context) *abstractBlock {
	templateName := "block_" + context.String()
	level := -1 // there is no 'nil' for an int
	if context == document {
		level = 0
	} else if parent != nil && context != section {
		level = parent.Level()
	}
	abstractBlock := &abstractBlock{newAbstractNode(parent, context), compound, []string{}, templateName, []*abstractBlock{}, level, "", "", "", 0, 1}
	return abstractBlock
}

/* The types of content that this block can accomodate */
func (ab *abstractBlock) ContentModel() contentModel {
	return ab.contentModel
}
func (ab *abstractBlock) SetContentModel(c contentModel) {
	ab.contentModel = c
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
func (ab *abstractBlock) SetContext(c context) {
	ab.context = c
	ab.templateName = "block_" + ab.Context().String()
}
