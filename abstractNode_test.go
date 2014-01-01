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

	Convey("An abstractNode can retrieve an attribute", t, func() {

		parentDocument := newAbstractNode(nil, document)
		parentDocument.setAttr("key", "val1", true)
		parent := newAbstractNode(parentDocument, document)
		an := newAbstractNode(nil, document)
		Convey("If inherited, it is the attribute if there, or the document attribute, or default value", func() {
			So(an.Attr("key", nil, true), ShouldBeNil)
			an.setAttr("key", "val", false)
			So(an.Attr("key", nil, true), ShouldEqual, "val")
			an.SetParent(parent)
			So(an.Attr("key", nil, true), ShouldEqual, "val")
			delete(an.attributes, "key")
			So(an.Attr("key", nil, true), ShouldEqual, "val1")
			// an should have for parent a child, which has an as a document
			// then an.document would be "parent".document, meaning an, when
			// setting an.setParent(child)
			an.document = an
			an.setAttr("key", "val", false)
			So(an.Attr("key", nil, true), ShouldEqual, "val")
		})
		Convey("If not inherited, it is the attribute if there, or default value", func() {
			an.SetParent(parent)
			So(an.Document(), ShouldEqual, parentDocument)
			So(an.Document(), ShouldEqual, parent.Document())
			So(parentDocument.Attr("key", nil, false), ShouldEqual, "val1")
			So(an.Attr("key", nil, false), ShouldEqual, "val")
		})
	})
	Convey("An abstractNode can set an attribute", t, func() {
		an := newAbstractNode(nil, document)
		an.setAttr("key", "val", true)
		So(an.Attr("key", nil, true), ShouldEqual, "val")
		Convey("If not override and already present, the value should not change", func() {
			res := an.setAttr("key", "val1", false)
			So(an.Attr("key", nil, true), ShouldEqual, "val")
			So(res, ShouldBeFalse)
		})
	})

	Convey("An abstractNode can set an option attribute", t, func() {
		an := newAbstractNode(nil, document)
		Convey("First option means options attributes has len 1", func() {
			an.SetOption("opt1")
			So(len(an.attributes["options"].(map[string]bool)), ShouldEqual, 1)
			So(an.Attr("opt1-option", nil, false), ShouldEqual, true)
		})
		Convey("Second option means options attributes has len 2", func() {
			an.SetOption("opt2")
			So(len(an.attributes["options"].(map[string]bool)), ShouldEqual, 2)
			So(an.Attr("opt2-option", nil, false), ShouldEqual, true)
		})
	})

	Convey("An abstractNode can get an option attribute", t, func() {
		an := newAbstractNode(nil, document)
		Convey("Zero option means Option returns false", func() {
			So(an.Option("opt1"), ShouldBeFalse)
			an.SetOption("opt1")
		})
		Convey("One option means Option returns true", func() {
			So(an.Option("opt1"), ShouldBeTrue)
		})
	})

	Convey("An abstractNode update option attributes with other attributes", t, func() {
		an := newAbstractNode(nil, document)
		an.setAttr("key1", "val1", true)
		an.setAttr("key2", "val2", true)
		Convey("New Attributes are added during an update", func() {
			attrs := map[string]interface{}{"key3": "val3", "key4": "val4"}
			an.Update(attrs)
			So(an.Attr("key1", nil, false), ShouldEqual, "val1")
			So(an.Attr("key2", nil, false), ShouldEqual, "val2")
			So(an.Attr("key3", nil, false), ShouldEqual, "val3")
			So(an.Attr("key4", nil, false), ShouldEqual, "val4")
		})
		Convey("Common Attributes are overrriden during an update", func() {
			attrs := map[string]interface{}{"key2": "val2b", "key3": "val3"}
			an.Update(attrs)
			So(an.Attr("key1", nil, false), ShouldEqual, "val1")
			So(an.Attr("key2", nil, false), ShouldEqual, "val2b")
			So(an.Attr("key3", nil, false), ShouldEqual, "val3")
		})
	})
}
