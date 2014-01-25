package asciidocgo

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestContext(t *testing.T) {

	Convey("An context has a string", t, func() {
		So(document.String(), ShouldEqual, "document")
		So(section.String(), ShouldEqual, "section")
		So(paragraph.String(), ShouldEqual, "paragraph")
		So(unknown.String(), ShouldEqual, "unknown")
	})

}
