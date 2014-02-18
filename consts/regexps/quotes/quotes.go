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

/* Results for QuoteSubRxres */
func NewQuoteSubRxres(s string, qs *QuoteSub) *QuoteSubRxres {
	return &QuoteSubRxres{regexps.NewReres(s, qs.rx), qs}
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
	res = addQuoteSub(res, Strong, false, `(?m)\\?(?:\[([^\]]+?)\])?\*\*(.+?)\*\*`)
	return res
}

func addQuoteSub(res []*QuoteSub, typeqs QuoteSubType, constrained bool, rxp string) []*QuoteSub {
	rx, _ := regexp.Compile(rxp)
	qs := &QuoteSub{typeqs, constrained, rx}
	res = append(res, qs)
	return res
}
