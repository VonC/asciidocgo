package asciidocgo

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAbstractBlock(t *testing.T) {

	Convey("An abstractNode can be initialized", t, func() {

		Convey("By default, an AbstractBlock can be created", func() {
			So(&abstractBlock{}, ShouldNotBeNil)
			So(newAbstractBlock(nil, document), ShouldNotBeNil)
		})

	})

}
