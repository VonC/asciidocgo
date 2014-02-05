package asciidocgo

import (
	"testing"

	"github.com/VonC/asciidocgo/consts/contentModel"
	"github.com/VonC/asciidocgo/consts/context"
	"github.com/VonC/asciidocgo/consts/safemode"
	. "github.com/smartystreets/goconvey/convey"
)

type testSectionAble struct {
	*abstractBlock
	index    int
	number   int
	numbered bool
	name     string
	caption  string
	level    int
	special  bool
}

type testBlockDocumentAble struct {
	*abstractBlock
}

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
		parent := newTestBlockDocumentAble(nil).abstractBlock
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
		ab.SetCaption("a caption")
		Convey("By default, no caption if there is a caption already there, but no title", func() {
			ab.AssignCaption("a new caption", "key")
			So(ab.CaptionedTitle(), ShouldEqual, "a caption")
		})
		Convey("Assign caption if one is passed and a title is there", func() {
			ab.setTitle("a title")
			ab.AssignCaption("a new caption", "key")
			So(ab.CaptionedTitle(), ShouldEqual, "a new captiona title")
		})
		Convey("If title, no caption and empty document, no caption assigned", func() {
			ab.setTitle("a title")
			ab = newAbstractBlock(nil, context.Document)
			ab.AssignCaption("", "key")
			So(ab.CaptionedTitle(), ShouldEqual, "")
		})
		Convey("If title, no caption and actual document, caption assigned for key 'caption'", func() {
			parent := newTestBlockDocumentAble(nil).abstractBlock
			ab = newAbstractBlock(parent, context.Document)
			ab.setTitle("a title")
			parent.setAttr("caption", "an attr caption", false)
			ab.AssignCaption("", "key")
			So(ab.CaptionedTitle(), ShouldEqual, "an attr captiona title")
		})
		Convey("If title, no caption and actual document and no key, caption assigned for key equals to 'document context-caption'", func() {
			parent := newTestBlockDocumentAble(nil)
			ab = newAbstractBlock(parent.abstractBlock, context.Document)
			ab.setTitle("a title2")

			parent.setAttr(ab.Context().String()+"-caption", "an attr doc caption", false)
			So(ab.Document().HasAttr("caption", nil, false), ShouldBeFalse)
			So(ab.Document().HasAttr(ab.Context().String()+"-caption", nil, false), ShouldBeTrue)

			ab.AssignCaption("", "")
			So(ab.CaptionedTitle(), ShouldEqual, "an attr doc caption . a title2")
		})

	})

	Convey("An abstractBlock can assign an index to a section", t, func() {
		ab := newAbstractBlock(nil, context.Document)
		parent := newTestBlockDocumentAble(nil).abstractBlock
		ts := &testSectionAble{}
		ab.assignIndex(ts)
		ab.assignIndex(ts)
		ab.assignIndex(ts)
		So(ts.index, ShouldEqual, 2)
		Convey("Section not apendix and not numbered has no number or caption", func() {
			ab = newAbstractBlock(parent, context.Document)
			ab.assignIndex(ts)
			ab.assignIndex(ts)
			So(ts.index, ShouldEqual, 1)
			So(ts.caption, ShouldEqual, "")
			So(ts.number, ShouldEqual, 0)
		})
		Convey("An appendix Section has number only if numbered", func() {
			ts.name = "appendix"
			ab.assignIndex(ts)
			So(ts.number, ShouldEqual, 0)
			ts.numbered = true
			ab.assignIndex(ts)
			So(ts.number, ShouldEqual, -1)
			So(ts.caption, ShouldEqual, "-1. ")
		})
		Convey("An appendix Section has caption only if document has appendix-caption attribute", func() {
			parent.setAttr("appendix-caption", "an appendix CAPTION", false)
			ab.assignIndex(ts)
			So(ts.caption, ShouldEqual, "an appendix CAPTION -1: ")
		})
		Convey("An non-appendix Section with non-book document has number equals to nextSectionNumber", func() {
			ts.name = ""
			So(ts.number, ShouldEqual, -1)
			ab.assignIndex(ts)
			So(ts.number, ShouldEqual, 1)
			ab.assignIndex(ts)
			So(ts.number, ShouldEqual, 2)
		})
		Convey("An non-appendix Section with level 1 and book document should have counter number", func() {
			ts.number = 0
			ts.level = 1
			testab = "test_doctypeBook_assignIndex"
			ab.assignIndex(ts)
			testab = ""
			So(ts.number, ShouldEqual, -1)
		})
	})

	Convey("An abstractBlock can reindex sections", t, func() {
		parent := newTestBlockDocumentAble(nil).abstractBlock
		ab := newAbstractBlock(parent, context.Document)
		So(len(ab.Sections()), ShouldEqual, 0)
		section1 := newTestSection(nil, context.Section)
		section2 := newTestSection(nil, context.Section)
		section2.numbered = true
		section3 := newTestSection(nil, context.Section)
		section3.numbered = true
		ab.AppendBlock(section1.abstractBlock)
		ab.AppendBlock(section2.abstractBlock)
		ab.AppendBlock(section3.abstractBlock)
		//ab.Section()
		ab.nextSectionIndex = -1
		ab.nextSectionNumber = -1
		ab.reindexSections()
		So(ab.nextSectionIndex, ShouldEqual, 3)
		So(ab.nextSectionNumber, ShouldEqual, 2)
	})
}

func newTestSection(parent *abstractBlock, c context.Context) *testSectionAble {
	ab := newAbstractBlock(parent, context.Section)
	testSectionAble := &testSectionAble{ab, 0, 0, false, "", "", 0, false}
	ab.MainSectionAble(testSectionAble)
	//fmt.Printf("testSectionAble '%v' => '%v' => '%v' vs. '%v'\n", testSectionAble, testSectionAble.abstractBlock, testSectionAble.Section(), testSectionAble.abstractBlock.Section())
	return testSectionAble
}

func (ts *testSectionAble) SetIndex(index int) {
	ts.index = index
}

func (ts *testSectionAble) SectName() string {
	return ts.name
}
func (ts *testSectionAble) SetNumber(number int) {
	ts.number = number
}
func (ts *testSectionAble) IsNumbered() bool {
	return ts.numbered
}
func (ts *testSectionAble) SetCaption(caption string) {
	ts.caption = caption
}
func (ts *testSectionAble) Level() int {
	return ts.level
}
func (ts *testSectionAble) IsSpecial() bool {
	return ts.special
}
func (ts *testSectionAble) Section() sectionAble {
	return ts
}

func newTestBlockDocumentAble(parent *abstractBlock) *testBlockDocumentAble {
	an := newAbstractBlock(parent, context.Document)
	tba := &testBlockDocumentAble{an}
	an.MainDocumentable(tba)
	return tba
}

func (tbd *testBlockDocumentAble) Safe() safemode.SafeMode {
	return safemode.UNSAFE
}

func (tbd *testBlockDocumentAble) BaseDir() string {
	return ""
}

func (tbd *testBlockDocumentAble) PlaybackAttributes(map[string]interface{}) {
	//
}

func (tbd *testBlockDocumentAble) CounterIncrement(counterName string, block *abstractNode) string {
	return ""
}

func (tbd *testBlockDocumentAble) Counter(name, seed string) int {
	return -1
}

func (tbd *testBlockDocumentAble) DocType() string {
	return ""
}
