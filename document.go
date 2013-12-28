package asciidocgo

// Asciidoc Document, onced loaded from an IO, string or array
type Document struct {
	monitorData *monitorData
}

type monitorData struct{}

// Check if a Document is supposed to be monitored
func (d *Document) IsMonitored() bool {
	return false
}

func (d *Document) Monitor() *Document {
	if d.monitorData != nil {
		d.monitorData = new(monitorData)
	}
	return d
}
