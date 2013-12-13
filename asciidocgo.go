/* Asciidocgo implements an AsciiDoc renderer in Go. */
package asciidocgo

import "io"

/*
Accepts input as an IO (or StringIO), String or String Array object.
If the input is a File, information about the file is stored in attributes on
the Document object.
*/
func Load(input io.Reader) *Document {
	return nil
}
