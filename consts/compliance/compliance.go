package compliance

// Flags to control compliance with the behavior of AsciiDoc
type compliance struct {
	block_terminates_paragraph     bool
	strict_verbatim_paragraphs     bool
	underline_style_section_titles bool
	unwrap_standalone_preamble     bool
	attribute_missing              string
	attribute_undefined            string
	markdown_syntax                bool
}

var cpl = &compliance{
	//congruent_block_delimiters: true,
	block_terminates_paragraph:     true,
	strict_verbatim_paragraphs:     true,
	underline_style_section_titles: true,
	unwrap_standalone_preamble:     true,
	attribute_missing:              "skip",
	attribute_undefined:            "drop-line",
	markdown_syntax:                true,
}

/* AsciiDoc terminates paragraphs adjacent to block content
(delimiter or block attribute list)
This option allows this behavior to be modified
TODO what about literal paragraph?
Compliance value: true */
func BlockTerminatesParagraph() bool {
	return cpl.block_terminates_paragraph
}

/* AsciiDoc does not treat paragraphs labeled with a verbatim style
(literal, listing, source, verse) as verbatim.
This options allows this behavior to be modified
Compliance value: false */
func StrictVerbatimParagraphs() bool {
	return cpl.strict_verbatim_paragraphs
}

/* NOT CURRENTLY USED
AsciiDoc allows start and end delimiters around a block to be different lengths
Enabling this option requires matching lengths
Compliance value: false
func CongruentBlockDelimiters() bool {
	return cpl.congruent_block_delimiters
}
*/

/* AsciiDoc supports both single-line and underlined section titles.
This option disables the underlined variant.
Compliance value: true */
func UnderlineStyleSectionTitles() bool {
	return cpl.underline_style_section_titles
}

/* Asciidoctor will unwrap the content in a preamble
if the document has a title and no sections.
Compliance value: false */
func UnwrapStandalonePreamble() bool {
	return cpl.unwrap_standalone_preamble
}

/* AsciiDoc drops lines that contain references to missing attributes.
This behavior is not intuitive to most writers
Compliance value: 'drop-line' */
func AttributeMissing() string {
	return cpl.attribute_missing
}

/* AsciiDoc drops lines that contain an attribute unassignemnt.
This behavior may need to be tuned depending on the circumstances.
Compliance value: 'drop-line' */
func AttributeUndefined() string {
	return cpl.attribute_undefined
}

/* Asciidoctor will recognize commonly-used Markdown syntax
to the degree it does not interfere with existing AsciiDoc syntax and behavior.
Compliance value: false */
func MarkdownSyntax() bool {
	return cpl.markdown_syntax
}
