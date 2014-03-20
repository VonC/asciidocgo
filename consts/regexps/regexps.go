package regexps

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/VonC/asciidocgo/utils"
)

var ADMONITION_STYLES utils.Arr = []string{"NOTE", "TIP", "IMPORTANT", "WARNING", "CAUTION"}

const (
	CC_ALPHA = `a-zA-Z`
	CC_ALNUM = `a-zA-Z0-9`
	CC_BLANK = `[ \t]`
	// non-blank character
	CC_GRAPH = `[\x21-\x7E]`
	CC_EOL   = `(?=\n|$)`
)

var ORDERED_LIST_KEYWORDS = map[string]rune{
	"loweralpha": 'a',
	"lowerroman": 'i',
	"upperalpha": 'A',
	"upperroman": 'I',
	//'lowergreek': 'a'
	//'arabic': '1'
	//'decimal': '1'
}

var EolRx, _ = regexp.Compile(`[\r\n]+`)

/* Encapsulate a regex and a string,
for managing results from FindAllStringSubmatchIndex */
type Reres struct {
	r        *regexp.Regexp
	s        string
	matches  [][]int
	i        int
	previous int
}

func (r *Reres) String() string {
	msg := fmt.Sprintf("Regexp res for '%v': (%v-%v; len %v) %v", r.r, r.i, r.previous, len(r.s), r.matches)
	return msg
}

/* Build new result from FindAllStringSubmatchIndex on a string */
func NewReres(s string, r *regexp.Regexp) *Reres {
	matches := r.FindAllStringSubmatchIndex(s, -1)
	return &Reres{r, s, matches, 0, 0}
}

type Qualifier func(lh string, match []int, s string) bool

/* Build new result from FindAllStringSubmatchIndex on a string,
validated by last group being a lookahead after each match */
func NewReresLAGroup(s string, r *regexp.Regexp) *Reres {
	return newReresLA(s, r, nil)
}

/* Build new result from FindAllStringSubmatchIndex on a string,
validated by last group being a lookahead after each match,
if that last group match qualifies */
func NewReresLAQual(s string, r *regexp.Regexp, q Qualifier) *Reres {
	return newReresLA(s, r, q)
}

/* Build new result from FindAllStringSubmatchIndex on a string,
validated by last group being a lookahead after each match */
func newReresLA(s string, r *regexp.Regexp, q Qualifier) *Reres {
	bf := bytes.NewBufferString(s)
	by := bf.Bytes()
	m := [][]int{}
	lg := []int{}
	res := &Reres{r, s, m, 0, 0}
	shift := 0
	for match := r.FindSubmatchIndex(by); match != nil && len(match) > 0; match = r.FindSubmatchIndex(by) {
		if len(match) > 0 {
			//fmt.Printf("\nMatch '%v' '%v'\n-------\n", match, string(by))
			match, lg = match[:len(match)-2], match[len(match)-2:]
			for i, mi := range match {
				match[i] = mi + shift
				/*
					ss := ""
					if mi > -1 {
						ss = s[mi+shift:]
					}
					fmt.Printf("\ni=%v: shift=%v match[i]=%v: s[match[%v]]='%v'\n", i, shift, mi, mi+shift, ss)
				*/
			}
			//fmt.Printf("\nAppend '%v'\n===\n", match)
			if lg[0] <= lg[1] && lg[0] > -1 {
				lh := string(by[lg[0]:lg[1]])
				by = by[lg[0]:]
				shift = shift + lg[0]
				delta := lg[1] - lg[0]
				match[1] = match[1] - delta
				if q == nil || q(lh, match, s) {
					m = append(m, match)
				}
			} else {
				m = append(m, match)
				//fmt.Printf("\nBY (%v-%v)'%v' '%v'\n'%v'\n-------\n", len(by), match[1]-shift, by, string(by), match)
				by = by[match[1]-shift:]
				shift = shift + (match[1] - shift)
			}
		}
	}
	res.matches = m
	//fmt.Printf("\n*** %v======\n", res.matches)
	return res
}

/* full initial text on which the regex was applied */
func (rr *Reres) Text() string {
	return rr.s
}

/* Check if there is any match */
func (rr *Reres) HasAnyMatch() bool {
	return len(rr.matches) > 0
}

/* Check if there is one more match */
func (rr *Reres) HasNext() bool {
	return rr.i < len(rr.matches)
}

/* Refers to the next match, for Group() to works with */
func (rr *Reres) Next() {
	rr.previous = rr.matches[rr.i][1]
	rr.i = rr.i + 1
}

/* Get back to the first match */
func (rr *Reres) ResetNext() {
	rr.i = 0
	rr.previous = 0
}

/* String from the last match to current one */
func (rr *Reres) Prefix() string {
	mi := rr.matches[rr.i]
	return rr.s[rr.previous:mi[0]]
}

/* String from current match to the end ofthe all string */
func (rr *Reres) Suffix() string {
	mi := rr.matches[rr.i]
	res := ""
	if len(rr.s) > mi[1] {
		res = rr.s[mi[1]:]
	}
	return res
}

/* First character of the current match */
func (rr *Reres) FirstChar() uint8 {
	mi := rr.matches[rr.i]
	return rr.s[mi[0]]
}

/* Test if first character of the current match is an escape */
func (rr *Reres) IsEscaped() bool {
	mi := rr.matches[rr.i]
	return rr.s[mi[0]] == '\\'
}

/* Full string matched for the current group */
func (rr *Reres) FullMatch() string {
	mi := rr.matches[rr.i]
	return rr.s[mi[0]:mi[1]]
}

/* Check if the ith group if present in the current match */
func (rr *Reres) HasGroup(j int) bool {
	res := false
	mi := rr.matches[rr.i]
	if len(mi) > (j*2)+1 {
		if mi[j*2] > -1 {
			if mi[j*2] < mi[(j*2)+1] {
				res = true
			}
		}
	}
	return res
}

/* return the ith group string, if present in the current match */
func (rr *Reres) Group(i int) string {
	res := ""
	if rr.HasGroup(i) {
		mi := rr.matches[rr.i]
		res = rr.s[mi[i*2]:mi[(i*2)+1]]
	}
	return res
}

/* Matches an inline attribute reference.
Examples
  {foo}
  {counter:pcount:1}
  {set:foo:bar}
  {set:name!}

   AttributeReferenceRx = /(\\)?\{((set|counter2?):.+?|\w+(?:[\-]\w+)*)(\\)?\}/ */
var AttributeReferenceRx, _ = regexp.Compile(`(\\)?\{((set|counter2?):.+?|\w+(?:[\-]\w+)*)(\\)?\}`)

type AttributeReferenceRxres struct {
	*Reres
}

/* Return true if first group is non empty and include an '\' */
func (arr *AttributeReferenceRxres) PreEscaped() bool {
	return arr.Group(1) == "\\"
}

/* Return true if last group is non empty and include an '\' */
func (arr *AttributeReferenceRxres) PostEscaped() bool {
	return arr.Group(4) == "\\"
}

/* Return directive of the reference, as 'counter' in '{counter:pcount:1}' */
func (arr *AttributeReferenceRxres) Directive() string {
	return arr.Group(3)
}

/* Return reference, as in 'counter:pcount:1' in {counter:pcount:1}' */
func (arr *AttributeReferenceRxres) Reference() string {
	return arr.Group(2)
}

/* Results for AttributeReferenceRx */
func NewAttributeReferenceRxres(s string) *AttributeReferenceRxres {
	return &AttributeReferenceRxres{NewReres(s, AttributeReferenceRx)}
}

/* The following pattern, which appears frequently, captures the contents
between square brackets, ignoring escaped closing brackets
(closing brackets prefixed with a backslash '\' character)

	Pattern:
	(?:\[((?:\\\]|[^\]])*?)\])
	Matches:
	[enclosed text here] or [enclosed [text\] here]
*/

/* Matches an admonition label at the start of a paragraph.
   Examples
     NOTE: Just a little note.
     TIP: Don't forget! */
var AdmonitionParagraphRx, _ = regexp.Compile(fmt.Sprintf("^(%v):%v", ADMONITION_STYLES.Mult("|"), CC_BLANK))

/* Inline macros */

/* Matches an anchor (i.e., id + optional reference text) in the flow of text.
 Examples
   [[idname]]
   [[idname,Reference Text]]
   anchor:idname[]
   anchor:idname[Reference Text]
InlineAnchorRx = /\\?(?:\[\[([#{CC_ALPHA}:_][\w:.-]*)(?:,#{CC_BLANK}*(\S.*?))?\]\]|anchor:(\S+)\[(.*?[^\\])?\])/ */

var InlineAnchorRx, _ = regexp.Compile(`\\?(?:\[\[([#{CC_ALPHA}:_][\w:.-]*)(?:,#{CC_BLANK}*(\S.*?))?\]\]|anchor:(\S+)\[(.*?[^\\])?\]`)

type InlineAnchorRxres struct {
	*Reres
}

/* Results for InlineAnchorRx */
func NewInlineAnchorRxres(s string) *InlineAnchorRxres {
	return &InlineAnchorRxres{NewReres(s, InlineAnchorRx)}
}

/* Return id of the macro in '[[idname]]' or 'anchor:idname[]' */
func (iar *InlineAnchorRxres) BibAnchorId() string {
	if iar.Group(1) != "" {
		return iar.Group(1)
	}
	return iar.Group(3)
}

/* Return text of the macro in '[[idname,Reference Text]]' or 'anchor:idname[Reference Text]' */
func (iar *InlineAnchorRxres) BibAnchorText() string {
	if iar.Group(2) != "" {
		return iar.Group(2)
	}
	return iar.Group(4)
}

/* Matches a bibliography anchor anywhere inline.
 Examples
   [[[Foo]]]
InlineBiblioAnchorRx = /\\?\[\[\[([\w:][\w:.-]*?)\]\]\]/ */

var InlineBiblioAnchorRx, _ = regexp.Compile(`\\?\[\[\[([\w:][\w:.-]*?)\]\]\]`)

type InlineBiblioAnchorRxres struct {
	*Reres
}

/* Results for InlineBiblioAnchorRx */
func NewInlineBiblioAnchorRxres(s string) *InlineBiblioAnchorRxres {
	return &InlineBiblioAnchorRxres{NewReres(s, InlineBiblioAnchorRx)}
}

/* Return id of the macro in '[[[id]]]' */
func (ibar *InlineBiblioAnchorRxres) BibId() string {
	return ibar.Group(1)
}

/* Matches an inline e-mail address.
   doc.writer@example.com
EmailInlineMacroRx = /([\\>:\/])?\w[\w.%+-]*@[#{CC_ALNUM}][#{CC_ALNUM}.-]*\.[#{CC_ALPHA}]{2,4}\b/ */

var EmailInlineMacroRx, _ = regexp.Compile(`([\\>:\/])?\w[\w.%+-]*@[a-zA-Z0-9][a-zA-Z0-9.-]*\.[a-zA-Z]{2,4}`)

type EmailInlineMacroRxres struct {
	*Reres
}

/* Results for EmailInlineMacroRx */
func NewEmailInlineMacroRxres(s string) *EmailInlineMacroRxres {
	return &EmailInlineMacroRxres{NewReres(s, EmailInlineMacroRx)}
}

/* Return lead of the macro in '>xx:@yyy.com' */
func (eimr *EmailInlineMacroRxres) EmailLead() string {
	return eimr.Group(1)
}

/* Matches an inline footnote macro, which is allowed to span multiple lines.
Examples
  footnote:[text]
  footnoteref:[id,text]
  footnoteref:[id]

   FootnoteInlineMacroRx = /\\?(footnote(?:ref)?):\[(.*?[^\\])\]/m */

var FootnoteInlineMacroRx, _ = regexp.Compile(`\\?(footnote(?:ref)?):\[(.*?[^\\])\]`)

type FootnoteInlineMacroRxres struct {
	*Reres
}

/* Results for FootnoteInlineMacroRx */
func NewFootnoteInlineMacroRxres(s string) *FootnoteInlineMacroRxres {
	return &FootnoteInlineMacroRxres{NewReres(s, FootnoteInlineMacroRx)}
}

/* Return prefix 'footnote' of the macro in 'footnote:[xxx]' */
func (fimr *FootnoteInlineMacroRxres) FootnotePrefix() string {
	return fimr.Group(1)
}

/* Return text 'xxx' of the macro in 'footnote:[xxx]' */
func (fimr *FootnoteInlineMacroRxres) FootnoteText() string {
	return fimr.Group(2)
}

/* Matches an image or icon inline macro.
Examples
   image:filename.png[Alt Text]
   image:http://example.com/images/filename.png[Alt Text]
   image:filename.png[More [Alt\] Text] (alt text becomes "More [Alt] Text")
   icon:github[large]
    ImageInlineMacroRx = /\\?(?:image|icon):([^:\[][^\[]*)\[((?:\\\]|[^\]])*?)\]/ */

var ImageInlineMacroRx, _ = regexp.Compile(`\\?(?:image|icon):([^:\[][^\[]*)\[((?:\\\]|[^\]])*?)\]`)

type ImageInlineMacroRxres struct {
	*Reres
}

/* Results for ImageInlineMacroRx */
func NewImageInlineMacroRxres(s string) *ImageInlineMacroRxres {
	return &ImageInlineMacroRxres{NewReres(s, ImageInlineMacroRx)}
}

/* Return target of the macro in 'image:target[attr1 attr2]' */
func (iimr *ImageInlineMacroRxres) ImageTarget() string {
	return iimr.Group(1)
}

/* Return attributes of the macro in 'image:target[attr1 attr2]' */
func (iimr *ImageInlineMacroRxres) ImageAttributes() string {
	return iimr.Group(2)
}

/* Matches an indexterm inline macro, which may span multiple lines.
Examples
  indexterm:[Tigers,Big cats]
  (((Tigers,Big cats)))
  indexterm2:[Tigers]
  ((Tigers))

   IndextermInlineMacroRx = /\\?(?:(indexterm2?):\[(.*?[^\\])\]|\(\((.+?)\)\)(?!\)))/m */

var IndextermInlineMacroRx, _ = regexp.Compile(`\\?(?:(indexterm2?):\[(.*?[^\\])\]|\(\((.+?)\)\)([^\)]|$))`)

type IndextermInlineMacroRxres struct {
	*Reres
}

/* Results for IndextermInlineMacroRx */
func NewIndextermInlineMacroRxres(s string) *IndextermInlineMacroRxres {
	return &IndextermInlineMacroRxres{NewReresLAGroup(s, IndextermInlineMacroRx)}
}

/* Return name indexterm of the macro in 'indexterm:[Tigers,Big cats]' */
func (itimr *IndextermInlineMacroRxres) IndextermMacroName() string {
	return itimr.Group(1)
}

/* Return terms of the macro in 'indexterm:[Tigers,Big cats]' */
func (itimr *IndextermInlineMacroRxres) IndextermTextOrTerms() string {
	return itimr.Group(2)
}

/* Return text in brackets of the macro in '((Tigers))' */
func (itimr *IndextermInlineMacroRxres) IndextermTextInBrackets() string {
	return itimr.Group(3)
}

/*
Matches either the kbd or btn inline macro.
Examples
   kbd:[F3]
   kbd:[Ctrl+Shift+T]
   kbd:[Ctrl+\]]
   kbd:[Ctrl,T]
   btn:[Save]
KbdBtnInlineMacroRx = /\\?(?:kbd|btn):\[((?:\\\]|[^\]])+?)\]/ */
var KbdBtnInlineMacroRx, _ = regexp.Compile(`\\?(?:kbd|btn):\[((?:\\\]|[^\]])+?)\]`)

type KbdBtnInlineMacroRxres struct {
	*Reres
}

/* Results for KbdBtnInlineMacroRx */
func NewKbdBtnInlineMacroRxres(s string) *KbdBtnInlineMacroRxres {
	return &KbdBtnInlineMacroRxres{NewReres(s, KbdBtnInlineMacroRx)}
}

/* Return key of the macro xxx in ':[xxx]' */
func (kbimr *KbdBtnInlineMacroRxres) Key() string {
	return kbimr.Group(1)
}

/* Matches the delimiter used for kbd value.
Examples
   Ctrl + Alt+T
   Ctrl,T
KbdDelimiterRx = /(?:\+|,)(?=#{CC_BLANK}*[^\1])/ */
var KbdDelimiterRx, _ = regexp.Compile(`(?:\+|,)([ \t]*[^ \t])`)

type KbdDelimiterRxres struct {
	*Reres
}

func kbdla(lh string, match []int, s string) bool {
	m := strings.TrimSpace(lh)
	g1 := string(s[match[0]:match[1]])
	//fmt.Printf("\nm='%v' vs. g1='%v'\n", m, g1)
	return m != g1
}

/* Results for KbdDelimiterRx */
func NewKbdDelimiterRxres(s string) *KbdDelimiterRxres {
	return &KbdDelimiterRxres{NewReresLAQual(s, KbdDelimiterRx, kbdla)}
}

/* Matches an implicit link and some of the link inline macro.
 Examples
   http://github.com
   http://github.com[GitHub]

 FIXME revisit! the main issue is we need different rules for implicit vs explicit
LinkInlineRx = %r{(^|link:|&lt;|[\s>\(\)\[\];])(\\?(?:https?|file|ftp|irc)://[^\s\[\]<]*[^\s.,\[\]<])(?:\[((?:\\\]|[^\]])*?)\])?} */
var LinkInlineRx, _ = regexp.Compile(`(^|link:|&lt;|[\s>\(\)\[\];])(\\?(?:https?|file|ftp|irc)://[^\s\[\]<]*[^\s.,\[\]<])(?:\[((?:\\\]|[^\]])*?)\])?`)

type LinkInlineRxres struct {
	*Reres
}

/* Results for LinkInlineRx */
func NewLinkInlineRxres(s string) *LinkInlineRxres {
	return &LinkInlineRxres{NewReres(s, LinkInlineRx)}
}

/* Return true if '\' found in 'link:\http:xxx[yyy]' */
func (lir *LinkInlineRxres) IsLinkEscaped() bool {
	return strings.HasPrefix(lir.Group(2), "\\")
}

/* Return 'link:' in 'link:http:xxx[yyy]' */
func (lir *LinkInlineRxres) LinkPrefix() string {
	return lir.Group(1)
}

/* Return 'xxx' in 'link:http:xxx[yyy]' */
func (lir *LinkInlineRxres) LinkTarget() string {
	return lir.Group(2)
}

/* Return 'yyy' in 'link:http:xxx[yyy]' */
func (lir *LinkInlineRxres) LinkText() string {
	return lir.Group(3)
}

/* Match a link or e-mail inline macro.
 Examples
   link:path[label]
   mailto:doc.writer@example.com[]

LinkInlineMacroRx = /\\?(?:link|mailto):([^\s\[]+)(?:\[((?:\\\]|[^\]])*?)\])/ */
var LinkInlineMacroRx, _ = regexp.Compile(`\\?(?:link|mailto):([^\s\[]+)(?:\[((?:\\\]|[^\]])*?)\])`)

type LinkInlineMacroRxres struct {
	*Reres
}

/* Results for LinkInlineMacroRx */
func NewLinkInlineMacroRxres(s string) *LinkInlineMacroRxres {
	return &LinkInlineMacroRxres{NewReres(s, LinkInlineMacroRx)}
}

/* Return 'xxx' in 'link:xxx[yyy]' */
func (limr *LinkInlineMacroRxres) LinkInlineTarget() string {
	return limr.Group(1)
}

/* Return 'yyy' in 'link:xxx[yyy]' */
func (limr *LinkInlineMacroRxres) LinkInlineText() string {
	return limr.Group(2)
}

/* Matches a math inline macro, which may span multiple lines.
Examples
  math:[x != 0]
  asciimath:[x != 0]
  latexmath:[\sqrt{4} = 2]

MathInlineMacroRx = /\\?((?:latex|ascii)?math):([a-z,]*)\[(.*?[^\\])\]/m */
var MathInlineMacroRx, _ = regexp.Compile(`(?sm)\\?((?:latex|ascii)?math):([a-z,]*)\[(.*?[^\\])\]`)

type MathInlineMacroRxres struct {
	*Reres
}

/* Results for MathInlineMacroRx */
func NewMathInlineMacroRxres(s string) *MathInlineMacroRxres {
	return &MathInlineMacroRxres{NewReres(s, MathInlineMacroRx)}
}

/* Return type 'math' in 'math:xx[yyy]' */
func (mimr *MathInlineMacroRxres) MathType() string {
	return mimr.Group(1)
}

/* Return sub 'xx' in 'math:xx[yyy]' */
func (mimr *MathInlineMacroRxres) MathSub() string {
	return mimr.Group(2)
}

/* Return text 'yyy' in 'math:xx[yyy]' */
func (mimr *MathInlineMacroRxres) MathText() string {
	return mimr.Group(3)
}

/* Matches a passthrough literal value, which may span multiple lines.
Examples
  `text`
*/
// PassInlineLiteralRx = /(^|[^`\w])(?:\[([^\]]+?)\])?(\\?`([^`\s]|[^`\s].*?\S)`)(?![`\w])/m

/* Matches a menu inline macro.
    # Examples
   menu:File[New...]
   menu:View[Page Style > No Style]
   menu:View[Page Style, No Style]
   MenuInlineMacroRx = /\\?menu:(\w|\w.*?\S)\[#{CC_BLANK}*(.+?)?\]/
*/
var MenuInlineMacroRx, _ = regexp.Compile(`(?sm)\\?menu:(\w|\w.*?\S)\[[ \t]*(.+?)?\]`)

type MenuInlineMacroRxres struct {
	*Reres
}

/* Results for KbdBtnInlineMacroRx */
func NewMenuInlineMacroRxres(s string) *MenuInlineMacroRxres {
	return &MenuInlineMacroRxres{NewReres(s, MenuInlineMacroRx)}
}

/* Return name of the macro in 'menu:name[xxx]' */
func (mimr *MenuInlineMacroRxres) MenuName() string {
	return mimr.Group(1)
}

/* Return items of the macro xxx in 'menu:name[xxx]' */
func (mimr *MenuInlineMacroRxres) MenuItems() string {
	return mimr.Group(2)
}

/* # Matches an implicit menu inline macro.
 Examples
   "File > New..."
MenuInlineRx = /\\?"(\w[^"]*?#{CC_BLANK}*&gt;#{CC_BLANK}*[^" \t][^"]*)"/ */
var MenuInlineRx, _ = regexp.Compile(`(?sm)\\?"(\w[^"]*?[ \t]*&gt;[ \t]*[^" \t][^"]*)"`)

type MenuInlineRxres struct {
	*Reres
}

/* Results for MenuInlineRx */
func NewMenuInlineRxres(s string) *MenuInlineRxres {
	return &MenuInlineRxres{NewReres(s, MenuInlineRx)}
}

/* Return input of the macro in '"File &gt; New"' */
func (mir *MenuInlineRxres) MenuInput() string {
	return mir.Group(1)
}

var PassInlineLiteralRx, _ = regexp.Compile("(?sm)(^|[^`\\w])(?:\\[([^\\]]+?)\\])?(\\\\?`([^`\\s]|[^`\\s].*?\\S)`)([^`\\w])")

type PassInlineLiteralRxres struct {
	*Reres
}

/* Results for PassInlineLiteralRx */
func NewPassInlineLiteralRxres(s string) *PassInlineLiteralRxres {
	res := &PassInlineLiteralRxres{NewReresLAGroup(s, PassInlineLiteralRx)}
	return res
}

func (pilr *PassInlineLiteralRxres) FirstChar() string {
	return pilr.Group(1)
}

func (pilr *PassInlineLiteralRxres) Attributes() string {
	return pilr.Group(2)
}

func (pilr *PassInlineLiteralRxres) Literal() string {
	return pilr.Group(3)
}

func (pilr *PassInlineLiteralRxres) LiteralText() string {
	return pilr.Group(4)
}

/* Matches several variants of the passthrough inline macro,
which may span multiple lines.

 Examples

   +++text+++
   $$text$$
   pass:quotes[text] */
// http://stackoverflow.com/questions/6770898/unknown-escape-sequence-error-in-go
var PassInlineMacroRx, _ = regexp.Compile(`(?s)\\?(?:(\+{3})(.*?)\+{3}|(\${2})(.*?)\${2}|pass:([a-z,]*)\[(.*?[^\\])\])`)

type PassInlineMacroRxres struct {
	*Reres
}

/* Results for PassInlineMacroRx */
func NewPassInlineMacroRxres(s string) *PassInlineMacroRxres {
	return &PassInlineMacroRxres{NewReres(s, PassInlineMacroRx)}
}

/* Check if has text 'yyy' in 'pass:xx[yyy]' */
func (pr *PassInlineMacroRxres) HasPassText() bool {
	return pr.HasGroup(6)
}

/* Return text 'yyy' in 'pass:xx[yyy]' */
func (pr *PassInlineMacroRxres) PassText() string {
	return pr.Group(6)
}

/* Return text 'yyy' in 'xxyyyxx' */
func (pr *PassInlineMacroRxres) InlineText() string {
	res := pr.Group(2)
	if res == "" {
		res = pr.Group(4)
	}
	return res
}

/* Check if has sub 'xx' in 'pass:xx[yyy]' */
func (pr *PassInlineMacroRxres) HasPassSub() bool {
	return pr.HasGroup(5)
}

/* Return sub 'xx' in 'pass:xx[yyy]' */
func (pr *PassInlineMacroRxres) PassSub() string {
	return pr.Group(5)
}

/* Return text 'xx' in 'xxyyyxx' */
func (pr *PassInlineMacroRxres) InlineSub() string {
	res := pr.Group(1)
	if res == "" {
		res = pr.Group(3)
	}
	return res
}

/* Matches an xref (i.e., cross-reference) inline macro,
which may span multiple lines.
 Examples
   <<id,reftext>>
   xref:id[reftext]
NOTE special characters have already been escaped, hence the entity references
XrefInlineMacroRx = /\\?(?:&lt;&lt;([\w":].*?)&gt;&gt;|xref:([\w":].*?)\[(.*?)\])/m */

var XrefInlineMacroRx, _ = regexp.Compile(`(?s)\\?(?:&lt;&lt;([\w":].*?)&gt;&gt;|xref:([\w":].*?)\[(.*?)\])`)

type XrefInlineMacroRxres struct {
	*Reres
}

/* Results for XrefInlineMacroRx */
func NewXrefInlineMacroRxres(s string) *XrefInlineMacroRxres {
	return &XrefInlineMacroRxres{NewReres(s, XrefInlineMacroRx)}
}

/* Return id of '<<id,reftext>>' or xref:id[reftext]' */
func (ximr *XrefInlineMacroRxres) XId() string {
	if ximr.Group(1) != "" {
		t := strings.Split(ximr.Group(1), ",")
		return t[0]
	} else {
		return ximr.Group(2)
	}
}

/* Return reftext of '<<id,reftext>>' or xref:id[reftext]' */
func (ximr *XrefInlineMacroRxres) XrefText() string {
	if ximr.Group(1) != "" {
		t := strings.Split(ximr.Group(1), ",")
		return t[1]
	} else {
		return ximr.Group(3)
	}
}

/* Detects strings that resemble URIs.

   Examples
     http://domain
     https://domain
     data:info */
var UriSniffRx, _ = regexp.Compile(fmt.Sprintf("^([%v][%v.+-]*:/{0,2}).*", CC_ALPHA, CC_ALNUM))

/* Detects the end of an implicit URI in the text
 Examples
   (http://google.com)
   &gt;http://google.com&lt;
   (See http://google.com):
UriTerminator = /[);:]$/ */
var UriTerminator, _ = regexp.Compile(`[);:]$`)

/* Detects escaped brackets */
var EscapedBracketRx, _ = regexp.Compile(`\\\]`)

type Replacement struct {
	rx                *regexp.Regexp
	leading           bool
	bounding          bool
	repl              string
	endsWithLookAhead bool
}

func (r *Replacement) Rx() *regexp.Regexp      { return r.rx }
func (r *Replacement) Leading() bool           { return r.leading }
func (r *Replacement) Bounding() bool          { return r.bounding }
func (r *Replacement) None() bool              { return !r.leading && !r.bounding }
func (r *Replacement) Repl() string            { return r.repl }
func (r *Replacement) EndsWithLookAhead() bool { return r.endsWithLookAhead }
func (r *Replacement) Reres(text string) *Reres {
	if r.EndsWithLookAhead() {
		return NewReresLAGroup(text, r.Rx())
	}
	return NewReres(text, r.Rx())
}

var Replacements []*Replacement = iniReplacements()

func Rtos(runes ...rune) string {
	res := ""
	for _, r := range runes {
		res = res + string(r)
	}
	return res
}
func iniReplacements() []*Replacement {
	res := []*Replacement{}
	var rx *regexp.Regexp = nil
	// (C)
	rx, _ = regexp.Compile(`\\?\(C\)`)
	res = append(res, &Replacement{rx, false, false, Rtos(169), false})
	// (R)
	rx, _ = regexp.Compile(`\\?\(R\)`)
	res = append(res, &Replacement{rx, false, false, Rtos(174), false})
	// (TM)
	rx, _ = regexp.Compile(`\\?\(TM\)`)
	res = append(res, &Replacement{rx, false, false, Rtos(8482), false})
	// foo -- bar
	rx, _ = regexp.Compile(`(^|\n| |\\)--( |\n|$)`)
	res = append(res, &Replacement{rx, false, false, Rtos(8201, 8212, 8201), false})
	// foo--bar
	rx, _ = regexp.Compile(`(\w)(\\?--)(\w)`)
	res = append(res, &Replacement{rx, true, false, Rtos(8212), true})
	// ellipsis
	rx, _ = regexp.Compile(`\\?\.\.\.`)
	res = append(res, &Replacement{rx, true, false, Rtos(8230), false})
	// apostrophe or a closing single quote (planned)
	rx, _ = regexp.Compile(`([a-zA-Z])(\\?')($|[^'])`)
	res = append(res, &Replacement{rx, true, false, Rtos(8217), true})
	// an opening single quote (planned)
	// #[/\B\\?'(?=[#{CC_ALPHA}])/, '&#8216;', :none],
	// right arrow ->
	rx, _ = regexp.Compile(`\\?-&gt;`)
	res = append(res, &Replacement{rx, false, false, Rtos(8594), false})
	// right double arrow =>
	rx, _ = regexp.Compile(`\\?=&gt;`)
	res = append(res, &Replacement{rx, false, false, Rtos(8658), false})
	// left arrow <-
	rx, _ = regexp.Compile(`\\?&lt;-`)
	res = append(res, &Replacement{rx, false, false, Rtos(8592), false})
	// left double arrow <=
	rx, _ = regexp.Compile(`\\?&lt;=`)
	res = append(res, &Replacement{rx, false, false, Rtos(8656), false})
	// restore entities
	rx, _ = regexp.Compile(`\\?(&)amp;((?:[a-zA-Z]+|#\d{2,5}|#x[a-fA-F0-9]{2,4});)`)
	res = append(res, &Replacement{rx, false, true, "", false})
	return res
}
