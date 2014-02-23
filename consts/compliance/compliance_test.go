package compliance

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRegexps(t *testing.T) {

	Convey("Complance has default values", t, func() {
		So(BlockTerminatesParagraph(), ShouldBeTrue)
		So(StrictVerbatimParagraphs(), ShouldBeTrue)
		So(UnderlineStyleSectionTitles(), ShouldBeTrue)
		So(UnwrapStandalonePreamble(), ShouldBeTrue)
		So(AttributeMissing(), ShouldEqual, "skip")
		So(AttributeUndefined(), ShouldEqual, "drop-line")
		So(MarkdownSyntax(), ShouldBeTrue)
	})
}
