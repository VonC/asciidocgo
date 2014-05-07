package asciidocgo

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSubstitutorSubs(t *testing.T) {

	Convey("A subEnum can be converted from string", t, func() {
		So(aToSE(string(subsBasic)), ShouldEqual, sub.basic)
		So(aToSE(string(subsNormal)), ShouldEqual, sub.normal)
		So(aToSE(string(subsVerbatim)), ShouldEqual, sub.verbatim)
		So(aToSE(string(subsTitle)), ShouldEqual, sub.title)
		So(aToSE(string(subsHeader)), ShouldEqual, sub.header)
		So(aToSE(string(subsPass)), ShouldEqual, sub.pass)
		So(aToSE(string(subsUnknown)), ShouldEqual, sub.unknown)
		So(aToSE("xxxtestxxx"), ShouldEqual, nil)
	})
	Convey("A subValue can be converted from string", t, func() {
		So(aToSEValues(string(subsSpecialCharacters)), ShouldEqual, subValue.specialcharacters)
		So(aToSEValues(string(subsQuotes)), ShouldEqual, subValue.quotes)
		So(aToSEValues(string(subsAttributes)), ShouldEqual, subValue.attributes)
		So(aToSEValues(string(subsReplacements)), ShouldEqual, subValue.replacements)
		So(aToSEValues(string(subsMacros)), ShouldEqual, subValue.macros)
		So(aToSEValues(string(subsPostReplacements)), ShouldEqual, subValue.replacements)
		So(aToSEValues(string(subsCallout)), ShouldEqual, subValue.callouts)
		So(aToSEValues("xxxtestxxx1"), ShouldEqual, nil)
	})
	Convey("A composite SE can be converted from string", t, func() {
		So(aToCompositeSE(string(subsNone)), ShouldEqual, compositeSub.none)
		So(aToCompositeSE(string(subsNormal)), ShouldEqual, compositeSub.normal)
		So(aToCompositeSE(string(subsVerbatim)), ShouldEqual, compositeSub.verbatim)
		So(aToCompositeSE(string(subsSpecialChars)), ShouldEqual, compositeSub.specialchars)
		So(aToCompositeSE("xxxtestxxx2"), ShouldEqual, nil)
	})
	Convey("A subSymbol can be converted from string", t, func() {
		So(aToSubSymbol(string(subsA)), ShouldEqual, subSymbol.a)
		So(aToSubSymbol(string(subsM)), ShouldEqual, subSymbol.m)
		So(aToSubSymbol(string(subsN)), ShouldEqual, subSymbol.n)
		So(aToSubSymbol(string(subsP)), ShouldEqual, subSymbol.p)
		So(aToSubSymbol(string(subsQ)), ShouldEqual, subSymbol.q)
		So(aToSubSymbol(string(subsR)), ShouldEqual, subSymbol.r)
		So(aToSubSymbol(string(subsC)), ShouldEqual, subSymbol.c)
		So(aToSubSymbol(string(subsV)), ShouldEqual, subSymbol.v)
		So(aToSubSymbol("xxxtestxxx3"), ShouldEqual, nil)
	})
	Convey("A subOption can be converted from string", t, func() {
		So(aToSubOption(string(subsBlock)), ShouldEqual, subOption.block)
		So(aToSubOption(string(subsInline)), ShouldEqual, subOption.inline)
		So(aToCompositeSE("xxxtestxxx4"), ShouldEqual, nil)
	})

	Convey("Subarrays with nil elements can be intersected or removed", t, func() {
		s1 := &subArray{nil, aToSE(string(subsNormal)), aToSE(string(subsVerbatim))}
		s2 := &subArray{aToSE(string(subsNormal)), nil}
		s1i2 := s1.Intersect(*s2)
		s1r2 := s1.Remove(*s2)
		So(fmt.Sprintf("%s", s1i2), ShouldEqual, "[%!s(*asciidocgo.subsEnum=&{normal})]")
		So(fmt.Sprintf("%s", s1r2), ShouldEqual, "[%!s(*asciidocgo.subsEnum=&{verbatim})]")
	})
	Convey("a sub can test for composite", t, func() {
		So(aToCompositeSE(string(subsNormal)).isCompositeSub(), ShouldBeTrue)
		So(aToSEValues(string(subsAttributes)).isCompositeSub(), ShouldBeFalse)
	})
}
