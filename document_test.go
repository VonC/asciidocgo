package asciidocgo

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDocument(t *testing.T) {
	Convey("A Document can be monitored", t, func() {
		Convey("By default, a Document is not monitored", func() {
			So(new(Document).isMonitored(), ShouldBeFalse)
		})
	})
}
