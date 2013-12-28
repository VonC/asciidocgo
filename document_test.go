package asciidocgo

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

var dm = new(Document).Monitor()
var dnm = new(Document)

func TestDocumentMonitor(t *testing.T) {
	Convey("A Document can be monitored", t, func() {
		Convey("By default, a Document is not monitored", func() {
			So(dnm.IsMonitored(), ShouldBeFalse)
		})
		Convey("A monitored Document is monitored", func() {
			So(dm.IsMonitored(), ShouldBeTrue)
		})
	})
	Convey("A non-monitored Document should return error when accessing times", t, func() {
		_, err := dnm.ReadTime()
		So(err, ShouldNotBeNil)
	})
}
