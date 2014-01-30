package asciidocgo

// Symbol name for the type of content (e.g., :paragraph).
type context int

const (
	document context = iota
	section
	paragraph
	unknown
)

func (c context) String() string {
	switch c {
	case document:
		return "document"
	case section:
		return "section"
	case paragraph:
		return "paragraph"
	}
	return "unknown"
}
