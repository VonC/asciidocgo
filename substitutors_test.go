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
}
