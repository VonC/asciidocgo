package quotes

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestQuotes(t *testing.T) {

	Convey("Quotes subs have a fixed number of regexps", t, func() {
		So(len(QuoteSubs), ShouldEqual, 1)
	})
	Convey("Quotes subs should detect quotes", t, func() {
		someQuotes := []string{
			"**Git**Hub",
		}
		for i, aquote := range someQuotes {
			reres := NewQuoteSubRxres(aquote, QuoteSubs[i])
			So(reres.HasAnyMatch(), ShouldBeTrue)
		}
	})
}
