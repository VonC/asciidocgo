package asciidocgo

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestContentModel(t *testing.T) {

	Convey("A content model has a string", t, func() {
		So(compound.String(), ShouldEqual, "compound")
		So(verse.String(), ShouldEqual, "verse")
		So(verbatim.String(), ShouldEqual, "verbatim")
		So(simple.String(), ShouldEqual, "simple")
		So(unknowncm.String(), ShouldEqual, "unknowncm")
	})

}
