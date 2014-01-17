package asciidocgo

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

/*
Handles all operations for resolving, cleaning and joining paths.
This class includes operations for handling both web paths (request URIs) and
system paths.

The main emphasis of the class is on creating clean and secure paths. Clean
paths are void of duplicate parent and current directory references in the
path name. Secure paths are paths which are restricted from accessing
directories outside of a jail root, if specified.

Since joining two paths can result in an insecure path, this class also
handles the task of joining a parent (start) and child (target) path.

This class makes no use of path utilities from the Ruby libraries. Instead,
it handles all aspects of path manipulation. The main benefit of
internalizing these operations is that the class is able to handle both posix
and windows paths independent of the operating system on which it runs. This
makes the class both deterministic and easier to test.

Examples:

    resolver = PathResolver.new

    Web Paths

    resolver.web_path('images')
    => 'images'

    resolver.web_path('./images')
    => './images'

    resolver.web_path('/images')
    => '/images'

    resolver.web_path('./images/../assets/images')
    => './assets/images'

    resolver.web_path('/../images')
    => '/images'

    resolver.web_path('images', 'assets')
    => 'assets/images'

    resolver.web_path('tiger.png', '../assets/images')
    => '../assets/images/tiger.png'

    System Paths

    resolver.working_dir
    => '/path/to/docs'

    resolver.system_path('images')
    => '/path/to/docs/images'

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
type PathResolver struct {
	fileSeparator byte
	workingDir    string
}

func (pr *PathResolver) FileSeparator() byte {
	return pr.fileSeparator
}
func (pr *PathResolver) WorkingDir() string {
	return pr.workingDir
}

/* Construct a new instance of PathResolver, optionally specifying
the file separator (to override the system default) and
the working directory (to override the present working directory).
The working directory will be expanded to an absolute path inside the constructor.
file_separator - the String file separator to use for path operations
(optional, default: File::FILE_SEPARATOR)
working_dir    - the String working directory (optional, default: Dir.pwd)
*/
func NewPathResolver(fileSeparator byte, workingDir string) *PathResolver {
	if fileSeparator == 0 {
		fileSeparator = os.PathSeparator
	}
	if workingDir == "" || workingDir == "panic on os.Getwd" {
		pwd, err := os.Getwd()
		if err != nil || workingDir == "panic on os.Getwd" {
			if workingDir == "panic on os.Getwd" {
				err = errors.New("test on bad os.Getwd")
			}
			panic(err)
		}
		workingDir = pwd
	} else {
		if IsRoot(workingDir) == false {
			wd, err := filepath.Abs(workingDir)
			if err != nil || workingDir == "panic on filepath.Abs" {
				if workingDir == "panic on filepath.Abs" {
					err = errors.New("test on bad filepath.Abs")
				}
				panic(err)
			}
			workingDir = wd
		}
	}
	return &PathResolver{fileSeparator, workingDir}
}

/*Check if the specified path is an absolute root path
his operation correctly handles both posix and windows paths.
returns a Boolean indicating whether the path is an absolute root path
*/
func IsRoot(apath string) bool {
	return filepath.IsAbs(apath)
}

/*Determine if the path is an absolute (root) web path.
Returns a Boolean indicating whether the path is an absolute (root) web path*/
func IsWebRoot(apath string) bool {
	return path.IsAbs(apath)
}

/*Normalize path by converting any backslashes to forward slashes
Returns a String path with any backslashes replaced with
forward slashes*/
func Posixfy(path string) string {
	return strings.Replace(path, "\\", "/", -1)
}

/* Expand the path by resolving any parent references (..)
and cleaning self references (.).
The result will be relative if the path is relative and
absolute if the path is absolute.
The file separator used  in the expanded path is the one specified
when the class was constructed.
path - the String path to expand
returns a String path with any parent or self references resolved.
*/
func ExpandPath(path string) string {
	pathSegments, pathRoot, _ := PartitionPath(path, false)
	return JoinPath(pathSegments, pathRoot)
}

/*Partition the path into path segments and remove any empty segments
or segments that are self references (.).
The path is split on either posix or windows file separators.
returns a 3-item Array containing the Array of String path segments,
the path root, if the path is absolute, and the posix version of the path.
*/
func PartitionPath(path string, webPath bool) (pathSegments []string, root string, posixPath string) {
	posixPath = Posixfy(path)
	isRoot := false
	if webPath {
		isRoot = IsWebRoot(posixPath)
	} else {
		isRoot = IsRoot(posixPath)
	}
	reg, _ := regexp.Compile("/+")
	posixPathCleaned := reg.ReplaceAllString(posixPath, "/")
	pathSegmentsWithDots := strings.Split(posixPathCleaned, "/")
	if pathSegmentsWithDots[0] == "." {
		root = "."
	} else {
		root = ""
	}
	pathSegments = []string{}
	if len(pathSegmentsWithDots) > 1 || pathSegmentsWithDots[0] != "" {
		for k := 0; k < len(pathSegmentsWithDots); k++ {
			if pathSegmentsWithDots[k] != "." {
				pathSegments = append(pathSegments, pathSegmentsWithDots[k])
			}
		}
	}
	if isRoot {
		root, pathSegments = pathSegments[0], pathSegments[1:len(pathSegments)]
		if root == "" {
			root, pathSegments = "/"+pathSegments[0], pathSegments[1:len(pathSegments)]
		}
	}
	return pathSegments, root, posixPath
}

/* Join the segments using the posix file separator.
Use the root, if specified, to construct an absolute path.
Otherwise join the segments as a relative path.
segments - a String Array of path segments
root     - a String path root (optional, default: nil)
returns a String path formed by joining the segments
using the posix file separator and prepending the root, if specified.
*/
func JoinPath(segments []string, root string) string {
	res := ""
	if segments != nil {
		for i := 0; i < len(segments); i++ {
			res = res + segments[i]
			if i < len(segments)-1 {
				res = res + "/"
			}
		}
	}
	if root != "" {
		res = root + "/" + res
	}
	return res
}

/*
Resolve a system path from the target and start paths.
If a jail path is specified, enforce that the resolved directory
is contained within the jail path.
If a jail path is not provided, the resolved path may be any location
on the system.
If the resolved path is absolute, use it as is.
If the resolved path is relative, resolve it relative to the working_dir
specified in the constructor.
target - the String target path
start  - the String start (i.e., parent) path
jail   - the String jail path to confine the resolved path
opts   - an optional Hash of options to control processing (default: {}):
  * :recover is used to control whether the processor should auto-recover
    when an illegal path is encountered
  * :target_name is used in messages to refer to the path being resolved
returns a String path that joins the target path with the start path with
any parent references resolved and self references removed and enforces
that the resolved path be contained within the jail, if provided
*/
func (pr *PathResolver) SystemPath(target, start, jail string, recover bool, targetName string) string {
	if jail != "" && !IsRoot(jail) {
		panic(fmt.Sprintf("Jail is not an absolute path: %v", jail))
	}
	jail = Posixfy(jail)
	targetSegments, targetRoot, _ := PartitionPath(target, false)
	if len(targetSegments) == 0 {
		if target == "a/b1" {
			panic(fmt.Sprintf("should not happen yet %v => %v", targetSegments, targetRoot))
		}
		if start == "" {
			if jail == "" {
				return pr.WorkingDir()
			}
		} else if IsRoot(start) {
			if jail == "" {
				return ExpandPath(start)
			}
		} else {
			// TODO return system_path(start, jail, jail)
			panic("should not happen yet")
		}
	}

	if targetRoot != "" && targetRoot != "." {
		resolvedTarget := JoinPath(targetSegments, targetRoot)
		// if target is absolute and a sub-directory of jail, or
		// a jail is not in place, let it slide
		if target == "a/b1" {
			panic(fmt.Sprintf("should not happen yet %v => %v", resolvedTarget, targetRoot))
		}
		if jail == "" || strings.HasPrefix(resolvedTarget, jail) {
			return resolvedTarget
		}
	}

	if start == "" {
		if jail == "" {
			start = pr.WorkingDir()
		} else {
			start = jail
		}
	} else if IsRoot(start) {
		start = Posixfy(start)
	} else {
		// TODO start = system_path(start, jail, jail)
		panic("should not happen yet (start)")
	}
	if Test == "test_SystemPath_start" {
		return start
	}

	jailSegments := []string{}
	jailRoot := ""
	startSegments := []string{}
	// both jail and start have been posixfied at this point
	if jail == start {
		jailSegments, jailRoot, _ = PartitionPath(jail, false)
		startSegments = make([]string, len(jailSegments))
		copy(startSegments, jailSegments)
	} else if jail != "" {
		if !strings.HasPrefix(start, jail) {
			panic(fmt.Sprintf("start '%v' is outside of jail: '%v' (disallowed in safe mode)", start, jail))
		}
	}

	if Test == "test_SystemPath_segments" {
		return fmt.Sprintf("jail='%v', jailRoot='%v', jailSegments '%v', startSegments '%v'", jail, jailRoot, jailSegments, startSegments)
	}

	return ""
}
