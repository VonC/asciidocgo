package asciidocgo

import (
	"os"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

var Test = "aaa"

func TestPathResolver(t *testing.T) {

	Convey("A pathResolver can be initialized", t, func() {

		Convey("By default, a pathResolver can be created", func() {
			So(NewPathResolver(0, ""), ShouldNotBeNil)
		})
		Convey("By default, a pathResolver has a system path separator", func() {
			So(NewPathResolver(0, "").FileSeparator(), ShouldEqual, os.PathSeparator)
			So(NewPathResolver('/', "").FileSeparator(), ShouldNotEqual, os.PathSeparator)
			So(NewPathResolver('/', "").FileSeparator(), ShouldEqual, '/')
		})

		Convey("By default, a pathResolver has a current working path", func() {
			pwd, err := os.Getwd()
			if err != nil {
				panic(err)
			}
			So(NewPathResolver(0, "").WorkingDir(), ShouldEqual, pwd)
			So(NewPathResolver(0, "C:\\").WorkingDir(), ShouldEqual, "C:\\")
			So(NewPathResolver(0, "test").WorkingDir(), ShouldEqual, pwd+string(os.PathSeparator)+"test")
			//So(NewPathResolver(0, "panicnoroot").WorkingDir(), ShouldEqual, pwd)

		})
		Convey("A pathResolver should not panic on getting pwd", func() {
			recovered := false
			defer func() {
				recover()
				recovered = true
				So(recovered, ShouldBeTrue)
			}()
			_ = NewPathResolver(0, "panic on os.Getwd")
		})
		Convey("A pathResolver should not panic on filepath.Abs", func() {
			recovered := false
			defer func() {
				recover()
				recovered = true
				So(recovered, ShouldBeTrue)
			}()
			_ = NewPathResolver(0, "panic on filepath.Abs")
		})
	})

	Convey("A pathResolver can test for a web path", t, func() {
		So(IsWebRoot(""), ShouldBeFalse)
		So(IsWebRoot("a"), ShouldBeFalse)
		So(IsWebRoot("\\a\\b/c"), ShouldBeFalse)
		So(IsWebRoot("/a/b/c"), ShouldBeTrue)
	})

	Convey("A pathResolver can replace backslash by slash", t, func() {
		So(Posixfy(""), ShouldEqual, "")
		So(Posixfy("a/b/c"), ShouldEqual, "a/b/c")
		So(Posixfy("a\\b\\c"), ShouldEqual, "a/b/c")
	})

	Convey("A pathResolver can test for root", t, func() {
		Convey("A Path starting with C:/ is root", func() {
			So(IsRoot(""), ShouldBeFalse)
			So(IsRoot("C:\\"), ShouldBeTrue)
			So(IsRoot("C:/"), ShouldBeTrue)
			So(IsRoot("C:\\a/b/"), ShouldBeTrue)
			So(IsRoot("c:/a/b/../c"), ShouldBeTrue)
			So(IsRoot("c:\\a/b/../c"), ShouldBeTrue)
		})
	})

	Convey("A pathResolver can test for web root", t, func() {
		Convey("A Path starting with / is web root", func() {
			So(IsWebRoot(""), ShouldBeFalse)
			So(IsWebRoot("C:\\"), ShouldBeFalse)
			So(IsWebRoot("\\"), ShouldBeFalse)
			So(IsWebRoot("/"), ShouldBeTrue)
			So(IsWebRoot("/a/b/"), ShouldBeTrue)
			So(IsWebRoot("/a\\b/./c"), ShouldBeTrue)
			So(IsWebRoot("/a/b/./c"), ShouldBeTrue)
		})
	})

	Convey("A pathResolver can expand a path", t, func() {
		Convey("empty path returns an empty string", func() {
			So(ExpandPath(""), ShouldEqual, "")
		})
		Convey("non-empty path returns an posix path", func() {
			So(ExpandPath("c:\\a/.\\b/../c"), ShouldEqual, "c:/a/b/../c")
		})
	})

	Convey("A pathResolver can partition a path", t, func() {
		pathSegments, root, posixPath := PartitionPath("", false)
		So(len(pathSegments), ShouldEqual, 0)
		So(root, ShouldEqual, "")
		So(posixPath, ShouldEqual, "")

		Convey("A Path starting with dot has a dot root", func() {
			pathSegments, root, posixPath := PartitionPath(".", false)
			So(len(pathSegments), ShouldEqual, 0)
			So(root, ShouldEqual, ".")
			So(posixPath, ShouldEqual, ".")

			pathSegments, root, posixPath = PartitionPath(".\\a/b", false)
			So(len(pathSegments), ShouldEqual, 2)
			So(root, ShouldEqual, ".")
			So(posixPath, ShouldEqual, "./a/b")

		})
		Convey("A Partition removes self-reference path", func() {
			pathSegments, root, posixPath := PartitionPath("a\\b/./c", false)
			So(len(pathSegments), ShouldEqual, 3)
			So(root, ShouldEqual, "")
			So(posixPath, ShouldEqual, "a/b/./c")

			pathSegments, root, posixPath = PartitionPath("C:/a\\b/./c", false)
			So(len(pathSegments), ShouldEqual, 3)
			So(root, ShouldEqual, "C:")
			So(posixPath, ShouldEqual, "C:/a/b/./c")

			pathSegments, root, posixPath = PartitionPath("/a\\b/./c", true)
			So(len(pathSegments), ShouldEqual, 2)
			So(root, ShouldEqual, "/a")
			So(posixPath, ShouldEqual, "/a/b/./c")
		})
		Convey("A Partition keep '..' paths", func() {
			pathSegments, root, posixPath = PartitionPath("a\\b/../c", true)
			So(len(pathSegments), ShouldEqual, 4)
			So(root, ShouldEqual, "")
			So(posixPath, ShouldEqual, "a/b/../c")

			pathSegments, root, posixPath = PartitionPath("\\a\\b/../c", true)
			So(len(pathSegments), ShouldEqual, 3)
			So(root, ShouldEqual, "/a")
			So(posixPath, ShouldEqual, "/a/b/../c")

			pathSegments, root, posixPath = PartitionPath("c:\\a\\b/../c", false)
			So(len(pathSegments), ShouldEqual, 4)
			So(root, ShouldEqual, "c:")
			So(posixPath, ShouldEqual, "c:/a/b/../c")

			pathSegments, root, posixPath = PartitionPath("a/b", false)
			So(len(pathSegments), ShouldEqual, 2)
			So(root, ShouldEqual, "")
			So(posixPath, ShouldEqual, "a/b")
		})
	})

	Convey("A pathResolver can join a path", t, func() {
		Convey("No segment and empty root returns an empty string", func() {
			So(JoinPath(nil, ""), ShouldEqual, "")
		})
		Convey("Segments with no root returns an slash-separated segments", func() {
			So(JoinPath([]string{"a"}, ""), ShouldEqual, "a")
			So(JoinPath([]string{"a", "b"}, ""), ShouldEqual, "a/b")
			So(JoinPath([]string{"a", "b", "c"}, ""), ShouldEqual, "a/b/c")
		})
		Convey("Segments with root returns an root plus slash-separated segments", func() {
			So(JoinPath([]string{"a"}, "c:"), ShouldEqual, "c:/a")
			So(JoinPath([]string{"a", "b"}, "d:"), ShouldEqual, "d:/a/b")
			So(JoinPath([]string{"a", "b", "c"}, "e:"), ShouldEqual, "e:/a/b/c")
		})
	})

	Convey("A Partition can resolve a system path from the target and start paths (internal tests)", t, func() {
		Test = ""
		pr := NewPathResolver(0, "C:/a/working/dir")
		Convey("A Non-absolute jail path should panic", func() {
			recovered := false
			defer func() {
				r := recover()
				recovered = true
				So(recovered, ShouldBeTrue)
				So(r, ShouldEqual, "Jail is not an absolute path: c")
			}()
			_ = pr.SystemPath("a", "b", "c", false, "")
		})
		/*
			Convey("A system path with no start resolves from the root", func() {
				So(pr.SystemPath("images", "", "", false, ""), ShouldEqual, "")
				So(pr.SystemPath("../images", "", "", false, ""), ShouldEqual, "")
				So(pr.SystemPath("/etc/images", "", "", false, ""), ShouldEqual, "")
			})*/
		Convey("Empty target segment and empty start and empty jail means working dir", func() {
			So(pr.SystemPath("", "", "", false, ""), ShouldEqual, "C:/a/working/dir")
		})
		Convey("Empty target segment, non-empty root start and empty jail means expanded start path", func() {
			So(pr.SystemPath("", "C:\\start/../b", "", false, ""), ShouldEqual, "C:/start/../b")
		})
		Convey("Empty target segment, non-empty non-root start means susyem path start", func() {
			recovered := false
			defer func() {
				r := recover()
				recovered = true
				So(recovered, ShouldBeTrue)
				So(r, ShouldEqual, "should not happen yet")
			}()
			SkipSo(pr.SystemPath("", "start/../b", "", false, ""), ShouldEqual, "start/../b")
			SkipSo(pr.SystemPath("", "start/../b", "C:\\", false, ""), ShouldEqual, "start/../b")
		})
		Convey("Non-Empty target segments starting with jail (or empty jail) returns target", func() {
			So(pr.SystemPath("C:/start/b", "", "", false, ""), ShouldEqual, "C:/start/b")
			So(pr.SystemPath("C:/start/b", "C:\\start", "", false, ""), ShouldEqual, "C:/start/b")
			So(pr.SystemPath("C:/start/b", "C:\\start/", "", false, ""), ShouldEqual, "C:/start/b")
		})

		Convey("Empty start and jail means start is working dir", func() {
			Test = "test_SystemPath_start"
			So(pr.SystemPath("a/b1", "", "", false, ""), ShouldEqual, pr.WorkingDir())
		})
		Convey("Empty start and non-empty jail means start is jail", func() {
			Test = "test_SystemPath_start"
			So(pr.SystemPath("a/b1", "", "C:/c/d", false, ""), ShouldEqual, "C:/c/d")
		})

		Convey("Non-Empty root start means posixfied start", func() {
			Test = "test_SystemPath_start"
			So(pr.SystemPath("a/b1", "C:\\a/b\\c", "C:/c/d", false, ""), ShouldEqual, "C:/a/b/c")
		})

		Convey("Non Empty target segment, non-empty non-root start means sustem path start with jail", func() {
			recovered := false
			defer func() {
				r := recover()
				recovered = true
				So(recovered, ShouldBeTrue)
				So(r, ShouldEqual, "should not happen yet (start)")
			}()
			SkipSo(pr.SystemPath("a/b2", "start/../b", "C:\\", false, ""), ShouldEqual, "start/../b")
		})

		Convey("Same jail and start means posixfied start", func() {
			Test = "test_SystemPath_segments"
			So(pr.SystemPath("a/b1", "C:\\a/b\\c", "C:\\a/b/c", false, ""), ShouldEqual, "jail='C:/a/b/c', jailRoot='C:', jailSegments '[a b c]', startRoot='', startSegments '[a b c]'")
		})

		Convey("Different jail and start means panic if start doesn't include jail", func() {
			recovered := false
			defer func() {
				r := recover()
				recovered = true
				So(recovered, ShouldBeTrue)
				So(r, ShouldEqual, "start 'C:/a/b/c' is outside of jail: 'C:/e/b/c' (disallowed in safe mode)")
			}()
			_ = pr.SystemPath("a/b1", "C:\\a/b\\c", "C:\\e/b/c", false, "")
		})

		Convey("Start must includes jail", func() {
			Test = "test_SystemPath_segments"
			So(pr.SystemPath("a/b1", "C:\\a/b\\c/e/f", "C:\\a/b/c", false, ""), ShouldEqual, "jail='C:/a/b/c', jailRoot='C:', jailSegments '[a b c]', startRoot='C:', startSegments '[a b c e f]'")
		})

		Convey("Start with empty jail", func() {
			Test = "test_SystemPath_segments"
			So(pr.SystemPath("a/b1", "C:\\a/b\\c/e/f", "", false, ""), ShouldEqual, "jail='', jailRoot='C:', jailSegments '[]', startRoot='C:', startSegments '[a b c e f]'")
		})
	})

	Convey("A Partition can resolve a system path from the target and start paths (Unit Tests)", t, func() {
		Test = ""
		pr := NewPathResolver(0, "C:/a/working/dir")

		Convey("Simple non-root target is append to current working dir", func() {
			// resolver.system_path('images')
			// => '/path/to/docs/images'
			So(pr.SystemPath("images", "", "", false, ""), ShouldEqual, "C:/a/working/dir/images")
		})
	})
	/*

	   resolver.system_path('../images')
	   => '/path/to/images'

	   resolver.system_path('/etc/images')
	   => '/etc/images'

	   resolver.system_path('images', '/etc')
	   => '/etc/images'

	   resolver.system_path('', '/etc/images')
	   => '/etc/images'

	   resolver.system_path(nil, nil, '/path/to/docs')
	   => '/path/to/docs'

	   resolver.system_path('..', nil, '/path/to/docs')
	   => '/path/to/docs'

	   resolver.system_path('../../../css', nil, '/path/to/docs')
	   => '/path/to/docs/css'

	   resolver.system_path('../../../css', '../../..', '/path/to/docs')
	   => '/path/to/docs/css'

	   resolver.system_path('..', 'C:\\data\\docs\\assets', 'C:\\data\\docs')
	   => 'C:/data/docs'

	   resolver.system_path('..\\..\\css', 'C:\\data\\docs\\assets', 'C:\\data\\docs')
	   => 'C:/data/docs/css'

	   begin
	     resolver.system_path('../../../css', '../../..', '/path/to/docs', :recover => false)
	   rescue SecurityError => e
	     puts e.message
	   end
	   => 'path ../../../../../../css refers to location outside jail: /path/to/docs (disallowed in safe mode)'

	   resolver.system_path('/path/to/docs/images', nil, '/path/to/docs')
	   => '/path/to/docs/images'

	   begin
	     resolver.system_path('images', '/etc', '/path/to/docs')
	   rescue SecurityError => e
	     puts e.message
	   end
	   => Start path /etc is outside of jail: /path/to/docs'
	*/
}
