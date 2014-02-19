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
	res := &QuoteSubRxres{nil, qs}
	if qs.constrained {
		res.Reres = regexps.NewReresLAGroup(s, qs.rx)
	} else {
		res.Reres = regexps.NewReres(s, qs.rx)
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
	res = addQuoteSub(res, Strong, true, `(?s)(^|[^\w;:}])(?:\[([^\]]+?)\])?\*(\S|\S.*?\S)\*($|\W)`)
	// ``double-quoted''
	res = addQuoteSub(res, Double, true, "(?s)(^|[^\\w;:}])(?:\\[([^\\]]+?)\\])?``(\\S|\\S.*?\\S)''(\\W|$)")
	// 'emphasis'
	res = addQuoteSub(res, Emphasis, true, `(?s)(^|[^\w;:}])(?:\[([^\]]+?)\])?'(\S|\S.*?\S)'(\W|$)`)
	// `single-quoted'
	res = addQuoteSub(res, Single, true, "(?s)(^|[^\\w;:}])(?:\\[([^\\]]+?)\\])?`(\\S|\\S.*?\\S)'(\\W|$)")
	// ++monospaced++
	res = addQuoteSub(res, Monospaced, false, `(?s)\\?(?:\[([^\]]+?)\])?\+\+(.+?)\+\+`)
	// +monospaced+
	res = addQuoteSub(res, Monospaced, true, `(?s)(^|[^\w;:}])(?:\[([^\]]+?)\])?\+(\S|\S.*?\S)\+(\W|$)`)
	// __emphasis__
	res = addQuoteSub(res, Emphasis, false, `(?s)\\?(?:\[([^\]]+?)\])?\_\_(.+?)\_\_`)
	// _emphasis_
	res = addQuoteSub(res, Emphasis, true, `(?s)(^|[^\w;:}])(?:\[([^\]]+?)\])?_(\S|\S.*?\S)_(\W|$)`)
	// ##unquoted##
	res = addQuoteSub(res, None, false, `(?s)\\?(?:\[([^\]]+?)\])?##(.+?)##`)
	// #unquoted#
	res = addQuoteSub(res, None, true, `(?s)(^|[^\w;:}])(?:\[([^\]]+?)\])?#(\S|\S.*?\S)#(\W|$)`)
	return res
}

func addQuoteSub(res []*QuoteSub, typeqs QuoteSubType, constrained bool, rxp string) []*QuoteSub {
	rx, _ := regexp.Compile(rxp)
	qs := &QuoteSub{typeqs, constrained, rx}
	res = append(res, qs)
	return res
}
