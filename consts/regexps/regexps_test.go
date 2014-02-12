package regexps

import (
	"regexp"
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
		So(PassInlineMacroRx.MatchString(`$$text
			mulutple
			line$$`), ShouldBeTrue)
		So(PassInlineMacroRx.MatchString(`pass:quotes[text]`), ShouldBeTrue)
		So(PassInlineMacroRx.MatchString(`pass:quotes[text
			line2
			line3]`), ShouldBeTrue)
	})

	Convey("Regexps can detect strings that resemble URIs", t, func() {
		So(UriSniffRx.MatchString("http://domain"), ShouldBeTrue)
		So(UriSniffRx.MatchString("https://domain"), ShouldBeTrue)
		So(UriSniffRx.MatchString("data:info"), ShouldBeTrue)
	})

	Convey("Regexps can detect escaped brackets", t, func() {
		So(EscapedBracketRx.MatchString(`\]`), ShouldBeTrue)
		So(EscapedBracketRx.MatchString(`a\\]a`), ShouldBeTrue)
	})

	Convey("Regexps can encapsulate results in a struct Reres", t, func() {
		testRx, _ := regexp.Compile("\\\\?a(b*)c")
		r := NewReres("xxxabbbbcyyy111aabbbcc222\\ac33", testRx)

		Convey("Regexps can create a Reres struct", func() {
			So(r, ShouldNotBeNil)
		})
		Convey("Regexps can test for matches", func() {
			So(r.HasAnyMatch(), ShouldBeTrue)
		})

		Convey("Regexps can iterate over matches", func() {
			So(r.HasNext(), ShouldBeTrue)
			r.Next()
			So(r.HasNext(), ShouldBeTrue)
			r.Next()
			So(r.HasNext(), ShouldBeTrue)
			r.Next()
			So(r.HasNext(), ShouldBeFalse)
			r.ResetNext()
		})

		Convey("Regexps can get the prefix, string before each match", func() {
			So(r.Prefix(), ShouldEqual, "xxx")
			r.Next()
			So(r.Prefix(), ShouldEqual, "yyy111a")
			r.ResetNext()
		})

		Convey("Regexps can get the suffix, string after each match", func() {
			So(r.Suffix(), ShouldEqual, "xxxabbbbcyyy111aabbbcc222\\ac33")
			r.Next()
			So(r.Suffix(), ShouldEqual, "yyy111aabbbcc222\\ac33")
			r.Next()
			So(r.Suffix(), ShouldEqual, "c222\\ac33")
			r.Next()
			So(r.Suffix(), ShouldEqual, "33")
			r.ResetNext()
		})

		Convey("Regexps can get the first character of the current match", func() {
			So(r.FirstChar(), ShouldEqual, 'a')
		})

		Convey("Regexps can get the full string of the current match", func() {
			So(r.FullMatch(), ShouldEqual, "abbbbc")
			r.Next()
			So(r.FullMatch(), ShouldEqual, "abbbc")
			r.Next()
			So(r.FullMatch(), ShouldEqual, "\\ac")
			r.ResetNext()
		})

		Convey("Regexps can test for a group within the current match", func() {
			//fmt.Printf("t %v (%d) %v => %d\n", r.matches, r.i, r.matches[r.i], r.previous)
			So(r.HasGroup(1), ShouldBeTrue)
			So(r.HasGroup(2), ShouldBeFalse)
			r.Next()
			So(r.HasGroup(1), ShouldBeTrue)
			So(r.HasGroup(2), ShouldBeFalse)
			r.Next()
			So(r.HasGroup(1), ShouldBeFalse)
			So(r.HasGroup(2), ShouldBeFalse)
			r.ResetNext()
		})

		Convey("Regexps can get a group within the current match", func() {
			So(r.Group(1), ShouldEqual, "bbbb")
			So(r.Group(2), ShouldEqual, "")
			r.Next()
			So(r.Group(1), ShouldEqual, "bbb")
			So(r.Group(2), ShouldEqual, "")
			r.Next()
			So(r.Group(1), ShouldEqual, "")
			So(r.Group(2), ShouldEqual, "")
			r.ResetNext()
		})
	})

}
