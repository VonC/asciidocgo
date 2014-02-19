package quotes

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestQuotes(t *testing.T) {

	Convey("Quotes subs have a fixed number of regexps", t, func() {
		So(len(QuoteSubs), ShouldEqual, 13)
	})

	Convey("Quotes subs should detect unconstrained **strong** quotes", t, func() {
		qs := QuoteSubs[0]
		So(qs.Rx().String(), ShouldEqual, `(?s)\\?(?:\[([^\]]+?)\])?\*\*(.+?)\*\*`)
		So(qs.TypeQS(), ShouldEqual, Strong)
		So(qs.Constrained(), ShouldBeFalse)
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
	Convey("Quotes subs should detect constrained 'emphasis' quotes", t, func() {
		qs := QuoteSubs[3]
		Convey("single-line constrained quote variation emphasized string", func() {
			reres := NewQuoteSubRxres("'a few emphasized words'", qs)
			So(reres.Prefix(), ShouldEqual, "")
			So(reres.PrefixQuote(), ShouldEqual, "")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "a few emphasized words")
		})
		Convey("escaped single-line constrained quote variation emphasized string", func() {
			reres := NewQuoteSubRxres("\\'a few emphasized words'", qs)
			So(reres.Prefix(), ShouldEqual, "")
			So(reres.PrefixQuote(), ShouldEqual, "\\")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "a few emphasized words")
		})
		Convey("multi-line constrained emphasized quote variation string", func() {
			reres := NewQuoteSubRxres("'a few\nemphasized words'", qs)
			So(reres.Prefix(), ShouldEqual, "")
			So(reres.PrefixQuote(), ShouldEqual, "")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "a few\nemphasized words")
		})
		Convey("single-quoted string containing an emphasized phrase", func() {
			reres := NewQuoteSubRxres("`I told him, 'Just go for it!''", qs)
			So(reres.Prefix(), ShouldEqual, "`I told him,")
			So(reres.PrefixQuote(), ShouldEqual, " ")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "Just go for it!")
		})
		Convey("escaped single-quotes inside emphasized words are restored", func() {
			reres := NewQuoteSubRxres("'Here\\'s Johnny!'", qs)
			So(reres.Prefix(), ShouldEqual, "")
			So(reres.PrefixQuote(), ShouldEqual, "")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "Here\\'s Johnny!")
		})
	})
	Convey("Quotes subs should detect constrained `single-quoted' quotes", t, func() {
		qs := QuoteSubs[4]
		Convey("single-line single-quoted string", func() {
			reres := NewQuoteSubRxres("`a few quoted words'", qs)
			So(reres.Prefix(), ShouldEqual, "")
			So(reres.PrefixQuote(), ShouldEqual, "")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "a few quoted words")
		})
		Convey("escaped single-line single-quoted string", func() {
			reres := NewQuoteSubRxres("\\`a few quoted words'", qs)
			So(reres.Prefix(), ShouldEqual, "")
			So(reres.PrefixQuote(), ShouldEqual, "\\")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "a few quoted words")
		})
		Convey("multi-line single-quoted string", func() {
			reres := NewQuoteSubRxres("`a few\nquoted words'", qs)
			So(reres.Prefix(), ShouldEqual, "")
			So(reres.PrefixQuote(), ShouldEqual, "")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "a few\nquoted words")
		})
		Convey("single-quoted string with inline single quote", func() {
			reres := NewQuoteSubRxres("`That isn't what I did.'", qs)
			So(reres.Prefix(), ShouldEqual, "")
			So(reres.PrefixQuote(), ShouldEqual, "")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "That isn't what I did.")
		})
		Convey("single-quoted string with inline backquote", func() {
			reres := NewQuoteSubRxres("`Here`s Johnny!'", qs)
			So(reres.Prefix(), ShouldEqual, "")
			So(reres.PrefixQuote(), ShouldEqual, "")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "Here`s Johnny!")
		})
	})
	Convey("Quotes subs should detect unconstrained ++monospaced++ quotes", t, func() {
		qs := QuoteSubs[5]
		Convey("single-line unconstrained monospaced chars", func() {
			reres := NewQuoteSubRxres("Git++Hub++", qs)
			So(reres.Prefix(), ShouldEqual, "Git")
			So(reres.PrefixQuote(), ShouldEqual, "")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "Hub")
		})
		Convey("escaped single-line unconstrained monospaced chars", func() {
			reres := NewQuoteSubRxres("Git\\++Hub++", qs)
			So(reres.Prefix(), ShouldEqual, "Git")
			So(reres.PrefixQuote(), ShouldEqual, "") // TOFIX sould be + here? :  assert_equal 'Git+<code>Hub</code>+'
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "Hub")
		})
		Convey("multi-line unconstrained monospaced chars", func() {
			reres := NewQuoteSubRxres("Git++\nH\nu\nb++", qs)
			So(reres.Prefix(), ShouldEqual, "Git")
			So(reres.PrefixQuote(), ShouldEqual, "")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "\nH\nu\nb")
		})
	})
	Convey("Quotes subs should detect constrained +monospaced+ quotes", t, func() {
		qs := QuoteSubs[6]
		Convey("single-line constrained monospaced chars", func() {
			reres := NewQuoteSubRxres("call +save()+ to persist the changes", qs)
			So(reres.Prefix(), ShouldEqual, "call")
			So(reres.PrefixQuote(), ShouldEqual, " ")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "save()")
		})
		Convey("single-line constrained monospaced chars with role", func() {
			reres := NewQuoteSubRxres("call [method]+save()+ to persist the changes", qs)
			So(reres.Prefix(), ShouldEqual, "call")
			So(reres.PrefixQuote(), ShouldEqual, " ")
			So(reres.Attribute(), ShouldEqual, "method")
			So(reres.Quote(), ShouldEqual, "save()")
		})
		Convey("escaped single-line constrained monospaced chars", func() {
			reres := NewQuoteSubRxres(`call \+save()+ to persist the changes`, qs)
			So(reres.Prefix(), ShouldEqual, "call ")
			So(reres.PrefixQuote(), ShouldEqual, "\\")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "save()")
		})
		Convey("escaped single-line constrained monospaced chars with role", func() {
			reres := NewQuoteSubRxres(`call [method]\+save()+ to persist the changes`, qs)
			So(reres.Prefix(), ShouldEqual, "call [method]")
			So(reres.PrefixQuote(), ShouldEqual, "\\")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "save()")
		})
		Convey("escaped role on single-line constrained monospaced chars", func() {
			reres := NewQuoteSubRxres(`call \[method]+save()+ to persist the changes`, qs)
			So(reres.Prefix(), ShouldEqual, "call ")
			So(reres.PrefixQuote(), ShouldEqual, "\\")
			So(reres.Attribute(), ShouldEqual, "method")
			So(reres.Quote(), ShouldEqual, "save()")
		})
		Convey("escaped role on escaped single-line constrained monospaced chars", func() {
			reres := NewQuoteSubRxres(`call \[method]\+save()+ to persist the changes`, qs)
			So(reres.Prefix(), ShouldEqual, "call \\[method]")
			So(reres.PrefixQuote(), ShouldEqual, "\\")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "save()")
		})
	})
	Convey("Quotes subs should detect unconstrained __emphasis__ quotes", t, func() {
		qs := QuoteSubs[7]
		Convey("single-line unconstrained emphasized chars", func() {
			reres := NewQuoteSubRxres("__Git__Hub", qs)
			So(reres.Prefix(), ShouldEqual, "")
			So(reres.PrefixQuote(), ShouldEqual, "")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "Git")
		})
		Convey("escaped single-line unconstrained emphasized chars", func() {
			reres := NewQuoteSubRxres(`\__Git__Hub`, qs)
			So(reres.Prefix(), ShouldEqual, "")
			So(reres.PrefixQuote(), ShouldEqual, "")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "Git")
		})
		Convey("multi-line unconstrained emphasized chars", func() {
			reres := NewQuoteSubRxres("__G\ni\nt\n__Hub", qs)
			So(reres.Prefix(), ShouldEqual, "")
			So(reres.PrefixQuote(), ShouldEqual, "")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "G\ni\nt\n")
		})
		Convey("unconstrained emphasis chars with role", func() {
			reres := NewQuoteSubRxres("[gray]__Git__Hub", qs)
			So(reres.Prefix(), ShouldEqual, "")
			So(reres.PrefixQuote(), ShouldEqual, "")
			So(reres.Attribute(), ShouldEqual, "gray")
			So(reres.Quote(), ShouldEqual, "Git")
		})
		Convey("escaped unconstrained emphasis chars with role", func() {
			reres := NewQuoteSubRxres(`\[gray]__Git__Hub`, qs)
			So(reres.Prefix(), ShouldEqual, "")
			So(reres.PrefixQuote(), ShouldEqual, "")
			So(reres.Attribute(), ShouldEqual, "gray")
			So(reres.Quote(), ShouldEqual, "Git")
		})
	})
	Convey("Quotes subs should detect constrained _emphasis_ quotes", t, func() {
		qs := QuoteSubs[8]
		Convey("single-line constrained emphasized underline variation string", func() {
			reres := NewQuoteSubRxres("_a few emphasized words_", qs)
			So(reres.Prefix(), ShouldEqual, "")
			So(reres.PrefixQuote(), ShouldEqual, "")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "a few emphasized words")
		})
		Convey("escaped single-line constrained emphasized underline variation string", func() {
			reres := NewQuoteSubRxres(`\_a few emphasized words_`, qs)
			So(reres.Prefix(), ShouldEqual, "")
			So(reres.PrefixQuote(), ShouldEqual, "\\")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "a few emphasized words")
		})
		Convey("multi-line constrained emphasized underline variation string", func() {
			reres := NewQuoteSubRxres("_a few\nemphasized words_", qs)
			So(reres.Prefix(), ShouldEqual, "")
			So(reres.PrefixQuote(), ShouldEqual, "")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "a few\nemphasized words")
		})
	})
	Convey("Quotes subs should detect unconstrained ##unquoted## quotes", t, func() {
		qs := QuoteSubs[9]
		Convey("single-line unconstrained unquoted string", func() {
			reres := NewQuoteSubRxres("##--anything goes ##", qs)
			So(reres.Prefix(), ShouldEqual, "")
			So(reres.PrefixQuote(), ShouldEqual, "")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "--anything goes ")
		})
		Convey("escaped single-line unconstrained unquoted string", func() {
			reres := NewQuoteSubRxres(`\##--anything goes ##`, qs)
			So(reres.Prefix(), ShouldEqual, "")
			So(reres.PrefixQuote(), ShouldEqual, "") // TOFIX? Why not \\, as in 365?
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "--anything goes ")
		})
		Convey("multi-line unconstrained unquoted string", func() {
			reres := NewQuoteSubRxres("##--anything\ngoes ##", qs)
			So(reres.Prefix(), ShouldEqual, "")
			So(reres.PrefixQuote(), ShouldEqual, "")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "--anything\ngoes ")
		})
	})
	Convey("Quotes subs should detect constrained #unquoted# quotes", t, func() {
		qs := QuoteSubs[10]
		Convey("single-line constrained unquoted string", func() {
			reres := NewQuoteSubRxres("#a few words#", qs)
			So(reres.Prefix(), ShouldEqual, "")
			So(reres.PrefixQuote(), ShouldEqual, "")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "a few words")
		})
		Convey("escaped single-line constrained unquoted string", func() {
			reres := NewQuoteSubRxres(`\#a few words#`, qs)
			So(reres.Prefix(), ShouldEqual, "")
			So(reres.PrefixQuote(), ShouldEqual, "\\")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "a few words")
		})
		Convey("multi-line constrained unquoted string", func() {
			reres := NewQuoteSubRxres("#a few\nwords#", qs)
			So(reres.Prefix(), ShouldEqual, "")
			So(reres.PrefixQuote(), ShouldEqual, "")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "a few\nwords")
		})
	})
	Convey("Quotes subs should detect unconstrained ^superscript^ quotes", t, func() {
		qs := QuoteSubs[11]
		Convey("single-line superscript chars", func() {
			reres := NewQuoteSubRxres("x^2^ = x * x, e = mc^2^, there's a 1^st^ time for everything", qs)
			So(reres.Prefix(), ShouldEqual, "x")
			So(reres.PrefixQuote(), ShouldEqual, "")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "2")
			reres.Next()
			So(reres.Prefix(), ShouldEqual, " = x * x, e = mc")
			So(reres.PrefixQuote(), ShouldEqual, "")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "2")
			reres.Next()
			So(reres.Prefix(), ShouldEqual, ", there's a 1")
			So(reres.PrefixQuote(), ShouldEqual, "")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "st")
		})
		Convey("escaped single-line superscript chars", func() {
			reres := NewQuoteSubRxres(`x\^2^ = x * x`, qs)
			So(reres.Prefix(), ShouldEqual, "x")
			So(reres.PrefixQuote(), ShouldEqual, "")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "2")
			So(reres.IsEscaped(), ShouldBeTrue)
		})
		Convey("multi-line superscript chars", func() {
			reres := NewQuoteSubRxres("x^(n\n-\n1)^", qs)
			So(reres.Prefix(), ShouldEqual, "x")
			So(reres.PrefixQuote(), ShouldEqual, "")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "(n\n-\n1)")
		})
	})
	Convey("Quotes subs should detect unconstrained ~subscript~ quotes", t, func() {
		qs := QuoteSubs[12]
		Convey("single-line subscript chars", func() {
			reres := NewQuoteSubRxres("H~2~O", qs)
			So(reres.Prefix(), ShouldEqual, "H")
			So(reres.PrefixQuote(), ShouldEqual, "")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "2")
		})
		Convey("escaped single-line subscript chars", func() {
			reres := NewQuoteSubRxres(`H\~2~O`, qs)
			So(reres.Prefix(), ShouldEqual, "H")
			So(reres.PrefixQuote(), ShouldEqual, "")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, "2")
			So(reres.IsEscaped(), ShouldBeTrue)
		})
		Convey("multi-line subscript chars", func() {
			reres := NewQuoteSubRxres("project~ view\non\nGitHub~", qs)
			So(reres.Prefix(), ShouldEqual, "project")
			So(reres.PrefixQuote(), ShouldEqual, "")
			So(reres.Attribute(), ShouldEqual, "")
			So(reres.Quote(), ShouldEqual, " view\non\nGitHub")
		})
	})
}
