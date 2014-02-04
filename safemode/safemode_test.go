package safemode

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSafeMode(t *testing.T) {

	Convey("A safemode is gradual", t, func() {
		So(UNSAFE < SAFE, ShouldBeTrue)
		So(SAFE < SERVER, ShouldBeTrue)
		So(SERVER < SECURE, ShouldBeTrue)
		So(SECURE < PARANOID, ShouldBeTrue)
	})

}
