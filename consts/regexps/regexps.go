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

type Reres struct {
	r        *regexp.Regexp
	s        string
	matches  [][]int
	i        int
	previous int
}

func NewReres(s string, r *regexp.Regexp) *Reres {
	matches := r.FindAllStringSubmatchIndex(s, -1)
	return &Reres{r, s, matches, 0, 0}
}

func (rr *Reres) HasAnyMatch() bool {
	return len(rr.matches) > 0
}

func (rr *Reres) HasNext() bool {
	return rr.i < len(rr.matches)
}

func (rr *Reres) Next() {
	rr.previous = rr.matches[rr.i][1]
	rr.i = rr.i + 1
}

func (rr *Reres) Previous() string {
	mi := rr.matches[rr.i]
	return rr.s[rr.previous:mi[0]]
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
