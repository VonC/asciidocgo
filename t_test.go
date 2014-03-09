package asciidocgo

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSubstitutor2(t *testing.T) {
	Convey("A substitutors can substitute extension index term inline macro references", t, func() {
		s := &substitutors{}
		testDocument := newTestSubstDocumentAble(s)
		tim := &testInlineMacro{}
		testDocument.te.inlineMacros = append(testDocument.te.inlineMacros, tim)
		s.document = testDocument
		s.attributeListMaker = &testAttributeListMaker{}
		Convey("Substitute escaped index term inline macro should return macro", func() {
			So(s.SubMacros("\\indexterm:[Tigers,Big cats]\n  \\(((Tigers,Big cats))) \n   \\indexterm2:[Tigers] \n \\((Tigers)))"), ShouldEqual, "indexterm:[Tigers,Big cats]\n  (((Tigers,Big cats))) \n   indexterm2:[Tigers] \n ((Tigers)))")
		})
	})
}
