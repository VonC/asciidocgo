package asciidocgo

/* Methods for rendering Asciidoc Documents, Sections, and Blocks
using <del>eRuby</del> Go templates */
type Renderer struct{}

/* Render an Asciidoc object with a specified view template.
view   - the String view template name.
object - the Object to be used as an evaluation scope.
locals - the optional Hash of locals to be passed to Tilt (default {})
(also ignored, really) */
func (r *Renderer) Render(view string, object interface{}, locals []interface{}) string { return "" }
