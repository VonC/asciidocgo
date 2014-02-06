package asciidocgo

type _subs string

const (
	subsBasic    _subs = "basic"
	subsNormal   _subs = "normal"
	subsVerbatim _subs = "verbatim"
	subsTitle    _subs = "title"
	subsHeader   _subs = "header"
	subsPass     _subs = "pass"
	subsUnknown  _subs = "unknown"
)

type subsEnum struct {
	value _subs
}

type subsEnums struct {
	basic    *subsEnum
	normal   *subsEnum
	verbatim *subsEnum
	title    *subsEnum
	header   *subsEnum
	pass     *subsEnum
	unknown  *subsEnum
}

func newSubsEnums() *subsEnums {
	return &subsEnums{
		&subsEnum{subsBasic},
		&subsEnum{subsNormal},
		&subsEnum{subsVerbatim},
		&subsEnum{subsTitle},
		&subsEnum{subsHeader},
		&subsEnum{subsPass},
		&subsEnum{subsUnknown}}
}

func (se *subsEnum) values() []string {
	switch se.value {
	case subsBasic:
		return []string{"specialcharacters"}
	case subsNormal:
		return []string{"specialcharacters", "quotes", "attributes", "replacements", "macros", "post_replacements"}
	case subsVerbatim:
		return []string{"specialcharacters", "callouts"}
	case subsTitle:
		return []string{"specialcharacters", "quotes", "replacements", "macros", "post_replacements"}
	case subsHeader:
		return []string{"specialcharacters", "attributes"}
	case subsPass:
		return []string{}
	}
	return []string{}
}

var subs = newSubsEnums()

/* Methods to perform substitutions on lines of AsciiDoc text.
This module is intented to be mixed-in to Section and Block to provide
operations for performing the necessary substitutions. */
type substitutors struct {
	// A String Array of passthough (unprocessed) text captured from this block
	passthroughs []string
}

/* Apply the specified substitutions to the lines of text

source  - The String or String Array of text to process
subs    - The substitutions to perform. Can be a Symbol or a Symbol Array (default: :normal)
expand -  A Boolean to control whether sub aliases are expanded (default: true)

returns Either a String or String Array, whichever matches the type of the first argument */
func (s *substitutors) ApplySubs(source []string, sub *subsEnum) []string {
	if sub == nil || sub == subs.pass {
		return source
	}
	return []string{}
}
