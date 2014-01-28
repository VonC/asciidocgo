package asciidocgo

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRenderer(t *testing.T) {

	Convey("A Renderer can be initialized", t, func() {
		// TODO r := newRenderer()
		Convey("By default, an Renderer can be created", func() {
			So(&Renderer{}, ShouldNotBeNil)
			// TODO So(r, ShouldNotBeNil)
		})
	})
	Convey("A Renderer can render a template", t, func() {
		Convey("Empty template means empty result", func() {
			r := &Renderer{}
			So(r.Render("", nil, nil), ShouldEqual, "")
		})
	})
}
