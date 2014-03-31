package asciidocgo

import (
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/VonC/asciidocgo/consts/compliance"
	"github.com/VonC/asciidocgo/consts/context"
	"github.com/VonC/asciidocgo/consts/regexps"
	"github.com/VonC/asciidocgo/consts/regexps/quotes"
	"github.com/VonC/asciidocgo/debug"
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

type InlineMacroable interface {
	IsShortFormat() bool
	IsContentModelAttributes() bool
	Regexp() *regexp.Regexp
	ProcessMethod(self interface{}, target string, attributes map[string]interface{}) string
	PosAttrs() []string
}

type Extensionables interface {
	HasInlineMacros() bool
	InlineMacros() []InlineMacroable
}

type SubstDocumentable interface {
	Attr(name string, defaultValue interface{}, inherit bool) interface{}
	Basebackend(base interface{}) bool
	SubAttributes(data string, opts *OptionsParseAttributes) string
	Counter(name string, seed int) string
	HasAttr(name string, expect interface{}, inherit bool) bool
	Extensions() Extensionables
	Register(typeDoc string, value []string)
	References() Referencable
}

type Referencable interface {
	HasId(id string) bool
	Get(id string) string
}

type Convertable interface {
	Convert() string
}
type AbstractNodable interface {
	IsAbstractNodable()
}

type OptionsInline struct {
	id         string
	typeInline string
	target     string
	attributes map[string]interface{}
}

/*func (oi *OptionsInline) Id() string                         { return oi.id }*/
func (oi *OptionsInline) TypeInline() string                 { return oi.typeInline }
func (oi *OptionsInline) Target() string                     { return oi.target }
func (oi *OptionsInline) Attributes() map[string]interface{} { return oi.attributes }

type InlineMaker interface {
	NewInline(parent AbstractNodable, c context.Context, text string, opts *OptionsInline) Convertable
}

type passthrough struct {
	text       string
	subs       subArray
	attributes map[string]interface{}
	typePT     string
}

type AttributeListable interface {
	ParseInto(into map[string]interface{}, posAttrs []string) map[string]interface{}
	Parse(posAttrs []string) map[string]interface{}
}

type ApplyNormalSubsable interface {
	ApplyNormalSubs(lines string) string
}

type AttributeListMaker interface {
	NewAttributeList(attrline string, block ApplyNormalSubsable, delimiter string) AttributeListable
}

/* Methods to perform substitutions on lines of AsciiDoc text.
This module is intented to be mixed-in to Section and Block to provide
operations for performing the necessary substitutions. */
type substitutors struct {
	// A String Array of passthough (unprocessed) text captured from this block
	passthroughs       []*passthrough
	document           SubstDocumentable
	inlineMaker        InlineMaker
	abstractNodable    AbstractNodable
	attributeListMaker AttributeListMaker
}

func (s *substitutors) Document() SubstDocumentable {
	return s.document
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
	for _, aSub := range allSubs {
		switch aSub.value {
		case "specialcharacters":
			text = subSpecialCharacters(text)
		case "quotes":
			text = subQuotes(text)
		case "attributes":
			text = s.SubAttributes(text, nil)
		case "replacements":
			text = subReplacements(text)
		case "macros":
			text = s.SubMacros(text)
			/*
				case "highlight":
				case "callouts":
				case "post_replacements":
			*/
		}
	}
	if testsub == "test_ApplySubs_applyAllsubs" {
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
			goto PassInlineLiteralRx
		}
		res = ""
		suffix := ""
		for reres.HasNext() {
			res = res + reres.Prefix()
			textOri := ""
			subsOri := subArray{}
			if reres.IsEscaped() {
				// honor the escape
				// meaning don't transform anything, but loose the escape
				res = res + reres.FullMatch()[1:]
				suffix = reres.Suffix()
				reres.Next()
				continue
			}
			if reres.HasPassText() {
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
				p := &passthrough{textOri, subsOri, make(map[string]interface{}), ""}
				s.passthroughs = append(s.passthroughs, p)
				index := len(s.passthroughs) - 1
				res = res + fmt.Sprintf("%s%d%s", subPASS_START, index, subPASS_END)
			}
			suffix = reres.Suffix()
			reres.Next()
		}
		res = res + suffix
	}
PassInlineLiteralRx:

	if strings.Contains(res, "`") {

		reres := regexps.NewPassInlineLiteralRxres(res)
		if !reres.HasAnyMatch() {
			goto MathInlineMacroRx
		}

		res = ""
		suffix := ""
		for reres.HasNext() {
			res = res + reres.Prefix()

			unescaped_attrs := ""
			// honor the escape
			if reres.Literal()[0] == '\\' {
				//fmt.Printf("======== %v=====\n", reres.FullMatch())
				res = res + reres.FirstChar() + reres.Attributes() + reres.Literal()[1:] + " : " + reres.FirstChar() + reres.Literal()[1:]
				suffix = reres.Suffix()
				reres.Next()
				continue
			}
			if reres.IsEscaped() && reres.Attributes() != "" {
				unescaped_attrs = "[" + reres.Attributes() + "]"
				res = res + unescaped_attrs
			} else {
				res = res + reres.FirstChar()
			}

			attributes := make(map[string]interface{})
			if unescaped_attrs != "" && reres.Attributes() != "" {
				attributes = s.parseAttributes(reres.Attributes(), []string{}, &OptionsParseAttributes{})
			}

			p := &passthrough{reres.LiteralText(), subArray{subValue.specialcharacters}, attributes, "monospaced"}
			s.passthroughs = append(s.passthroughs, p) //TODO attributes, type (later, to make them type safe instead of hash)
			index := len(s.passthroughs) - 1
			res = res + fmt.Sprintf("%s%d%s", subPASS_START, index, subPASS_END)

			suffix = reres.Suffix()
			reres.Next()
		}
		res = res + suffix

	}

MathInlineMacroRx:

	if strings.Contains(res, "math:") {
		reres := regexps.NewMathInlineMacroRxres(res)
		if !reres.HasAnyMatch() {
			goto ExtractPassthroughsRes
		}

		res = ""
		suffix := ""
		for reres.HasNext() {
			res = res + reres.Prefix()

			if reres.IsEscaped() {
				// honor the escape
				// meaning don't transform anything, but loose the escape
				res = res + reres.FullMatch()[1:]
				suffix = reres.Suffix()
				reres.Next()
				continue
			}

			mathType := reres.MathType()
			if mathType == "math" {
				defaultType := "asciimath"
				if s.Document() != nil {
					defaultTypeI := s.Document().Attr("math", nil, false)
					if defaultTypeI != nil && defaultTypeI.(string) != "" {
						defaultType = defaultTypeI.(string)
					}
				}
				mathType = defaultType
			}
			mathText := unescapeBrackets(reres.MathText())
			mathSubs := subArray{}
			if reres.MathSub() == "" {
				if s.Document() != nil && s.Document().Basebackend("html") {
					mathSubs = subArray{subValue.specialcharacters}
				} else {
					mathSubs = resolvePassSubs(reres.MathSub())
				}
			}
			attributes := make(map[string]interface{})
			p := &passthrough{mathText, mathSubs, attributes, mathType}
			s.passthroughs = append(s.passthroughs, p)
			index := len(s.passthroughs) - 1
			res = res + fmt.Sprintf("%s%d%s", subPASS_START, index, subPASS_END)

			suffix = reres.Suffix()
			reres.Next()
		}
		res = res + suffix
	}

ExtractPassthroughsRes:

	return res
}

var PASS_MATCHRx, _ = regexp.Compile("\u0096" + `(\d+)` + "\u0097")

/* Internal: Restore the passthrough text by reinserting into the placeholder positions
text - The String text into which to restore the passthrough text
returns The String text with the passthrough text restored */
func (s *substitutors) restorePassthroughs(text string) string {
	res := text
	if s.passthroughs == nil || len(s.passthroughs) == 0 || !strings.Contains(text, subPASS_START) {
		return res
	}
	fmt.Printf("\n%v => %v\n", s.passthroughs, len(s.passthroughs))
	res = ""
	suffix := ""
	reres := regexps.NewReres(text, PASS_MATCHRx)
	for reres.HasNext() {
		res = res + reres.Prefix()
		index, _ := strconv.Atoi(reres.Group(1))
		pass := s.passthroughs[index]
		subs := pass.subs
		subbedText := pass.text
		//fmt.Printf("\nrestorePassthroughs subs '%v', index '%v' text '%v'\n", subs, index, subbedText)
		if subs != nil {
			subbedText = s.ApplySubs(subbedText, subs)
		}
		//fmt.Printf("subbedText '%v'\n", subbedText)
		typePT := pass.typePT
		if typePT != "" {
			optsInline := &OptionsInline{attributes: pass.attributes}
			optsInline.typeInline = typePT
			inline := s.inlineMaker.NewInline(s.abstractNodable, context.Quoted, subbedText, optsInline)
			res = res + inline.Convert()
		}
		suffix = reres.Suffix()
		reres.Next()
	}
	res = res + suffix
	return res
}

var specialCharacterPatternRx, _ = regexp.Compile(`[&<>]`)

type specialCharacterPatternRxRes struct {
	*regexps.Reres
}

/* Substitute special characters (i.e., encode XML)
Special characters are defined in the Asciidoctor::SPECIAL_CHARS Array constant

 text - The String text to process
 returns The String text with special characters replaced */
func subSpecialCharacters(text string) string {
	reres := &specialCharacterPatternRxRes{regexps.NewReres(text, specialCharacterPatternRx)}

	if !reres.HasAnyMatch() {
		return text
	}
	res := ""
	suffix := ""
	for reres.HasNext() {
		res = res + reres.Prefix()
		switch reres.FullMatch() {
		case "&":
			res = res + "&amp;"
		case "<":
			res = res + "&lt;"
		case ">":
			res = res + "&gt;"
		}
		suffix = reres.Suffix()
		reres.Next()
	}
	res = res + suffix
	return res
}

/* Substitute quoted text (includes emphasis, strong, monospaced, etc)
 text - The String text to process
returns The String text with quoted text rendered using
the backend templates */
func subQuotes(text string) string {
	result := text
	//fmt.Printf("subQuotes result='%v'\n", result)
	for _, qs := range quotes.QuoteSubs {
		//fmt.Printf("subQuotes rx='%v' on '%v' (%v)\n", qs.Rx(), result, qs.Constrained())
		match := quotes.NewQuoteSubRxres(result, qs)
		result = transformQuotedText(match, qs.TypeQS(), qs.Constrained())
	}
	return result
}

/* Substitute replacement characters (e.g., copyright, trademark, etc)
 text - The String text to process
returns The String text with the replacement characters substituted */
func subReplacements(text string) string {
	result := text
	for _, repl := range regexps.Replacements {
		reres := repl.Reres(result)
		if reres.HasAnyMatch() {
			result = ""
			suffix := ""
			for reres.HasNext() {
				result = result + reres.Prefix()
				result = result + doReplacement(reres, repl)
				suffix = reres.Suffix()
				reres.Next()
			}
			result = result + suffix
		}
	}
	return result
}

func doReplacement(reres *regexps.Reres, repl *regexps.Replacement) string {
	res := ""
	if reres.IsEscaped() {
		res = reres.FullMatch()[1:]
	} else if reres.HasGroup(2) && reres.Group(2)[0] == '\\' {
		res = reres.Group(1) + reres.Group(2)[1:]
	} else if repl.None() {
		res = repl.Repl()
	} else if repl.Leading() {
		res = reres.Group(1) + repl.Repl()
	} else if repl.Bounding() {
		res = reres.Group(1) + repl.Repl() + reres.Group(2)
	}
	return res
}

var intrinsicAttributes = map[string]rune{
	"startsb":        '[',
	"endsb":          ']',
	"vbar":           '|',
	"caret":          '^',
	"asterisk":       '*',
	"tilde":          '~',
	"plus":           43,
	"apostrophe":     '\'',
	"backslash":      '\\',
	"backtick":       '`',
	"empty":          0,
	"sp":             ' ',
	"space":          ' ',
	"two-colons":     ':',
	"two-semicolons": ';',
	"nbsp":           160,
	"deg":            176,
	"zwsp":           8203,
	"quot":           34,
	"apos":           39,
	"lsquo":          8216,
	"rsquo":          8217,
	"ldquo":          8220,
	"rdquo":          8221,
	"wj":             8288,
	"brvbar":         166,
	"amp":            '&',
	"lt":             '<',
	"gt":             '>',
}

/* Public: Substitute attribute references
Attribute references are in the format +{name}+.
If an attribute referenced in the line is missing, the line is dropped.
# text     - The String text to process
returns The String text with the attribute references replaced with attribute values
--
NOTE it's necessary to perform this substitution line-by-line
so that a missing key doesn't wipe out the whole block of data */
func (s *substitutors) SubAttributes(data string, opts *OptionsParseAttributes) string {
	if data == "" {
		return data
	}
	lines := strings.Split(data, "\n")
	res := ""
	for i, line := range lines {
		reject := false
		reject_if_empty := false
		lineres := line
		if strings.Contains(line, "{") {
			reres := regexps.NewAttributeReferenceRxres(line)
			if !reres.HasAnyMatch() {
				if i > 0 {
					res = res + "\n"
				}
				res = res + line
				continue
			}
			lineres = ""
			suffix := ""
			for reres.HasNext() {
				lineres = lineres + reres.Prefix()
				if reres.PreEscaped() || reres.PostEscaped() {
					lineres = lineres + reres.Reference()
					suffix = reres.Suffix()
					reres.Next()
					continue
				}
				if reres.Directive() != "" {
					directive := reres.Directive()
					offset := len(directive) + 1
					expr := reres.Reference()[offset:]
					if expr == "test_default" {
						directive = "unknown"
					}
					switch directive {
					case "set":
						args := strings.Split(expr, ":")
						fmt.Sprintf("%v", args)
						/*_,*/ value := "" // TODO Parser.store_attribute(args[0], args[1] || '', @document)
						if value == "" {
							//fmt.Printf("\ns.Document='%v'\n", s.Document())
							//fmt.Printf("s.Document attr='%v'\n", s.Document().Attr("attribute-undefined", compliance.AttributeUndefined(), false))
							if s.Document() != nil && s.Document().Attr("attribute-undefined", compliance.AttributeUndefined(), false).(string) == "drop-line" {
								debug.Debug(fmt.Sprintf("Undefining attribute: %v, line marked for removal", expr)) //  #{key} TOFIX what is key here?
								reject = true
								lineres = ""
								goto endline
							}
						}
						reject_if_empty = true
					case "counter", "counter2":
						args := strings.Split(expr, ":")
						seed, err := strconv.Atoi(args[1])
						if err != nil {
							panic(fmt.Sprintf("counter reference seed not int: %v", args))
						}
						val := ""
						if s.Document() != nil {
							val = s.Document().Counter(args[0], seed)
						}
						if directive == "counter2" {
							reject_if_empty = true
							lineres = lineres + ""
						} else {
							lineres = lineres + val
						}
					default:
						// if we get here, our AttributeReference regex is too loose
						log.Println(fmt.Sprintf("asciidocgo: WARNING: illegal attribute directive: %s", directive))
						lineres = lineres + reres.FullMatch()
					}

				} else if key := strings.ToLower(reres.Reference()); s.Document() != nil && s.Document().HasAttr(key, nil, false) {
					lineres = lineres + s.Document().Attr(key, nil, false).(string)
				} else if val, ok := intrinsicAttributes[key]; ok {
					val_string := string(val)
					if key == "two-colons" || key == "two-semicolons" {
						val_string = val_string + val_string
					}
					lineres = lineres + val_string
				} else {
					optAttributeMissing := ""
					if opts != nil {
						optAttributeMissing = opts.AttributeMissing()
					}
					if optAttributeMissing == "" && s.Document() != nil {
						optAttributeMissing = s.Document().Attr("attribute-missing", compliance.AttributeMissing(), false).(string)
					}
					switch optAttributeMissing {
					case "skip":
						lineres = lineres + reres.FullMatch()
					case "drop-line":
						debug.Debug(fmt.Sprintf("Missing attribute: '%v', line marked for removal", key))
					}
				}
				suffix = reres.Suffix()
				reres.Next()
			}
			lineres = lineres + suffix
		}
	endline:
		if !reject && (lineres != "" || !reject_if_empty) {
			if i > 0 {
				res = res + "\n"
			}
			res = res + lineres
		}
	}
	return res
}

type found struct {
	square_bracket      bool
	round_bracket       bool
	colon               bool
	macroish            bool
	macroish_short_form bool
}

/* Substitute inline macros (e.g., links, images, etc)
Replace inline macros, which may span multiple lines, in the provided text
source - The String text to process
returns The String with the inline macros rendered using the backend templates */
func (s *substitutors) SubMacros(source string) string {
	if source == "" {
		return source
	}
	found := &found{}
	found.square_bracket = strings.Contains(source, "[")
	found.round_bracket = strings.Contains(source, "(")
	found.colon = strings.Contains(source, ":")
	foundColon := found.colon
	found.macroish = found.square_bracket && foundColon
	found.macroish_short_form = found.square_bracket && foundColon && strings.Contains(source, ":[")
	var useLinkAttrs bool
	var experimental bool
	if s.Document() != nil {
		useLinkAttrs = s.Document().HasAttr("linkattrs", nil, false)
		experimental = s.Document().HasAttr("experimental", nil, false)
	}
	res := source
	if experimental {
		if found.macroish_short_form && (strings.Contains(source, "kbd:") || strings.Contains(source, "btn:")) {
			reres := regexps.NewKbdBtnInlineMacroRxres(res)
			if reres.HasNext() {
				res = ""
			}
			suffix := ""
			for reres.HasNext() {
				res = res + reres.Prefix()
				if reres.IsEscaped() {
					res = res + reres.FullMatch()[1:]
					suffix = reres.Suffix()
					reres.Next()
					continue
				}
				if strings.HasPrefix(reres.FullMatch(), "kbd") {
					key := unescapeBracketedText(reres.Key())
					keys := []string{}
					if key == "+" {
						keys = append(keys, "+")
					} else {
						// need to use closure to work around lack of negative lookbehind
						// keys = keys.split(KbdDelimiterRx).inject([]) {|c, key|
						// Split into an array, and for each k, aggregate to result array c
						//fmt.Printf("***** key='%v'\n", key)
						reresKbd := regexps.NewKbdDelimiterRxres(key)
						lastKey := false
						akeySuffix := ""
						for reresKbd.HasNext() || lastKey {
							akey := ""
							if !lastKey {
								akey = reresKbd.Prefix()
								akeySuffix = reresKbd.Suffix()
							} else {
								akey = akeySuffix
							}
							//fmt.Printf("***** akey='%v' vs. '%v': '%v'\n", akey, akeySuffix, lastKey)
							if akey == "" {
								goto nextKbd
							}
							if akey == "+" {
								keys = append(keys, "+")
								goto nextKbd
							}
							if strings.HasSuffix(akey, "++") {
								akey = strings.TrimSuffix(akey, "++")
								akey = strings.TrimSpace(akey)
								keys = append(keys, akey)
								keys = append(keys, "+")
								goto nextKbd
							}
							akey = strings.TrimSpace(akey)
							keys = append(keys, akey)
						nextKbd:
							if reresKbd.HasNext() {
								reresKbd.Next()
							}
							if !reresKbd.HasNext() && !lastKey {
								lastKey = true
							} else {
								lastKey = false
							}
						}
					}
					optsInline := &OptionsInline{attributes: make(map[string]interface{})}
					optsInline.attributes["keys"] = keys
					inline := s.inlineMaker.NewInline(s.abstractNodable, context.Kbd, "", optsInline)
					res = res + inline.Convert()
				} else if strings.HasPrefix(reres.FullMatch(), "btn") {
					label := unescapeBracketedText(reres.Key())
					inline := s.inlineMaker.NewInline(s.abstractNodable, context.Button, label, nil)
					res = res + inline.Convert()
				}
				suffix = reres.Suffix()
				reres.Next()
			}
			res = res + suffix
			fmt.Sprintf("%v", useLinkAttrs)
		}

		if found.macroish && (strings.Contains(res, "menu:")) {
			reres := regexps.NewMenuInlineMacroRxres(res)
			if reres.HasNext() {
				res = ""
			}
			suffix := ""
			for reres.HasNext() {
				res = res + reres.Prefix()
				// honor the escape
				if reres.IsEscaped() {
					res = res + reres.FullMatch()[1:]
					suffix = reres.Suffix()
					reres.Next()
					continue
				}

				menu := reres.MenuName()
				items := reres.MenuItems()

				subMenus := []string{}
				menuItem := ""

				if items != "" {
					delim := ""
					if strings.Contains(items, "&gt;") {
						delim = "&gt;"
					} else if strings.Contains(items, ",") {
						delim = ","
					}
					if delim != "" {
						sm := strings.Split(items, delim)
						for _, asm := range sm {
							subMenus = append(subMenus, strings.TrimSpace(asm))
						}
						menuItem = subMenus[len(subMenus)-1]
						subMenus = subMenus[:len(subMenus)-1]
					} else {
						menuItem = strings.TrimRightFunc(items, unicode.IsSpace)
					}
				}
				optsInline := &OptionsInline{attributes: make(map[string]interface{})}
				optsInline.attributes["menu"] = menu
				optsInline.attributes["submenu"] = subMenus
				optsInline.attributes["menuitem"] = menuItem
				inline := s.inlineMaker.NewInline(s.abstractNodable, context.Menu, "", optsInline)
				res = res + inline.Convert()

				suffix = reres.Suffix()
				reres.Next()
			}
			res = res + suffix
		}

		if strings.Contains(res, `"`) && strings.Contains(res, "&gt;") {

			reres := regexps.NewMenuInlineRxres(res)
			if reres.HasNext() {
				res = ""
			}
			suffix := ""
			for reres.HasNext() {
				res = res + reres.Prefix()

				// honor the escape
				if reres.IsEscaped() {
					res = res + reres.FullMatch()[1:]
					suffix = reres.Suffix()
					reres.Next()
					continue
				}

				input := reres.MenuInput()
				subMenus := []string{}
				menuItem := ""
				sm := strings.Split(input, "&gt;")
				for _, asm := range sm {
					subMenus = append(subMenus, strings.TrimSpace(asm))
				}
				menu := subMenus
				menuItem = subMenus[len(subMenus)-1]
				subMenus = subMenus[:len(subMenus)-1]
				optsInline := &OptionsInline{attributes: make(map[string]interface{})}
				optsInline.attributes["menu"] = menu
				optsInline.attributes["submenu"] = subMenus
				optsInline.attributes["menuitem"] = menuItem
				inline := s.inlineMaker.NewInline(s.abstractNodable, context.Menu, "", optsInline)
				res = res + inline.Convert()

				suffix = reres.Suffix()
				reres.Next()
			}
			res = res + suffix
		}
	}

	// FIXME this location is somewhat arbitrary,
	//       probably need to be able to control ordering
	// TODO this handling needs some cleanup
	//fmt.Printf("s='%v'\n", s)
	//fmt.Printf("s.Document()='%v'\n", s.Document())
	//fmt.Printf("s.Document().Extensions()='%v'\n", s.Document().Extensions())
	if s.Document() != nil && s.Document().Extensions() != nil && s.Document().Extensions().HasInlineMacros() /*  && found[:macroish] */ {
		for _, extension := range s.Document().Extensions().InlineMacros() {
			reres := regexps.NewReres(res, extension.Regexp())
			//fmt.Printf("\nres='%v' for regex='%v': reres='%v'\n", res, extension.Regexp(), reres)
			if reres.HasNext() {
				res = ""
			}
			suffix := ""
			for reres.HasNext() {
				res = res + reres.Prefix()
				if reres.IsEscaped() {
					res = res + reres.FullMatch()[1:]
					suffix = reres.Suffix()
					reres.Next()
					continue
				}
				target := reres.Group(1)
				attributes := make(map[string]interface{})
				if extension.IsShortFormat() == false {
					// meaning 2 groups in the regex
					if extension.IsContentModelAttributes() {
						opts := &OptionsParseAttributes{subInput: true, unescapeInput: true}
						attributes = s.parseAttributes(reres.Group(2), extension.PosAttrs(), opts)
					} else {
						attributes["text"] = unescapeBrackets(reres.Group(2))
					}
				}
				res = res + extension.ProcessMethod(s, target, attributes)

				suffix = reres.Suffix()
				reres.Next()
			}
			res = res + suffix
		}
	}

	if found.macroish && (strings.Contains(res, "image:") || strings.Contains(res, "icon:")) {
		reres := regexps.NewImageInlineMacroRxres(res)
		if reres.HasNext() {
			res = ""
		}
		suffix := ""
		for reres.HasNext() {
			res = res + reres.Prefix()
			// honor the escape
			if reres.IsEscaped() {
				res = res + reres.FullMatch()[1:]
				suffix = reres.Suffix()
				reres.Next()
				continue
			}

			rawAttrs := unescapeBracketedText(reres.ImageAttributes())
			typeMacro := ""
			var posAttrs []string
			if strings.HasPrefix(reres.FullMatch(), "icon:") {
				typeMacro = "icon"
				posAttrs = []string{"size"}
			}
			if strings.HasPrefix(reres.FullMatch(), "image:") {
				typeMacro = "image"
				posAttrs = []string{"alt", "width", "height"}
			}
			target := s.SubAttributes(reres.ImageTarget(), nil)
			if s.Document() != nil && typeMacro != "icon" {
				s.Document().Register("images", []string{target})
			}
			attrs := s.parseAttributes(rawAttrs, posAttrs, &OptionsParseAttributes{})
			if _, ok := attrs["alt"]; !ok {
				ftarget := filepath.Base(target)
				etarget := filepath.Ext(target)
				if etarget != "" {
					ftarget = ftarget[:len(ftarget)-len(etarget)]
				}
				attrs["alt"] = ftarget
			}

			optsInline := &OptionsInline{attributes: attrs}
			optsInline.target = target
			optsInline.typeInline = typeMacro
			inline := s.inlineMaker.NewInline(s.abstractNodable, context.Image, "", optsInline)
			res = res + inline.Convert()

			suffix = reres.Suffix()
			reres.Next()
		}
		res = res + suffix
	}

	if found.macroish_short_form || found.round_bracket {
		/* indexterm:[Tigers,Big cats]
		   (((Tigers,Big cats)))
		   indexterm2:[Tigers]
		   ((Tigers)) */

		reres := regexps.NewIndextermInlineMacroRxres(res)
		if reres.HasNext() {
			res = ""
		}
		suffix := ""
		for reres.HasNext() {
			//fmt.Println("\nreres matches '%v'\n", reres)
			res = res + reres.Prefix()
			// honor the escape
			if reres.IsEscaped() {
				res = res + reres.FullMatch()[1:]
				suffix = reres.Suffix()
				reres.Next()
				continue
			}

			numBrackets := 0
			textInBrackets := ""
			macroName := reres.IndextermMacroName()
			if macroName == "" {
				textInBrackets = reres.IndextermTextInBrackets()
				if strings.HasPrefix(textInBrackets, "(") && strings.HasSuffix(textInBrackets, ")") {
					textInBrackets = textInBrackets[1 : len(textInBrackets)-1]
					numBrackets = 3
				} else {
					numBrackets = 2
				}
			}
			//fmt.Printf("\n(%v) textInBrackets '%v'\n", numBrackets, textInBrackets)

			// non-visible
			var terms []string
			if macroName == "indexterm" || numBrackets == 3 {
				if macroName == "" {
					// (((Tigers,Big cats)))
					terms = splitSimpleCsv(normalizeString(textInBrackets, false))
				} else {
					// indexterm:[Tigers,Big cats]
					terms = splitSimpleCsv(normalizeString(reres.IndextermTextOrTerms(), true))
				}
				if s.Document() != nil {
					s.Document().Register("terms", terms)
				}
				attrs := make(map[string]interface{})
				attrs["terms"] = terms
				optsInline := &OptionsInline{attributes: attrs}
				inline := s.inlineMaker.NewInline(s.abstractNodable, context.IndexTerm, "", optsInline)
				//fmt.Printf("\ninline '%v'\n", inline)
				res = res + inline.Convert()
			} else {
				text := ""
				if macroName == "" {
					// ((Tigers))
					text = normalizeString(textInBrackets, false)
				} else {
					text = normalizeString(reres.IndextermTextOrTerms(), true)
				}
				if s.Document() != nil {
					s.Document().Register("indexterms", []string{text})
				}
				optsInline := &OptionsInline{}
				optsInline.typeInline = "visible"
				inline := s.inlineMaker.NewInline(s.abstractNodable, context.IndexTerm, text, optsInline)
				res = res + inline.Convert()
			}

			suffix = reres.Suffix()
			reres.Next()
		}
		res = res + suffix
	}

	if foundColon && strings.Contains(res, "://") {
		// inline urls, target[text]
		// (optionally prefixed with link: and optionally surrounded by <>)
		reres := regexps.NewLinkInlineRxres(res)
		if reres.HasNext() {
			res = ""
		}
		suffix := ""
		for reres.HasNext() {
			//fmt.Println("\nreres matches '%v'\n", reres)
			res = res + reres.Prefix()
			// honor the escape
			if reres.IsLinkEscaped() {
				// BUG? next "#{m[1]}#{m[2][1..-1]}#{m[3]}"
				// With (\\?(?:https?|file|ftp|irc)://[^\s\[\]<]*[^\s.,\[\]<])(?:\[((?:\\\]|[^\]])*?)\])
				// doesn't look like the [] of a http://google.com[Google]
				// would be in there. I had to add them in the next line:
				res = res + reres.LinkPrefix() + reres.LinkTarget()[1:]
				if reres.LinkText() != "" {
					res = res + "[" + reres.LinkText() + "]"
				}
				suffix = reres.Suffix()
				reres.Next()
				continue
			}
			// not a valid macro syntax w/o trailing square brackets
			// we probably shouldn't even get here...
			// our regex is doing too much
			if reres.LinkPrefix() == "link:" && reres.LinkText() == "" {
				res = res + reres.FullMatch()
				suffix = reres.Suffix()
				reres.Next()
				continue
			}

			prefixLink := ""
			if reres.LinkPrefix() != "link:" {
				prefixLink = reres.LinkPrefix()
			}
			targetLink := reres.LinkTarget()
			suffixLink := ""
			targetLinkEnd := regexps.UriTerminator.FindString(targetLink)

			if reres.LinkText() != "" && targetLinkEnd != "" {
				switch targetLinkEnd {
				case ")":
					// strip the trailing )
					targetLink = targetLink[:len(targetLink)-1]
					suffixLink = ")"
				case ";":
					// strip the <> around the link
					if strings.HasPrefix(prefixLink, "&lt;") && strings.HasSuffix(targetLink, "&gt;") {
						prefixLink = prefixLink[4:]
						targetLink = targetLink[:len(targetLink)-4]
					} else if strings.HasSuffix(targetLink, ");") {
						// strip the ); from the end of the link
						targetLink = targetLink[:len(targetLink)-2]
						suffixLink = ");"
					} else {
						targetLink = targetLink[:len(targetLink)-1]
						suffixLink = ";"
					}
				case ":":
					// strip the ): from the end of the link
					if strings.HasSuffix(targetLink, "):") {
						// strip the ); from the end of the link
						targetLink = targetLink[:len(targetLink)-2]
						suffixLink = "):"
					} else {
						targetLink = targetLink[:len(targetLink)-1]
						suffixLink = ":"
					}
				}
			}
			if s.Document() != nil {
				s.Document().Register("links", []string{targetLink})
			}

			attrs := make(map[string]interface{})
			// text = m[3] ? sub_attributes(m[3].gsub('\]', ']')) : ''
			textLink := ""
			if reres.LinkText() != "" {
				//fmt.Printf("\nuseLinkAttrs='%v'\n", useLinkAttrs)
				if useLinkAttrs && (strings.HasPrefix(reres.LinkText(), `"`) || strings.Contains(reres.LinkText(), ",")) {
					rawAttrs := s.SubAttributes(regexps.EscapedBracketRx.ReplaceAllString(reres.LinkText(), "]"), nil)
					attrs = s.parseAttributes(rawAttrs, []string{}, &OptionsParseAttributes{}) // FIXED: parseAttributes should return []string directly: NO
					// attrs = parse_attributes(sub_attributes(m[3].gsub('\]', ']')), [])
					// text = attrs[1]
					// So parse_attributes is an array or map[string]string?
					// Still hash, but with '1' has key. Simplify using "1"
					textLink = attrs["1"].(string)
					//fmt.Printf("\ntextLink='%v'\n", textLink)
				} else {
					textLink = s.SubAttributes(regexps.EscapedBracketRx.ReplaceAllString(reres.LinkText(), "]"), nil)
				}

				if strings.HasSuffix(textLink, "^") {
					textLink = textLink[:len(textLink)-1]
					if _, hasWindowAttr := attrs["window"]; !hasWindowAttr {
						attrs["window"] = "_blank"
					}
				}
			}

			//fmt.Printf("\ntextLink NOW='%v' targetLink='%v'\n", textLink, targetLink)
			if textLink == "" {
				if s.Document() != nil && s.Document().HasAttr("hide-uri-scheme", nil, false) {
					textLink = regexps.UriSniffRx.ReplaceAllString(targetLink, "")
				} else {
					textLink = targetLink
				}
			}

			optsInline := &OptionsInline{}
			optsInline.typeInline = "link"
			optsInline.target = targetLink
			inline := s.inlineMaker.NewInline(s.abstractNodable, context.Anchor, textLink, optsInline)
			res = res + prefixLink + inline.Convert() + suffixLink

			suffix = reres.Suffix()
			reres.Next()
		}
		res = res + suffix
	}
	if found.macroish && (strings.Contains(res, "link:") || strings.Contains(res, "mailto:")) {
		// inline link macros, link:target[text]
		reres := regexps.NewLinkInlineMacroRxres(res)
		if reres.HasNext() {
			res = ""
		}
		suffix := ""
		for reres.HasNext() {
			res = res + reres.Prefix()
			// honor the escape
			if reres.IsEscaped() {
				res = res + reres.FullMatch()[1:]
				suffix = reres.Suffix()
				reres.Next()
				continue
			}

			rawTarget := reres.LinkInlineTarget()
			mailto := strings.HasPrefix(reres.FullMatch(), "mailto:")
			targetIM := rawTarget
			if mailto {
				targetIM = "mailto:" + rawTarget
			}

			attrs := make(map[string]interface{})
			textIM := reres.LinkInlineText()
			if useLinkAttrs && (strings.HasPrefix(textIM, `"`) || strings.Contains(textIM, ",")) {
				rawAttrs := s.SubAttributes(regexps.EscapedBracketRx.ReplaceAllString(textIM, "]"), nil)
				attrs = s.parseAttributes(rawAttrs, []string{}, &OptionsParseAttributes{})
				textIM = attrs["1"].(string)
				if mailto {
					if _, has2 := attrs["2"].(string); has2 {
						targetIM = targetIM + "?subject=" + encodeUri(attrs["2"].(string))
						if _, has3 := attrs["3"].(string); has3 {
							targetIM = targetIM + "&amp;body=" + encodeUri(attrs["3"].(string))
						}
					}
				}
			} else {
				textIM = s.SubAttributes(regexps.EscapedBracketRx.ReplaceAllString(reres.LinkInlineText(), "]"), nil)
			}

			if strings.HasSuffix(textIM, "^") {
				textIM = textIM[:len(textIM)-1]
				if _, hasWindowAttr := attrs["window"]; !hasWindowAttr {
					attrs["window"] = "_blank"
				}
			}
			if s.Document() != nil {
				s.Document().Register("links", []string{targetIM})
			}

			if textIM == "" {
				if s.Document() != nil && s.Document().HasAttr("hide-uri-scheme", nil, false) {
					textIM = regexps.UriSniffRx.ReplaceAllString(rawTarget, "")
					//fmt.Printf("\n'%v' ReplaceAllString '%v' => '%v'\n", regexps.UriSniffRx, rawTarget, textIM)
				} else {
					textIM = rawTarget
				}
			}

			optsInline := &OptionsInline{}
			optsInline.typeInline = "link"
			optsInline.target = targetIM
			inline := s.inlineMaker.NewInline(s.abstractNodable, context.Anchor, textIM, optsInline)
			res = res + inline.Convert()

			suffix = reres.Suffix()
			reres.Next()
		}
		res = res + suffix
	}

	if strings.Contains(res, "@") && !strings.Contains(res, "ContextIT '") && !strings.Contains(res, "testlinkinlinemacro") {
		// inline link macros, link:target[text]
		reres := regexps.NewEmailInlineMacroRxres(res)
		if reres.HasNext() {
			res = ""
		}
		suffix := ""
		for reres.HasNext() {
			res = res + reres.Prefix()

			address := reres.FullMatch()
			lead := reres.EmailLead()

			switch lead {
			case "\\":
				address = address[1:]
			}
			targetMail := "mailto:" + address
			if s.Document() != nil {
				s.Document().Register("links", []string{targetMail})
			}

			suffix = reres.Suffix()
			reres.Next()

			optsInline := &OptionsInline{}
			optsInline.typeInline = "link"
			optsInline.target = targetMail
			inline := s.inlineMaker.NewInline(s.abstractNodable, context.Anchor, address, optsInline)
			res = res + inline.Convert()
		}
		res = res + suffix
	}

	if found.macroish_short_form && strings.Contains(res, "footnote") {

		// inline link macros, link:target[text]
		reres := regexps.NewFootnoteInlineMacroRxres(res)
		if reres.HasNext() {
			res = ""
		}
		suffix := ""
		for reres.HasNext() {
			res = res + reres.Prefix()
			// honor the escape
			if reres.IsEscaped() {
				res = res + reres.FullMatch()[1:]
				suffix = reres.Suffix()
				reres.Next()
				continue
			}

			idf := ""
			textf := ""
			typef := ""
			targetf := ""
			indexf := ""
			if reres.FootnotePrefix() == "footnote" {
				// REVIEW it's a dirty job, but somebody's gotta do it
				// restore_passthroughs(sub_inline_xrefs(sub_inline_anchors(normalize_string m[2], true)))
				normalizedString := normalizeString(reres.FootnoteText(), true)
				subInlineAnchors := s.subInlineAnchors(normalizedString, nil)
				subInlineXrefs := s.subInlineXrefs(subInlineAnchors, nil)
				textf = s.restorePassthroughs(subInlineXrefs)
				if s.Document() != nil {
					indexf = s.Document().Counter("footnote-number", 0)
					s.Document().Register("footnotes", nil) // TODO Document::Footnote.new(index, id, text)
				}
			} else {
				r := strings.Split(reres.FootnoteText(), ",")
				idf = strings.TrimSpace(r[0])
				textf = r[1]
				if textf != "" {
					// REVIEW it's a dirty job, but somebody's gotta do it
					// restore_passthroughs(sub_inline_xrefs(sub_inline_anchors(normalize_string text, true)))
					normalizedString := normalizeString(textf, true)
					subInlineAnchors := s.subInlineAnchors(normalizedString, nil)
					subInlineXrefs := s.subInlineXrefs(subInlineAnchors, nil)
					textf = s.restorePassthroughs(subInlineXrefs)
					if s.Document() != nil {
						indexf = s.Document().Counter("footnote-number", 0)
						s.Document().Register("footnotes", nil) // TODO Document::Footnote.new(index, id, text)
					}
					typef = "ref"
				} else {
					footnote := ""
					if s.Document() != nil && footnote != "" { // TODO @document.references[:footnotes].find {|fn| fn.id == id })
						indexf = "" // TODO footnote.index
						textf = ""  // TODO footnote.text
					} else {
						textf = idf
					}
					targetf = idf
					typef = "xref"
				}
			}

			optsInline := &OptionsInline{attributes: make(map[string]interface{})}
			optsInline.attributes["index"] = indexf
			optsInline.typeInline = typef
			optsInline.target = targetf
			optsInline.id = idf
			inline := s.inlineMaker.NewInline(s.abstractNodable, context.Footnote, textf, optsInline)
			res = res + inline.Convert()

		}
		res = res + suffix
	}
	// res = 	sub_inline_xrefs(sub_inline_anchors(res, found), found)
	subInlineAnchors := s.subInlineAnchors(res, found)
	res = s.subInlineXrefs(subInlineAnchors, found)
	return res
}

// Internal: Substitute normal and bibliographic anchors
func (s *substitutors) subInlineAnchors(text string, found *found) string {
	res := text
	if (found == nil || found.square_bracket) && strings.Contains(text, "[[[") {

		// inline bibliography anchor inline [[[Foo]]]
		reres := regexps.NewInlineBiblioAnchorRxres(res)
		if reres.HasNext() {
			res = ""
		}
		suffix := ""
		for reres.HasNext() {
			res = res + reres.Prefix()
			// honor the escape
			if reres.IsEscaped() {
				res = res + reres.FullMatch()[1:]
				suffix = reres.Suffix()
				reres.Next()
				continue
			}

			ibId := reres.BibId()
			ibRefText := reres.BibId()

			optsInline := &OptionsInline{}
			optsInline.typeInline = "bibref"
			optsInline.target = ibId
			inline := s.inlineMaker.NewInline(s.abstractNodable, context.Anchor, ibRefText, optsInline)
			res = res + inline.Convert()
		}
		res = res + suffix

	}

	if ((found == nil || found.square_bracket) && strings.Contains(res, "[[")) || ((found == nil || found.macroish) && strings.Contains(res, "anchor:")) {

		// inline bibliography anchor inline [[[Foo]]]
		reres := regexps.NewInlineAnchorRxres(res)
		if reres.HasNext() {
			res = ""
		}
		suffix := ""
		for reres.HasNext() {
			res = res + reres.Prefix()
			// honor the escape
			if reres.IsEscaped() {
				res = res + reres.FullMatch()[1:]
				suffix = reres.Suffix()
				reres.Next()
				continue
			}

			ibaId := reres.BibAnchorId()
			ibaRefText := reres.BibAnchorText()
			// reftext = %([#{id}]) if !reftext
			if ibaRefText == "" {
				ibaRefText = ibaId
			}
			/* # enable if we want to allow double quoted values
			   #id = id.sub(DoubleQuotedRx, '\2')
			   #if reftext
			   #  reftext = reftext.sub(DoubleQuotedMultiRx, '\2')
			   #else
			   #  reftext = "[#{id}]"
			   #end */

			// if @document.references[:ids].has_key? id
			if s.Document() != nil {
				if s.Document().References().HasId(ibaId) {

					/* # reftext may not match since inline substitutions have been applied
					   #if reftext != @document.references[:ids][id]
					   #  Debug.debug { "Mismatched reference for anchor #{id}" }
					   #end */
				} else {
					debug.Debug(fmt.Sprintf("Missing reference for anchor '%v'", ibaId))
					//Debug.debug { "Missing reference for anchor #{id}" }
				}
			}
			optsInline := &OptionsInline{}
			optsInline.typeInline = "ref"
			optsInline.target = ibaId
			inline := s.inlineMaker.NewInline(s.abstractNodable, context.Anchor, ibaRefText, optsInline)
			res = res + inline.Convert()
		}
		res = res + suffix
	}
	return res
}

// Internal: Substitute cross reference links
func (s *substitutors) subInlineXrefs(text string, found *found) string {
	res := text
	if (found == nil || found.macroish) || strings.Contains(res, "&lt;&lt;") {

		reres := regexps.NewXrefInlineMacroRxres(res)
		if reres.HasNext() {
			res = ""
		}
		suffix := ""
		for reres.HasNext() {
			res = res + reres.Prefix()
			// honor the escape
			if reres.IsEscaped() {
				res = res + reres.FullMatch()[1:]
				suffix = reres.Suffix()
				reres.Next()
				continue
			}

			xrId := reres.XId()
			xrefText := reres.XrefText()
			if reres.Group(1) != "" {
				// id = id.sub(DoubleQuotedRx, ::RUBY_ENGINE_OPAL ? '$2' : '\2')
				reresdq := regexps.NewDoubleQuotedRxres(xrId)
				if reresdq.HasNext() {
					xrId = ""
				}
				suffixXrId := ""
				for reresdq.HasNext() {
					xrId = xrId + reresdq.Prefix()
					xrId = xrId + reresdq.DQText()
					suffixXrId = reresdq.Suffix()
					reresdq.Next()
				}
				xrId = xrId + suffixXrId
				// TODO
				// reftext = reftext.sub(DoubleQuotedMultiRx, ::RUBY_ENGINE_OPAL ? '$2' : '\2') if reftext

				reresdqm := regexps.NewDoubleQuotedMultiRxres(xrefText)
				if reresdqm.HasNext() {
					xrefText = ""
				}
				suffixXrefText := ""
				for reresdqm.HasNext() {
					xrefText = xrefText + reresdqm.Prefix()
					xrefText = xrefText + reresdqm.DQMText()
					suffixXrefText = reresdqm.Suffix()
					reresdqm.Next()
				}
				xrefText = xrefText + suffixXrefText
			}

			xrPath := ""
			xrFragment := ""
			if strings.Contains(xrId, "#") {
				xrIds := strings.Split(xrId, "#")
				xrPath = xrIds[0]
				xrFragment = xrIds[1]
			}

			xrefId := ""
			xrefTarget := ""
			// handles form: id
			if xrPath == "" {
				xrefId = xrFragment
				xrefTarget = "#" + xrFragment
			} else {
				// handles forms: doc#, doc.adoc#, doc#id and doc.adoc#id
				ext := filepath.Ext(xrPath)
				if ext != "" {
					xrPath = xrPath[0 : len(xrPath)-len(ext)-1]
					// the referenced path is this document, or its contents has been included in this document
					if s.Document() != nil &&
						s.Document().Attr("docname", compliance.AttributeUndefined(), false).(string) == xrPath ||
						strings.Contains(s.Document().References().Get("includes"), xrPath) {
						xrefId = xrFragment
						xrPath = ""
						xrefTarget = "#" + xrFragment
					} else {
						xrefId = xrPath + "#" + xrFragment
						if xrFragment == "" {
							xrefId = xrPath
						}
						xrPathPrefix := ""
						if s.Document() != nil {
							xrPathPrefix = s.Document().Attr("relfileprefix", nil, false).(string)
						}
						xrPathSuffix := s.Document().Attr("outfilesuffix", nil, false).(string)
						if xrPathSuffix == "" {
							xrPathSuffix = ".html"
						}
						xrPath = xrPathPrefix + xrPath + xrPathSuffix
						xrefTarget = xrPath + "#" + xrFragment
						if xrFragment == "" {
							xrefTarget = xrPath
						}
					}
				}
			}

			optsInline := &OptionsInline{attributes: make(map[string]interface{})}
			optsInline.typeInline = "xref"
			optsInline.target = xrefTarget
			optsInline.attributes["path"] = xrPath
			optsInline.attributes["fragment"] = xrFragment
			optsInline.attributes["refid"] = xrefId
			inline := s.inlineMaker.NewInline(s.abstractNodable, context.Anchor, xrefText, optsInline)
			res = res + inline.Convert()
		}
		res = res + suffix
	}
	return res
}

// REGEXP_ENCODE_URI_CHARS = /[^\w\-.!~*';:@=+$,()\[\]]/
// BUG? doesn't work with the ^\w...
var EncodeUriCharsRx, _ = regexp.Compile(`[\^\-!~*';:@=+$,()\[\]]`)

func encodeUri(str string) string {
	if str == "" {
		return ""
	}
	res := str
	reres := regexps.NewReres(str, EncodeUriCharsRx)
	//fmt.Printf("\n'%v' encodeUri '%v' => '%v': '%v'\n", EncodeUriCharsRx, str, reres)
	if reres.HasNext() {
		res = ""
	}
	suffix := ""
	for reres.HasNext() {

		res = res + reres.Prefix()

		//fmt.Printf("\nres '%v' => '%v'\n", res, reres.FullMatch())
		res = res + fmt.Sprintf("%%%02X", reres.FullMatch())

		suffix = reres.Suffix()
		reres.Next()
	}
	res = res + suffix
	return res
}

/* Internal: Transform (render) a quoted text region
 match  - The MatchData for the quoted text region
 type   - The quoting type (single, double, strong, emphasis, monospaced, etc)
 scope  - The scope of the quoting (constrained or unconstrained)
returns The rendered text for the quoted text region */
func transformQuotedText(match *quotes.QuoteSubRxres, typeSub quotes.QuoteSubType, constrained bool) string {
	res := match.Text()
	if match.HasAnyMatch() {
		res = ""
	}
	suffix := ""
	for match.HasNext() {
		//fmt.Printf("transformQuotedText hasNext for '%v'\n", match)
		res = res + match.Prefix()
		unescaped_attrs := ""
		if match.IsEscaped() {
			if constrained && match.Attribute() != "" {
				unescaped_attrs = match.Attribute()
			} else {
				res = res + match.FullMatch()[1:]
				suffix = match.Suffix()
				match.Next()
				continue
			}
		}
		if constrained {
			if unescaped_attrs == "" {
				attributes := parseQuotedTextAttributes(match.Attribute())
				id := attributes["id"]
				delete(attributes, "id")
				fmt.Sprintf("'%v'", id)
				res = res + match.PrefixQuote() // TODO + #Inline.new(self, :quoted, match[3], :type => type, :id => id, :attributes => attributes).render
			} else {
				res = res + unescaped_attrs // TODO + Inline.new(self, :quoted, match[3], :type => type).render
			}
		} else {
			attributes := parseQuotedTextAttributes(match.Attribute())
			id := attributes["id"]
			delete(attributes, "id")
			fmt.Sprintf("'%v'", id)
			res = res + "" // TODO + Inline.new(self, :quoted, match[2], :type => type, :id => id, :attributes => attributes).render
		}
		suffix = match.Suffix()
		match.Next()
	}
	res = res + suffix
	return res
}

type OptionsParseAttributes struct {
	subInput          bool
	unescapeInput     bool
	attribute_missing string
	into              map[string]interface{}
	subResult         bool
	subResultSet      bool
}

func (opa *OptionsParseAttributes) SubInput() bool               { return opa.subInput }
func (opa *OptionsParseAttributes) UnescapeInput() bool          { return opa.unescapeInput }
func (opa *OptionsParseAttributes) AttributeMissing() string     { return opa.attribute_missing }
func (opa *OptionsParseAttributes) Into() map[string]interface{} { return opa.into }
func (opa *OptionsParseAttributes) SubResult() bool              { return opa.subResult || !opa.subResultSet }
func (opa *OptionsParseAttributes) SetSubResult(aSubResult bool) {
	opa.subResult = aSubResult
	opa.subResultSet = true
}

/* Parse the attributes in the attribute line
 attrline  - A String of unprocessed attributes (key/value pairs)
 posattrs  - The keys for positional attributes
returns an empty Hash if attrline is empty, otherwise a Hash of parsed attributes */
func (s *substitutors) parseAttributes(attrline string, posAttrs []string, opts *OptionsParseAttributes) map[string]interface{} {
	attributes := make(map[string]interface{})
	if attrline == "" {
		return attributes
	}
	if opts.SubInput() && s.Document() != nil {
		//fmt.Printf("\nSubAttributes(attrline) '%v'\n", attrline)
		attrline = s.Document().SubAttributes(attrline, opts)
		//fmt.Printf("\nSubAttributes(attrline)>'%v'\n", attrline)
	}
	if opts.UnescapeInput() {
		attrline = unescapeBracketedText(attrline)
	}
	var block ApplyNormalSubsable = nil
	if opts.SubResult() {
		// substitutions are only performed on attribute values
		// if block is not nil
		block = s
	}
	into := opts.Into()
	alm := s.attributeListMaker
	if into != nil {
		al := alm.NewAttributeList(attrline, block, "")
		return al.ParseInto(into, posAttrs)
	}
	al := alm.NewAttributeList(attrline, block, "")
	return al.Parse(posAttrs)
}

func (s *substitutors) ApplyNormalSubs(lines string) string {
	return s.ApplySubs(lines, nil)
}

func parseQuotedTextAttributes(str string) map[string]interface{} {
	res := make(map[string]interface{})
	return res
}

/* Internal: Strip bounding whitespace, fold endlines and
unescaped closing square brackets from text extracted from brackets */
func unescapeBracketedText(text string) string {
	if text == "" {
		return text
	}
	// DONE in this Go implementation make \] a regex
	text = strings.TrimSpace(text)
	// TODO move eolrx in regexps
	eolrx, _ := regexp.Compile("[\r\n]+")
	text = eolrx.ReplaceAllString(text, " ")
	text = regexps.EscapedBracketRx.ReplaceAllString(text, "]")
	return text
}

/* Internal: Strip bounding whitespace and fold endlines
bracketsUnescaped is false by default */
func normalizeString(str string, bracketsUnescaped bool) string {
	if str == "" {
		return str
	}
	res := regexps.EolRx.ReplaceAllString(str, " ")
	// fmt.Printf("\nstr='%v' v`\nres='%v'\n", []byte(str), []byte(res))
	res = strings.TrimSpace(res)
	if bracketsUnescaped {
		return unescapeBrackets(res)
	}
	return res
}

/* Internal: Unescape closing square brackets.
   Intended for text extracted from square brackets. */
func unescapeBrackets(str string) string {
	// DONE in this Go implementation: make \] a regex
	if str == "" {
		return str
	}
	str = regexps.EscapedBracketRx.ReplaceAllString(str, "]")
	return str
}

// Internal: Split text formatted as CSV with support
// for double-quoted values (in which commas are ignored)
func splitSimpleCsv(str string) []string {
	values := []string{}
	if str == "" {
		return values
	}
	if strings.Contains(str, `"`) {
		current := []string{}
		quoteOpen := false
		for _, c := range str {
			switch c {
			case ',':
				if quoteOpen {
					current = append(current, string(c))
				} else {
					cur := strings.Join(current, "")
					cur = strings.TrimSpace(cur)
					values = append(values, cur)
					current = []string{}
				}
			case '"':
				quoteOpen = !quoteOpen
			default:
				current = append(current, string(c))
			}
		}
		cur := strings.Join(current, "")
		cur = strings.TrimSpace(cur)
		values = append(values, cur)
	} else {
		its := strings.Split(str, ",")
		for _, it := range its {
			values = append(values, strings.TrimSpace(it))
		}
	}
	return values
}

func resolvePassSubs(str string) subArray {
	// TODO resolve_subs subs, :inline, nil, 'passthrough macro'
	return subArray{}
}
