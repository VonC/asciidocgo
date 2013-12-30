package asciidocgo

// Symbol name for the type of content (e.g., :paragraph).
type context int

const (
	document context = iota
	section
	paragraph
)
