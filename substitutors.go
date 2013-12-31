package asciidocgo

/* Methods to perform substitutions on lines of AsciiDoc text.
This module is intented to be mixed-in to Section and Block to provide
operations for performing the necessary substitutions. */
type substitutors struct {
	// A String Array of passthough (unprocessed) text captured from this block
	passthroughs []string
}
