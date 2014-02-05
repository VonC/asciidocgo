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

	Convey("Regexps can detect strings that resemble URIs", t, func() {
		So(UriSniffRx.MatchString("http://domain"), ShouldBeTrue)
		So(UriSniffRx.MatchString("https://domain"), ShouldBeTrue)
		So(UriSniffRx.MatchString("data:info"), ShouldBeTrue)
	})
}
