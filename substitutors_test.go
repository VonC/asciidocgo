package asciidocgo

import (
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
		subs := newSubsEnums()
		So(len(subs.basic.values()), ShouldEqual, 1)
		So(len(subs.normal.values()), ShouldEqual, 6)
		So(len(subs.verbatim.values()), ShouldEqual, 2)
		So(len(subs.title.values()), ShouldEqual, 5)
		So(len(subs.header.values()), ShouldEqual, 2)
		So(len(subs.pass.values()), ShouldEqual, 0)
		So(len(subs.unknown.values()), ShouldEqual, 0)
	})

	Convey("A substitutors can aaply substitutions", t, func() {

		source := []string{"test"}
		s := &substitutors{}

		Convey("By default, no substitution or a pass subs will return source unchanged", func() {
			So(s.ApplySubs(source, nil), ShouldResemble, source)
			So(s.ApplySubs(source, []*subsEnum{subs.pass}), ShouldResemble, source)
			So(len(s.ApplySubs(source, []*subsEnum{subs.unknown})), ShouldEqual, 0)
		})

		Convey("A normal substition will use normal substitution modes", func() {
			testsub = "test_ApplySubs_allsubs"
			So(s.ApplySubs(source, []*subsEnum{subs.normal}), ShouldResemble, []string{"specialcharacters", "quotes", "attributes", "replacements", "macros", "post_replacements"})
			testsub = ""
		})

	})
}
