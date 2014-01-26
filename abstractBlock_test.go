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
	})

}
