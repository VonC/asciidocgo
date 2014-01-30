package context

// Symbol name for the type of content (e.g., :paragraph).
type Context int

const (
	Document Context = iota
	Section
	Paragraph
	Unknown
)

func (c Context) String() string {
	switch c {
	case Document:
		return "document"
	case Section:
		return "section"
	case Paragraph:
		return "paragraph"
	}
	return "unknown"
}
