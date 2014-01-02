package utils

import "testing"
import . "github.com/smartystreets/goconvey/convey"

func TestUtils(t *testing.T) {

	Convey("An array of string can be multipled with a separator", t, func() {
		Convey("An empty array returns an empty string", func() {
			So(Arr{}.Mult("|"), ShouldEqual, "")
		})
		Convey("An array with only one element returns its element", func() {
			So(Arr{"a"}.Mult("|"), ShouldEqual, "a")
		})
		Convey("An array with multiple element returns those elements separated by the separator", func() {
			So(Arr{"a", "b"}.Mult("|"), ShouldEqual, "a|b")
		})
	})
}
