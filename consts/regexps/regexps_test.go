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
	Convey("Regexps can simulate a lookahead at the end of a regexp choice", t, func() {
		rx, _ := regexp.Compile(`a(b*)c|d(e*)([^f])`)
		r := NewReresLAGroup("aabbbbcdaefdeeeabbc", rx)

		So(r.HasAnyMatch(), ShouldBeTrue)
		So(len(r.matches), ShouldEqual, 4)
		//fmt.Println(r)

		So(r.Prefix(), ShouldEqual, "a")
		So(r.FullMatch(), ShouldEqual, "abbbbc")
		So(r.Group(1), ShouldEqual, "bbbb")
		So(r.Suffix(), ShouldEqual, "daefdeeeabbc")

		r.Next()

		So(r.Prefix(), ShouldEqual, "")
		So(r.FullMatch(), ShouldEqual, "d")
		So(r.Group(1), ShouldEqual, "")
		So(r.Suffix(), ShouldEqual, "aefdeeeabbc")

		r.Next()

		So(r.Prefix(), ShouldEqual, "aef")
		So(r.FullMatch(), ShouldEqual, "deee")
		So(r.Group(2), ShouldEqual, "eee")
		So(r.Suffix(), ShouldEqual, "abbc")

		r.Next()

		So(r.Prefix(), ShouldEqual, "")
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
		text1 := "a -- b"
		text2 := text1
		for _, repl := range Replacements {
			text1 = repl.Rx().ReplaceAllString(text1, repl.Repl())
			fmt.Sprintf("%v %v %v %v", repl.Leading(), repl.Bounding(), repl.None(), repl.EndsWithLookAhead())
			reres := repl.Reres(text2)
			if reres.HasAnyMatch() {
				text2 = reres.Prefix() + repl.Repl() + reres.Suffix()
			}
		}
		So(text1, ShouldEqual, "a"+Rtos(8201, 8212, 8201)+"b")
		So(text2, ShouldEqual, "a"+Rtos(8201, 8212, 8201)+"b")
	})

	Convey("Regexps can encapsulate KbdBtnInlineMacroRx results in a struct KbdBtnInlineMacroRxres", t, func() {
		r := NewKbdBtnInlineMacroRxres(`
   kbd:[F3]
   kbd:[Ctrl+Shift+T]
   kbd:[Ctrl+\]]
   kbd:[Ctrl,T]
   btn:[Save]`)

		So(r.HasAnyMatch(), ShouldBeTrue)
		So(len(r.matches), ShouldEqual, 5)

		So(r.Key(), ShouldEqual, "F3")
		r.Next()
		So(r.Key(), ShouldEqual, "Ctrl+Shift+T")
		r.Next()
		So(r.Key(), ShouldEqual, `Ctrl+\]`)
		r.Next()
		So(r.Key(), ShouldEqual, "Ctrl,T")
		r.Next()
		So(r.Key(), ShouldEqual, "Save")
		r.Next()
	})

	Convey("Regexps can encapsulate KbdDelimiterRx results in a struct KbdDelimiterRxres", t, func() {
		r := NewKbdDelimiterRxres(`Ctrl + Alt+T`)

		So(r.HasAnyMatch(), ShouldBeTrue)
		So(len(r.matches), ShouldEqual, 2)

		So(r.FullMatch(), ShouldEqual, "+")
		r.Next()
		So(r.FullMatch(), ShouldEqual, "+")

		r = NewKbdDelimiterRxres(`
   Ctrl,T`)

		So(r.HasAnyMatch(), ShouldBeTrue)
		So(len(r.matches), ShouldEqual, 1)
		So(r.FullMatch(), ShouldEqual, ",")

		r = NewKbdDelimiterRxres(`Ctrl,  ,`)
		So(r.HasAnyMatch(), ShouldBeFalse)
		r = NewKbdDelimiterRxres(`Ctrl +  +a`)
		So(len(r.matches), ShouldEqual, 1)
		So(r.FullMatch(), ShouldEqual, "+")
	})

	Convey("Regexps can encapsulate MenuInlineMacroRx results in a struct MenuInlineMacroRxres", t, func() {
		r := NewMenuInlineMacroRxres(`menu:File[New...]
   menu:View[Page Style > No Style]
   menu:View[Page Style, No Style]`)
		So(r.HasAnyMatch(), ShouldBeTrue)
		So(len(r.matches), ShouldEqual, 3)

		So(r.MenuName(), ShouldEqual, "File")
		So(r.MenuItems(), ShouldEqual, "New...")
		r.Next()
		So(r.MenuName(), ShouldEqual, "View")
		So(r.MenuItems(), ShouldEqual, "Page Style > No Style")
		r.Next()
		So(r.MenuName(), ShouldEqual, "View")
		So(r.MenuItems(), ShouldEqual, "Page Style, No Style")
	})

	Convey("Regexps can encapsulate MenuInlineRx results in a struct MenuInlineRxres", t, func() {
		r := NewMenuInlineRxres(`menu \"File &gt; New"
			menu "File1 &gt; New1" test
			menu "File2 &gt; New2   &gt;    Item2"`)
		So(r.HasAnyMatch(), ShouldBeTrue)
		So(len(r.matches), ShouldEqual, 3)

		So(r.MenuInput(), ShouldEqual, "File &gt; New")
		r.Next()
		So(r.MenuInput(), ShouldEqual, "File1 &gt; New1")
		r.Next()
		So(r.MenuInput(), ShouldEqual, "File2 &gt; New2   &gt;    Item2")
	})

	Convey("Regexps can encapsulate ImageInlineMacroRx results in a struct ImageInlineMacroRxres", t, func() {
		r := NewImageInlineMacroRxres(`\image:filename1.png[Alt Text]
   image:filename2.png[Alt2 Text2]
   image:http://example.com/images/filename3.png[Alt3 Text3]
   image:filename4.png[More4 [Alt4\] Text4] (alt text becomes "More [Alt] Text")
   icon:github[large]`)
		So(r.HasAnyMatch(), ShouldBeTrue)
		So(len(r.matches), ShouldEqual, 5)

		So(r.ImageTarget(), ShouldEqual, "filename1.png")
		So(r.ImageAttributes(), ShouldEqual, "Alt Text")
		r.Next()
		So(r.ImageTarget(), ShouldEqual, "filename2.png")
		So(r.ImageAttributes(), ShouldEqual, "Alt2 Text2")
		r.Next()
		So(r.ImageTarget(), ShouldEqual, "http://example.com/images/filename3.png")
		So(r.ImageAttributes(), ShouldEqual, "Alt3 Text3")
		r.Next()
		So(r.ImageTarget(), ShouldEqual, "filename4.png")
		So(r.ImageAttributes(), ShouldEqual, `More4 [Alt4\] Text4`)
		r.Next()
		So(r.ImageTarget(), ShouldEqual, "github")
		So(r.ImageAttributes(), ShouldEqual, "large")
	})
	Convey("Regexps can encapsulate IndextermInlineMacroRx results in a struct IndextermInlineMacroRxres", t, func() {
		Convey("Escaped IndextermInlineMacroRxres should be escaped", func() {
			r := NewIndextermInlineMacroRxres(`\indexterm:[Tigers,Big cats]
				\(((Tigers,Big cats)))
				\((Tigers))`)
			So(r.HasAnyMatch(), ShouldBeTrue)
			So(len(r.matches), ShouldEqual, 3)

			//fmt.Printf("\nr='%v'\n", r.matches)
			So(r.IsEscaped(), ShouldBeTrue)
			So(r.IndextermMacroName(), ShouldEqual, "indexterm")
			So(r.IndextermTextOrTerms(), ShouldEqual, "Tigers,Big cats")
			So(r.IndextermTextInBrackets(), ShouldEqual, "")
			r.Next()
			So(r.IsEscaped(), ShouldBeTrue)
			So(r.IndextermMacroName(), ShouldEqual, "")
			So(r.IndextermTextOrTerms(), ShouldEqual, "")
			So(r.IndextermTextInBrackets(), ShouldEqual, "(Tigers,Big cats)")
			r.Next()
			So(r.IsEscaped(), ShouldBeTrue)
			So(r.IndextermMacroName(), ShouldEqual, "")
			So(r.IndextermTextOrTerms(), ShouldEqual, "")
			So(r.IndextermTextInBrackets(), ShouldEqual, "Tigers")
			r.Next()
		})
		Convey("IndextermInlineMacroRxres shouldn't be followed by closing bracket", func() {
			r := NewIndextermInlineMacroRxres(`(((Tigers,Big cats))))
				((Tigers)))`)
			So(r.HasAnyMatch(), ShouldBeTrue)
			So(len(r.matches), ShouldEqual, 2)
			So(r.IndextermMacroName(), ShouldEqual, "")
			So(r.IndextermTextOrTerms(), ShouldEqual, "")
			So(r.IndextermTextInBrackets(), ShouldEqual, "(Tigers,Big cats))")
			r.Next()
			So(r.IsEscaped(), ShouldBeFalse)
			So(r.IndextermMacroName(), ShouldEqual, "")
			So(r.IndextermTextOrTerms(), ShouldEqual, "")
			So(r.IndextermTextInBrackets(), ShouldEqual, "Tigers)")
		})
	})

	Convey("Regexps can encapsulate LinkInlineRx results in a struct LinkInlineRxres", t, func() {
		Convey("Escaped LinkInlineRxres should be escaped", func() {
			r := NewLinkInlineRxres("\\http://google.com")
			So(r.HasAnyMatch(), ShouldBeTrue)
			So(r.IsLinkEscaped(), ShouldBeTrue)
			So(r.LinkPrefix(), ShouldEqual, "")
			So(r.LinkTarget(), ShouldEqual, "\\http://google.com")
			So(r.LinkText(), ShouldEqual, "")
		})
		Convey("a single-line raw url should be interpreted as a link", func() {
			r := NewLinkInlineRxres("http://google.com")
			So(r.HasAnyMatch(), ShouldBeTrue)
			So(r.IsLinkEscaped(), ShouldBeFalse)
			So(r.LinkPrefix(), ShouldEqual, "")
			So(r.LinkTarget(), ShouldEqual, "http://google.com")
			So(r.LinkText(), ShouldEqual, "")
		})
		Convey("a single-line raw url with text should be interpreted as a link", func() {
			r := NewLinkInlineRxres("http://google.com[Google]")
			So(r.HasAnyMatch(), ShouldBeTrue)
			So(r.IsLinkEscaped(), ShouldBeFalse)
			So(r.LinkPrefix(), ShouldEqual, "")
			So(r.LinkTarget(), ShouldEqual, "http://google.com")
			So(r.LinkText(), ShouldEqual, "Google")
		})
		Convey("a multi-line raw url with text should be interpreted as a link", func() {
			r := NewLinkInlineRxres("http://google.com[Google\nHomepage]")
			So(r.HasAnyMatch(), ShouldBeTrue)
			So(r.IsLinkEscaped(), ShouldBeFalse)
			So(r.LinkPrefix(), ShouldEqual, "")
			So(r.LinkTarget(), ShouldEqual, "http://google.com")
			So(r.LinkText(), ShouldEqual, "Google\nHomepage")
		})
	})

	Convey("Regexps can encapsulate LinkInlineMacroRx results in a struct LinkInlineMacroRxres", t, func() {
		Convey("Escaped LinkInlineMacroRx should be escaped", func() {
			r := NewLinkInlineMacroRxres(`\link:path[label]
			 \mailto:doc.writer@example.com[]
			 link:path2[label2]
			 mailto:doc2.writer@example2.com[xxx]`)
			So(r.HasAnyMatch(), ShouldBeTrue)
			So(len(r.matches), ShouldEqual, 4)

			So(r.IsEscaped(), ShouldBeTrue)
			So(r.LinkInlineTarget(), ShouldEqual, "path")
			So(r.LinkInlineText(), ShouldEqual, "label")

			r.Next()

			So(r.IsEscaped(), ShouldBeTrue)
			So(r.LinkInlineTarget(), ShouldEqual, "doc.writer@example.com")
			So(r.LinkInlineText(), ShouldEqual, "")

			r.Next()

			So(r.IsEscaped(), ShouldBeFalse)
			So(r.LinkInlineTarget(), ShouldEqual, "path2")
			So(r.LinkInlineText(), ShouldEqual, "label2")

			r.Next()

			So(r.IsEscaped(), ShouldBeFalse)
			So(r.LinkInlineTarget(), ShouldEqual, "doc2.writer@example2.com")
			So(r.LinkInlineText(), ShouldEqual, "xxx")
		})
	})

	Convey("Regexps can encapsulate EmailInlineMacroRx results in a struct EmailInlineMacroRxres", t, func() {
		Convey("EmailInlineMacroRx should detect lead", func() {
			r := NewEmailInlineMacroRxres(`doc.writer@example.com
				:doc2.writer@example.com`)
			So(r.HasAnyMatch(), ShouldBeTrue)
			So(len(r.matches), ShouldEqual, 2)

			So(r.IsEscaped(), ShouldBeFalse)
			So(r.EmailLead(), ShouldEqual, "")

			r.Next()

			So(r.IsEscaped(), ShouldBeFalse)
			So(r.EmailLead(), ShouldEqual, ":")
		})
	})

	Convey("Regexps can encapsulate FootnoteInlineMacroRx results in a struct FootnoteInlineMacroRxres", t, func() {
		Convey("FootnoteInlineMacroRx should detect lead", func() {
			r := NewFootnoteInlineMacroRxres(`\footnote:[text]
  footnoteref:[text]
  footnoteref:[id,text]
  footnoteref:[id]`)

			So(r.HasAnyMatch(), ShouldBeTrue)
			So(len(r.matches), ShouldEqual, 4)

			So(r.IsEscaped(), ShouldBeTrue)
			So(r.FootnotePrefix(), ShouldEqual, "footnote")
			So(r.FootnoteText(), ShouldEqual, "text")

			r.Next()
			So(r.IsEscaped(), ShouldBeFalse)
			So(r.FootnotePrefix(), ShouldEqual, "footnoteref")
			So(r.FootnoteText(), ShouldEqual, "text")

			r.Next()
			So(r.IsEscaped(), ShouldBeFalse)
			So(r.FootnotePrefix(), ShouldEqual, "footnoteref")
			So(r.FootnoteText(), ShouldEqual, "id,text")

			r.Next()
			So(r.IsEscaped(), ShouldBeFalse)
			So(r.FootnotePrefix(), ShouldEqual, "footnoteref")
			So(r.FootnoteText(), ShouldEqual, "id")

		})
	})

	Convey("Regexps can encapsulate InlineBiblioAnchorRx results in a struct InlineBiblioAnchorRxres", t, func() {
		Convey("InlineBiblioAnchorRx should detect id", func() {
			r := NewInlineBiblioAnchorRxres(`\[[[Foo]]]
  [[[Bar]]]`)

			So(r.HasAnyMatch(), ShouldBeTrue)
			So(len(r.matches), ShouldEqual, 2)

			So(r.IsEscaped(), ShouldBeTrue)
			So(r.BibId(), ShouldEqual, "Foo")

			r.Next()
			So(r.IsEscaped(), ShouldBeFalse)
			So(r.BibId(), ShouldEqual, "Bar")

		})
	})

	Convey("Regexps can encapsulate double quoted text results in a struct DoubleQuotedRxres", t, func() {
		Convey("InlineBiblioAnchorRx should detect text", func() {
			r := NewDoubleQuotedRxres(`"Who goes there?"
Who goes there2?
"notext
notext"`)

			So(r.HasAnyMatch(), ShouldBeTrue)
			So(len(r.matches), ShouldEqual, 2)

			So(r.DQQuote(), ShouldEqual, `"`)
			So(r.DQText(), ShouldEqual, `Who goes there?`)

			r.Next()
			So(r.DQQuote(), ShouldEqual, ``)
			So(r.DQText(), ShouldEqual, `Who goes there2?`)
		})
	})

	Convey("Regexps can encapsulate inline anchor text results in a struct InlineAnchorRxres", t, func() {
		Convey("InlineAnchorRxres should detect id and text", func() {
			r := NewInlineAnchorRxres(`\[[idname]]
   \[[idname2,Reference2 Text]]
   \anchor:idname3[]
   \anchor:idname4[Reference4 Text]
   [[idname5]]
   [[idname6,Reference6 Text]]
   anchor:idname7[]
   anchor:idname8[Reference8 Text]`)

			So(r.HasAnyMatch(), ShouldBeTrue)
			So(len(r.matches), ShouldEqual, 8)

			So(r.IsEscaped(), ShouldBeTrue)
			So(r.BibAnchorId(), ShouldEqual, `idname`)
			So(r.BibAnchorText(), ShouldEqual, ``)

			r.Next()
			So(r.IsEscaped(), ShouldBeTrue)
			So(r.BibAnchorId(), ShouldEqual, `idname2`)
			So(r.BibAnchorText(), ShouldEqual, `Reference2 Text`)

			r.Next()
			So(r.IsEscaped(), ShouldBeTrue)
			So(r.BibAnchorId(), ShouldEqual, `idname3`)
			So(r.BibAnchorText(), ShouldEqual, ``)

			r.Next()
			So(r.IsEscaped(), ShouldBeTrue)
			So(r.BibAnchorId(), ShouldEqual, `idname4`)
			So(r.BibAnchorText(), ShouldEqual, `Reference4 Text`)

			r.Next()
			So(r.IsEscaped(), ShouldBeFalse)
			So(r.BibAnchorId(), ShouldEqual, `idname5`)
			So(r.BibAnchorText(), ShouldEqual, ``)

			r.Next()
			So(r.IsEscaped(), ShouldBeFalse)
			So(r.BibAnchorId(), ShouldEqual, `idname6`)
			So(r.BibAnchorText(), ShouldEqual, `Reference6 Text`)

			r.Next()
			So(r.IsEscaped(), ShouldBeFalse)
			So(r.BibAnchorId(), ShouldEqual, `idname7`)
			So(r.BibAnchorText(), ShouldEqual, ``)

			r.Next()
			So(r.IsEscaped(), ShouldBeFalse)
			So(r.BibAnchorId(), ShouldEqual, `idname8`)
			So(r.BibAnchorText(), ShouldEqual, `Reference8 Text`)
		})
	})

	Convey("Regexps can encapsulate xref inline id and text results in a struct XrefInlineMacroRxres", t, func() {
		Convey("XrefInlineMacroRxres should detect id and text", func() {
			r := NewXrefInlineMacroRxres(`\&lt;&lt;id,reftext&gt;&gt;
   \xref:id2[reftext2]
   &lt;&lt;id3,reftext3&gt;&gt;
   xref:id4[reftext4]`)

			So(r.HasAnyMatch(), ShouldBeTrue)
			So(len(r.matches), ShouldEqual, 4)

			So(r.IsEscaped(), ShouldBeTrue)
			So(r.XId(), ShouldEqual, `id`)
			So(r.XrefText(), ShouldEqual, `reftext`)

			r.Next()
			So(r.IsEscaped(), ShouldBeTrue)
			So(r.XId(), ShouldEqual, `id2`)
			So(r.XrefText(), ShouldEqual, `reftext2`)

			r.Next()
			So(r.IsEscaped(), ShouldBeFalse)
			So(r.XId(), ShouldEqual, `id3`)
			So(r.XrefText(), ShouldEqual, `reftext3`)

			r.Next()
			So(r.IsEscaped(), ShouldBeFalse)
			So(r.XId(), ShouldEqual, `id4`)
			So(r.XrefText(), ShouldEqual, `reftext4`)

		})
	})
}
