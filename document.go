package asciidocgo

// Asciidoc Document, onced loaded from an IO, string or array
type Document struct {
	monitorData *monitorData
}

type monitorData struct {
	readTime  int
	parseTime int
}

// Error returned when accessing times on a Document not monitored
type NotMonitoredError struct {
	msg string // description of error
}

// Print description of a non-monitored error
func (e *NotMonitoredError) Error() string { return e.msg }

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
		return 0, &NotMonitoredError{"No readTime: current document is not monitored"}
	}
	return d.monitorData.readTime, nil
}

func (d *Document) ParseTime() (parseTime int, err error) {
	if d.IsMonitored() == false {
		return 0, &NotMonitoredError{"No parseTime: current document is not monitored"}
	}
	return d.monitorData.parseTime, nil
}

func (d *Document) LoadTime() (loadTime int, err error) {
	if d.IsMonitored() == false {
		return 0, &NotMonitoredError{"No loadTime: current document is not monitored"}
	}
	readTime, _ := d.ReadTime()
	parseTime, _ := d.ParseTime()
	return readTime + parseTime, nil
}
