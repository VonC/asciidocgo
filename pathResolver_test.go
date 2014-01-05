package asciidocgo

import (
	"os"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPathResolver(t *testing.T) {

	Convey("A pathResolver can be initialized", t, func() {

		Convey("By default, a pathResolver can be created", func() {
			So(newPathResolver(0, ""), ShouldNotBeNil)
		})
		Convey("By default, a pathResolver has a system path separator", func() {
			So(newPathResolver(0, "").FileSeparator(), ShouldEqual, os.PathSeparator)
			So(newPathResolver('/', "").FileSeparator(), ShouldNotEqual, os.PathSeparator)
			So(newPathResolver('/', "").FileSeparator(), ShouldEqual, '/')
		})

		Convey("By default, a pathResolver has a current working path", func() {
			pwd, err := os.Getwd()
			if err != nil {
				panic(err)
			}
			So(newPathResolver(0, "").WorkingDir(), ShouldEqual, pwd)
			So(newPathResolver(0, "C:\\").WorkingDir(), ShouldEqual, "C:\\")
			So(newPathResolver(0, "test").WorkingDir(), ShouldEqual, pwd+string(os.PathSeparator)+"test")
			//So(newPathResolver(0, "panicnoroot").WorkingDir(), ShouldEqual, pwd)

		})
		Convey("A pathResolver should not panic on getting pwd", func() {
			recovered := false
			defer func() {
				recover()
				recovered = true
				So(recovered, ShouldBeTrue)
			}()
			_ = newPathResolver(0, "panic on os.Getwd")
		})
		Convey("A pathResolver should not panic on filepath.Abs", func() {
			recovered := false
			defer func() {
				recover()
				recovered = true
				So(recovered, ShouldBeTrue)
			}()
			_ = newPathResolver(0, "panic on filepath.Abs")
		})
	})
}
