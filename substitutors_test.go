package asciidocgo

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/VonC/asciidocgo/consts/context"
	"github.com/VonC/asciidocgo/consts/regexps"
	. "github.com/smartystreets/goconvey/convey"
)

type testSubstDocumentAble struct {
	s *substitutors
}

func (tsd *testSubstDocumentAble) Attr(name string, defaultValue interface{}, inherit bool) interface{} {
	if name == "attribute-undefined" {
		return "drop-line"
	}
	if name == "attribute-missing" {
		return "skip"
	}
	return "mathtest"
}
func (tsd *testSubstDocumentAble) Basebackend(base interface{}) bool {
	return true
}
func (tsd *testSubstDocumentAble) SubAttributes(data string, opts *OptionsParseAttributes) string {
	if tsd.s != nil {
		return tsd.s.SubAttributes(data, opts)
	}
	return ""
}
func (tsd *testSubstDocumentAble) HasAttr(name string, expect interface{}, inherit bool) bool {
	if name == "test_attr_value" {
		return true
	}
	if name == "experimental" {
		return true
	}
	return false
}

func (tsd *testSubstDocumentAble) Counter(name string, seed int) string {
	seed = seed + 1
	return strconv.Itoa(seed)
}

type testConvertable struct {
	data interface{}
}

func (tc *testConvertable) Convert() string {
	return fmt.Sprintf("%v", tc.data)
}

type testInlineMaker struct {
}

func (tim *testInlineMaker) NewInline(parent AbstractNodable, c context.Context, text string, opts *OptionsInline) Convertable {
	switch c {
	case context.Kbd:
		return &testConvertable{opts.Attributes()["keys"]}
	case context.Button:
		return &testConvertable{text}
	case context.Menu:
		return &testConvertable{opts.Attributes()}
	}
	return &testConvertable{"unknown context"}
}

func TestSubstitutor(t *testing.T) {

	Convey("A substitutors can be initialized", t, func() {

		Convey("By default, a substitutors can be created", func() {
			So(&substitutors{}, ShouldNotBeNil)
		})

		Convey("A substitutors has an empty passthroughs array", func() {
			s := substitutors{}
			So(len(s.passthroughs), ShouldEqual, 0)
		})
	})

	Convey("A substitutors has subs type", t, func() {
		So(len(subs[sub.basic]), ShouldEqual, 1)
		So(len(subs[sub.normal]), ShouldEqual, 6)
		So(len(subs[sub.verbatim]), ShouldEqual, 2)
		So(len(subs[sub.title]), ShouldEqual, 6)
		So(len(subs[sub.header]), ShouldEqual, 2)
		So(len(subs[sub.pass]), ShouldEqual, 0)
		So(len(subs[sub.unknown]), ShouldEqual, 0)
	})

	Convey("A substitutors can apply substitutions", t, func() {

		source := "test"
		s := &substitutors{}

		Convey("By default, no substitution or a pass subs will return source unchanged", func() {
			So(s.ApplySubs(source, nil), ShouldEqual, source)
			So(s.ApplySubs(source, subArray{sub.pass}), ShouldResemble, source)
			So(len(s.ApplySubs(source, subArray{sub.unknown})), ShouldEqual, 0)
			So(s.ApplySubs(source, subArray{sub.title}), ShouldEqual, "test")
		})

		Convey("A normal substition will use normal substitution modes", func() {
			testsub = "test_ApplySubs_allsubs"
			So(s.ApplySubs(source, subArray{sub.normal}), ShouldEqual, "[specialcharacters quotes attributes replacements macros post_replacements]")
			So(s.ApplySubs(source, subArray{sub.title}), ShouldEqual, "[title]")
			testsub = ""
		})
		Convey("A macros substition will call extractPassthroughs", func() {
			testsub = "test_ApplySubs_extractPassthroughs"
			So(s.ApplySubs(source, subArray{subValue.macros}), ShouldEqual, "test")
			testsub = ""
		})

	})

	Convey("A substitutors can Extract the passthrough text from the document for reinsertion without processing if escaped", t, func() {
		source := `test \+++for
		a
		passthrough+++ by test2 \$$text
			multiple
			line$$ for
			test3 \pass:quotes[text
			line2
			line3] end test4`
		s := &substitutors{}
		testsub = "test_ApplySubs_extractPassthroughs"
		So(s.ApplySubs(source, subArray{subValue.macros}), ShouldEqual, `test +++for
		a
		passthrough+++ by test2 $$text
			multiple
			line$$ for
			test3 pass:quotes[text
			line2
			line3] end test4`)
		testsub = ""
	})

	Convey("A substitutors can Extract the passthrough text from the document for reinsertion after processing", t, func() {
		source := `test +++for
		a
		passthrough+++ by test2 $$text
			multiple
			line$$ for
			test3 pass:quotes[text
			line2
			line3] end test4`
		s := &substitutors{}
		testsub = "test_ApplySubs_extractPassthroughs"

		Convey("If no inline macros substitution detected, return text unchanged", func() {
			So(s.ApplySubs("test ++ nosub", subArray{subValue.macros}), ShouldEqual, "test ++ nosub")
		})

		So(s.ApplySubs(source, subArray{subValue.macros}), ShouldEqual, fmt.Sprintf(`test %s0%s by test2 %s1%s for
			test3 %s2%s end test4`, subPASS_START, subPASS_END, subPASS_START, subPASS_END, subPASS_START, subPASS_END))
		testsub = ""
	})
	Convey("A substitutors can unescape escaped branckets", t, func() {
		So(unescapeBrackets(""), ShouldEqual, "")
		So(unescapeBrackets(`a\]b]c\]`), ShouldEqual, `a]b]c]`)
	})

	Convey("A substitutors can Extract inline text", t, func() {
		source := "`a few <\\{monospaced\\}> words`" +
			"[input]`A few <\\{monospaced\\}> words`\n" +
			"\\[input]`a few <monospaced> words`\n" +
			"\\[input]\\`a few <monospaced> words`\n" +
			"`a few\n<\\{monospaced\\}> words`" +
			"\\[input]`a few &lt;monospaced&gt; words`\n" +
			"the text `asciimath:[x = y]` should be passed through as `literal` text\n" +
			"`Here`s Johnny!"
		s := &substitutors{}
		testsub = "test_ApplySubs_extractPassthroughs"

		So(s.ApplySubs(source, subArray{subValue.macros}), ShouldEqual, fmt.Sprintf(`%s0%s%s1%s
[input]%s2%s
\input`+"`"+`a few <monospaced> words`+"`"+` : \`+"`"+`a few <monospaced> words`+"`"+`
%s3%s[input]%s4%s
the text %s5%s should be passed through as %s6%s text
`+"`"+`Here`+"`"+`s Johnny!`, subPASS_START, subPASS_END, subPASS_START, subPASS_END, subPASS_START, subPASS_END, subPASS_START, subPASS_END, subPASS_START, subPASS_END, subPASS_START, subPASS_END, subPASS_START, subPASS_END))

		Convey("If no literal text substitution detected, return text unchanged", func() {
			So(s.ApplySubs("test`nosub", subArray{subValue.macros}), ShouldEqual, "test`nosub")
		})
		testsub = ""
	})

	Convey("A substitutors can Extract math inline text", t, func() {
		source := `math:[x != 0]
   \math:[x != 0]
   asciimath:[x != 0]
   latexmath:abc[\sqrt{4} = 2]`
		s := &substitutors{}
		testsub = "test_ApplySubs_extractPassthroughs"

		So(s.ApplySubs(source, subArray{subValue.macros}), ShouldEqual, fmt.Sprintf(`%s0%s
   math:[x != 0]
   %s1%s
   %s2%s`, subPASS_START, subPASS_END, subPASS_START, subPASS_END, subPASS_START, subPASS_END))

		s.document = &testSubstDocumentAble{}

		So(s.ApplySubs(source, subArray{subValue.macros}), ShouldEqual, fmt.Sprintf(`%s3%s
   math:[x != 0]
   %s4%s
   %s5%s`, subPASS_START, subPASS_END, subPASS_START, subPASS_END, subPASS_START, subPASS_END))

		Convey("If no math literal substitution detected, return text unchanged", func() {
			So(s.ApplySubs("math:nosub", subArray{subValue.macros}), ShouldEqual, "math:nosub")
		})

		Convey("If no math literal substitution detected, return text unchanged", func() {
			So(s.ApplySubs("asciimath:[x <> 0]", subArray{subValue.specialcharacters}), ShouldEqual, "asciimath:[x &lt;&gt; 0]")
		})

		testsub = ""
	})
	Convey("A substitutors can be substitute special characters", t, func() {

		Convey("If none, return text unchanged", func() {
			So(subSpecialCharacters("abcd"), ShouldEqual, "abcd")
		})

		Convey("All special characters should be replaced", func() {
			So(subSpecialCharacters("&"), ShouldEqual, "&amp;")
			So(subSpecialCharacters("<"), ShouldEqual, "&lt;")
			So(subSpecialCharacters(">"), ShouldEqual, "&gt;")
			So(subSpecialCharacters(">&<"), ShouldEqual, "&gt;&amp;&lt;")
			So(subSpecialCharacters(">a&bc<"), ShouldEqual, "&gt;a&amp;bc&lt;")
		})
	})
	Convey("A substitutors can substitute special characters", t, func() {

		Convey("If none, return text unchanged", func() {
			s := &substitutors{}
			source := "\\[input]`a few <monospaced> words`"
			testsub = "test_ApplySubs_applyAllsubs"
			So(s.ApplySubs(source, subArray{subValue.specialcharacters}), ShouldEqual, fmt.Sprintf("\\[input]`a few &lt;monospaced&gt; words`"))
			testsub = ""
		})
	})

	Convey("A substitutors can Extract inline quotes", t, func() {
		s := &substitutors{}
		testsub = "test_ApplySubs_applyAllsubs"
		Convey("test inline quote, constrained, no attribute", func() {
			So(s.ApplySubs("test 'quote'", subArray{subValue.quotes}), ShouldEqual, "test ")
			testsub = ""
		})
		Convey("test inline quote, unconstrained, escaped, no attribute", func() {
			So(s.ApplySubs(`\[gray]__Git__Hub`, subArray{subValue.quotes}), ShouldEqual, "[gray]__Git__Hub")
			testsub = ""
		})
		Convey("test inline quote, constrained, escaped, with attribute", func() {
			So(s.ApplySubs(`\[gray]_Git_ Hub`, subArray{subValue.quotes}), ShouldEqual, "gray Hub")
			testsub = ""
		})
		Convey("test inline quote, unconstrained, unescaped, attribute", func() {
			So(s.ApplySubs(`[gray]__Git__Hub`, subArray{subValue.quotes}), ShouldEqual, "Hub")
			testsub = ""
		})

	})
	Convey("A substitutors can parse attributes", t, func() {
		s := &substitutors{}
		s.document = &testSubstDocumentAble{s}
		opts := &OptionsParseAttributes{}
		Convey("Parsing no attributes returns empty map", func() {
			So(len(s.parseAttributes("", opts)), ShouldEqual, 0)
		})
		Convey("Parsing attributes with SubInput means calling document SubAttributes", func() {
			opts.subInput = true
			So(len(s.parseAttributes("test", opts)), ShouldEqual, 0)
		})
	})
	Convey("A substitutors can substitute attribute references", t, func() {
		s := &substitutors{}
		testDocument := &testSubstDocumentAble{s}
		s.document = testDocument
		opts := &OptionsParseAttributes{}
		Convey("Substitute empty attribute references returns empty empty string", func() {
			So(s.SubAttributes("", opts), ShouldEqual, "")
		})
		Convey("Substitute attribute references detect references '{'", func() {
			So(s.SubAttributes("a {test1} b\nc {test2} d\n{noref", opts), ShouldEqual, "a {test1} b\nc {test2} d\n{noref")
		})
		Convey("Pre or Post escaped reference returns only the reference", func() {
			So(s.SubAttributes("a \\{test1} b\nc {test2\\} d\n{noref", opts), ShouldEqual, "a test1 b\nc test2 d\n{noref")
		})
		Convey("Reference with set directive drops the line if Parser.store_attribute returns empty", func() {
			So(s.SubAttributes("a {set:foo:bar} b", opts), ShouldEqual, "")
			s.document = nil
		})
		Convey("Reference with set directive and no document don't drops the line", func() {
			s.document = nil
			So(s.SubAttributes("a {set:foo:bar} b", opts), ShouldEqual, "a  b")
		})
		Convey("Reference with counter directive returns incremented counter", func() {
			s.document = testDocument
			So(s.SubAttributes("a {counter:aaa:2} b", opts), ShouldEqual, "a 3 b")
		})
		Convey("Reference with non-integer counter directive panic", func() {
			s.document = testDocument
			recovered := false
			defer func() {
				recover()
				recovered = true
				So(recovered, ShouldBeTrue)
			}()
			s.SubAttributes("a {counter:aaa:bbb} b", opts)
		})
		Convey("Reference with unknown directive warns and returns the all reference", func() {
			s.document = testDocument
			So(s.SubAttributes("a {counter:test_default} b", opts), ShouldEqual, "a {counter:test_default} b")
		})

		Convey("Reference with counter2 directive skip the counter", func() {
			s.document = testDocument
			So(s.SubAttributes("a {counter2:aaa:3} b", opts), ShouldEqual, "a  b")
		})
		Convey("Reference with no directive look for key", func() {
			s.document = testDocument
			So(s.SubAttributes("a {test_attr_value} b", opts), ShouldEqual, "a mathtest b")
		})
		Convey("Reference with intrinsec attributes returns translated string", func() {
			s.document = testDocument
			So(s.SubAttributes("a {caret} b {quot} c {ldquo} d", opts), ShouldEqual, "a ^ b \" c "+string(rune(8220))+" d")
			So(s.SubAttributes("a{space}b {two-colons} c {two-semicolons}ddd", opts), ShouldEqual, "a b :: c ;;ddd")
		})
		Convey("Reference with custom value look for attribute-missing attribute", func() {
			s.document = testDocument // meaning "skip"
			So(s.SubAttributes("a {test} b", opts), ShouldEqual, "a {test} b")
			opts.attribute_missing = "drop-line"
			So(s.SubAttributes("a {test} b2", opts), ShouldEqual, "a  b2")
		})
	})
	Convey("A substitutors can Extract Reference attributes", t, func() {
		s := &substitutors{}
		testDocument := &testSubstDocumentAble{s}
		s.document = testDocument
		testsub = "test_ApplySubs_applyAllsubs"
		Convey("reference counter attribute", func() {
			So(s.ApplySubs("a1 {counter:aaa:2} b1", subArray{subValue.attributes}), ShouldEqual, "a1 3 b1")
			testsub = ""
		})
	})

	Convey("A substitutors can Extract replaced text", t, func() {
		s := &substitutors{}
		testsub = "test_ApplySubs_applyAllsubs"
		Convey("(C) copyright sign is replaced", func() {
			So(s.ApplySubs("text with (C) copyright", subArray{subValue.replacements}), ShouldEqual, "text with "+regexps.Rtos(169)+" copyright")
			So(s.ApplySubs(`text with \(C) escaped copyright`, subArray{subValue.replacements}), ShouldEqual, "text with (C) escaped copyright")

			So(s.ApplySubs("text with (R) Registered Trademark", subArray{subValue.replacements}), ShouldEqual, "text with "+regexps.Rtos(174)+" Registered Trademark")
			So(s.ApplySubs(`text with \(R) escaped Registered Trademark`, subArray{subValue.replacements}), ShouldEqual, "text with (R) escaped Registered Trademark")

			So(s.ApplySubs("text with (TM) Trademark", subArray{subValue.replacements}), ShouldEqual, "text with "+regexps.Rtos(8482)+" Trademark")
			So(s.ApplySubs(`text with \(TM) escaped Trademark`, subArray{subValue.replacements}), ShouldEqual, "text with (TM) escaped Trademark")

			So(s.ApplySubs("text with -- dash-dash", subArray{subValue.replacements}), ShouldEqual, "text with"+regexps.Rtos(8201, 8212, 8201)+"dash-dash")
			So(s.ApplySubs(`text with \-- escaped dash-dash`, subArray{subValue.replacements}), ShouldEqual, "text with -- escaped dash-dash")

			So(s.ApplySubs("text with linked a--b--c dash-dash", subArray{subValue.replacements}), ShouldEqual, "text with linked a"+regexps.Rtos(8212)+"b"+regexps.Rtos(8212)+"c dash-dash")
			So(s.ApplySubs(`text with linked a\--b\--c escaped dash-dash`, subArray{subValue.replacements}), ShouldEqual, "text with linked a--b--c escaped dash-dash")

			So(s.ApplySubs("text with ... ellipsis", subArray{subValue.replacements}), ShouldEqual, "text with "+regexps.Rtos(8230)+" ellipsis")
			So(s.ApplySubs(`text with \... escaped ellipsis`, subArray{subValue.replacements}), ShouldEqual, "text with ... escaped ellipsis")

			So(s.ApplySubs("text with a'b'c' apostrophe or a closing single quote", subArray{subValue.replacements}), ShouldEqual, "text with a"+regexps.Rtos(8217)+"b"+regexps.Rtos(8217)+"c"+regexps.Rtos(8217)+" apostrophe or a closing single quote")
			So(s.ApplySubs(`text with a\'b\'c\' apostrophe or a closing single quote`, subArray{subValue.replacements}), ShouldEqual, "text with a'b'c' apostrophe or a closing single quote")

			So(s.ApplySubs("text with a-&gt;b -&gt; right arrow", subArray{subValue.replacements}), ShouldEqual, "text with a"+regexps.Rtos(8594)+"b "+regexps.Rtos(8594)+" right arrow")
			So(s.ApplySubs(`text with a\-&gt;b \-&gt; escaped right arrow`, subArray{subValue.replacements}), ShouldEqual, "text with a-&gt;b -&gt; escaped right arrow")

			So(s.ApplySubs("text with a=&gt;b =&gt; right double arrow", subArray{subValue.replacements}), ShouldEqual, "text with a"+regexps.Rtos(8658)+"b "+regexps.Rtos(8658)+" right double arrow")
			So(s.ApplySubs(`text with a\=&gt;b \=&gt; escaped right double arrow`, subArray{subValue.replacements}), ShouldEqual, "text with a=&gt;b =&gt; escaped right double arrow")

			So(s.ApplySubs("text with a&lt;-b &lt;- left arrow", subArray{subValue.replacements}), ShouldEqual, "text with a"+regexps.Rtos(8592)+"b "+regexps.Rtos(8592)+" left arrow")
			So(s.ApplySubs(`text with a\&lt;-b \&lt;- escaped left arrow`, subArray{subValue.replacements}), ShouldEqual, "text with a&lt;-b &lt;- escaped left arrow")

			So(s.ApplySubs("text with a&lt;=b &lt;= left double arrow", subArray{subValue.replacements}), ShouldEqual, "text with a"+regexps.Rtos(8656)+"b "+regexps.Rtos(8656)+" left double arrow")
			So(s.ApplySubs(`text with a\&lt;=b \&lt;= escaped left double arrow`, subArray{subValue.replacements}), ShouldEqual, "text with a&lt;=b &lt;= escaped left double arrow")

			So(s.ApplySubs("text with &amp;abc; &amp;#123; &amp;#123456; &amp;#xA1b2; &amp;#xA1b2c3; restore entities", subArray{subValue.replacements}), ShouldEqual, "text with &abc; &#123; &amp;#123456; &#xA1b2; &amp;#xA1b2c3; restore entities")
			So(s.ApplySubs(`text with \&amp;abc; \&amp;#123; \&amp;#123456; \&amp;#xA1b2; \&amp;#xA1b2c3; escaped restore entities`, subArray{subValue.replacements}), ShouldEqual, "text with &amp;abc; &amp;#123; \\&amp;#123456; &amp;#xA1b2; \\&amp;#xA1b2c3; escaped restore entities")

			testsub = ""
		})
	})

	Convey("A substitutors can substitute macros references", t, func() {
		s := &substitutors{}
		testDocument := &testSubstDocumentAble{s}
		s.document = testDocument
		s.inlineMaker = &testInlineMaker{}

		Convey("Substitute empty macros references returns empty empty string", func() {
			So(s.SubMacros(""), ShouldEqual, "")
		})
		Convey("Substitute non-empty without macros references returns same text", func() {
			So(s.SubMacros("test"), ShouldEqual, "test")
		})
		Convey("Substitute kbd macro with single key", func() {
			So(s.SubMacros("kbd:[F3]"), ShouldEqual, "[]")
		})
		Convey("Substitute kbd macro with escaped single key", func() {
			So(s.SubMacros(`\kbd:[F3]`), ShouldEqual, "kbd:[F3]")
		})
		Convey("Substitute kbd macro with single '+' key", func() {
			So(s.SubMacros("kbd:[+]"), ShouldEqual, "[+]")
		})
		Convey("Substitute kbd macro ignores first empty key, detects others", func() {
			So(s.SubMacros("kbd:[+ Alt+T]"), ShouldEqual, "[Alt T]")
		})
		Convey("Substitute kbd macro detects '+' key", func() {
			So(s.SubMacros("kbd:[Ctrl,+]"), ShouldEqual, "[Ctrl +]")
		})
		Convey("Substitute kbd macro detects '++' suffixed key", func() {
			So(s.SubMacros("kbd:[Ctrl,abc++]"), ShouldEqual, "[Ctrl abc +]")
		})

		Convey("Substitute btn macro detects the label", func() {
			So(s.SubMacros("btn:[alabel]"), ShouldEqual, "alabel")
		})

		Convey("Substitute menu macro with escape return menu macro", func() {
			So(s.SubMacros(`\menu:name0[items0]`), ShouldEqual, "menu:name0[items0]")
		})
		Convey("Substitute menu macro detects the item", func() {
			So(s.SubMacros("menu:name[items]"), ShouldEqual, "map[menu:name submenu:[] menuitem:items]")
		})

		Convey("Substitute menu macro detects the items", func() {
			So(s.SubMacros("menu:name[item1 item2 item3  ]"), ShouldEqual, "map[menu:name submenu:[] menuitem:item1 item2 item3]")
			So(s.SubMacros("menu:name[item1b ,  item2b,  item3b]"), ShouldEqual, "map[menu:name submenu:[item1b item2b] menuitem:item3b]")
			So(s.SubMacros("menu:name[item1c  &gt; item2c &gt;  item3c]"), ShouldEqual, "map[menu:name submenu:[item1c item2c] menuitem:item3c]")
		})
		Convey("Substitute menu macro detects the inline items with &gt;", func() {
			So(s.SubMacros(`menu \"File &gt; New" test`), ShouldEqual, `menu "File &gt; New" test`)
			So(s.SubMacros(`menu "File1 &gt; New1" test1`), ShouldEqual, "menu map[menu:[File1 New1] submenu:[File1] menuitem:New1] test1")
			So(s.SubMacros(`menu "File2 &gt; New2   &gt;    Item2" test2`), ShouldEqual, "menu map[menu:[File2 New2 Item2] submenu:[File2 New2] menuitem:Item2] test2")
		})
	})
}
