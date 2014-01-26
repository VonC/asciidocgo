package asciidocgo

/* An abstract class that provides state and methods for managing
a block of AsciiDoc content, which is a node. */
type abstractBlock struct {
	*abstractNode
	contentModel
	subs         []string
	templateName string
	blocks       []*abstractBlock
	level        int
	title        string
}

func newAbstractBlock(parent Documentable, context context) *abstractBlock {
	templateName := "block_" + context.String()
	level := -1 // there is no 'nil' for an int
	if context == document {
		level = 0
	} else if parent != nil && context != section {
		level = parent.Level()
	}
	abstractBlock := &abstractBlock{newAbstractNode(parent, context), compound, []string{}, templateName, []*abstractBlock{}, level, ""}
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
