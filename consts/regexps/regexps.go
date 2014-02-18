package regexps

import (
	"fmt"
	"regexp"

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

/* Encapsulate a regex and a string,
for managing results from FindAllStringSubmatchIndex */
type Reres struct {
	r        *regexp.Regexp
	s        string
	matches  [][]int
	i        int
	previous int
}

type lookahead func(suffix string) bool

/* Build new result from FindAllStringSubmatchIndex on a string */
func NewReres(s string, r *regexp.Regexp) *Reres {
	matches := r.FindAllStringSubmatchIndex(s, -1)
	return &Reres{r, s, matches, 0, 0}
}

/* Build new result from FindAllStringSubmatchIndex on a string,
validated by a lookahead after each match */
func NewReresLA(s string, r *regexp.Regexp, la lookahead) *Reres {
	res := NewReres(s, r)
	if la != nil && res.HasAnyMatch() {
		m := [][]int{}
		for res.HasNext() {
			if la(res.Suffix()) {
				m = append(m, res.matches[res.i])
			}
			res.Next()
		}
		res.ResetNext()
		res.matches = m
	}
	return res
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

var PassInlineLiteralRx, _ = regexp.Compile("(?sm)(^|[^`\\w])(?:\\[([^\\]]+?)\\])?(\\\\?`([^`\\s]|[^`\\s].*?\\S)`)")
var PassInlineMacroRxLookahead, _ = regexp.Compile("[`\\w]")

type PassInlineLiteralRxres struct {
	*Reres
}

func passInlineMacroRxLookahead(suffix string) bool {
	c := ""
	if suffix != "" {
		c = suffix[0:1]
	}
	b := PassInlineMacroRxLookahead.FindAllString(c, -1) == nil
	return b
}

/* Results for PassInlineLiteralRx */
func NewPassInlineLiteralRxres(s string) *PassInlineLiteralRxres {
	res := &PassInlineLiteralRxres{NewReresLA(s, PassInlineLiteralRx, passInlineMacroRxLookahead)}
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

/* Detects strings that resemble URIs.

   Examples
     http://domain
     https://domain
     data:info */
var UriSniffRx, _ = regexp.Compile(fmt.Sprintf("^([%v][%v.+-]*:/{0,2}).*", CC_ALPHA, CC_ALNUM))

/* Detects escaped brackets */
var EscapedBracketRx, _ = regexp.Compile(`\\\]`)
