package contentmodel

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestContentModel(t *testing.T) {

	Convey("A content model has a string", t, func() {
		So(Compound.String(), ShouldEqual, "compound")
		So(Verse.String(), ShouldEqual, "verse")
		So(Verbatim.String(), ShouldEqual, "verbatim")
		So(Simple.String(), ShouldEqual, "simple")
		So(UnknownCM.String(), ShouldEqual, "unknowncm")
	})

}
