package asciidocgo

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestPathResolver(t *testing.T) {

	Convey("A pathResolver can be initialized", t, func() {

		Convey("By default, a pathResolver can be created", func() {
			So(&abstractNode{}, ShouldNotBeNil)
		})
	})
}
