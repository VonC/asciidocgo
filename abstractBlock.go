package asciidocgo

/* An abstract class that provides state and methods for managing
a block of AsciiDoc content, which is a node. */
type abstractBlock struct {
	*abstractNode
	contentModel
	subs         []string
	templateName string
	blocks       []*abstractBlock
}

func newAbstractBlock(parent Documentable, context context) *abstractBlock {
	templateName := "block_" + context.String()
	abstractBlock := &abstractBlock{newAbstractNode(parent, context), compound, []string{}, templateName, []*abstractBlock{}}
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
