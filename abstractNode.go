package asciidocgo

/* An abstract base class that provides state and methods for managing
a node of AsciiDoc content.
The state and methods on this class are comment to all content segments
in an AsciiDoc document. */
type abstractNode struct {
	parent   *abstractNode
	context  context
	document *abstractNode
}

func newAbstractNode(parent *abstractNode, context context) *abstractNode {
	abstractNode := &abstractNode{parent, context, nil}
	if context == document {
		abstractNode.parent = nil
		abstractNode.document = parent
	} else if parent != nil {
		abstractNode.document = parent.Document()
	}
	return abstractNode
}

//  Get the element which is the parent of this node
func (an *abstractNode) Parent() *abstractNode {
	return an.parent
}

//  Get the Asciidoctor::Document to which this node belongs
func (an *abstractNode) Document() *abstractNode {
	return an.document
}

// Get the Symbol context for this node
func (an *abstractNode) Context() context {
	return an.context
}
