package asciidocgo

// Symbol name for the type of content (e.g., :paragraph).
type contentModel int

const (
	compound contentModel = iota
	verse
	verbatim
	simple
)
