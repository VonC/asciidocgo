package asciidocgo

import (
	"fmt"
	"reflect"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

var dm = new(Document).Monitor()
var dnm = new(Document)
var notMonitoredError = &NotMonitoredError{"test"}
var monitorFNames = [1]string{"ReadTime"}

func TestDocumentMonitor(t *testing.T) {
	Convey("A Document can be monitored", t, func() {
		Convey("By default, a Document is not monitored", func() {
			So(dnm.IsMonitored(), ShouldBeFalse)
		})
		Convey("A monitored Document is monitored", func() {
			So(dm.IsMonitored(), ShouldBeTrue)
		})
	})
	Convey("A non-monitored Document should return error when accessing times", t, func() {
		_, err := dnm.ReadTime()
		So(err, ShouldNotBeNil)
		So(err, ShouldHaveSameTypeAs, notMonitoredError)
		So(err.Error(), ShouldContainSubstring, "not monitored")
	})
	Convey("A monitored empty Document should return 0 when accessing times", t, func() {
		dtype := reflect.ValueOf(dm)
		for _, fname := range monitorFNames {
			dfunc := dtype.MethodByName(fname)
			ret := dfunc.Call([]reflect.Value{})
			So(ret[1], shouldBeNilReflectValue)
			So(ret[0].Int(), ShouldBeZeroValue)
		}
	})
}

func shouldBeNilReflectValue(actual interface{}, expected ...interface{}) string {
	if actual.(reflect.Value).IsNil() {
		return ""
	}
	return "Value " + fmt.Sprintf("%v", actual) + " should be nil"
}
