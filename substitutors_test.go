package asciidocgo

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSubstitutor(t *testing.T) {

	Convey("A substitutors can be initialized", t, func() {

		Convey("By default, a substitutors can be created", func() {
			So(&substitutors{}, ShouldNotBeNil)
		})

		Convey("A substitutors has an empty passthroughs array", func() {
			s := substitutors{}
			So(len(s.passthroughs), ShouldEqual, 0)
		})
	})

	Convey("A substitutors has subs type", t, func() {
		subs := newSubs()
		So(len(subs.basic), ShouldEqual, 1)
		So(len(subs.normal), ShouldEqual, 6)
		So(len(subs.verbatim), ShouldEqual, 2)
		So(len(subs.title), ShouldEqual, 5)
		So(len(subs.header), ShouldEqual, 2)
		So(len(subs.pass), ShouldEqual, 0)
	})
}
