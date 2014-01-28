package asciidocgo

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAbstractBlock(t *testing.T) {

	Convey("An abstractBlock can be initialized", t, func() {
		ab := newAbstractBlock(nil, document)
		Convey("By default, an AbstractBlock can be created", func() {
			So(&abstractBlock{}, ShouldNotBeNil)
			So(ab, ShouldNotBeNil)
		})
		Convey("By default, an AbstractBlock has a 'compound' content model", func() {
			So(ab.ContentModel(), ShouldEqual, compound)
			ab.SetContentModel(simple)
			So(ab.ContentModel(), ShouldEqual, simple)
		})
		Convey("By default, an AbstractBlock has no subs", func() {
			So(len(ab.Subs()), ShouldEqual, 0)
		})
		Convey("By default, an AbstractBlock has a template name equals to block_#{context}", func() {
			So(ab.TemplateName(), ShouldEqual, "block_"+ab.Context().String())
			ab.SetTemplateName("aTemplateName")
			So(ab.TemplateName(), ShouldEqual, "aTemplateName")
		})
		Convey("By default, an AbstractBlock has no blocks", func() {
			So(len(ab.Blocks()), ShouldEqual, 0)
		})
		Convey("By default, an AbstractBlock with no document context and no parent has a level of -1", func() {
			So(newAbstractBlock(nil, section).Level(), ShouldEqual, -1)
		})
		Convey("By default, an AbstractBlock with document context has a level of 0", func() {
			So(ab.Level(), ShouldEqual, 0)
		})
		Convey("By default, an AbstractBlock with parent of non-section context has a level of the parent", func() {
			parent := newAbstractBlock(nil, document)
			parent.SetLevel(2)
			ablock := newAbstractBlock(parent, paragraph)
			So(ablock.Level(), ShouldEqual, 2)
		})
		Convey("By default, an AbstractBlock has an empty title", func() {
			So(ab.title, ShouldEqual, "")
			ab.setTitle("a title")
			So(ab.title, ShouldEqual, "a title")
		})
		Convey("By default, an AbstractBlock has an empty style", func() {
			So(ab.Style(), ShouldEqual, "")
			ab.SetStyle("a style")
			So(ab.Style(), ShouldEqual, "a style")
		})
		Convey("By default, an AbstractBlock has an empty caption", func() {
			So(ab.Caption(), ShouldEqual, "")
			ab.SetCaption("a caption")
			So(ab.Caption(), ShouldEqual, "a caption")
		})
	})

	Convey("An abstractBlock can set its context", t, func() {
		ab := newAbstractBlock(nil, document)
		So(ab.Context(), ShouldEqual, document)
		So(ab.TemplateName(), ShouldEqual, "block_document")
		ab.SetContext(paragraph)
		So(ab.Context(), ShouldEqual, paragraph)
		So(ab.TemplateName(), ShouldEqual, "block_paragraph")
	})

	Convey("An abstractBlock can render its content", t, func() {
		parent := newAbstractBlock(nil, paragraph)
		ab := newAbstractBlock(parent, document)
		So(ab.Render(), ShouldEqual, "")
		// TODO complete
	})

	Convey("An abstractBlock can add blocks", t, func() {
		ab := newAbstractBlock(nil, document)
		So(len(ab.Blocks()), ShouldEqual, 0)
		ab1 := newAbstractBlock(nil, document)
		ab.AppendBlock(ab1)
		So(len(ab.Blocks()), ShouldEqual, 1)
	})
}
