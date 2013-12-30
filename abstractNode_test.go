package asciidocgo

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAbstractNode(t *testing.T) {

	Convey("An abstractNode can be initialized", t, func() {

		Convey("By default, an AbstractNode can be created", func() {
			So(&abstractNode{}, ShouldNotBeNil)
		})
		Convey("An AbstractNode takes a parent and a context", func() {
			So(newAbstractNode(nil, document), ShouldNotBeNil)
		})
		Convey("If context is a document, then parent is nil and document is parent", func() {
			parent := &abstractNode{}
			an := newAbstractNode(parent, document)
			So(an.Context(), ShouldEqual, document)
			So(an.Parent(), ShouldBeNil)
			So(an.Document(), ShouldEqual, parent)
		})
	})
}
