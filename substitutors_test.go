package asciidocgo

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/VonC/asciidocgo/consts/context"
	"github.com/VonC/asciidocgo/consts/regexps"
	. "github.com/smartystreets/goconvey/convey"
)

type testSubstDocumentAble struct {
	s               *substitutors
	te              *testExtensionables
	linkAttrs       bool
	hideUriScheme   bool
	references      *testReferencable
	footnotes       []Footnotable
	counterFootnote int
}

type testReferencable struct {
}

func (tr *testReferencable) HasId(id string) bool {
	if id == "testref2" {
		return true
	}
	return false
}
func (tr *testReferencable) Get(id string) string {
	if id == "includes" {
		return "xxdoc8yy"
	}
	return ""
}

func newTestSubstDocumentAble(s *substitutors) *testSubstDocumentAble {
	tsd := &testSubstDocumentAble{s: s}
	tsd.te = &testExtensionables{}
	tsd.footnotes = []Footnotable{}
	return tsd
}

func (tsd *testSubstDocumentAble) References() Referencable {
	return tsd.references
}

func (tsd *testSubstDocumentAble) Attr(name string, defaultValue interface{}, inherit bool) interface{} {
	if name == "attribute-undefined" {
		return "drop-line"
	}
	if name == "attribute-missing" {
		return "skip"
	}
	if name == "linkattrs" {
		if tsd.linkAttrs {
			return "linkattrs"
		}
		return ""
	}
	if name == "hide-uri-scheme" {
		if tsd.hideUriScheme {
			return "hide-uri-scheme"
		}
		return ""
	}
	if name == "relfileprefix" {
		return "relfileprefixAttr"
	}
	if name == "outfilesuffix" {
		return ""
	}
	if name == "docname" {
		return "doc8"
	}
	return "mathtest"
}
func (tsd *testSubstDocumentAble) Basebackend(base interface{}) bool {
	return true
}

type testFootnotable struct {
	index int
	id    int
	text  string
}

func (tf *testFootnotable) Index() int {
	return tf.index
}
func (tf *testFootnotable) Id() int {
	return tf.id
}
func (tf *testFootnotable) Text() string {
	return tf.text
}

func (tf *testFootnotable) String() string {
	return fmt.Sprintf("footnote(%v,%v): '%v'", tf.Id(), tf.Index(), tf.Text())
}

func (tsd *testSubstDocumentAble) NewFootnote(index int, id int, text string) Footnotable {
	return &testFootnotable{index: index, id: id, text: text}
}
func (tsd *testSubstDocumentAble) RegisterFootnote(f Footnotable) {
	tsd.footnotes = append(tsd.footnotes, f)
}
func (tsd *testSubstDocumentAble) FindFootnote(id int) Footnotable {
	var footnote Footnotable
	for _, f := range tsd.footnotes {
		if f.Id() == id {
			footnote = f
			break
		}
	}
	return footnote
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
	if name == "linkattrs" {
		return tsd.linkAttrs
	}
	if name == "hide-uri-scheme" {
		return tsd.hideUriScheme
	}
	return false
}

func (tsd *testSubstDocumentAble) Counter(name string, seed int) string {
	if name == "footnote-number" {
		tsd.counterFootnote = tsd.counterFootnote + 1
		return strconv.Itoa(tsd.counterFootnote)
	}
	seed = seed + 1
	return strconv.Itoa(seed)
}
func (tsd *testSubstDocumentAble) Register(typeDoc string, value []string) {
}

type testGlobalParsable struct {
}

func (tgp *testGlobalParsable) storeAttribute(name string, value string, doc SubstDocumentable, attrs map[string]interface{}) (nameRes string, valueRes string) {
	if name == "undefined" {
		return "undefined_tgpnres", ""
	}
	return name + "_tgpnres", value + "_tgpvres"
}

type testInlineMacro struct {
	rx                     *regexp.Regexp
	rxshort                *regexp.Regexp
	shortFormat            bool
	contentModelAttributes bool
	posAttrs               []string
}

func (tim *testInlineMacro) IsShortFormat() bool            { return tim.shortFormat }
func (tim *testInlineMacro) IsContentModelAttributes() bool { return tim.contentModelAttributes }
func (tim *testInlineMacro) Regexp() *regexp.Regexp {
	if tim.rx == nil {
		tim.rx, _ = regexp.Compile(`\\?test:(\S+?)\[((?:\\\]|[^\]])*?)\]`)
	}
	if tim.rxshort == nil {
		tim.rxshort, _ = regexp.Compile(`\\?testShort:\[((?:\\\]|[^\]])*?)\]`)
	}
	if tim.shortFormat {
		return tim.rxshort
	}
	return tim.rx
}
func (tim *testInlineMacro) ProcessMethod(self interface{}, target string, attributes map[string]interface{}) string {
	return fmt.Sprintf("%v", attributes)
}
func (tim *testInlineMacro) PosAttrs() []string { return tim.posAttrs }

type testExtensionables struct {
	inlineMacros []InlineMacroable
}

func (te *testExtensionables) HasInlineMacros() bool {
	return len(te.inlineMacros) > 0
}
func (te *testExtensionables) InlineMacros() []InlineMacroable {
	return te.inlineMacros
}

func (tsd *testSubstDocumentAble) Extensions() Extensionables {
	return tsd.te
}

type testConvertable struct {
	data interface{}
}

func (tc *testConvertable) Convert() string {
	//fmt.Printf("\ntc.data: '%v'\n", tc.data)
	return fmt.Sprintf("%v", tc.data)
}

type testAbstractNodable struct {
}

func (tan *testAbstractNodable) IsAbstractNodable() {}

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
	case context.Image:
		msg := fmt.Sprintf("Context '%v': target '%v' type '%v' attrs: '%v'", c, opts.Target(), opts.TypeInline(), opts.Attributes())
		return &testConvertable{msg}
	case context.IndexTerm:
		msg := fmt.Sprintf("ContextIT '%v': text '%v' ===> type '%v' attrs: '%v'", c, text, opts.TypeInline(), opts.Attributes())
		//fmt.Printf("\n msg='%v'", msg)
		return &testConvertable{msg}
	case context.Anchor:
		msg := fmt.Sprintf("ContextAn '%v': text '%v' ===> type '%v' target '%v' attrs: '%v'", c, text, opts.TypeInline(), opts.Target(), opts.Attributes())
		return &testConvertable{msg}
	case context.Quoted:
		msg := fmt.Sprintf("ContextQt '%v': text '%v' ===> type '%v' target '%v' id '%v' attrs: '%v'", c, text, opts.TypeInline(), opts.Target(), opts.id, opts.Attributes())
		return &testConvertable{msg}
	case context.Footnote:
		msg := fmt.Sprintf("ContextFt '%v': text '%v' ===> type '%v' target '%v' id '%v' attrs: '%v'", c, text, opts.TypeInline(), opts.Target(), opts.id, opts.Attributes())
		return &testConvertable{msg}
	}
	return &testConvertable{"unknown context"}
}

type testAttributeListable struct {
	attrline string
	block    ApplyNormalSubsable
}

func (tal *testAttributeListable) ParseInto(into map[string]interface{}, posAttrs []string) map[string]interface{} {
	into["*testpi: "+tal.attrline+"*"] = posAttrs
	//fmt.Printf("\ninto='%v'\n", into)
	return into
}

func (tal *testAttributeListable) Parse(posAttrs []string) map[string]interface{} {
	res := make(map[string]interface{})
	b := ""
	if tal.block != nil {
		b = tal.block.ApplyNormalSubs("block")
	}
	res["*testp: "+tal.attrline+"*"+b+"*"] = posAttrs
	if tal.attrline == "\"text,url" {
		res["1"] = tal.attrline[1:] + "*" + b + "*"
	}
	if tal.attrline == "\"text2,url2^" {
		res["1"] = "*" + b + "*" + tal.attrline[1:]
	}
	if tal.attrline == "\"label, b^" {
		res["1"] = tal.attrline[1:]
	}
	if tal.attrline == "\"a,b,c" {
		res["1"] = "a"
		res["2"] = "b"
		res["3"] = "c"
	}
	if tal.attrline == "\"a,,c=(d)" {
		res["1"] = "a"
		res["2"] = ""
		res["3"] = "c=(d)"
	}
	//fmt.Printf("\npars '%v'='%v'\n", tal.attrline, res)
	return res
}

type testAttributeListMaker struct {
}

func (talm *testAttributeListMaker) NewAttributeList(attrline string, block ApplyNormalSubsable, delimiter string) AttributeListable {
	tal := &testAttributeListable{attrline: attrline, block: block}
	return tal
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
			So(s.ApplySubs(source, nil, false), ShouldEqual, source)
			So(s.ApplySubs(source, subArray{sub.pass}, false), ShouldResemble, source)
			So(len(s.ApplySubs(source, subArray{sub.unknown}, false)), ShouldEqual, 0)
			So(s.ApplySubs(source, subArray{sub.title}, false), ShouldEqual, "test")
		})

		Convey("A normal substition will use normal substitution modes", func() {
			testsub = "test_ApplySubs_allsubs"
			So(s.ApplySubs(source, subArray{sub.normal}, true), ShouldEqual, "[specialcharacters quotes attributes replacements macros post_replacements]")
			So(s.ApplySubs(source, subArray{sub.title}, false), ShouldEqual, "[title]")
			So(s.ApplySubs(source, subArray{sub.title}, true), ShouldEqual, "[title]")
			So(s.ApplySubs(source, subArray{subSymbol.v}, true), ShouldEqual, "[specialcharacters callouts]")
			testsub = ""
		})
		Convey("A macros substition will call extractPassthroughs", func() {
			testsub = "test_ApplySubs_extractPassthroughs"
			So(s.ApplySubs(source, subArray{subValue.macros}, false), ShouldEqual, "test")
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
		So(s.ApplySubs(source, subArray{subValue.macros}, false), ShouldEqual, `test +++for
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
			So(s.ApplySubs("test ++ nosub", subArray{subValue.macros}, false), ShouldEqual, "test ++ nosub")
		})

		So(s.ApplySubs(source, subArray{subValue.macros}, false), ShouldEqual, fmt.Sprintf(`test %s0%s by test2 %s1%s for
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
		s.attributeListMaker = &testAttributeListMaker{}

		So(s.ApplySubs(source, subArray{subValue.macros}, false), ShouldEqual, fmt.Sprintf(`%s0%s%s1%s
[input]%s2%s
\input`+"`"+`a few <monospaced> words`+"`"+` : \`+"`"+`a few <monospaced> words`+"`"+`
%s3%s[input]%s4%s
the text %s5%s should be passed through as %s6%s text
`+"`"+`Here`+"`"+`s Johnny!`, subPASS_START, subPASS_END, subPASS_START, subPASS_END, subPASS_START, subPASS_END, subPASS_START, subPASS_END, subPASS_START, subPASS_END, subPASS_START, subPASS_END, subPASS_START, subPASS_END))

		Convey("If no literal text substitution detected, return text unchanged", func() {
			So(s.ApplySubs("test`nosub", subArray{subValue.macros}, false), ShouldEqual, "test`nosub")
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

		So(s.ApplySubs(source, subArray{subValue.macros}, false), ShouldEqual, fmt.Sprintf(`%s0%s
   math:[x != 0]
   %s1%s
   %s2%s`, subPASS_START, subPASS_END, subPASS_START, subPASS_END, subPASS_START, subPASS_END))

		s.document = newTestSubstDocumentAble(nil)

		So(s.ApplySubs(source, subArray{subValue.macros}, false), ShouldEqual, fmt.Sprintf(`%s3%s
   math:[x != 0]
   %s4%s
   %s5%s`, subPASS_START, subPASS_END, subPASS_START, subPASS_END, subPASS_START, subPASS_END))

		Convey("If no math literal substitution detected, return text unchanged", func() {
			So(s.ApplySubs("math:nosub", subArray{subValue.macros}, false), ShouldEqual, "math:nosub")
		})

		Convey("If no math literal substitution detected, return text unchanged", func() {
			So(s.ApplySubs("asciimath:[x <> 0]", subArray{subValue.specialcharacters}, false), ShouldEqual, "asciimath:[x &lt;&gt; 0]")
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
			So(s.ApplySubs(source, subArray{subValue.specialcharacters}, false), ShouldEqual, fmt.Sprintf("\\[input]`a few &lt;monospaced&gt; words`"))
			testsub = ""
		})
	})

	Convey("A substitutors can Extract inline quotes", t, func() {
		s := &substitutors{}
		s.inlineMaker = &testInlineMaker{}
		testsub = "test_ApplySubs_applyAllsubs"
		Convey("test inline quote, constrained, no attribute", func() {
			So(s.ApplySubs("test 'quote'", subArray{subValue.quotes}, false), ShouldEqual, "test ContextQt 'quoted': text 'quote' ===> type 'Emphasis' target '' id '' attrs: 'map[]'")
			testsub = ""
		})
		Convey("test inline quote, unconstrained, escaped, no attribute", func() {
			So(s.ApplySubs(`\[gray]__Git__Hub`, subArray{subValue.quotes}, false), ShouldEqual, "[gray]__Git__Hub")
			testsub = ""
		})
		Convey("test inline quote, constrained, escaped or not, with attribute", func() {
			So(s.ApplySubs(`\[gray]_Git_ Hub`, subArray{subValue.quotes}, false), ShouldEqual, "grayContextQt 'quoted': text 'Git' ===> type 'Emphasis' target '' id '' attrs: 'map[]' Hub")
			So(s.ApplySubs(`[.bbb#gray3.gray4]_Git2_ Hub2`, subArray{subValue.quotes}, false), ShouldEqual, "ContextQt 'quoted': text 'Git2' ===> type 'Emphasis' target '' id 'gray3' attrs: 'map[roles:[bbb gray4]]' Hub2")
			testsub = ""
		})
		Convey("test inline quote, unconstrained, unescaped, attribute", func() {
			So(s.ApplySubs(`[gray]__Git__Hub`, subArray{subValue.quotes}, false), ShouldEqual, "ContextQt 'quoted': text 'Git' ===> type 'Emphasis' target '' id '' attrs: 'map[roles:[gray]]'Hub")
			So(s.ApplySubs(`[.aaa#gray1.gray2]__Git__Hub`, subArray{subValue.quotes}, false), ShouldEqual, "ContextQt 'quoted': text 'Git' ===> type 'Emphasis' target '' id 'gray1' attrs: 'map[roles:[aaa gray2]]'Hub")
			testsub = ""
		})

	})
	Convey("A substitutors can parse attributes", t, func() {
		s := &substitutors{}
		s.document = newTestSubstDocumentAble(s)
		s.attributeListMaker = &testAttributeListMaker{}
		opts := &OptionsParseAttributes{}
		Convey("Parsing no attributes returns empty map", func() {
			So(len(s.parseAttributes("", []string{}, opts)), ShouldEqual, 0)
		})
		Convey("Parsing attributes with SubInput means calling document SubAttributes", func() {
			opts.subInput = true
			attrs := s.parseAttributes("test", []string{}, opts)
			So(len(attrs), ShouldEqual, 1)
			So(fmt.Sprintf("%v", attrs), ShouldEqual, "map[*testp: test*block*:[]]")
		})
		Convey("Parsing attributes with SubInput means calling document ParseInto", func() {
			into := make(map[string]interface{})
			intoAttrs := []string{"intoa1", "intoa2"}
			into["into1"] = intoAttrs
			opts.into = into
			attrs := s.parseAttributes("test2", []string{}, opts)
			So(len(attrs), ShouldEqual, 2)
			So(fmt.Sprintf("%v", attrs), ShouldEqual, "map[into1:[intoa1 intoa2] *testpi: test2*:[]]")
		})
		Convey("Parsing attributes with SubResult means using substitutor as block", func() {
			opts.SetSubResult(true)
			attrs := s.parseAttributes("test3", []string{}, opts)
			So(len(attrs), ShouldEqual, 3)
			So(fmt.Sprintf("%v", attrs), ShouldEqual, "map[into1:[intoa1 intoa2] *testpi: test2*:[] *testpi: test3*:[]]")
		})
	})
	Convey("A substitutors can parse quoted text attributes", t, func() {
		s := &substitutors{}
		testDocument := newTestSubstDocumentAble(s)
		s.document = testDocument
		s.parser = &testGlobalParsable{}
		Convey("simple quoted text means no id, and quoted text as role", func() {
			So(fmt.Sprintf("%v", s.parseQuotedTextAttributes("test")), ShouldEqual, "map[roles:[test]]")
		})
		Convey("simple quoted text with {} means substituting attributes", func() {
			So(fmt.Sprintf("%v", s.parseQuotedTextAttributes("{set:foo:bar}")), ShouldEqual, "map[]")
			So(fmt.Sprintf("%v", s.parseQuotedTextAttributes("{counter:aaa:2}")), ShouldEqual, "map[roles:[3]]")
		})
		Convey("simple quoted text with ',' means only first part is considered", func() {
			So(fmt.Sprintf("%v", s.parseQuotedTextAttributes("test1,test2")), ShouldEqual, "map[roles:[test1]]")
		})
		Convey("simple quoted text starting with '.' means first part is skipped for roles", func() {
			So(fmt.Sprintf("%v", s.parseQuotedTextAttributes(".test3.test4.test5")), ShouldEqual, "map[roles:[test3 test4 test5]]")
		})
		Convey("simple quoted text starting with '#' means first part is id", func() {
			So(fmt.Sprintf("%v", s.parseQuotedTextAttributes("#aa.test6.test7.test8")), ShouldEqual, "map[id:aa roles:[test6 test7 test8]]")
			So(fmt.Sprintf("%v", s.parseQuotedTextAttributes(".testA.testB#ccc.ddd")), ShouldEqual, "map[id:ccc roles:[testA testB ddd]]")
		})
	})
	Convey("A substitutors can substitute attribute references", t, func() {
		s := &substitutors{}
		testDocument := newTestSubstDocumentAble(s)
		s.document = testDocument
		s.parser = &testGlobalParsable{}
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
			So(s.SubAttributes("a {set:foo:bar} b", opts), ShouldEqual, "a  b")
			So(s.SubAttributes("a {set:undefined:bar} b2", opts), ShouldEqual, "")
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
		testDocument := newTestSubstDocumentAble(s)
		s.document = testDocument
		testsub = "test_ApplySubs_applyAllsubs"
		Convey("reference counter attribute", func() {
			So(s.ApplySubs("a1 {counter:aaa:2} b1", subArray{subValue.attributes}, false), ShouldEqual, "a1 3 b1")
			testsub = ""
		})
	})

	Convey("A substitutors can Extract replaced text", t, func() {
		s := &substitutors{}
		testsub = "test_ApplySubs_applyAllsubs"
		Convey("(C) copyright sign is replaced", func() {
			So(s.ApplySubs("text with (C) copyright", subArray{subValue.replacements}, false), ShouldEqual, "text with "+regexps.Rtos(169)+" copyright")
			So(s.ApplySubs(`text with \(C) escaped copyright`, subArray{subValue.replacements}, false), ShouldEqual, "text with (C) escaped copyright")

			So(s.ApplySubs("text with (R) Registered Trademark", subArray{subValue.replacements}, false), ShouldEqual, "text with "+regexps.Rtos(174)+" Registered Trademark")
			So(s.ApplySubs(`text with \(R) escaped Registered Trademark`, subArray{subValue.replacements}, false), ShouldEqual, "text with (R) escaped Registered Trademark")

			So(s.ApplySubs("text with (TM) Trademark", subArray{subValue.replacements}, false), ShouldEqual, "text with "+regexps.Rtos(8482)+" Trademark")
			So(s.ApplySubs(`text with \(TM) escaped Trademark`, subArray{subValue.replacements}, false), ShouldEqual, "text with (TM) escaped Trademark")

			So(s.ApplySubs("text with -- dash-dash", subArray{subValue.replacements}, false), ShouldEqual, "text with"+regexps.Rtos(8201, 8212, 8201)+"dash-dash")
			So(s.ApplySubs(`text with \-- escaped dash-dash`, subArray{subValue.replacements}, false), ShouldEqual, "text with -- escaped dash-dash")

			So(s.ApplySubs("text with linked a--b--c dash-dash", subArray{subValue.replacements}, false), ShouldEqual, "text with linked a"+regexps.Rtos(8212)+"b"+regexps.Rtos(8212)+"c dash-dash")
			So(s.ApplySubs(`text with linked a\--b\--c escaped dash-dash`, subArray{subValue.replacements}, false), ShouldEqual, "text with linked a--b--c escaped dash-dash")

			So(s.ApplySubs("text with ... ellipsis", subArray{subValue.replacements}, false), ShouldEqual, "text with "+regexps.Rtos(8230)+" ellipsis")
			So(s.ApplySubs(`text with \... escaped ellipsis`, subArray{subValue.replacements}, false), ShouldEqual, "text with ... escaped ellipsis")

			So(s.ApplySubs("text with a'b'c' apostrophe or a closing single quote", subArray{subValue.replacements}, false), ShouldEqual, "text with a"+regexps.Rtos(8217)+"b"+regexps.Rtos(8217)+"c"+regexps.Rtos(8217)+" apostrophe or a closing single quote")
			So(s.ApplySubs(`text with a\'b\'c\' apostrophe or a closing single quote`, subArray{subValue.replacements}, false), ShouldEqual, "text with a'b'c' apostrophe or a closing single quote")

			So(s.ApplySubs("text with a-&gt;b -&gt; right arrow", subArray{subValue.replacements}, false), ShouldEqual, "text with a"+regexps.Rtos(8594)+"b "+regexps.Rtos(8594)+" right arrow")
			So(s.ApplySubs(`text with a\-&gt;b \-&gt; escaped right arrow`, subArray{subValue.replacements}, false), ShouldEqual, "text with a-&gt;b -&gt; escaped right arrow")

			So(s.ApplySubs("text with a=&gt;b =&gt; right double arrow", subArray{subValue.replacements}, false), ShouldEqual, "text with a"+regexps.Rtos(8658)+"b "+regexps.Rtos(8658)+" right double arrow")
			So(s.ApplySubs(`text with a\=&gt;b \=&gt; escaped right double arrow`, subArray{subValue.replacements}, false), ShouldEqual, "text with a=&gt;b =&gt; escaped right double arrow")

			So(s.ApplySubs("text with a&lt;-b &lt;- left arrow", subArray{subValue.replacements}, false), ShouldEqual, "text with a"+regexps.Rtos(8592)+"b "+regexps.Rtos(8592)+" left arrow")
			So(s.ApplySubs(`text with a\&lt;-b \&lt;- escaped left arrow`, subArray{subValue.replacements}, false), ShouldEqual, "text with a&lt;-b &lt;- escaped left arrow")

			So(s.ApplySubs("text with a&lt;=b &lt;= left double arrow", subArray{subValue.replacements}, false), ShouldEqual, "text with a"+regexps.Rtos(8656)+"b "+regexps.Rtos(8656)+" left double arrow")
			So(s.ApplySubs(`text with a\&lt;=b \&lt;= escaped left double arrow`, subArray{subValue.replacements}, false), ShouldEqual, "text with a&lt;=b &lt;= escaped left double arrow")

			So(s.ApplySubs("text with &amp;abc; &amp;#123; &amp;#123456; &amp;#xA1b2; &amp;#xA1b2c3; restore entities", subArray{subValue.replacements}, false), ShouldEqual, "text with &abc; &#123; &amp;#123456; &#xA1b2; &amp;#xA1b2c3; restore entities")
			So(s.ApplySubs(`text with \&amp;abc; \&amp;#123; \&amp;#123456; \&amp;#xA1b2; \&amp;#xA1b2c3; escaped restore entities`, subArray{subValue.replacements}, false), ShouldEqual, "text with &amp;abc; &amp;#123; \\&amp;#123456; &amp;#xA1b2; \\&amp;#xA1b2c3; escaped restore entities")

			testsub = ""
		})
	})

	Convey("A substitutors can substitute macros references", t, func() {
		s := &substitutors{}
		testDocument := newTestSubstDocumentAble(s)
		s.document = testDocument
		s.inlineMaker = &testInlineMaker{}
		s.attributeListMaker = &testAttributeListMaker{}

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
	Convey("A substitutors can substitute extension inline macro references", t, func() {
		s := &substitutors{}
		testDocument := newTestSubstDocumentAble(s)
		tim := &testInlineMacro{}
		testDocument.te.inlineMacros = append(testDocument.te.inlineMacros, tim)
		s.document = testDocument
		s.attributeListMaker = &testAttributeListMaker{}
		/*
			fmt.Printf("\ns.document '%v'\n", s.document)
			fmt.Printf("\ntestDocument.te '%v'\n", testDocument.te)
			fmt.Printf("\ntestDocument.te.inlineMacros '%v'\n", testDocument.te.inlineMacros)
			fmt.Printf("\ntestDocument.te.HasInlineMacros() '%v'\n", testDocument.te.HasInlineMacros())
			fmt.Printf("\ntestDocument.te.inlineMacros[0] '%v'\n", testDocument.te.inlineMacros[0])
			fmt.Printf("\ntestDocument.te.inlineMacros[0].Regexp() '%v'\n", testDocument.te.inlineMacros[0].Regexp())
		*/
		Convey("Substitute escaped test inline macro should return macro", func() {
			So(s.SubMacros(`\test:target1[attr1 attr2]`), ShouldEqual, "test:target1[attr1 attr2]")
		})
		Convey("Substitute non-escaped test inline macro should return attributes", func() {
			So(s.SubMacros(`test:target1[attr1 attr2]`), ShouldEqual, "map[text:attr1 attr2]")
		})
		Convey("Substitute non-escaped test inline macro with 'attributes' content model should return attributes", func() {
			tim.contentModelAttributes = true
			So(s.SubMacros(`test:target2[attr21 attr22]`), ShouldEqual, "map[*testp: attr21 attr22*block*:[]]")
			So(unescapeBracketedText(""), ShouldEqual, "")
		})
	})

	Convey("A substitutors can substitute image or icon macros references", t, func() {
		s := &substitutors{}

		testDocument := newTestSubstDocumentAble(s)
		s.document = testDocument
		s.inlineMaker = &testInlineMaker{}
		s.attributeListMaker = &testAttributeListMaker{}

		Convey("Substitute escaped image macros references returns same text", func() {
			So(s.SubMacros(`\image:filename1.png[Alt Text]`), ShouldEqual, "image:filename1.png[Alt Text]")
		})
		Convey("Substitute non-escaped image macros references returns target and attributes", func() {
			So(s.SubMacros(`image:filename2.png[Alt2 Text2]`), ShouldEqual, "Context 'image': target 'filename2.png' type 'image' attrs: 'map[*testp: Alt2 Text2*block*:[alt width height] alt:filename2]'")
			So(s.SubMacros(`icon:filename3.png[Alt3 Text3]`), ShouldEqual, "Context 'image': target 'filename3.png' type 'icon' attrs: 'map[*testp: Alt3 Text3*block*:[size] alt:filename3]'")
		})
	})

	Convey("A substitutors can normalizeString", t, func() {
		So(normalizeString("", false), ShouldEqual, "")
		So(normalizeString(" aaa  ", false), ShouldEqual, "aaa")
		So(normalizeString(` abcaa  
			  def       
		ghi  `, false), ShouldEqual, `abcaa   			  def        		ghi`)
		So(normalizeString(` ab\]aa  
			  d]f       
		\]hi  `, true), ShouldEqual, `ab]aa   			  d]f        		]hi`)
	})

	Convey("A substitutors can splitSimpleCsv", t, func() {
		So(len(splitSimpleCsv("")), ShouldEqual, 0)
		So(fmt.Sprintf("%v", splitSimpleCsv("aaa")), ShouldEqual, "[aaa]")
		var res []string

		res = splitSimpleCsv("aaa  ,    bb,cc,  ddd  ,eee , ")
		So(fmt.Sprintf("%v", res), ShouldEqual, "[aaa bb cc ddd eee ]")
		So(len(res), ShouldEqual, 6)

		res = splitSimpleCsv(`aa " 12 3 " a  ,    bb,c"c1"c,  d"dd  ,ee"e , `)
		So(fmt.Sprintf("%v", res), ShouldEqual, "[aa  12 3  a bb cc1c ddd  ,eee ]")
		So(len(res), ShouldEqual, 5)
		So(res[0], ShouldEqual, "aa  12 3  a")
		So(res[1], ShouldEqual, "bb")
		So(res[2], ShouldEqual, "cc1c")
		So(res[3], ShouldEqual, "ddd  ,eee")
		So(res[4], ShouldEqual, "")
	})

	Convey("A substitutors can substitute extension index term inline macro references", t, func() {
		s := &substitutors{}
		testDocument := newTestSubstDocumentAble(s)
		tim := &testInlineMacro{}
		testDocument.te.inlineMacros = append(testDocument.te.inlineMacros, tim)
		s.document = testDocument
		s.inlineMaker = &testInlineMaker{}
		s.attributeListMaker = &testAttributeListMaker{}
		Convey("Substitute escaped index term inline macro should return macro", func() {
			So(s.SubMacros("\\indexterm:[Tigers,Big cats]\n  \\(((Tigers,Big cats))) \n   \\indexterm2:[Tigers] \n \\((Tigers)))"), ShouldEqual, "indexterm:[Tigers,Big cats]\n  (((Tigers,Big cats))) \n   indexterm2:[Tigers] \n ((Tigers)))")
		})
		Convey("Substitute index term inline macro with text in brackets should return substituted macro", func() {
			So(s.SubMacros("(((Tigers,Big cats))) "), ShouldEqual, "ContextIT 'indexterm': text '' ===> type '' attrs: 'map[terms:[Tigers Big cats]]' ")
			So(s.SubMacros("((Tigers2,Big2 Tig, cats2)) "), ShouldEqual, "ContextIT 'indexterm': text 'Tigers2,Big2 Tig, cats2' ===> type 'visible' attrs: 'map[]' ")
		})
		Convey("Substitute index term inline macro with text or termsshould return substituted macro", func() {
			So(s.SubMacros("indexterm:[Tigers,Big cats]"), ShouldEqual, "ContextIT 'indexterm': text '' ===> type '' attrs: 'map[terms:[Tigers Big cats]]'")
			So(s.SubMacros("indexterm2:[Tigers]"), ShouldEqual, "ContextIT 'indexterm': text 'Tigers' ===> type 'visible' attrs: 'map[]'")
		})
	})
	Convey("A substitutors can substitute raw url macro references", t, func() {
		s := &substitutors{}
		testDocument := newTestSubstDocumentAble(s)
		tim := &testInlineMacro{}
		testDocument.te.inlineMacros = append(testDocument.te.inlineMacros, tim)
		s.document = testDocument
		s.inlineMaker = &testInlineMaker{}
		s.attributeListMaker = &testAttributeListMaker{}
		Convey("Substitute escaped raw url macro should return macro unescaped", func() {
			So(s.SubMacros("\\http://google.com[Google]\n  \\http://google.com[Google\nHomepage]"), ShouldEqual, "http://google.com[Google]\n  http://google.com[Google\nHomepage]")
		})

		Convey("Substitute invalid raw url macro should return macro unchanged", func() {
			So(s.SubMacros("link:http://google.com"), ShouldEqual, "link:http://google.com")
		})

		Convey("Substitute valid raw url macro without text should return target link", func() {
			So(s.SubMacros("&lt;http://google.com"), ShouldEqual, "&lt;ContextAn 'anchor': text 'http://google.com' ===> type 'link' target 'http://google.com' attrs: 'map[]'")
		})
		Convey("Substitute raw url macro with text and uri terminator should return target link and suffix", func() {
			So(s.SubMacros("&lt;http://google.com)[texturl]"), ShouldEqual, "&lt;ContextAn 'anchor': text 'texturl' ===> type 'link' target 'http://google.com' attrs: 'map[]')")
			So(s.SubMacros("&lt;http://google.com;[texturl2]"), ShouldEqual, "&lt;ContextAn 'anchor': text 'texturl2' ===> type 'link' target 'http://google.com' attrs: 'map[]';")
			So(s.SubMacros("&lt;http://google.com:[texturl3]"), ShouldEqual, "&lt;ContextAn 'anchor': text 'texturl3' ===> type 'link' target 'http://google.com' attrs: 'map[]':")
		})
		Convey("Substitute raw url macro with text and uri terminator ';' should return target link updated and suffix", func() {
			So(s.SubMacros("&lt;http://google.com&gt;[texturl2]"), ShouldEqual, "ContextAn 'anchor': text 'texturl2' ===> type 'link' target 'http://google.com' attrs: 'map[]'")
			So(s.SubMacros("&lt;http://google.com);[texturl3]"), ShouldEqual, "&lt;ContextAn 'anchor': text 'texturl3' ===> type 'link' target 'http://google.com' attrs: 'map[]');")
		})
		Convey("Substitute raw url macro with text and uri terminator ':' should return target link updated and suffix", func() {
			So(s.SubMacros("&lt;http://google.com):[texturl3]"), ShouldEqual, "&lt;ContextAn 'anchor': text 'texturl3' ===> type 'link' target 'http://google.com' attrs: 'map[]'):")
		})
		Convey("Substitute raw url macro with a document using link should return modified target link", func() {
			s.Document().(*testSubstDocumentAble).linkAttrs = true
			So(s.SubMacros("&lt;http://google.com[\"text,url]"), ShouldEqual, "&lt;ContextAn 'anchor': text 'text,url*block*' ===> type 'link' target 'http://google.com' attrs: 'map[]'")
			So(s.SubMacros("&lt;http://google2.com[\"text2,url2^]"), ShouldEqual, "&lt;ContextAn 'anchor': text '*block*text2,url2' ===> type 'link' target 'http://google2.com' attrs: 'map[]'")
		})
		Convey("Substitute raw url macro with text having uri inside: should return modified link text without uri", func() {
			s.Document().(*testSubstDocumentAble).linkAttrs = false
			s.Document().(*testSubstDocumentAble).hideUriScheme = true
			So(s.SubMacros("&lt;http://google.com:test"), ShouldEqual, "&lt;ContextAn 'anchor': text '' ===> type 'link' target 'http://google.com:test' attrs: 'map[]'")
		})
	})
	Convey("A substitutors can substitute link inline macro references", t, func() {
		s := &substitutors{}
		testDocument := newTestSubstDocumentAble(s)
		tim := &testInlineMacro{}
		testDocument.te.inlineMacros = append(testDocument.te.inlineMacros, tim)
		s.document = testDocument
		s.inlineMaker = &testInlineMaker{}
		s.attributeListMaker = &testAttributeListMaker{}
		Convey("Substitute escaped link inline macro should return macro unescaped", func() {
			So(s.SubMacros("\\link:path[label] \n \\mailto:doc.writer@testlinkinlinemacro.com[]"), ShouldEqual, "link:path[label] \n mailto:doc.writer@testlinkinlinemacro.com[]")
		})
		Convey("Substitute link inline macro with mailto: should return mailto: target", func() {
			So(s.SubMacros("mailto:doc.writer@example.com[] "), ShouldEqual, "ContextAn 'anchor': text 'ContextAn 'anchor': text 'doc.writer@example.com' ===> type 'link' target 'mailto:doc.writer@example.com' attrs: 'map[]'' ===> type 'link' target 'mailtoContextAn 'anchor': text ':doc.writer@example.com' ===> type 'link' target 'mailto::doc.writer@example.com' attrs: 'map[]'' attrs: 'map[]' ")
		})
		Convey("Substitute link inline macro with a document using link and quoted text should return modified target link", func() {
			s.Document().(*testSubstDocumentAble).linkAttrs = true
			So(s.SubMacros("link:path[\"label, b^] \n mailto:doc.writer@example.com[\"a,b,c] \n mailto:doc2.writer2@example2.com[\"a,,c=(d)]"), ShouldEqual, "ContextAn 'anchor': text 'label, b' ===> type 'link' target 'path' attrs: 'map[]' \n ContextAn 'anchor': text 'a' ===> type 'link' target 'mailtoContextAn 'anchor': text ':doc.writer@example.com' ===> type 'link' target 'mailto::doc.writer@example.com' attrs: 'map[]'?subject=b&amp;body=c' attrs: 'map[]' \n ContextAn 'anchor': text 'a' ===> type 'link' target 'mailtoContextAn 'anchor': text ':doc2.writer2@example2.com' ===> type 'link' target 'mailto::doc2.writer2@example2.com' attrs: 'map[]'?subject=&amp;body=c%3D%28d%29' attrs: 'map[]'")
		})

		Convey("Substitute link inline macro with text having uri inside: should return modified link text without uri", func() {
			s.Document().(*testSubstDocumentAble).linkAttrs = false
			s.Document().(*testSubstDocumentAble).hideUriScheme = true
			So(s.SubMacros("link:http://a[]"), ShouldEqual, "ContextAn 'anchor': text '' ===> type 'link' target 'http://a' attrs: 'map[]'")
		})
	})

	Convey("A substitutors can substitute email inline macro references", t, func() {
		s := &substitutors{}
		testDocument := newTestSubstDocumentAble(s)
		tim := &testInlineMacro{}
		testDocument.te.inlineMacros = append(testDocument.te.inlineMacros, tim)
		s.document = testDocument
		s.inlineMaker = &testInlineMaker{}
		s.attributeListMaker = &testAttributeListMaker{}
		Convey("Substitute escaped email link inline macro should ignore the escape", func() {
			So(s.SubMacros("\\doc.writer@test.com[]"), ShouldEqual, "ContextAn 'anchor': text 'doc.writer@test.com' ===> type 'link' target 'mailto:doc.writer@test.com' attrs: 'map[]'[]")
		})
		Convey("Substitute email inline macro with mailto: should return mailto: target", func() {
			So(s.SubMacros("doc.writer@test2.com[] "), ShouldEqual, "ContextAn 'anchor': text 'doc.writer@test2.com' ===> type 'link' target 'mailto:doc.writer@test2.com' attrs: 'map[]'[] ")
		})
	})

	Convey("A substitutors can restore passthrough", t, func() {
		s := &substitutors{}
		Convey("By default, empty passthrough means text is returned unchanged", func() {
			So(s.restorePassthroughs("test"), ShouldEqual, "test")
		})
		Convey("non-empty passthrough apply subs", func() {
			p := &passthrough{}
			// []*subsEnum
			p.subs = append(p.subs, sub.title)
			p.text = "test"
			p.typePT = "visible"
			// not really needed because testInlineMaker doesn't use parent
			s.abstractNodable = &testAbstractNodable{}
			s.inlineMaker = &testInlineMaker{}
			s.passthroughs = append(s.passthroughs, p)
			So(s.restorePassthroughs("abc\u00960\u0097def"), ShouldEqual, "abcContextQt 'quoted': text 'test' ===> type 'visible' target '' id '' attrs: 'map[]'def")
		})
	})

	Convey("A substitutors can Substitute normal and bibliographic anchors", t, func() {
		s := &substitutors{}
		Convey("Substitute normal anchor '[[['", func() {
			s.inlineMaker = &testInlineMaker{}
			So(s.subInlineAnchors(`\[[[test]]]`, nil), ShouldEqual, "[ContextAn 'anchor': text 'test' ===> type 'ref' target 'test' attrs: 'map[]']")
			So(s.subInlineAnchors("[[[test]]]", nil), ShouldEqual, "ContextAn 'anchor': text 'test' ===> type 'bibref' target 'test' attrs: 'map[]'")
		})
		Convey("Substitute ref anchor '[['", func() {
			So(s.subInlineAnchors(`\[[testref]]`, nil), ShouldEqual, "[[testref]]")
			testDocument := newTestSubstDocumentAble(s)
			testDocument.references = &testReferencable{}
			s.document = testDocument
			So(s.subInlineAnchors(`[[testref]]`, nil), ShouldEqual, "ContextAn 'anchor': text 'testref' ===> type 'ref' target 'testref' attrs: 'map[]'")
			So(s.subInlineAnchors(`[[testref2]]`, nil), ShouldEqual, "ContextAn 'anchor': text 'testref2' ===> type 'ref' target 'testref2' attrs: 'map[]'")
		})
	})

	Convey("A substitutors can Substitute cross reference links", t, func() {
		s := &substitutors{}
		Convey("Substitute <<id,reftext>>", func() {
			s.inlineMaker = &testInlineMaker{}
			So(s.subInlineXrefs(`\&lt;&lt;id1,reftext&gt;&gt;`, nil), ShouldEqual, "&lt;&lt;id1,reftext&gt;&gt;")
			So(s.subInlineXrefs(`&lt;&lt;id2,reftext2&gt;&gt;`, nil), ShouldEqual, "ContextAn 'anchor': text 'reftext2' ===> type 'xref' target '#' attrs: 'map[path: fragment: refid:]'")
		})
		Convey("Substitute xref:id[reftext]", func() {
			So(s.subInlineXrefs(`\xref:id3[reftext3]`, nil), ShouldEqual, "xref:id3[reftext3]")
			So(s.subInlineXrefs(`xref:id4[reftext4]`, nil), ShouldEqual, "ContextAn 'anchor': text 'reftext4' ===> type 'xref' target '#' attrs: 'map[path: fragment: refid:]'")
		})
		Convey("Substitute xref:id#xx[reftext]", func() {
			So(s.subInlineXrefs(`xref:id5#xxx5[reftext5]`, nil), ShouldEqual, "ContextAn 'anchor': text 'reftext5' ===> type 'xref' target '' attrs: 'map[path:id5 fragment:xxx5 refid:]'")
		})
		Convey("Substitute xref:doc.adoc#xx[reftext]", func() {
			testDocument := newTestSubstDocumentAble(s)
			testDocument.references = &testReferencable{}
			s.document = testDocument
			So(s.subInlineXrefs(`xref:doc6.adoc6#xxx6[reftext6]`, nil), ShouldEqual, "ContextAn 'anchor': text 'reftext6' ===> type 'xref' target 'relfileprefixAttrdoc6.html#xxx6' attrs: 'map[path:relfileprefixAttrdoc6.html fragment:xxx6 refid:doc6#xxx6]'")

			So(s.subInlineXrefs(`xref:doc7.adoc7#[reftext7]`, nil), ShouldEqual, "ContextAn 'anchor': text 'reftext7' ===> type 'xref' target 'relfileprefixAttrdoc7.html' attrs: 'map[path:relfileprefixAttrdoc7.html fragment: refid:doc7]'")
			So(s.subInlineXrefs(`xref:doc8.adoc8#frag8[reftext8]`, nil), ShouldEqual, "ContextAn 'anchor': text 'reftext8' ===> type 'xref' target '#frag8' attrs: 'map[path: fragment:frag8 refid:frag8]'")
		})
	})

	Convey("A substitutors can substitute footnote inline macro references", t, func() {
		s := &substitutors{}
		testDocument := newTestSubstDocumentAble(s)
		tim := &testInlineMacro{}
		testDocument.te.inlineMacros = append(testDocument.te.inlineMacros, tim)
		s.document = testDocument
		s.inlineMaker = &testInlineMaker{}
		s.attributeListMaker = &testAttributeListMaker{}

		p := &passthrough{}
		// []*subsEnum
		p.subs = append(p.subs, sub.title)
		p.text = "test6"
		p.typePT = "visible"
		s.passthroughs = append(s.passthroughs, p)

		Convey("Substitute escaped footnote link inline macro should ignore the escape", func() {
			So(s.SubMacros("test \\footnoteref:[id,text] ww\n \\footnote:[text]hh\nq q\\footnoteref:[id] ww"), ShouldEqual, "test footnoteref:[id,text] ww\n footnote:[text]hh\nq qfootnoteref:[id] ww")
		})
		Convey("Substitute footnote:xx", func() {
			So(s.SubMacros("test footnote:[text2] ww\n ss"), ShouldEqual, "test ContextFt 'footnote': text 'text2' ===> type '' target '' id '' attrs: 'map[index:1]' ww\n ss")
			// Normalize string: strips outboud spaces
			So(s.SubMacros("test footnote:[  text3  aa ] ww\n ss"), ShouldEqual, "test ContextFt 'footnote': text 'text3  aa' ===> type '' target '' id '' attrs: 'map[index:2]' ww\n ss")
			// TOFIX? Is it normal that normalizedString='anchor:idname4[Reference4 Text4' and subInlineAnchors='anchor:idname4[Reference4 Text4'? Last ] is missing, which should prevent anchor to be detected... Yet it appears to besubstitute anyway
			So(s.SubMacros("test footnote:[anchor:idname4[Reference4 Text4]] ww\n ss"), ShouldEqual, "test ContextFt 'footnote': text 'ContextAn 'anchor': text 'Reference4 Text4' ===> type '' target '' id '' attrs: 'map[index:3' ===> type 'ref' target 'idname4' attrs: 'map[]''] ww\n ss")
			// At least, subInlineXrefs is detectable
			So(s.SubMacros("test footnote:[&lt;&lt;id5,reftext5&gt;&gt;] ww\n ss"), ShouldEqual, "test ContextFt 'footnote': text 'ContextAn 'anchor': text 'reftext5' ===> type 'xref' target '#' attrs: 'map[path: fragment: refid:]'' ===> type '' target '' id '' attrs: 'map[index:4]' ww\n ss")
			// Restore passthrough works too
			So(s.SubMacros("test footnote:[abc6\u00960\u0097def6] ww\n ss"), ShouldEqual, "test ContextFt 'footnote': text 'abc6ContextQt 'quoted': text 'test6' ===> type 'visible' target '' id '' attrs: 'map[]'def6' ===> type '' target '' id '' attrs: 'map[index:5]' ww\n ss")
		})
		Convey("Substitute footnote:xx", func() {
			So(s.SubMacros("test footnote:[text2] ww\n ss"), ShouldEqual, "test ContextFt 'footnote': text 'text2' ===> type '' target '' id '' attrs: 'map[index:6]' ww\n ss")
			So(len(testDocument.footnotes), ShouldEqual, 6)
			footnote := testDocument.footnotes[len(testDocument.footnotes)-1]
			So(footnote.String(), ShouldEqual, "footnote(0,6): 'text2'")
		})
		Convey("Substitute footnoteref:[id]", func() {
			So(s.SubMacros("test footnoteref:[4] ww\n ss"), ShouldEqual, "test ContextFt 'footnote': text '4' ===> type 'xref' target '4' id '4' attrs: 'map[index:]' ww\n ss")
		})
		Convey("Substitute footnoteref:[id,text]", func() {
			So(s.SubMacros("test footnoteref:[4,text4] ww4\n ss"), ShouldEqual, "test ContextFt 'footnote': text 'text4' ===> type 'ref' target '' id '4' attrs: 'map[index:7]' ww4\n ss")
			So(len(testDocument.footnotes), ShouldEqual, 7)
			footnote := testDocument.footnotes[len(testDocument.footnotes)-1]
			So(footnote.String(), ShouldEqual, "footnote(4,7): 'text4'")
			So(s.SubMacros("test footnoteref:[4] ww\n ss"), ShouldEqual, "test ContextFt 'footnote': text 'text4' ===> type 'xref' target '4' id '4' attrs: 'map[index:7]' ww\n ss")
		})
	})
	Convey("A substitutors can resolve Subs references", t, func() {

		Convey("Empty subs means empty candidates", func() {
			So(len(resolveSubs("", nil, nil, "")), ShouldEqual, 0)
		})
		Convey("Empty subs means empty candidates", func() {
			So(len(resolveSubs(",", nil, nil, "")), ShouldEqual, 0)
		})
		Convey("No subs means empty candidates", func() {
			So(len(resolveSubs("test", nil, nil, "")), ShouldEqual, 0)
		})
		Convey("No defaults, no subOptions[key] means no subs returned, invalid warning", func() {
			So(len(resolveSubs("basic", nil, nil, "")), ShouldEqual, 0)
		})
		Convey("No defaults, compositeSub key means same composite sub value returned", func() {
			c := resolveSubs("specialchars", nil, nil, "")
			So(len(c), ShouldEqual, 1)
			So(fmt.Sprintf("%s", c), ShouldEqual, "[%!s(*asciidocgo.subsEnum=&{specialcharacters})]")
			c = resolveSubs("verbatim", nil, nil, "")
			So(len(c), ShouldEqual, 2)
			So(fmt.Sprintf("%s", c), ShouldEqual, "[%!s(*asciidocgo.subsEnum=&{specialcharacters}) %!s(*asciidocgo.subsEnum=&{callouts})]")
		})
		Convey("No defaults, inline type for verbatim means returned", func() {
			c := resolveSubs("verbatim", subOption.inline, nil, "")
			So(len(c), ShouldEqual, 1)
			So(fmt.Sprintf("%s", c), ShouldEqual, "[%!s(*asciidocgo.subsEnum=&{specialcharacters})]")
		})
		Convey("Defaults, no operation", func() {
			c := resolveSubs("verbatim", subOption.inline, subArray{sub.title}, "")
			So(len(c), ShouldEqual, 1)
			So(fmt.Sprintf("%s", c), ShouldEqual, "[%!s(*asciidocgo.subsEnum=&{specialcharacters})]")
		})
		Convey("Defaults, multiple candidates, no operation, no intersaction with subOptions", func() {
			c := resolveSubs("basic", subOption.inline, subArray{sub.title}, "")
			So(len(c), ShouldEqual, 0)
			So(fmt.Sprintf("%s", c), ShouldEqual, "[]")
		})

		Convey("Defaults, multiple candidates, no operation, no intersaction with subOptions", func() {
			c := resolveSubs("normal", subOption.inline, subArray{sub.title}, "testx")
			So(len(c), ShouldEqual, 6)
			So(fmt.Sprintf("%s", c), ShouldEqual, "[%!s(*asciidocgo.subsEnum=&{specialcharacters}) %!s(*asciidocgo.subsEnum=&{quotes}) %!s(*asciidocgo.subsEnum=&{attributes}) %!s(*asciidocgo.subsEnum=&{replacements}) %!s(*asciidocgo.subsEnum=&{macros}) %!s(*asciidocgo.subsEnum=&{post_replacements})]")
		})

		Convey("Subs with append, followed by invalid entry, Defaults, multiple candidates, no operation", func() {
			c := resolveSubs("+normal,title", subOption.inline, subArray{sub.title}, "aSubject")
			So(len(c), ShouldEqual, 6)
			So(fmt.Sprintf("%s", c), ShouldEqual, "[%!s(*asciidocgo.subsEnum=&{specialcharacters}) %!s(*asciidocgo.subsEnum=&{quotes}) %!s(*asciidocgo.subsEnum=&{attributes}) %!s(*asciidocgo.subsEnum=&{replacements}) %!s(*asciidocgo.subsEnum=&{macros}) %!s(*asciidocgo.subsEnum=&{post_replacements})]")
		})
		Convey("Subs with append, Defaults, multiple candidates, no operation", func() {
			c := resolveSubs("+normal", subOption.inline, subArray{sub.title}, "Subs with append")
			So(len(c), ShouldEqual, 6)
			So(fmt.Sprintf("%s", c), ShouldEqual, "[%!s(*asciidocgo.subsEnum=&{specialcharacters}) %!s(*asciidocgo.subsEnum=&{quotes}) %!s(*asciidocgo.subsEnum=&{attributes}) %!s(*asciidocgo.subsEnum=&{replacements}) %!s(*asciidocgo.subsEnum=&{macros}) %!s(*asciidocgo.subsEnum=&{post_replacements})]")
		})

		Convey("Subs with remove, Defaults, multiple candidates, no operation", func() {
			c := resolveSubs("-v", subOption.block, subArray{sub.normal}, "Subs with remove")
			So(len(c), ShouldEqual, 1)
			So(fmt.Sprintf("%s", c), ShouldEqual, "[%!s(*asciidocgo.subsEnum=&{normal})]")
			c = resolveBlockSubs("-v", subArray{sub.normal}, "Subs with remove 2")
			So(len(c), ShouldEqual, 1)
			So(fmt.Sprintf("%s", c), ShouldEqual, "[%!s(*asciidocgo.subsEnum=&{normal})]")
		})
		Convey("Subs with prepend, Defaults, multiple candidates, no operation", func() {
			c := resolveSubs("normal+", subOption.inline, subArray{sub.title, sub.pass}, "Subs with prepend")
			So(len(c), ShouldEqual, 6)
			So(fmt.Sprintf("%s", c), ShouldEqual, "[%!s(*asciidocgo.subsEnum=&{specialcharacters}) %!s(*asciidocgo.subsEnum=&{quotes}) %!s(*asciidocgo.subsEnum=&{attributes}) %!s(*asciidocgo.subsEnum=&{replacements}) %!s(*asciidocgo.subsEnum=&{macros}) %!s(*asciidocgo.subsEnum=&{post_replacements})]")
		})
		Convey("Subs with append, with resolved symbol", func() {
			c := resolveSubs("+C", subOption.inline, subArray{sub.normal}, "Subs with remove")
			So(len(c), ShouldEqual, 2)
			So(fmt.Sprintf("%s", c), ShouldEqual, "[%!s(*asciidocgo.subsEnum=&{normal}) %!s(*asciidocgo.subsEnum=&{specialcharacters})]")
		})
	})
}
