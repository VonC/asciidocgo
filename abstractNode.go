package asciidocgo

/* An abstract base class that provides state and methods for managing
a node of AsciiDoc content.
The state and methods on this class are comment to all content segments
in an AsciiDoc document. */
type abstractNode struct {
}

func newAbstractNode(data []string, options map[string]string) *abstractNode {
	abstractNode := &abstractNode{}
	return abstractNode
}
