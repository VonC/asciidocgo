package asciidocgo

import (
	"testing"

	"github.com/VonC/asciidocgo/contentModel"
	"github.com/VonC/asciidocgo/context"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAbstractBlock(t *testing.T) {

	Convey("An abstractBlock can be initialized", t, func() {
		ab := newAbstractBlock(nil, context.Document)
		Convey("By default, an AbstractBlock can be created", func() {
			So(&abstractBlock{}, ShouldNotBeNil)
			So(ab, ShouldNotBeNil)
		})
		Convey("By default, an AbstractBlock has a 'compound' content model", func() {
			So(ab.ContentModel(), ShouldEqual, contentmodel.Compound)
			ab.SetContentModel(contentmodel.Simple)
			So(ab.ContentModel(), ShouldEqual, contentmodel.Simple)
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
			So(newAbstractBlock(nil, context.Section).Level(), ShouldEqual, -1)
		})
		Convey("By default, an AbstractBlock with document context has a level of 0", func() {
			So(ab.Level(), ShouldEqual, 0)
		})
		Convey("By default, an AbstractBlock with parent of non-section context has a level of the parent", func() {
			parent := newAbstractBlock(nil, context.Document)
			parent.SetLevel(2)
			ablock := newAbstractBlock(parent, context.Paragraph)
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
		ab := newAbstractBlock(nil, context.Document)
		So(ab.Context(), ShouldEqual, context.Document)
		So(ab.TemplateName(), ShouldEqual, "block_document")
		ab.SetContext(context.Paragraph)
		So(ab.Context(), ShouldEqual, context.Paragraph)
		So(ab.TemplateName(), ShouldEqual, "block_paragraph")
	})

	Convey("An abstractBlock can render its content", t, func() {
		parent := newAbstractBlock(nil, context.Paragraph)
		ab := newAbstractBlock(parent, context.Document)
		So(ab.Render(), ShouldEqual, "")
		// TODO complete
	})

	Convey("An abstractBlock can get its content", t, func() {
		ab := newAbstractBlock(nil, context.Document)
		So(ab.Content(), ShouldEqual, "")
		ab1 := newAbstractBlock(nil, context.Document)
		ab.AppendBlock(ab1)
		So(ab.Content(), ShouldEqual, "\n")
	})

	Convey("An abstractBlock can test for blocks", t, func() {
		ab := newAbstractBlock(nil, context.Document)
		So(ab.HasBlocks(), ShouldEqual, false)
		ab1 := newAbstractBlock(nil, context.Document)
		ab.AppendBlock(ab1)
		So(ab.HasBlocks(), ShouldEqual, true)
	})

	Convey("An abstractBlock can add blocks", t, func() {
		ab := newAbstractBlock(nil, context.Document)
		So(len(ab.Blocks()), ShouldEqual, 0)
		ab1 := newAbstractBlock(nil, context.Document)
		ab.AppendBlock(ab1)
		So(len(ab.Blocks()), ShouldEqual, 1)
	})

	Convey("An abstractBlock can test for sub", t, func() {
		ab := newAbstractBlock(nil, context.Document)
		So(ab.HasSub("test"), ShouldBeFalse)
		ab.subs = []string{"a", "test", "c"}
		So(ab.HasSub("test"), ShouldBeTrue)
	})

	Convey("An abstractBlock can test for title", t, func() {
		ab := newAbstractBlock(nil, context.Document)
		So(ab.HasTitle(), ShouldBeFalse)
		ab.setTitle("a title")
		So(ab.HasTitle(), ShouldBeTrue)
	})

	Convey("An abstractBlock can get its title", t, func() {
		ab := newAbstractBlock(nil, context.Document)
		So(ab.Title(), ShouldEqual, "")
		ab.setTitle("a title")
		So(ab.Title(), ShouldEqual, "a title")
	})

	Convey("An abstractBlock can get its captioned title", t, func() {
		ab := newAbstractBlock(nil, context.Document)
		So(ab.CaptionedTitle(), ShouldEqual, "")
		ab.setTitle("a title")
		So(ab.CaptionedTitle(), ShouldEqual, "a title")
		ab.SetCaption("a caption ")
		So(ab.CaptionedTitle(), ShouldEqual, "a caption a title")
	})

	Convey("An abstractBlock can get its Sections", t, func() {
		ab := newAbstractBlock(nil, context.Document)
		So(len(ab.Sections()), ShouldEqual, 0)
		section := newAbstractBlock(nil, context.Section)
		ab.AppendBlock(section)
		So(len(ab.Sections()), ShouldEqual, 1)
	})

	Convey("An abstractBlock can remove a sub", t, func() {
		ab := newAbstractBlock(nil, context.Document)
		ab.subs = append(ab.subs, "test1")
		ab.subs = append(ab.subs, "test2")
		ab.subs = append(ab.subs, "test3")
		ab.subs = append(ab.subs, "test4")
		So(len(ab.Subs()), ShouldEqual, 4)
		ab.RemoveSub("test2")
		So(len(ab.Subs()), ShouldEqual, 3)
		So(ab.HasSub("test3"), ShouldBeTrue)
		ab.RemoveSub("test4")
		So(len(ab.Subs()), ShouldEqual, 2)
		So(ab.HasSub("test3"), ShouldBeTrue)
		So(ab.HasSub("test4"), ShouldBeFalse)
		ab.RemoveSub("test1")
		So(len(ab.Subs()), ShouldEqual, 1)
		So(ab.HasSub("test3"), ShouldBeTrue)
		So(ab.HasSub("test4"), ShouldBeFalse)
		So(ab.HasSub("test1"), ShouldBeFalse)
		ab.RemoveSub("test1")
		So(ab.HasSub("test1"), ShouldBeFalse)
	})
	Convey("An abstractBlock can assign a caption", t, func() {
		ab := newAbstractBlock(nil, context.Document)
		Convey("By default, no caption if none passed", func() {
			ab.AssignCaption("", "key")
			So(ab.CaptionedTitle(), ShouldEqual, "")
		})
		Convey("By default, no caption if title already there", func() {
			ab.setTitle("a title")
			ab.AssignCaption("a caption", "key")
			So(ab.CaptionedTitle(), ShouldEqual, "a title")
		})
	})
}
