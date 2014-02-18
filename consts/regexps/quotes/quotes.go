package quotes

import (
	"regexp"

	"github.com/VonC/asciidocgo/consts/regexps"
)

type QuoteSubType int

const (
	Strong QuoteSubType = iota
	Double
	Emphasis
	Single
	Monospaced
	None
	Superscript
	Subscript
)

type QuoteSub struct {
	typeqs      QuoteSubType
	constrained bool
	rx          *regexp.Regexp
}

type QuoteSubRxres struct {
	*regexps.Reres
	qs *QuoteSub
}

var constraintRx, _ = regexp.Compile(`\W`)

func quoteSubLookAhead(suffix string) bool {
	if suffix == "" {
		return true
	}
	c := suffix[0:1]
	//fmt.Printf("c%v====\n", c)
	r := constraintRx.FindStringSubmatch(c)
	//fmt.Printf("c%v====%d\n", r, len(r))
	return r != nil && len(r) > 0
}

/* Results for QuoteSubRxres */
func NewQuoteSubRxres(s string, qs *QuoteSub) *QuoteSubRxres {
	res := &QuoteSubRxres{regexps.NewReres(s, qs.rx), qs}
	if qs.constrained {
		res = &QuoteSubRxres{regexps.NewReresLA(s, qs.rx, quoteSubLookAhead), qs}
	}
	return res
}

func (qsr *QuoteSubRxres) PrefixQuote() string {
	if !qsr.HasAnyMatch() || !qsr.qs.constrained {
		return ""
	}
	return qsr.Group(1)
}

func (qsr *QuoteSubRxres) Attribute() string {
	if !qsr.HasAnyMatch() {
		return ""
	}
	if qsr.qs.constrained {
		return qsr.Group(2)
	}
	return qsr.Group(1)
}

func (qsr *QuoteSubRxres) Quote() string {
	if !qsr.HasAnyMatch() {
		return ""
	}
	if qsr.qs.constrained {
		return qsr.Group(3)
	}
	return qsr.Group(2)
}

/* unconstrained quotes:: can appear anywhere
   constrained quotes:: must be bordered by non-word characters
   NOTE these substitutions are processed in the order they appear here and
   the order in which they are replaced is important
*/
//QUOTE_SUBS = [

var QuoteSubs []*QuoteSub = iniQuoteSubs()

func iniQuoteSubs() []*QuoteSub {
	res := []*QuoteSub{}
	// **strong**
	res = addQuoteSub(res, Strong, false, `(?s)\\?(?:\[([^\]]+?)\])?\*\*(.+?)\*\*`)
	// *strong*
	res = addQuoteSub(res, Strong, true, `(?s)(^|[^\w;:}])(?:\[([^\]]+?)\])?\*(\S|\S.*?\S)\*`)
	return res
}

func addQuoteSub(res []*QuoteSub, typeqs QuoteSubType, constrained bool, rxp string) []*QuoteSub {
	rx, _ := regexp.Compile(rxp)
	qs := &QuoteSub{typeqs, constrained, rx}
	res = append(res, qs)
	return res
}
