package asciidocgo

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSubstitutorSubs(t *testing.T) {

	Convey("A sunEnum can be converted from string", t, func() {
		So(aToSE(string(subsBasic)), ShouldEqual, sub.basic)
		So(aToSE(string(subsNormal)), ShouldEqual, sub.normal)
		So(aToSE(string(subsVerbatim)), ShouldEqual, sub.verbatim)
		So(aToSE(string(subsTitle)), ShouldEqual, sub.title)
		So(aToSE(string(subsHeader)), ShouldEqual, sub.header)
		So(aToSE(string(subsPass)), ShouldEqual, sub.pass)
		So(aToSE(string(subsUnknown)), ShouldEqual, sub.unknown)
		So(aToSE("xxxtestxxx"), ShouldEqual, nil)
	})
}
