package quotes

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestQuotes(t *testing.T) {

	Convey("Quotes subs have a fixed number of regexps", t, func() {
		So(len(QuoteSubs), ShouldEqual, 1)
	})

	Convey("Quotes subs should detect unconstrained **strong** quotes", t, func() {
		qs := QuoteSubs[0]
		Convey("single-line unconstrained strong chars", func() {
			reres := NewQuoteSubRxres("**Git**Hub", qs)
			So(reres.Prefix(), ShouldEqual, "")
			So(reres.PrefixQuote(), ShouldEqual, "")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "Git")
		})
		Convey("escaped single-line unconstrained strong chars", func() {
			reres := NewQuoteSubRxres(`\**Git**Hub`, qs)
			So(reres.Prefix(), ShouldEqual, "")
			So(reres.PrefixQuote(), ShouldEqual, "")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "Git")
		})
		Convey("multi-line unconstrained strong chars", func() {
			reres := NewQuoteSubRxres("**G\ni\nt\n**Hub", qs)
			So(reres.Prefix(), ShouldEqual, "")
			So(reres.PrefixQuote(), ShouldEqual, "")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "G\ni\nt\n")
		})
		Convey("unconstrained strong chars with inline asterisk", func() {
			reres := NewQuoteSubRxres("**bl*ck**-eye", qs)
			So(reres.Prefix(), ShouldEqual, "")
			So(reres.PrefixQuote(), ShouldEqual, "")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "bl*ck")
		})
		Convey("unconstrained strong chars with role", func() {
			reres := NewQuoteSubRxres("Git[blue]**Hub**", qs)
			So(reres.Prefix(), ShouldEqual, "Git")
			So(reres.PrefixQuote(), ShouldEqual, "")
			So(reres.Attribute(), ShouldEqual, "blue")
			So(reres.Quote(), ShouldEqual, "Hub")
		})
		Convey("escaped unconstrained strong chars with role", func() {
			reres := NewQuoteSubRxres(`Git\[blue]**Hub**`, qs)
			So(reres.Prefix(), ShouldEqual, "Git")
			So(reres.PrefixQuote(), ShouldEqual, "")
			So(reres.Attribute(), ShouldEqual, "blue")
			So(reres.Quote(), ShouldEqual, "Hub")
		})

	})
}
