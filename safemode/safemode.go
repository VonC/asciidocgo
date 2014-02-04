package safemode

// Symbol name for the type of content (e.g., :paragraph).
type SafeMode int

const (
	/* A safe mode level that disables any of the security features enforced
	   by Asciidocgo (Go is still subject to its own restrictions). */
	UNSAFE SafeMode = iota
	/* A safe mode level that closely parallels safe mode in AsciiDoc.
	   This value prevents access to files which reside outside of the
	   parent directory of the source file and disables any macro other
	   than the include::[] macro. */
	SAFE
	/*A safe mode level that disallows the document from setting attributes
	  that would affect the rendering of the document, in addition to all the
	  security features of SafeMode::SAFE. For instance, this level disallows
	  changing the backend or the source-highlighter using an attribute defined
	  in the source document. This is the most fundamental level of security
	  for server-side deployments (hence the name).*/
	SERVER
	/*A safe mode level that disallows the document from attempting to read
	  files from the file system and including the contents of them into the
	  document, in additional to all the security features of SafeMode::SERVER.
	  For instance, this level disallows use of the include::[] macro and the
	  embedding of binary content (data uri), stylesheets and JavaScripts
	  referenced by the document.(Asciidoctor and trusted extensions may still
	  be allowed to embed trusted content into the document).

	  Since Asciidocgo is aiming for wide adoption, this level is the default
	  and is recommended for server-side deployments.*/
	SECURE
	/*A planned safe mode level that disallows the use of passthrough macros and
	  prevents the document from setting any known attributes,
	  in addition to all the security features of SafeMode::SECURE.

	  Please note that this level is not currently implemented
	  (and therefore not enforced)! */
	PARANOID
)
