package context

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestContext(t *testing.T) {

	Convey("An context has a string", t, func() {
		So(Document.String(), ShouldEqual, "document")
		So(Section.String(), ShouldEqual, "section")
		So(Paragraph.String(), ShouldEqual, "paragraph")
		So(Unknown.String(), ShouldEqual, "unknown")
	})

}
