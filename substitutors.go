package asciidocgo

type subsDef struct {
	basic    []string
	normal   []string
	verbatim []string
	title    []string
	header   []string
	pass     []string
}

func newSubs() *subsDef {
	return &subsDef{
		/* basic    */ []string{"specialcharacters"},
		/* normal   */ []string{"specialcharacters", "quotes", "attributes", "replacements", "macros", "post_replacements"},
		/* verbatim */ []string{"specialcharacters", "attributes"},
		/* title    */ []string{"specialcharacters", "quotes", "replacements", "macros", "post_replacements"},
		/* header   */ []string{"specialcharacters", "attributes"},
		// by default, Asciidocgo performs :attributes and :macros on a pass block
		/* pass     */ []string{}}
}

var subs = newSubs()

/* Methods to perform substitutions on lines of AsciiDoc text.
This module is intented to be mixed-in to Section and Block to provide
operations for performing the necessary substitutions. */
type substitutors struct {
	// A String Array of passthough (unprocessed) text captured from this block
	passthroughs []string
}
