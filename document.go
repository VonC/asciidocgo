package asciidocgo

// Asciidoc Document, onced loaded from an IO, string or array
type Document struct {
	monitorData *monitorData
}

// Check if a Document is supposed to be monitored
func (d *Document) IsMonitored() bool {
	return false
}
