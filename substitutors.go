package asciidocgo

import (
	"fmt"
	"strings"

	"github.com/VonC/asciidocgo/consts/regexps"
)

type _sub string

const (
	subsBasic    _sub = "basic"
	subsNormal   _sub = "normal"
	subsVerbatim _sub = "verbatim"
	subsTitle    _sub = "title"
	subsHeader   _sub = "header"
	subsPass     _sub = "pass"
	subsUnknown  _sub = "unknown"
)
const (
	subsSpecialCharacters _sub = "specialcharacters"
	subsQuotes            _sub = "quotes"
	subsAttributes        _sub = "attributes"
	subsReplacements      _sub = "replacements"
	subsMacros            _sub = "macros"
	subsPostReplacements  _sub = "post_replacements"
	subsCallout           _sub = "callouts"
)
const (
	subsNone         _sub = "none"
	subsSpecialChars _sub = "specialchars"
)
const (
	subsA _sub = "a"
	subsM _sub = "m"
	subsN _sub = "n"
	subsP _sub = "p"
	subsQ _sub = "q"
	subsR _sub = "R"
	subsC _sub = "C"
	subsV _sub = "V"
)
const (
	subsBlock  _sub = "block"
	subsInline _sub = "inline"
)

var testsub = ""

type subsEnum struct {
	value _sub
}

type subsEnums struct {
	basic    *subsEnum
	normal   *subsEnum
	verbatim *subsEnum
	title    *subsEnum
	header   *subsEnum
	pass     *subsEnum
	unknown  *subsEnum
}

type subsEnumsValues struct {
	specialcharacters *subsEnum
	quotes            *subsEnum
	attributes        *subsEnum
	replacements      *subsEnum
	macros            *subsEnum
	postReplacements  *subsEnum
	callouts          *subsEnum
}

type compositeSubsEnums struct {
	none         *subsEnum
	normal       *subsEnum
	verbatim     *subsEnum
	specialchars *subsEnum
}

type subSymbolsEnums struct {
	a *subsEnum
	m *subsEnum
	n *subsEnum
	p *subsEnum
	q *subsEnum
	r *subsEnum
	c *subsEnum
	v *subsEnum
}

type subOptionsEnums struct {
	block  *subsEnum
	inline *subsEnum
}

func newSubsEnums() *subsEnums {
	return &subsEnums{
		&subsEnum{subsBasic},
		&subsEnum{subsNormal},
		&subsEnum{subsVerbatim},
		&subsEnum{subsTitle},
		&subsEnum{subsHeader},
		&subsEnum{subsPass},
		&subsEnum{subsUnknown}}
}

func newSubsEnumsValues() *subsEnumsValues {
	return &subsEnumsValues{
		&subsEnum{subsSpecialCharacters},
		&subsEnum{subsQuotes},
		&subsEnum{subsAttributes},
		&subsEnum{subsReplacements},
		&subsEnum{subsMacros},
		&subsEnum{subsPostReplacements},
		&subsEnum{subsCallout}}
}

func newCompositeSubsEnums() *compositeSubsEnums {
	return &compositeSubsEnums{
		&subsEnum{subsNone},
		&subsEnum{subsNormal},
		&subsEnum{subsVerbatim},
		&subsEnum{subsSpecialChars}}
}

func newSubSymbolsEnums() *subSymbolsEnums {
	return &subSymbolsEnums{
		&subsEnum{subsA},
		&subsEnum{subsM},
		&subsEnum{subsN},
		&subsEnum{subsP},
		&subsEnum{subsQ},
		&subsEnum{subsR},
		&subsEnum{subsC},
		&subsEnum{subsV}}
}

func newSubOptionsEnums() *subOptionsEnums {
	return &subOptionsEnums{
		&subsEnum{subsBlock},
		&subsEnum{subsInline}}
}

var sub = newSubsEnums()
var subValue = newSubsEnumsValues()
var compositeSub = newCompositeSubsEnums()
var subSymbol = newSubSymbolsEnums()
var subOption = newSubOptionsEnums()

type subArray []*subsEnum

func (cses *compositeSubsEnums) keys() subArray {
	res := subArray{}
	res = append(res, cses.none)
	res = append(res, cses.normal)
	res = append(res, cses.verbatim)
	res = append(res, cses.specialchars)
	return res
}

var subs = map[*subsEnum]subArray{
	sub.basic:    subArray{subValue.specialcharacters},
	sub.normal:   subArray{subValue.specialcharacters, subValue.quotes, subValue.attributes, subValue.replacements, subValue.macros, subValue.postReplacements},
	sub.verbatim: subArray{subValue.specialcharacters, subValue.callouts},
	sub.title:    subArray{subValue.specialcharacters, subValue.quotes, subValue.replacements, subValue.macros, subValue.attributes, subValue.postReplacements},
	sub.header:   subArray{subValue.specialcharacters, subValue.attributes},
	sub.pass:     subArray{},
}
var compositeSubs = map[*subsEnum]subArray{
	compositeSub.none:         subArray{},
	compositeSub.normal:       subs[sub.normal],
	sub.normal:                subs[sub.normal],
	compositeSub.verbatim:     subs[sub.verbatim],
	compositeSub.specialchars: subArray{subValue.specialcharacters},
}
var subSymbols = map[*subsEnum]subArray{
	subSymbol.a: subArray{subValue.attributes},
	subSymbol.m: subArray{subValue.macros},
	subSymbol.n: subArray{sub.normal},
	subSymbol.p: subArray{subValue.postReplacements},
	subSymbol.q: subArray{subValue.quotes},
	subSymbol.r: subArray{subValue.replacements},
	subSymbol.c: subArray{subValue.specialcharacters},
	subSymbol.v: subArray{sub.verbatim},
}
var subOptions = map[*subsEnum]subArray{
	subOption.block:  append(append(compositeSub.keys(), subs[sub.normal]...), subValue.callouts),
	subOption.inline: append(compositeSub.keys(), subs[sub.normal]...),
}

func (se *subsEnum) isCompositeSub() bool {
	if _, ok := compositeSubs[se]; ok {
		return true
	}
	return false
}

func values(someSubs subArray) []string {
	res := []string{}
	for _, aSub := range someSubs {
		res = append(res, string(aSub.value))
	}
	return res
}

func (sa subArray) include(s *subsEnum) bool {
	for _, aSub := range sa {
		if aSub == s {
			return true
		}
	}
	return false
}

type passthrough struct {
	text string
	subs subArray
}

/* Methods to perform substitutions on lines of AsciiDoc text.
This module is intented to be mixed-in to Section and Block to provide
operations for performing the necessary substitutions. */
type substitutors struct {
	// A String Array of passthough (unprocessed) text captured from this block
	passthroughs []passthrough
}

/* Apply the specified substitutions to the lines of text

source  - The String or String Array of text to process
subs    - The substitutions to perform. Can be a Symbol or a Symbol Array (default: :normal)
expand -  A Boolean to control whether sub aliases are expanded (default: true)

returns Either a String or String Array, whichever matches the type of the first argument */
func (s *substitutors) ApplySubs(source string, someSubs subArray) string {
	text := ""
	var allSubs subArray
	if len(someSubs) == 1 {
		if someSubs[0] == sub.pass {
			return source
		}
		if someSubs[0] == sub.unknown {
			return text
		}
	}
	for _, aSub := range someSubs {
		if aSub.isCompositeSub() {
			allSubs = append(allSubs, compositeSubs[aSub]...)
		} else {
			allSubs = append(allSubs, aSub)
		}
	}
	if testsub == "test_ApplySubs_allsubs" {
		return fmt.Sprintf("%v", values(allSubs))
	}
	if len(allSubs) == 0 {
		return source
	}
	text = source
	if allSubs.include(subValue.macros) {
		text = s.extractPassthroughs(text)
	}
	if testsub == "test_ApplySubs_extractPassthroughs" {
		return text
	}
	// TODO complete (s *substitutors) ApplySubs after extractPassthroughs
	return text
}

// Delimiters and matchers for the passthrough placeholder
// See http://www.aivosto.com/vbtips/control-characters.html#listabout
// for characters to use

const (
	// SPA, start of guarded protected area (\u0096)
	subPASS_START = "\u0096"

	// EPA, end of guarded protected area (\u0097)
	subPASS_END = "\u0097"
)

/* Extract the passthrough text from the document for reinsertion after processing.
text - The String from which to extract passthrough fragements
returns - The text with the passthrough region substituted with placeholders */
func (s *substitutors) extractPassthroughs(text string) string {
	res := text
	if strings.Contains(res, "++") || strings.Contains(res, "$$") || strings.Contains(res, "ss:") {
		reres := regexps.NewPassInlineMacroRxres(res)
		if !reres.HasAnyMatch() {
			goto Next
		}
		res = ""
		for reres.HasNext() {
			res = res + reres.Prefix()
			textOri := ""
			subsOri := subArray{}
			if reres.IsEscaped() {
				// honor the escape
				// meaning don't transform anything, but loose the escape
				res = res + reres.FullMatch()[1:]
			} else if reres.HasPassText() {
				textOri = unescapeBrackets(reres.PassText())
				if reres.HasPassSub() {
					subsOri = resolvePassSubs(reres.PassSub())
				}
			} else {
				textOri = reres.InlineText()
				if reres.InlineSub() == "$$" {
					subsOri = subArray{subValue.specialcharacters}
				}
			}
			if textOri != "" {
				p := passthrough{textOri, subsOri}
				s.passthroughs = append(s.passthroughs, p)
				index := len(s.passthroughs) - 1
				res = res + fmt.Sprintf("%s%d%s", subPASS_START, index, subPASS_END)
			}
			reres.Next()
		}
		res = res + reres.Suffix()
	}
Next:
	return res
}

/* Internal: Unescape closing square brackets.
   Intended for text extracted from square brackets. */
func unescapeBrackets(str string) string {
	// FIXME make \] a regex
	if str == "" {
		return str
	}
	str = regexps.EscapedBracketRx.ReplaceAllString(str, "]")
	return str
}

func resolvePassSubs(str string) subArray {
	// resolve_subs subs, :inline, nil, 'passthrough macro'
	return subArray{}
}
