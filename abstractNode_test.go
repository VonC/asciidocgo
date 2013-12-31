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
		Convey("If context is not document, then parent is parent and document is parent document", func() {
			parent := &abstractNode{nil, document, &abstractNode{}, nil, &substitutors{}}
			an := newAbstractNode(parent, section)
			So(an.Context(), ShouldEqual, section)
			So(an.Parent(), ShouldEqual, parent)
			So(an.Document(), ShouldEqual, parent.Document())
		})
		Convey("If context is not document, and parent is nil, then document is nil", func() {
			an := newAbstractNode(nil, section)
			So(an.Context(), ShouldEqual, section)
			So(an.Parent(), ShouldBeNil)
			So(an.Document(), ShouldBeNil)
		})

		Convey("An abstractNode has an empty attributes map", func() {
			an := newAbstractNode(nil, section)
			So(len(an.Attributes()), ShouldEqual, 0)
		})
	})

	Convey("An abstractNode can be associated to a parent", t, func() {
		an := newAbstractNode(nil, section)
		documentParent := &abstractNode{}
		parent := &abstractNode{nil, document, documentParent, nil, &substitutors{}}
		an.SetParent(parent)
		So(an.Parent(), ShouldEqual, parent)
		So(an.Document(), ShouldEqual, parent.Document())
	})
}
