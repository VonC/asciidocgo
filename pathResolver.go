package asciidocgo

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/VonC/asciidocgo/consts/regexps"
)

var testpr = ""

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
	}
	wd, err := filepath.Abs(workingDir)
	if err != nil || workingDir == "panic on filepath.Abs" {
		if workingDir == "panic on filepath.Abs" {
			err = errors.New("test on bad filepath.Abs")
		}
		panic(err)
	}
	workingDir = wd
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
		if !isRoot && IsRoot(posixPath) {
			panic(fmt.Sprintf("path '%v' is a root path, but not a web root path", path))
		}
	} else {
		isRoot = IsRoot(posixPath)
		if !isRoot && IsWebRoot(posixPath) {
			panic(fmt.Sprintf("path '%v' is a root path, but not a windows root path", path))
		}
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
			if pathSegmentsWithDots[k] != "." && (pathSegmentsWithDots[k] != "" || k == 0) {
				pathSegments = append(pathSegments, pathSegmentsWithDots[k])
			}
		}
	}
	if testpr == "test_PartitionPath_segments" {
		return pathSegments, root, fmt.Sprintf("pathSegments=(%v)'%v'- root='%v'(%v), posixPath='%v'", len(pathSegments), pathSegments, root, isRoot, posixPath)
	}
	if isRoot {
		root, pathSegments = pathSegments[0], pathSegments[1:len(pathSegments)]
	}
	if testpr == "test_PartitionPath_rootSegments" {
		return pathSegments, root, fmt.Sprintf("pathSegments=(%v)'%v'- root='%v'(%v), posixPath='%v'", len(pathSegments), pathSegments, root, isRoot, posixPath)
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
		if len(segments) > 0 {
			res = root + "/" + res
		} else {
			res = root
		}
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
func (pr *PathResolver) SystemPath(target, start, jail string, canrecover bool, targetName string) string {
	if jail != "" && !IsRoot(jail) {
		panic(fmt.Sprintf("Jail is not an absolute path: %v", jail))
	}
	jail = Posixfy(jail)
	targetSegments, targetRoot, _ := PartitionPath(target, false)
	if len(targetSegments) == 0 {
		if start == "" {
			if jail == "" {
				return Posixfy(pr.WorkingDir())
			}
		} else if IsRoot(start) {
			if jail == "" {
				return ExpandPath(start)
			}
		} else {
			return pr.SystemPath(start, jail, jail, canrecover, targetName)
		}
	}

	if targetRoot != "" && targetRoot != "." {
		resolvedTarget := JoinPath(targetSegments, targetRoot)
		// if target is absolute and a sub-directory of jail, or
		// a jail is not in place, let it slide
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
		start = pr.SystemPath(start, jail, jail, true, targetName)
	}
	if testpr == "test_SystemPath_start" {
		return start
	}

	jailSegments := []string{}
	jailRoot := ""
	startSegments := []string{}
	startRoot := ""
	// both jail and start have been posixfied at this point
	if jail == start {
		jailSegments, jailRoot, _ = PartitionPath(jail, false)
		startSegments = make([]string, len(jailSegments))
		copy(startSegments, jailSegments)
	} else if jail != "" {
		if !strings.HasPrefix(start, jail) {
			aTargetName := targetName
			if targetName == "" {
				aTargetName = "Start path"
			}
			panic(fmt.Sprintf("%v '%v' is outside of jail: '%v' (disallowed in safe mode)", aTargetName, start, jail))
		}
		startSegments, startRoot, _ = PartitionPath(start, false)
		jailSegments, jailRoot, _ = PartitionPath(jail, false)
	} else {
		startSegments, startRoot, _ = PartitionPath(start, false)
		jailRoot = startRoot
	}

	if testpr == "test_SystemPath_segments" {
		return fmt.Sprintf("jail='%v', jailRoot='%v', jailSegments '%v', startRoot='%v', startSegments '%v'", jail, jailRoot, jailSegments, startRoot, startSegments)
	}

	resolvedSegments := make([]string, len(startSegments))
	copy(resolvedSegments, startSegments)
	warned := false

	for _, segment := range targetSegments {
		if segment == ".." {
			lr := len(resolvedSegments)
			if jail != "" {
				aTargetName := targetName
				if targetName == "" {
					aTargetName = "path"
				}
				if lr > len(jailSegments) {
					resolvedSegments = resolvedSegments[:lr-1]
				} else if !canrecover {
					panic(fmt.Sprintf("%v '%v' refers to location outside jail: '%v' (disallowed in safe mode)", aTargetName, target, jail))
				} else if !warned {
					fmt.Errorf("asciidoctor: WARNING: %v '%v' has illegal reference to ancestor of jail, auto-recovering", aTargetName, target)
					warned = true
				}
			} else {
				resolvedSegments = resolvedSegments[:lr-1]
			}
		} else {
			resolvedSegments = append(resolvedSegments, segment)
		}
	}

	return JoinPath(resolvedSegments, jailRoot)
}

/* Resolve a web path from the target and start paths.
The main function of this operation is to resolve any parent references
and remove any self references.
start  - the String start (i.e., parent) path
returns a String path that joins the target path with
the start path with any parent references resolved
and self references removed
*/
func WebPath(target, start string) string {
	target = Posixfy(target)
	start = Posixfy(start)
	uriPrefix := ""
	isWebroot := IsWebRoot(target)

	if !isWebroot && start != "" {
		target = start + "/" + target
		if strings.Contains(target, ":") {
			if res := regexps.UriSniffRx.FindStringSubmatchIndex(target); len(res) == 4 {
				uriPrefix = target[:res[3]]
				target = target[res[3]:]
			}
		}
	}
	// BUG? slash seems to be lost in https://github.com/asciidoctor/asciidoctor/blob/ab1e0b9c45e5138394b089dac205fb6d854e15e6/lib/asciidoctor/path_resolver.rb#L352-L366
	isWebroot = IsWebRoot(target)
	if testpr == "test_Webath_uriPrefix" {
		return fmt.Sprintf("target='%v', uriPrefix='%v'", target, uriPrefix)
	}
	targetSegments, targetRoot, _ := PartitionPath(target, true)
	if testpr == "test_Webath_partitionTarget" {
		return fmt.Sprintf("targetSegments=(%v)'%v', targetRoot='%v'", len(targetSegments), targetSegments, targetRoot)
	}
	accum := []string{}
	for _, segment := range targetSegments {
		if segment == ".." {
			if len(accum) == 0 {
				if targetRoot == "" || targetRoot == "." {
					accum = append(accum, segment)
				}
			} else if accum[len(accum)-1] == ".." {
				accum = append(accum, segment)
			} else {
				accum = accum[:len(accum)-1]
			}
		} else {
			accum = append(accum, segment)
		}
	}
	resolvedSegments := accum

	joinPath := JoinPath(resolvedSegments, targetRoot)
	if uriPrefix != "" {
		return uriPrefix + joinPath
	}
	if isWebroot {
		return "/" + joinPath
	}
	return joinPath
}

/*
Calculate the relative path to this absolute filename from
the specified base directory.
If either the filename or the base_directory are not absolute paths,
no work is done.
filename       - An absolute file name as a String
base_directory - An absolute base directory as a String
Return the relative path String of the filename calculated
from the base directory
*/
func (pr *PathResolver) RelativePath(filename, baseDirectory string) string {
	if IsRoot(filename) && IsRoot(baseDirectory) {
		offset := baseDirectory
		if strings.HasSuffix(baseDirectory, string(pr.FileSeparator())) {
			offset = baseDirectory[:len(baseDirectory)-1]
		}
		filename = filename[len(offset)+1:]
	}
	return filename
}
