package asciidocgo

/* An abstract class that provides state and methods for managing
a block of AsciiDoc content, which is a node. */
type abstractBlock struct {
	*abstractNode
}

func newAbstractBlock(parent Documentable, context context) *abstractBlock {
	abstractBlock := &abstractBlock{&abstractNode{}}
	return abstractBlock
}
