package regexps

import (
	"fmt"
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
			So(r.Text(), ShouldEqual, "xxxabbbbcyyy111aabbbcc222\\ac33")
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

		Convey("Regexps can test for escape as first charater in the current match", func() {
			So(r.IsEscaped(), ShouldBeFalse)
			r.Next()
			So(r.IsEscaped(), ShouldBeFalse)
			r.Next()
			So(r.IsEscaped(), ShouldBeTrue)
			r.ResetNext()
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

	Convey("Regexps can encapsulate PassInlineMacroRx results in a struct PassInlineMacroRxres", t, func() {
		r := NewPassInlineMacroRxres(`test \+++for
		a
		passthrough+++ by test2 $$text
			multiple
			line$$ for
			test3 pass:quotes[text
			line2
			line3] end test4`)
		So(r.HasAnyMatch(), ShouldBeTrue)

		Convey("PassInlineMacroRx can test for pass text", func() {
			So(r.HasPassText(), ShouldBeFalse)
			r.Next()
			So(r.HasPassText(), ShouldBeFalse)
			r.Next()
			So(r.HasPassText(), ShouldBeTrue)
			r.ResetNext()
		})

		Convey("PassInlineMacroRx can get pass text", func() {
			So(r.PassText(), ShouldEqual, "")
			r.Next()
			So(r.PassText(), ShouldEqual, "")
			r.Next()
			So(r.PassText(), ShouldEqual, `text
			line2
			line3`)
			r.ResetNext()
		})

		Convey("PassInlineMacroRx can test for pass sub", func() {
			So(r.HasPassSub(), ShouldBeFalse)
			r.Next()
			So(r.HasPassSub(), ShouldBeFalse)
			r.Next()
			So(r.HasPassSub(), ShouldBeTrue)
			r.ResetNext()
		})

		Convey("PassInlineMacroRx can get pass sub", func() {
			So(r.PassSub(), ShouldEqual, "")
			r.Next()
			So(r.PassSub(), ShouldEqual, "")
			r.Next()
			So(r.PassSub(), ShouldEqual, "quotes")
			r.ResetNext()
		})

		Convey("PassInlineMacroRx can get inline text", func() {
			So(r.InlineText(), ShouldEqual, `for
		a
		passthrough`)
			r.Next()
			So(r.InlineText(), ShouldEqual, `text
			multiple
			line`)
			r.Next()
			So(r.InlineText(), ShouldEqual, "")
			r.ResetNext()
		})

		Convey("PassInlineMacroRx can get inline sub", func() {
			So(r.InlineSub(), ShouldEqual, "+++")
			r.Next()
			So(r.InlineSub(), ShouldEqual, "$$")
			r.Next()
			So(r.InlineSub(), ShouldEqual, "")
			r.ResetNext()
		})

	})
	Convey("Regexps can encapsulate PassInlineLiteralRx results in a struct PassInlineLiteralRxRes", t, func() {

		r := NewPassInlineLiteralRxres(
			"`a few <\\{monospaced\\}> words`" +
				"[input]`A few <\\{monospaced\\}> words`\n" +
				"\\[input]`a few <monospaced> words`\n" +
				"\\[input]\\`a few <monospaced> words`\n" +
				"`a few\n<\\{monospaced\\}> words`" +
				"\\[input]`a few &lt;monospaced&gt; words`\n" +
				"the text `asciimath:[x = y]` should be passed through as `literal` text\n" +
				"`Here`s Johnny!")

		So(r.HasAnyMatch(), ShouldBeTrue)
		So(len(r.matches), ShouldEqual, 8)
		r.Next()
		So(r.FirstChar(), ShouldEqual, "")
		So(r.Attributes(), ShouldEqual, "input")
		So(r.Literal(), ShouldEqual, "`A few <\\{monospaced\\}> words`")
		So(r.LiteralText(), ShouldEqual, "A few <\\{monospaced\\}> words")

		r.Next()
		r.Next()
		So(r.FirstChar(), ShouldEqual, "\\")
		So(r.Attributes(), ShouldEqual, "input")
		So(r.Literal(), ShouldEqual, "\\`a few <monospaced> words`")
		So(r.LiteralText(), ShouldEqual, "a few <monospaced> words")

		r.Next()
		r.Next()
		r.Next()
		So(r.FirstChar(), ShouldEqual, " ")
		So(r.Attributes(), ShouldEqual, "")
		So(r.Literal(), ShouldEqual, "`asciimath:[x = y]`")
		So(r.LiteralText(), ShouldEqual, "asciimath:[x = y]")

		r.Next()
		So(r.FirstChar(), ShouldEqual, " ")
		So(r.Attributes(), ShouldEqual, "")
		So(r.Literal(), ShouldEqual, "`literal`")
		So(r.LiteralText(), ShouldEqual, "literal")
	})

	Convey("Regexps can encapsulate MathInlineMacroRx results in a struct MathInlineMacroRxRes", t, func() {
		r := NewMathInlineMacroRxres(`
			math:[x != 0]
   asciimath:[x != 0]
   latexmath:abc[\sqrt{4} = 2]`)

		So(r.HasAnyMatch(), ShouldBeTrue)
		So(len(r.matches), ShouldEqual, 3)
		So(r.MathType(), ShouldEqual, "math")
		So(r.MathSub(), ShouldEqual, "")
		So(r.MathText(), ShouldEqual, "x != 0")
		r.Next()
		So(r.MathType(), ShouldEqual, "asciimath")
		So(r.MathSub(), ShouldEqual, "")
		So(r.MathText(), ShouldEqual, "x != 0")
		r.Next()
		So(r.MathType(), ShouldEqual, "latexmath")
		So(r.MathSub(), ShouldEqual, "abc")
		So(r.MathText(), ShouldEqual, "\\sqrt{4} = 2")
	})

	Convey("Regexps can simulate a lookahead at the end of a regexp", t, func() {
		rx, _ := regexp.Compile(`a(b*)c($|de)`)
		r := NewReresLAGroup("aabbbbcdefabbcabbbcdeabcdabbc", rx)

		So(r.HasAnyMatch(), ShouldBeTrue)
		So(len(r.matches), ShouldEqual, 3)

		So(r.Prefix(), ShouldEqual, "a")
		So(r.FullMatch(), ShouldEqual, "abbbbc")
		So(r.Group(1), ShouldEqual, "bbbb")
		So(r.Suffix(), ShouldEqual, "defabbcabbbcdeabcdabbc")

		r.Next()

		So(r.Prefix(), ShouldEqual, "defabbc")
		So(r.FullMatch(), ShouldEqual, "abbbc")
		So(r.Group(1), ShouldEqual, "bbb")
		So(r.Suffix(), ShouldEqual, "deabcdabbc")

		r.Next()

		So(r.Prefix(), ShouldEqual, "deabcd")
		So(r.FullMatch(), ShouldEqual, "abbc")
		So(r.Group(1), ShouldEqual, "bb")
		So(r.Suffix(), ShouldEqual, "")
	})

	Convey("Regexps can encapsulate AttributeReferenceRx results in a struct AttributeReferenceRxres", t, func() {
		r := NewAttributeReferenceRxres(`
			{foo}
  {counter:pcount:1}
  {set:foo:bar}
  {set:name!}
  a\{counter:pcount:1}
  {set:foo:bar\}b
  a\{set:name!\}b`)

		So(r.HasAnyMatch(), ShouldBeTrue)
		So(len(r.matches), ShouldEqual, 7)

		So(r.PreEscaped(), ShouldBeFalse)
		So(r.Directive(), ShouldEqual, "")
		So(r.Reference(), ShouldEqual, "foo")
		So(r.PostEscaped(), ShouldBeFalse)
		r.Next()
		So(r.PreEscaped(), ShouldBeFalse)
		So(r.Directive(), ShouldEqual, "counter")
		So(r.Reference(), ShouldEqual, "counter:pcount:1")
		So(r.PostEscaped(), ShouldBeFalse)
		r.Next()
		So(r.PreEscaped(), ShouldBeFalse)
		So(r.Directive(), ShouldEqual, "set")
		So(r.Reference(), ShouldEqual, "set:foo:bar")
		So(r.PostEscaped(), ShouldBeFalse)
		r.Next()
		So(r.PreEscaped(), ShouldBeFalse)
		So(r.Directive(), ShouldEqual, "set")
		So(r.Reference(), ShouldEqual, "set:name!")
		So(r.PostEscaped(), ShouldBeFalse)
		r.Next()
		So(r.PreEscaped(), ShouldBeTrue)
		So(r.Directive(), ShouldEqual, "counter")
		So(r.Reference(), ShouldEqual, "counter:pcount:1")
		So(r.PostEscaped(), ShouldBeFalse)
		r.Next()
		So(r.PreEscaped(), ShouldBeFalse)
		So(r.Directive(), ShouldEqual, "set")
		So(r.Reference(), ShouldEqual, "set:foo:bar")
		So(r.PostEscaped(), ShouldBeTrue)
		r.Next()
		So(r.PreEscaped(), ShouldBeTrue)
		So(r.Directive(), ShouldEqual, "set")
		So(r.Reference(), ShouldEqual, "set:name!")
		So(r.PostEscaped(), ShouldBeTrue)
	})

	Convey("Regexps can replace special html characters", t, func() {
		text := "a -- b"
		for _, repl := range Replacements {
			text = repl.Rx().ReplaceAllString(text, repl.Repl())
			fmt.Sprintf("%v %v %v %v", repl.Leading(), repl.Bounding(), repl.None(), repl.EndsWithLookAhead())
		}
		So(text, ShouldEqual, "a"+rtos(8201, 8212, 8201)+"b")
	})
}
