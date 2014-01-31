package asciidocgo

import (
	"os"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

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
			So(len(pathSegments), ShouldEqual, 3)
			So(root, ShouldEqual, "")
			So(posixPath, ShouldEqual, "/a/b/./c")
		})
		Convey("A Partition keep '..' paths", func() {
			pathSegments, root, posixPath = PartitionPath("a\\b/../c", true)
			So(len(pathSegments), ShouldEqual, 4)
			So(root, ShouldEqual, "")
			So(posixPath, ShouldEqual, "a/b/../c")

			pathSegments, root, posixPath = PartitionPath("\\a\\b/../c", true)
			So(len(pathSegments), ShouldEqual, 4)
			So(root, ShouldEqual, "")
			So(posixPath, ShouldEqual, "/a/b/../c")

			pathSegments, root, posixPath = PartitionPath("c:\\a\\b/../c", false)
			So(len(pathSegments), ShouldEqual, 4)
			So(root, ShouldEqual, "c:")
			So(posixPath, ShouldEqual, "c:/a/b/../c")

			pathSegments, root, posixPath = PartitionPath("a/b", false)
			So(len(pathSegments), ShouldEqual, 2)
			So(root, ShouldEqual, "")
			So(posixPath, ShouldEqual, "a/b")

			testpr = "test_PartitionPath_segments"
			pathSegments, root, posixPath = PartitionPath("/../images", true)
			So(posixPath, ShouldEqual, "pathSegments=(3)'[ .. images]'- root=''(true), posixPath='/../images'")
			testpr = ""

			testpr = "test_PartitionPath_rootSegments"
			pathSegments, root, posixPath = PartitionPath("/../images", true)
			So(posixPath, ShouldEqual, "pathSegments=(2)'[.. images]'- root=''(true), posixPath='/../images'")
			testpr = ""

		})

		Convey("A Windows root path should panic if partitioned as web", func() {
			recovered := false
			defer func() {
				r := recover()
				recovered = true
				So(recovered, ShouldBeTrue)
				So(r, ShouldEqual, "path 'C:\\a/b' is a root path, but not a web root path")
			}()
			pathSegments, root, posixPath = PartitionPath("C:\\a/b", true)
		})
		Convey("A Web root path should panic if partitioned as windows", func() {
			recovered := false
			defer func() {
				r := recover()
				recovered = true
				So(recovered, ShouldBeTrue)
				So(r, ShouldEqual, "path '/a/b' is a root path, but not a windows root path")
			}()
			pathSegments, root, posixPath = PartitionPath("/a/b", false)
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
			So(JoinPath([]string{}, "e:"), ShouldEqual, "e:")
			So(JoinPath(nil, "e:"), ShouldEqual, "e:")
		})
	})

	Convey("A Partition can resolve a system path from the target and start paths (internal tests)", t, func() {
		testpr = ""
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
		Convey("Empty target segment, non-empty non-root start means system path start", func() {
			So(pr.SystemPath("", "start/../b", "", false, ""), ShouldEqual, "C:/a/working/dir/b")
			So(pr.SystemPath("", "start/../b", "C:\\", false, ""), ShouldEqual, "C:/b")
			So(pr.SystemPath("start/../b", "C:\\", "C:\\", false, ""), ShouldEqual, "C:/b")
		})
		Convey("Non-Empty target segments starting with jail (or empty jail) returns target", func() {
			So(pr.SystemPath("C:/start/b", "", "", false, ""), ShouldEqual, "C:/start/b")
			So(pr.SystemPath("C:/start/b", "C:\\start", "", false, ""), ShouldEqual, "C:/start/b")
			So(pr.SystemPath("C:/start/b", "C:\\start/", "", false, ""), ShouldEqual, "C:/start/b")
		})

		Convey("Empty start and jail means start is working dir", func() {
			testpr = "test_SystemPath_start"
			So(pr.SystemPath("a/b1", "", "", false, ""), ShouldEqual, pr.WorkingDir())
		})
		Convey("Empty start and non-empty jail means start is jail", func() {
			testpr = "test_SystemPath_start"
			So(pr.SystemPath("a/b1", "", "C:/c/d", false, ""), ShouldEqual, "C:/c/d")
		})

		Convey("Non-Empty root start means posixfied start", func() {
			testpr = "test_SystemPath_start"
			So(pr.SystemPath("a/b1", "C:\\a/b\\c", "C:/c/d", false, ""), ShouldEqual, "C:/a/b/c")
		})

		Convey("Non Empty target segment, non-empty non-root start means system path start with jail", func() {
			So(pr.SystemPath("a/b2", "start/../b", "C:\\", false, ""), ShouldEqual, "C:/b/a/b2")
		})

		Convey("Same jail and start means posixfied start", func() {
			testpr = "test_SystemPath_segments"
			So(pr.SystemPath("a/b1", "C:\\a/b\\c", "C:\\a/b/c", false, ""), ShouldEqual, "jail='C:/a/b/c', jailRoot='C:', jailSegments '[a b c]', startRoot='', startSegments '[a b c]'")
		})

		Convey("Different jail and start means panic if start doesn't include jail", func() {
			recovered := false
			defer func() {
				r := recover()
				recovered = true
				So(recovered, ShouldBeTrue)
				So(r, ShouldEqual, "Start path 'C:/a/b/c' is outside of jail: 'C:/e/b/c' (disallowed in safe mode)")
			}()
			_ = pr.SystemPath("a/b1", "C:\\a/b\\c", "C:\\e/b/c", false, "")
		})

		Convey("Start must includes jail", func() {
			testpr = "test_SystemPath_segments"
			So(pr.SystemPath("a/b1", "C:\\a/b\\c/e/f", "C:\\a/b/c", false, ""), ShouldEqual, "jail='C:/a/b/c', jailRoot='C:', jailSegments '[a b c]', startRoot='C:', startSegments '[a b c e f]'")
		})

		Convey("Start with empty jail", func() {
			testpr = "test_SystemPath_segments"
			So(pr.SystemPath("a/b1", "C:\\a/b\\c/e/f", "", false, ""), ShouldEqual, "jail='', jailRoot='C:', jailSegments '[]', startRoot='C:', startSegments '[a b c e f]'")
		})
	})

	Convey("A Partition can resolve a system path from the target and start paths (Unit Tests)", t, func() {
		testpr = ""
		pr := NewPathResolver(0, "C:/a/working/dir")

		Convey("Simple non-root target is append to current working dir", func() {
			// resolver.system_path('images')
			// => '/path/to/docs/images'
			So(pr.SystemPath("images", "", "", false, ""), ShouldEqual, "C:/a/working/dir/images")
		})

		Convey("dot-dot target current working dir back one folder up", func() {
			// resolver.system_path('../images')
			// => '/path/to/images'
			So(pr.SystemPath("../images", "", "", false, ""), ShouldEqual, "C:/a/working/images")
		})

		Convey("dot-dot target current working dir back one folder up", func() {
			// resolver.system_path('/etc/images')
			// => '/etc/images'
			So(pr.SystemPath("C:/etc/images", "", "", false, ""), ShouldEqual, "C:/etc/images")
		})

		Convey("non-empty target is appended to non-empty start", func() {
			// resolver.system_path('images', '/etc')
			// => '/etc/images'
			So(pr.SystemPath("images", "C:/etc", "", false, ""), ShouldEqual, "C:/etc/images")
		})

		Convey("empty target returns non-empty start", func() {
			// resolver.system_path('', '/etc/images')
			// => '/etc/images'
			So(pr.SystemPath("", "C:/etc/images", "", false, ""), ShouldEqual, "C:/etc/images")
		})

		Convey("empty target and empty start returns jail", func() {
			// resolver.system_path(nil, nil, '/path/to/docs')
			// => '/path/to/docs'
			So(pr.SystemPath("", "", "C:/etc/images", false, ""), ShouldEqual, "C:/etc/images")
		})

		Convey("dot_dot target and empty start, returns non-empty jail", func() {
			// resolver.system_path('..', nil, '/path/to/docs')
			// => '/path/to/docs'
			So(pr.SystemPath("..", "", "C:/etc/images", true, ""), ShouldEqual, "C:/etc/images")
		})

		Convey("dot_dot path target and empty start, returns non-empty jail plus path", func() {
			// resolver.system_path('../../../css', nil, '/path/to/docs')
			// => '/path/to/docs/css'
			So(pr.SystemPath("../../../css", "", "C:/etc/to/images", true, ""), ShouldEqual, "C:/etc/to/images/css")
		})

		Convey("dot_dot path target and empty start, returns non-empty jail plus path", func() {
			// resolver.system_path('../../../css', '../../..', '/path/to/docs')
			// => '/path/to/docs/css'
			So(pr.SystemPath("../../../css", "", "C:/etc/to/images", true, ""), ShouldEqual, "C:/etc/to/images/css")
		})

		Convey("dot_dot path target, different start and jail returns jail", func() {
			// resolver.system_path('..', 'C:\\data\\docs\\assets', 'C:\\data\\docs')
			//=> 'C:/data/docs'
			So(pr.SystemPath("..", "C:\\data\\docs\\assets", "C:\\data\\docs", true, ""), ShouldEqual, "C:/data/docs")
		})

		Convey("dot_dot path target, start including jail returns start+target", func() {
			// resolver.system_path('..\\..\\css', 'C:\\data\\docs\\assets', 'C:\\data\\docs')
			// => 'C:/data/docs/css'
			So(pr.SystemPath("..\\..\\css", "C:\\data\\docs\\assets", "C:\\data\\docs", true, ""), ShouldEqual, "C:/data/docs/css")
		})

		Convey("Different jail and start means panic if start doesn't include jail", func() {
			/*
					begin
				     resolver.system_path('../../../css', '../../..', '/path/to/docs', :recover => false)
						rescue SecurityError => e
						puts e.message
					end
					=> 'path ../../../../../../css refers to location outside jail: /path/to/docs (disallowed in safe mode)'

			*/
			recovered := false
			defer func() {
				r := recover()
				recovered = true
				So(recovered, ShouldBeTrue)
				So(r, ShouldEqual, "path '../../../css' refers to location outside jail: 'C:/path/to/docs' (disallowed in safe mode)")
			}()
			_ = pr.SystemPath("../../../css", "../../..", "C:\\path/to/docs", false, "")
		})

		Convey("target, including jail but empty start returns target", func() {
			// resolver.system_path('/path/to/docs/images', nil, '/path/to/docs')
			//	=> '/path/to/docs/images'
			So(pr.SystemPath("C:/path/to/docs/images", "", "C:/path/to/docs", false, ""), ShouldEqual, "C:/path/to/docs/images")
		})

		Convey("start outside of jail means panic", func() {
			/*

				begin
				  resolver.system_path('images', '/etc', '/path/to/docs')
				  rescue SecurityError => e
				    puts e.message
				  end
				=> Start path /etc is outside of jail: /path/to/docs'
			*/
			recovered := false
			defer func() {
				r := recover()
				recovered = true
				So(recovered, ShouldBeTrue)
				So(r, ShouldEqual, "Start path 'C:/etc' is outside of jail: 'C:/path/to/docs' (disallowed in safe mode)")
			}()
			_ = pr.SystemPath("images", "C:/etc", "C:/path/to/docs", false, "")
		})
	})

	Convey("A PathResolver can compute a web path from the target and start paths (internal tests)", t, func() {
		testpr = ""

		Convey("Empty target and start returns empty web path", func() {
			So(WebPath("", ""), ShouldEqual, "")
		})

		Convey("target and start with http means non-empty uriPrefix", func() {
			testpr = "test_Webath_uriPrefix"
			So(WebPath("b/c", "http://a"), ShouldEqual, "target='a/b/c', uriPrefix='http://'")
			So(WebPath("/images", ""), ShouldEqual, "target='/images', uriPrefix=''")
			So(WebPath("/../images", ""), ShouldEqual, "target='/../images', uriPrefix=''")
			So(WebPath("images", "/assets"), ShouldEqual, "target='/assets/images', uriPrefix=''")
		})

		Convey("target and start with http means non-empty target segments", func() {
			testpr = "test_Webath_partitionTarget"
			So(WebPath("b/c", "http://a"), ShouldEqual, "targetSegments=(3)'[a b c]', targetRoot=''")
			So(WebPath("b/c", "/a//"), ShouldEqual, "targetSegments=(3)'[a b c]', targetRoot=''")
			So(WebPath("/b/c", "/a/"), ShouldEqual, "targetSegments=(2)'[b c]', targetRoot=''")
			So(WebPath("/images", ""), ShouldEqual, "targetSegments=(1)'[images]', targetRoot=''")
			So(WebPath("/../images", ""), ShouldEqual, "targetSegments=(2)'[.. images]', targetRoot=''")
			So(WebPath("images", "/assets"), ShouldEqual, "targetSegments=(2)'[assets images]', targetRoot=''")
		})
	})

	Convey("A PathResolver can compute a web path from the target and start paths (unit tests)", t, func() {
		testpr = ""

		Convey("Simple target and empty start returns simple target", func() {
			/* resolver.web_path('images')
			=> 'images'*/
			So(WebPath("images", ""), ShouldEqual, "images")
		})
		Convey("Simple dot target and empty start returns simple dot target", func() {
			/* resolver.web_path('./images')
			=> './images'*/
			So(WebPath("./images", ""), ShouldEqual, "./images")
		})

		Convey("Simple dot target and empty start returns simple dot target", func() {
			/* resolver.web_path('/images')
			=> '/images'*/
			So(WebPath("/images", ""), ShouldEqual, "/images")
		})

		Convey("Target with dot and dots and empty start returns resolved target", func() {
			/* resolver.web_path('./images/../assets/images')
			=> './assets/images'*/
			So(WebPath("./images/../assets/images", ""), ShouldEqual, "./assets/images")
		})

		Convey("Target with dots and empty start returns resolved target", func() {
			/* resolver.web_path('/../images')
			   => '/../images' (not /images/ as commented. BUG?) */
			So(WebPath("/../images", ""), ShouldEqual, "/../images")
			So(WebPath("/../../images", ""), ShouldEqual, "/../../images")
			So(WebPath("a/../../images", ""), ShouldEqual, "../images")
			So(WebPath("a/../b/../images", ""), ShouldEqual, "images")
			So(WebPath("a/../b/../../images", ""), ShouldEqual, "../images")
		})

		Convey("Target with start returns start plus target", func() {
			/* resolver.web_path('images', 'assets')
			   => 'assets/images' */
			So(WebPath("images", "assets"), ShouldEqual, "assets/images")
			So(WebPath("images", "/assets"), ShouldEqual, "/assets/images")
			So(WebPath("images", "http://assets"), ShouldEqual, "http://assets/images")
		})

		Convey("Target with start with dots returns start plus target", func() {
			/* resolver.web_path('tiger.png', '../assets/images')
			   => '../assets/images/tiger.png' */
			So(WebPath("tiger.png", "../assets/images"), ShouldEqual, "../assets/images/tiger.png")
		})
	})
	Convey("A PathResolver can compute relative path out of two absolute paths", t, func() {

		pr := NewPathResolver(0, "C:/a/working/dir")

		Convey("2 absolute paths, one begins with the other", func() {
			So(pr.RelativePath("C:\\a/b\\c/d", "C:\\a\\b\\"), ShouldEqual, "c/d")
			So(pr.RelativePath("C:\\z/b\\c/d", "C:\\a\\b\\"), ShouldEqual, "c/d")
			So(pr.RelativePath("C:\\z/b\\c/d", "a\\b\\"), ShouldEqual, "C:\\z/b\\c/d")
			//So(pr.RelativePath("C:\\z", "C:\\a\\b\\"), ShouldEqual, "c/d")
			//So(pr.RelativePath("", "C:\\a\\b\\"), ShouldEqual, "c/d")
		})
		testpr = ""
	})

}
