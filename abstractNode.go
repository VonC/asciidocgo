package asciidocgo

/* An abstract base class that provides state and methods for managing
a node of AsciiDoc content.
The state and methods on this class are comment to all content segments
in an AsciiDoc document. */
type abstractNode struct {
	parent     *abstractNode
	context    context
	document   *abstractNode
	attributes map[string]interface{}
	*substitutors
}

func newAbstractNode(parent *abstractNode, context context) *abstractNode {
	abstractNode := &abstractNode{parent, context, nil, make(map[string]interface{}), &substitutors{}}
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

func (an *abstractNode) Attributes() map[string]interface{} {
	return an.attributes
}

// Associate this Block with a new parent Block
// parent: The Block to set as the parent of this Block
func (an *abstractNode) SetParent(parent *abstractNode) {
	an.parent = parent
	an.document = parent.Document()
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
