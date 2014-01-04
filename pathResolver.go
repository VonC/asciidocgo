package asciidocgo

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
