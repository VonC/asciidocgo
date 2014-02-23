package debug

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRegexps(t *testing.T) {

	Convey("Debug can print message", t, func() {
		Debug("test")
		// TODO really test with a test logger instead of the default one
		So(debugLog, ShouldNotBeNil)
		Switch()
		Debug("test2")
		Switch()
		Debug("test3")
	})
}
