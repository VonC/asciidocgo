package asciidocgo

import "testing"
import . "github.com/smartystreets/goconvey/convey"

func TestAsciidocgo(t *testing.T) {
	Load(nil)
	Convey("Asciidocgo takes a Reader and return a Document", t, func() {
		Convey("A nil Reader must returns a nil Document", func() {
			So(Load(nil), ShouldBeNil)
		})
	})
}
