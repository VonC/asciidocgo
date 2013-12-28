package asciidocgo

import "errors"

// Asciidoc Document, onced loaded from an IO, string or array
type Document struct {
	monitorData *monitorData
}

type monitorData struct {
	readTime int
}

// Check if a Document is supposed to be monitored
func (d *Document) IsMonitored() bool {
	return (d.monitorData != nil)
}

// Setup a monitor for the document
// (or does nothing if the monitor already exists).
// Returns self, for easy composition
func (d *Document) Monitor() *Document {
	if d.monitorData == nil {
		d.monitorData = new(monitorData)
	}
	return d
}

func (d *Document) ReadTime() (readTime int, err error) {
	if d.IsMonitored() == false {
		return 0, errors.New("z")
	}
	return d.monitorData.readTime, nil
}
