package asciidocgo

import (
	"fmt"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

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
		So(s.ApplySubs(source, subArray{subValue.macros}), ShouldEqual, `test +++for
		a
		passthrough+++ by test2 $$text
			multiple
			line$$ for
			test3 pass:quotes[text
			line2
			line3] end test4`)
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

		Convey("If no substitution detected, return text unchanged", func() {
			So(s.ApplySubs("test ++ nosub", subArray{subValue.macros}), ShouldEqual, "test ++ nosub")
		})

		So(s.ApplySubs(source, subArray{subValue.macros}), ShouldEqual, fmt.Sprintf(`test %s0%s by test2 %s1%s for
			test3 %s2%s end test4`, subPASS_START, subPASS_END, subPASS_START, subPASS_END, subPASS_START, subPASS_END))
	})
	Convey("A substitutors can unescape escaped branckets", t, func() {
		So(unescapeBrackets(""), ShouldEqual, "")
		So(unescapeBrackets(`a\]b]c\]`), ShouldEqual, `a]b]c]`)
	})
}
