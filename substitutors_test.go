package asciidocgo

import (
	"fmt"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

type testSubstDocumentAble struct {
}

func (tsd *testSubstDocumentAble) Attr(name string, defaultValue interface{}, inherit bool) interface{} {
	return "mathtest"
}
func (tsd *testSubstDocumentAble) Basebackend(base interface{}) bool {
	return true
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
}
