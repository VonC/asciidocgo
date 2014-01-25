package asciidocgo

/* An abstract class that provides state and methods for managing
a block of AsciiDoc content, which is a node. */
type abstractBlock struct {
	*abstractNode
	contentModel
}

func newAbstractBlock(parent Documentable, context context) *abstractBlock {
	abstractBlock := &abstractBlock{newAbstractNode(parent, context), compound}
	return abstractBlock
}
