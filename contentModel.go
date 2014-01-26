package asciidocgo

// Symbol name for the type of content (e.g., :paragraph).
type contentModel int

const (
	compound contentModel = iota
	verse
	verbatim
	simple
	unknowncm
)

func (cm contentModel) String() string {
	switch cm {
	case compound:
		return "compound"
	case verse:
		return "verse"
	case verbatim:
		return "verbatim"
	case simple:
		return "simple"
	}
	return "unknowncm"
}
