package asciidocgo

import (
	"reflect"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

var dm = new(Document).Monitor()
var dnm = new(Document)
var dnmini = NewDocument([]string{}, map[string]string{})
var notMonitoredError = &NotMonitoredError{"test"}
var monitorFNames = []string{"ReadTime", "ParseTime", "LoadTime", "RenderTime", "LoadRenderTime", "WriteTime", "TotalTime"}

func TestDocumentMonitor(t *testing.T) {
	Convey("A Document can be monitored", t, func() {
		Convey("By default, a Document is not monitored", func() {
			So(dnm.IsMonitored(), ShouldBeFalse)
			So(dnmini.IsMonitored(), ShouldBeFalse)
		})
		Convey("A monitored Document is monitored", func() {
			So(dm.IsMonitored(), ShouldBeTrue)
		})
	})
	Convey("A non-monitored Document should return error when accessing times", t, func() {
		defer func() {
			if x := recover(); x != nil {
				So(x, ShouldBeNil)
			}
		}()
		dtype := reflect.ValueOf(dnm)
		for _, fname := range monitorFNames {
			dfunc := dtype.MethodByName(fname)
			ret := dfunc.Call([]reflect.Value{})
			err := ret[1].Interface().(error)
			So(err, ShouldNotBeNil)
			So(err, ShouldHaveSameTypeAs, notMonitoredError)
			So(err.Error(), ShouldContainSubstring, "not monitored")
		}
	})
	Convey("A monitored empty Document should return 0 when accessing times", t, func() {
		defer func() {
			if x := recover(); x != nil {
				So(x, ShouldBeNil)
			}
		}()
		dtype := reflect.ValueOf(dm)
		for _, fname := range monitorFNames {
			dfunc := dtype.MethodByName(fname)
			ret := dfunc.Call([]reflect.Value{})
			time := ret[0].Int()
			err := ret[1].Interface()
			So(err, ShouldBeNil)
			So(time, ShouldBeZeroValue)
		}
	})
	Convey("Load time equals read time + parse time", t, func() {
		loadTime, _ := dm.LoadTime()
		readTime, _ := dm.ReadTime()
		parseTime, _ := dm.ParseTime()
		So(loadTime, ShouldEqual, readTime+parseTime)
	})
	Convey("LoadRender time equals load time + render time", t, func() {
		loadRenderTime, _ := dm.LoadRenderTime()
		loadTime, _ := dm.LoadTime()
		renderTime, _ := dm.ReadTime()
		So(loadRenderTime, ShouldEqual, loadTime+renderTime)
	})
	Convey("Total time equals LoadRender time + write time", t, func() {
		totalTime, _ := dm.TotalTime()
		loadRenderTime, _ := dm.LoadRenderTime()
		writeTime, _ := dm.WriteTime()
		So(totalTime, ShouldEqual, loadRenderTime+writeTime)
	})
}

func TestDocumentInitialization(t *testing.T) {

	Convey("A Document can be initialized", t, func() {

		Convey("By default, a Document has no data, no options", func() {
			So(&Document{}, ShouldNotBeNil)
		})
		Convey("A Document take an array of strings as data, and a map as options", func() {
			So(NewDocument([]string{}, map[string]string{}), ShouldNotBeNil)
		})
	})
}
