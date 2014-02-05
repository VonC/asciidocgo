package contentmodel

// Symbol name for the type of content (e.g., :paragraph).
type ContentModel int

const (
	Compound ContentModel = iota
	Verse
	Verbatim
	Simple
	UnknownCM
)

func (cm ContentModel) String() string {
	switch cm {
	case Compound:
		return "compound"
	case Verse:
		return "verse"
	case Verbatim:
		return "verbatim"
	case Simple:
		return "simple"
	}
	return "unknowncm"
}
