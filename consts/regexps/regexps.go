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

/* Build new result from FindAllStringSubmatchIndex on a string */
func NewReres(s string, r *regexp.Regexp) *Reres {
	matches := r.FindAllStringSubmatchIndex(s, -1)
	return &Reres{r, s, matches, 0, 0}
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
	return rr.s[rr.previous:]
}

/* First character of the current match */
func (rr *Reres) FirstChar() uint8 {
	mi := rr.matches[rr.i]
	return rr.s[mi[0]]
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

/* Matches several variants of the passthrough inline macro,
which may span multiple lines.

 Examples

   +++text+++
   $$text$$
   pass:quotes[text] */
// http://stackoverflow.com/questions/6770898/unknown-escape-sequence-error-in-go
var PassInlineMacroRx, _ = regexp.Compile(`(?s)\\?(?:(\+{3})(.*?)\+{3}|(\${2})(.*?)\${2}|pass:([a-z,]*)\[(.*?[^\\])\])`)

/* Detects strings that resemble URIs.

   Examples
     http://domain
     https://domain
     data:info */
var UriSniffRx, _ = regexp.Compile(fmt.Sprintf("^([%v][%v.+-]*:/{0,2}).*", CC_ALPHA, CC_ALNUM))

/* Detects escaped brackets */
var EscapedBracketRx, _ = regexp.Compile(`\\\]`)
