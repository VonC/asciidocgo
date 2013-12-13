/* Asciidocgo implements an AsciiDoc renderer in Go.

Methods for parsing Asciidoc input files and rendering documents using eRuby
templates.

Asciidoc documents comprise a header followed by zero or more sections.
Sections are composed of blocks of content.  For example:

  = Doc Title

  == Section 1

  This is a paragraph block in the first section.

  == Section 2

  This section has a paragraph block and an olist block.

  . Item 1
  . Item 2

Examples:

Use built-in templates:

  lines = File.readlines("your_file.asc")
  doc = Asciidoctor::Document.new(lines)
  html = doc.render
  File.open("your_file.html", "w+") do |file|
    file.puts html
  end

Use custom (Tilt-supported) templates:

  lines = File.readlines("your_file.asc")
  doc = Asciidoctor::Document.new(lines, :template_dir => 'templates')
  html = doc.render
  File.open("your_file.html", "w+") do |file|
    file.puts html
  end

*/
package asciidocgo

import "io"

// Accepts input as an IO (or StringIO), String or String Array object.
// If the input is a File, information about the file is stored in attributes on
// the Document object.
func Load(input io.Reader) *Document {
	return nil
}
