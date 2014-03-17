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
		So(Kbd.String(), ShouldEqual, "kbd")
		So(Button.String(), ShouldEqual, "button")
		So(Menu.String(), ShouldEqual, "menu")
		So(Image.String(), ShouldEqual, "image")
		So(IndexTerm.String(), ShouldEqual, "indexterm")
		So(Anchor.String(), ShouldEqual, "anchor")
		So(Footnote.String(), ShouldEqual, "footnote")
		So(Unknown.String(), ShouldEqual, "unknown")
	})

}
