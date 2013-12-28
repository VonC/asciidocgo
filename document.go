package asciidocgo

// Asciidoc Document, onced loaded from an IO, string or array
type Document struct {
}

func (d *Document) isMonitored() bool {
	return false
}
