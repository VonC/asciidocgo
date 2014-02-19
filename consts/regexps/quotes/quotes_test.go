package quotes

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestQuotes(t *testing.T) {

	Convey("Quotes subs have a fixed number of regexps", t, func() {
		So(len(QuoteSubs), ShouldEqual, 3)
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

	Convey("Quotes subs should detect constrained *strong* quotes", t, func() {
		qs := QuoteSubs[1]
		Convey("single-line constrained strong string", func() {
			reres := NewQuoteSubRxres(`*a few strong words*`, qs)
			So(reres.Prefix(), ShouldEqual, "")
			So(reres.PrefixQuote(), ShouldEqual, "")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "a few strong words")
		})
		Convey("single-line constrained strong string", func() {
			reres := NewQuoteSubRxres(`*a few strong failed words*a`, qs)
			So(reres.HasAnyMatch(), ShouldBeFalse)
			So(reres.PrefixQuote(), ShouldEqual, "")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "")
		})
		Convey("escaped single-line constrained strong string", func() {
			reres := NewQuoteSubRxres(`\*a few strong words*`, qs)
			So(reres.Prefix(), ShouldEqual, "")
			So(reres.PrefixQuote(), ShouldEqual, "\\")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "a few strong words")
		})
		Convey("multi-line constrained strong string", func() {
			reres := NewQuoteSubRxres("*a few\nstrong words*", qs)
			So(reres.Prefix(), ShouldEqual, "")
			So(reres.PrefixQuote(), ShouldEqual, "")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "a few\nstrong words")
		})
		Convey("constrained strong string containing an asterisk", func() {
			reres := NewQuoteSubRxres("*bl*ck*-eye-*2d*word*--", qs)
			So(reres.Prefix(), ShouldEqual, "")
			So(reres.PrefixQuote(), ShouldEqual, "")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "bl*ck")
			reres.Next()
			So(reres.Prefix(), ShouldEqual, "-eye")
			So(reres.PrefixQuote(), ShouldEqual, "-")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "2d*word")
		})

		Convey("consecutive constrained strong string containing an asterisk", func() {
			reres := NewQuoteSubRxres("*bl*ck*-*2d*word*--", qs)
			So(reres.Prefix(), ShouldEqual, "")
			So(reres.PrefixQuote(), ShouldEqual, "")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "bl*ck")
			So(reres.Suffix(), ShouldEqual, "-*2d*word*--")
			reres.Next()
			So(reres.Prefix(), ShouldEqual, "")
			So(reres.PrefixQuote(), ShouldEqual, "-")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "2d*word")
			So(reres.Suffix(), ShouldEqual, "--")

			reres = NewQuoteSubRxres("*bl*ck**2d*word*--", qs)
			So(reres.Prefix(), ShouldEqual, "")
			So(reres.PrefixQuote(), ShouldEqual, "")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "bl*ck")
			So(reres.Suffix(), ShouldEqual, "*2d*word*--")
			reres.Next()
			So(reres.Prefix(), ShouldEqual, "")
			So(reres.PrefixQuote(), ShouldEqual, "")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "2d*word")
			So(reres.Suffix(), ShouldEqual, "--")
		})

	})
	Convey("Quotes subs should detect constrained ``double-quoted'' quotes", t, func() {
		qs := QuoteSubs[2]
		Convey("single-line double-quoted string", func() {
			reres := NewQuoteSubRxres("``a few quoted words''", qs)
			So(reres.Prefix(), ShouldEqual, "")
			So(reres.PrefixQuote(), ShouldEqual, "")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "a few quoted words")
		})
		Convey("escaped single-line double-quoted string", func() {
			reres := NewQuoteSubRxres("\\``a few quoted words''", qs)
			So(reres.Prefix(), ShouldEqual, "")
			So(reres.PrefixQuote(), ShouldEqual, "\\")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "a few quoted words")
		})
		Convey("multi-line double-quoted string", func() {
			reres := NewQuoteSubRxres("``a few\nquoted words''", qs)
			So(reres.Prefix(), ShouldEqual, "")
			So(reres.PrefixQuote(), ShouldEqual, "")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "a few\nquoted words")
		})
		Convey("double-quoted string with inline single quote", func() {
			reres := NewQuoteSubRxres("``Here's Johnny!''", qs)
			So(reres.Prefix(), ShouldEqual, "")
			So(reres.PrefixQuote(), ShouldEqual, "")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "Here's Johnny!")
		})
		Convey("double-quoted string with inline backquote", func() {
			reres := NewQuoteSubRxres("``Here`s Johnny!''", qs)
			So(reres.Prefix(), ShouldEqual, "")
			So(reres.PrefixQuote(), ShouldEqual, "")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "Here`s Johnny!")
		})
	})
}
