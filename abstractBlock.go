package asciidocgo

/* An abstract class that provides state and methods for managing
a block of AsciiDoc content, which is a node. */
type abstractBlock struct {
	*abstractNode
	contentModel
	subs         []string
	templateName string
}

func newAbstractBlock(parent Documentable, context context) *abstractBlock {
	templateName := "block_" + context.String()
	abstractBlock := &abstractBlock{newAbstractNode(parent, context), compound, []string{}, templateName}
	return abstractBlock
}
