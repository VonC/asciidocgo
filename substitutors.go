package asciidocgo

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

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

type SubstDocumentable interface {
	Attr(name string, defaultValue interface{}, inherit bool) interface{}
	Basebackend(base interface{}) bool
	SubAttributes(data string, opts *OptionsParseAttributes) string
	Counter(name string, seed int) string
	HasAttr(name string, expect interface{}, inherit bool) bool
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

/*func (oi *OptionsInline) Id() string                         { return oi.id }
func (oi *OptionsInline) TypeInline() string                 { return oi.typeInline }
func (oi *OptionsInline) Target() string                     { return oi.target }*/
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

/* Methods to perform substitutions on lines of AsciiDoc text.
This module is intented to be mixed-in to Section and Block to provide
operations for performing the necessary substitutions. */
type substitutors struct {
	// A String Array of passthough (unprocessed) text captured from this block
	passthroughs    []passthrough
	document        SubstDocumentable
	inlineMaker     InlineMaker
	abstractNodable AbstractNodable
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
				p := passthrough{textOri, subsOri, make(map[string]interface{}), ""}
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
				attributes = s.parseAttributes(reres.Attributes(), &OptionsParseAttributes{})
			}

			p := passthrough{reres.LiteralText(), subArray{subValue.specialcharacters}, attributes, "monospaced"}
			s.passthroughs = append(s.passthroughs, p) //TODO attributes, type
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
			p := passthrough{mathText, mathSubs, attributes, mathType}
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
					key := unescapeBrackets(reres.Key())
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
					label := unescapeBrackets(reres.Key())
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
			fmt.Sprintf("test")
		}
	}
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
	attribute_missing string
}

func (opa *OptionsParseAttributes) SubInput() bool           { return opa.subInput }
func (opa *OptionsParseAttributes) AttributeMissing() string { return opa.attribute_missing }

/* Parse the attributes in the attribute line
 attrline  - A String of unprocessed attributes (key/value pairs)
 posattrs  - The keys for positional attributes
returns an empty Hash if attrline is empty, otherwise a Hash of parsed attributes */
func (s *substitutors) parseAttributes(attrline string, opts *OptionsParseAttributes) map[string]interface{} {
	attributes := make(map[string]interface{})
	if attrline == "" {
		return attributes
	}
	if opts.SubInput() && s.Document() != nil {
		attrline = s.Document().SubAttributes(attrline, opts)
	}
	// TODO implement parseAttributes posattrs and opt map[string]interface{}
	return attributes
}

func parseQuotedTextAttributes(str string) map[string]interface{} {
	res := make(map[string]interface{})
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
	// TODO resolve_subs subs, :inline, nil, 'passthrough macro'
	return subArray{}
}
