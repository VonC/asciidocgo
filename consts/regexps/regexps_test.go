package regexps

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRegexps(t *testing.T) {

	Convey("Regexps can match an admonition label at the start of a paragraph", t, func() {
		So(AdmonitionParagraphRx.MatchString("NOTE: Just a little note."), ShouldBeTrue)
		So(AdmonitionParagraphRx.MatchString("TIP: Don't forget!"), ShouldBeTrue)
	})

	Convey("Regexps can match several variants of the passthrough inline macro, which may span multiple lines", t, func() {
		So(PassInlineMacroRx.MatchString("+++text+++"), ShouldBeTrue)
		So(PassInlineMacroRx.MatchString(`+++text
			line2
			line3+++`), ShouldBeTrue)
		So(PassInlineMacroRx.MatchString("$$text$$"), ShouldBeTrue)
		So(PassInlineMacroRx.MatchString("$$text$$"), ShouldBeTrue)
		So(PassInlineMacroRx.MatchString(`pass:quotes[text
			line2
			line3]`), ShouldBeTrue)
		So(PassInlineMacroRx.MatchString(`pass:quotes[text
			line2
			line3]`), ShouldBeTrue)
	})

	Convey("Regexps can detect strings that resemble URIs", t, func() {
		So(UriSniffRx.MatchString("http://domain"), ShouldBeTrue)
		So(UriSniffRx.MatchString("https://domain"), ShouldBeTrue)
		So(UriSniffRx.MatchString("data:info"), ShouldBeTrue)
	})
}
