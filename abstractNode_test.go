package asciidocgo

import (
	"strconv"
	"strings"
	"testing"

	"github.com/VonC/asciidocgo/consts/context"
	"github.com/VonC/asciidocgo/consts/safemode"
	. "github.com/smartystreets/goconvey/convey"
)

type testDocumentAble struct {
	*abstractNode
}

func TestAbstractNode(t *testing.T) {

	Convey("An abstractNode can be initialized", t, func() {

		Convey("By default, an AbstractNode can be created", func() {
			So(&abstractNode{}, ShouldNotBeNil)
		})
		Convey("An AbstractNode takes a parent and a context", func() {
			So(newAbstractNode(nil, context.Document), ShouldNotBeNil)
		})
		Convey("If context is a document, then parent is nil and document is parent", func() {
			parent := newTestDocumentAble(nil)
			an := newAbstractNode(parent.abstractNode, context.Document)
			So(an.Context(), ShouldEqual, context.Document)
			So(an.Parent(), ShouldBeNil)
			So(an.Document(), ShouldEqual, parent)
			So(an.Id(), ShouldEqual, "")
		})
		Convey("If context is not document, then parent is parent and document is parent document", func() {
			doc := newTestDocumentAble(nil)
			parent := newAbstractNode(doc.abstractNode, context.Document)
			an := newAbstractNode(parent, context.Section)
			So(an.Context(), ShouldEqual, context.Section)
			So(an.Parent(), ShouldEqual, parent)
			So(an.Document(), ShouldEqual, parent.Document())
		})
		Convey("If context is not document, and parent is nil, then document is nil", func() {
			an := newAbstractNode(nil, context.Section)
			So(an.Context(), ShouldEqual, context.Section)
			So(an.Parent(), ShouldBeNil)
			So(an.Document(), ShouldBeNil)
		})

		Convey("An abstractNode has an empty attributes map", func() {
			an := newAbstractNode(nil, context.Section)
			So(len(an.Attributes()), ShouldEqual, 0)
		})
	})

	Convey("An abstractNode can be associated to a parent", t, func() {
		an := newAbstractNode(nil, context.Section)
		documentParent := newTestDocumentAble(nil)
		parent := newAbstractNode(documentParent.abstractNode, context.Document)
		an.SetParent(parent)
		So(an.Parent(), ShouldEqual, parent)
		So(an.Document(), ShouldEqual, parent.Document())
	})

	Convey("An abstractNode can retrieve and set an id", t, func() {
		an := newAbstractNode(nil, context.Section)
		So(an.Id(), ShouldEqual, "")
		an.SetId("test")
		So(an.Id(), ShouldEqual, "test")
	})

	Convey("An abstractNode can retrieve an attribute", t, func() {

		parentDocument := newTestDocumentAble(nil)
		parentDocument.setAttr("key", "val1", true)
		parent := newTestDocumentAble(parentDocument.abstractNode).abstractNode
		an := newTestDocumentAble(nil).abstractNode
		Convey("If inherited, it is the attribute if there, or the document attribute, or default value", func() {
			So(an.Attr("key", nil, true), ShouldBeNil)
			an.setAttr("key", "val", false)
			So(an.Attr("key", nil, true), ShouldEqual, "val")
			an.SetParent(parent)
			So(an.Attr("key", nil, true), ShouldEqual, "val")
			delete(an.attributes, "key")
			So(an.Attr("key", nil, true), ShouldEqual, "val1")
			// an should have for parent a child, which has an as a document
			// then an.document would be "parent".document, meaning an, when
			// setting an.setParent(child)
			an.document = an._doc
			an.setAttr("key", "val", false)
			So(an.Attr("key", nil, true), ShouldEqual, "val")
		})
		Convey("If not inherited, it is the attribute if there, or default value", func() {
			an.SetParent(parent)
			So(an.Document(), ShouldEqual, parentDocument)
			So(an.Document(), ShouldEqual, parent.Document())
			So(parentDocument.Attr("key", nil, false), ShouldEqual, "val1")
			So(an.Attr("key", nil, false), ShouldEqual, "val")
		})
	})

	Convey("An abstractNode can check for an attribute", t, func() {
		parentDocument := newTestDocumentAble(nil)
		parentDocument.setAttr("key1", "val1", true)
		parent := newAbstractNode(parentDocument.abstractNode, context.Document)
		an := newTestDocumentAble(nil).abstractNode

		Convey("If expect nil, check if key is there", func() {
			So(an.HasAttr("key1", nil, false), ShouldBeFalse)
			an.SetParent(parent)
			So(an.Document(), ShouldEqual, parentDocument)
			So(an.HasAttr("key1", nil, false), ShouldBeFalse)
			So(an.HasAttr("key1", nil, true), ShouldBeTrue)
			an.SetParent(nil)
			So(an.HasAttr("key1", nil, false), ShouldBeFalse)
			an.setAttr("key1", nil, true)
			So(an.HasAttr("key1", nil, false), ShouldBeTrue)
		})
		Convey("If expect not nil, check if key is there and the value matches", func() {
			an.SetParent(parent)
			So(an.HasAttr("key1", "val1", false), ShouldBeFalse)
			So(an.HasAttr("key1", "val1", true), ShouldBeTrue)
			an.setAttr("key1", "val1", true)
			an.SetParent(nil)
			So(an.HasAttr("key1", "val1", false), ShouldBeTrue)
			an.document = an._doc
			So(an.HasAttr("key1", "val1", true), ShouldBeTrue)
		})
	})

	Convey("An abstractNode can set an attribute", t, func() {
		an := newAbstractNode(nil, context.Document)
		an.setAttr("key", "val", true)
		So(an.Attr("key", nil, true), ShouldEqual, "val")
		Convey("If not override and already present, the value should not change", func() {
			res := an.setAttr("key", "val1", false)
			So(an.Attr("key", nil, true), ShouldEqual, "val")
			So(res, ShouldBeFalse)
		})
	})

	Convey("An abstractNode can set an option attribute", t, func() {
		an := newAbstractNode(nil, context.Document)
		Convey("First option means options attributes has len 1", func() {
			an.SetOption("opt1")
			So(len(an.attributes["options"].(map[string]bool)), ShouldEqual, 1)
			So(an.Attr("opt1-option", nil, false), ShouldEqual, true)
		})
		Convey("Second option means options attributes has len 2", func() {
			an.SetOption("opt2")
			So(len(an.attributes["options"].(map[string]bool)), ShouldEqual, 2)
			So(an.Attr("opt2-option", nil, false), ShouldEqual, true)
		})
	})

	Convey("An abstractNode can get an option attribute", t, func() {
		an := newAbstractNode(nil, context.Document)
		Convey("Zero option means Option returns false", func() {
			So(an.HasOption("opt1"), ShouldBeFalse)
			an.SetOption("opt1")
		})
		Convey("One option means Option returns true", func() {
			So(an.HasOption("opt1"), ShouldBeTrue)
		})
	})

	Convey("An abstractNode update option attributes with other attributes", t, func() {
		an := newAbstractNode(nil, context.Document)
		an.setAttr("key1", "val1", true)
		an.setAttr("key2", "val2", true)
		Convey("New Attributes are added during an update", func() {
			attrs := map[string]interface{}{"key3": "val3", "key4": "val4"}
			an.UpdateAttributes(attrs)
			So(an.Attr("key1", nil, false), ShouldEqual, "val1")
			So(an.Attr("key2", nil, false), ShouldEqual, "val2")
			So(an.Attr("key3", nil, false), ShouldEqual, "val3")
			So(an.Attr("key4", nil, false), ShouldEqual, "val4")
		})
		Convey("Common Attributes are overrriden during an update", func() {
			attrs := map[string]interface{}{"key2": "val2b", "key3": "val3"}
			an.UpdateAttributes(attrs)
			So(an.Attr("key1", nil, false), ShouldEqual, "val1")
			So(an.Attr("key2", nil, false), ShouldEqual, "val2b")
			So(an.Attr("key3", nil, false), ShouldEqual, "val3")
		})
	})
	Convey("An abstractNode can check for a role", t, func() {
		an := newAbstractNode(nil, context.Document)
		parentDocument := newTestDocumentAble(nil).abstractNode
		parentDocument.setAttr("role", "roleFromParentDocument", true)
		parent := newTestDocumentAble(parentDocument).abstractNode
		Convey("A role can be checked, whatever its value is", func() {
			So(an.HasRole(nil), ShouldBeFalse)
			So(parentDocument.HasRole(nil), ShouldBeTrue)
			an.SetParent(parent)
			So(an.HasRole(nil), ShouldBeTrue)
		})
		Convey("A role can be checked against an expected value", func() {
			an := newAbstractNode(nil, context.Document)
			an.SetParent(parent)
			So(an.HasRole("roleFromAN"), ShouldBeFalse)
			So(an.HasRole("roleFromParentDocument"), ShouldBeTrue)
			an.setAttr("role", "roleFromAN", true)
			So(an.HasRole("roleFromAN"), ShouldBeTrue)
			So(an.HasRole("roleFromParentDocument"), ShouldBeFalse)
		})
	})
	Convey("An abstractNode can check for a role name", t, func() {
		an := newAbstractNode(nil, context.Document)
		parentDocument := newTestDocumentAble(nil)
		parentDocument.setAttr("role", "role1FromParentDocument role2FromParentDocument role3FromParentDocument", true)
		parent := newTestDocumentAble(parentDocument.abstractNode).abstractNode

		Convey("A role name can be checked on the document", func() {
			So(an.HasARole("role3FromAN"), ShouldBeFalse)
			So(an.HasARole("role2FromParentDocument"), ShouldBeFalse)
			an.SetParent(parent)
			So(an.Document(), ShouldEqual, parentDocument)
			//So(an.Document().Attr("role", nil, true), ShouldEqual, "test")
			So(an.HasARole("role2FromParentDocument"), ShouldBeTrue)
			So(an.HasARole("role4FromParentDocument"), ShouldBeFalse)
		})
		Convey("A role name can be checked on the abstractNode itself", func() {
			an.setAttr("role", "role1FromAN role2FromAN role3FromAN", true)
			So(an.HasARole("role3FromAN"), ShouldBeTrue)
		})
		Convey("An empty role name is always false=", func() {
			So(an.HasARole(""), ShouldBeFalse)
		})
	})
	Convey("An abstractNode can access role", t, func() {
		an := newAbstractNode(nil, context.Document)
		parentDocument := newTestDocumentAble(nil).abstractNode
		parentDocument.setAttr("role", "roleFromParentDocument", true)
		parent := newTestDocumentAble(parentDocument).abstractNode
		Convey("A role can be access from an document, when an has no role", func() {
			So(an.Role(), ShouldBeNil)
			So(parentDocument.Role(), ShouldEqual, "roleFromParentDocument")
			an.SetParent(parent)
			So(an.Role(), ShouldEqual, "roleFromParentDocument")
		})
		Convey("A role can be access from an itself", func() {
			an.setAttr("role", "roleFromAN", true)
			So(an.Role(), ShouldEqual, "roleFromAN")
		})
	})
	Convey("An abstractNode can access role names", t, func() {
		an := newAbstractNode(nil, context.Document)
		parentDocument := newTestDocumentAble(nil)
		parentDocument.setAttr("role", "role1FromParentDocument role2FromParentDocument role3FromParentDocument", true)
		parent := newTestDocumentAble(parentDocument.abstractNode).abstractNode

		Convey("A role name can be accessed on the document", func() {
			So(len(an.RoleNames()), ShouldBeZeroValue)
			an.SetParent(parent)
			So(an.Document(), ShouldEqual, parentDocument)
			So(len(an.RoleNames()), ShouldEqual, 3)
		})
		Convey("A role name can be accessed on the abstractNode itself", func() {
			an.setAttr("role", "role1FromAN role2FromAN role3FromAN role5FromAN role4FromAN", true)
			So(len(an.RoleNames()), ShouldEqual, 5)
		})
	})

	Convey("An abstractNode can check for a reftext", t, func() {
		an := newAbstractNode(nil, context.Document)
		parentDocument := newTestDocumentAble(nil).abstractNode
		parentDocument.setAttr("reftext", "reftextFromParentDocument", true)
		parent := newTestDocumentAble(parentDocument).abstractNode
		Convey("A reftext can be checked on the document", func() {
			So(an.HasReftext(), ShouldBeFalse)
			So(parentDocument.HasReftext(), ShouldBeTrue)
			an.SetParent(parent)
			So(an.HasReftext(), ShouldBeTrue)
		})
		Convey("A reftext can be checked directly on the abstractNode", func() {
			an := newAbstractNode(nil, context.Document)
			parentDocument := newTestDocumentAble(nil).abstractNode
			parent := newTestDocumentAble(parentDocument).abstractNode
			So(an.HasReftext(), ShouldBeFalse)
			an.SetParent(parent)
			So(an.HasReftext(), ShouldBeFalse)
			an.setAttr("reftext", "reftextFromAN", true)
			So(an.HasReftext(), ShouldBeTrue)
			So(an.Document().HasReftext(), ShouldBeFalse)
		})
	})

	Convey("An abstractNode can access reftext", t, func() {
		an := newAbstractNode(nil, context.Document)
		parentDocument := newTestDocumentAble(nil).abstractNode
		parentDocument.setAttr("reftext", "reftextFromParentDocument", true)
		parent := newTestDocumentAble(parentDocument).abstractNode
		Convey("A reftext can be access from an document, when an has no reftext", func() {
			So(an.Reftext(), ShouldBeNil)
			So(parentDocument.Reftext(), ShouldEqual, "reftextFromParentDocument")
			an.SetParent(parent)
			So(an.Reftext(), ShouldEqual, "reftextFromParentDocument")
		})
		Convey("A reftext can be access from an itself", func() {
			an.setAttr("reftext", "reftextFromAN", true)
			So(an.Reftext(), ShouldEqual, "reftextFromAN")
		})
	})

	Convey("An abstractNode can check for Slash Usage", t, func() {
		an := newAbstractNode(nil, context.Section)
		parentDocument := newTestDocumentAble(nil).abstractNode
		Convey("A abstractNode without document won't use slash", func() {
			So(an.ShortTagSlash(), ShouldBeNil)
		})
		Convey("A abstractNode with document htmlsyntax set to not xml won't use slash", func() {
			parentDocument.setAttr("htmlsyntax", "notxml", true)
			parent := newTestDocumentAble(parentDocument).abstractNode
			an.SetParent(parent)
			So(an.ShortTagSlash(), ShouldBeNil)
		})
		Convey("A abstractNode with document htmlsyntax set to not xml will use slash", func() {
			parentDocument.setAttr("htmlsyntax", "xml", true)
			So(strings.Trim(strconv.QuoteRune(*an.ShortTagSlash()), "'"), ShouldEqual, "/")
		})
	})
	Convey("An abstractNode can build media uri", t, func() {
		an := newAbstractNode(nil, context.Section)
		target := ""
		Convey("If the target media is a URI reference, then leave it untouched.", func() {
			target = "data:info"
			// So(REGEXP[":uri_sniff"].String(), ShouldEqual, "^[a-zA-Z][a-zA-Z0-9.+-]*:/{0,2}.*")
			// So(REGEXP[":uri_sniff"].MatchString(target), ShouldBeTrue)
			So(an.MediaUri(target, ""), ShouldEqual, target)
		})
		Convey("If the assetDirKey attribute is not there, normalize the target.", func() {
			target = "data"
			So(an.MediaUri(target, ""), ShouldEqual, target)
		})
		Convey("If the assetDirKey attribute is there, normalize the target with it.", func() {
			an.setAttr("imagesdir", "/images", true)
			target = "data"
			So(an.MediaUri(target, ""), ShouldEqual, "/images/data")
			So(an.MediaUri(target, "dummy"), ShouldEqual, target)
			So(an.MediaUri("http://a/b", ""), ShouldEqual, "http://a/b")
		})

	})

	Convey("An abstractNode can build icon uri", t, func() {
		parent := newTestDocumentAble(nil).abstractNode
		an := newAbstractNode(parent, context.Document)
		target := ""
		Convey("If the icon attribute is not there, imageUri with icon", func() {
			So(an.IconUri(target), ShouldEqual, ".png")
		})
		Convey("If the icon attribute is there, imageUri with name plus icontype (default png)", func() {
			parent.setAttr("icontype", "jpg", true)
			So(an.IconUri("a"), ShouldEqual, "a.jpg")
			an.setAttr("icon", "/icon", true)
			So(an.IconUri("a.anext"), ShouldEqual, "/icon")
		})
	})

	Convey("An abstractNode can normalize system paths", t, func() {
		parent := newTestDocumentAble(nil).abstractNode
		an := newAbstractNode(parent, context.Document)
		pr := NewPathResolver(0, "")
		wd := Posixfy(pr.WorkingDir())
		Convey("Empty target means working dir", func() {
			So(an.normalizeSystemPath("", "", "", false, ""), ShouldEqual, Posixfy(pr.WorkingDir()))
		})
		Convey("Empty start and jail means working dir", func() {
			So(an.normalizeSystemPath("a/b", "", "", false, ""), ShouldEqual, wd+"/a/b")
		})
		Convey("Empty start and jail and safe document means working dir", func() {
			testan = "test_normalizeSystemPath_safeDocument"
			So(an.normalizeSystemPath("a/b", "", "", false, ""), ShouldEqual, wd+"/a/b")
			testan = ""
		})
	})

	Convey("An abstractNode can generate data uri", t, func() {
		parent := newTestDocumentAble(nil).abstractNode
		an := newAbstractNode(parent, context.Document)
		pr := NewPathResolver(0, "")
		wd := Posixfy(pr.WorkingDir())

		Convey("Empty target and assetDir means working dir, meaning defaut data uri content", func() {
			So(an.generateDataUri("", ""), ShouldEqual, "data:image/:base64,")
			So(an.generateDataUri("a/b.exe", ""), ShouldEqual, "data:image/exe:base64,")
		})
		Convey("Svg non-existing target and empty assetDir means data: with svg+xml mimetype", func() {
			So(an.generateDataUri("a/b.svg", ""), ShouldEqual, "data:image/svg+xml:base64,")
		})
		Convey("Svg target and non-empty assetDir imagePath", func() {
			testan = "test_generateDataUri_imagePath"
			So(an.generateDataUri("a/b.svg", "akey"), ShouldEqual, "imagePath='"+wd+"/a/b.svg'")
			parent.setAttr("akey", "c:/x", true)
			So(an.generateDataUri("a/b.svg", "akey"), ShouldEqual, "imagePath='c:/x/a/b.svg'")
			testan = ""
		})
		Convey("Existing target and empty assetDir means data content", func() {
			So(an.generateDataUri("test/t.txt", ""), ShouldEqual, "test data")
		})
	})

	Convey("An abstractNode can build image uri", t, func() {
		parent := newTestDocumentAble(nil)
		an := newAbstractNode(parent.abstractNode, context.Document)
		Convey("If the data-uri attribute is not there, imageUri with icon", func() {
			So(an.ImageUri("http://a/b", ""), ShouldEqual, "http://a/b")
		})
		Convey("If the assetDirKey attribute is there, imageUri with start", func() {
			assetDirKey := "start-uri"
			an.Document().setAttr("start-uri", "/a/b", true)
			So(an.Document(), ShouldEqual, parent)
			So(an.HasAttr("start-uri", nil, true), ShouldBeTrue)
			So(an.Document().HasAttr(assetDirKey, nil, true), ShouldBeTrue)
			So(an.Attr(assetDirKey, nil, true).(string), ShouldEqual, "/a/b")
			So(an.Document().Attr(assetDirKey, nil, true).(string), ShouldEqual, "/a/b")
			So(an.ImageUri("c/d", assetDirKey), ShouldEqual, "/a/b/c/d")
			an.Document().setAttr("start-uri", "a/b", true)
			So(an.Attr(assetDirKey, nil, true).(string), ShouldEqual, "a/b")
			So(an.ImageUri("c/d", assetDirKey), ShouldEqual, "a/b/c/d")
			So(an.ImageUri("c/d", assetDirKey+"2"), ShouldEqual, "c/d")
		})

		Convey("If the data-uri attribute is on the Document, generate data uri", func() {
			an.Document().setAttr("data-uri", "anything", true)
			an.ImageUri("c/d", "")
			So(an.ImageUri("c/d.anext", ""), ShouldEqual, "data:image/anext:base64,")
		})
	})
	Convey("An abstractNode can read asset", t, func() {
		Convey("It can warn on an non-existing asset", func() {
			So(ReadAsset("a/b.txt", true), ShouldEqual, "")
		})
		Convey("It read an existing asset", func() {
			So(ReadAsset("test/t.txt", true), ShouldEqual, "test data")
		})
	})
	Convey("An abstractNode can normalize asset path", t, func() {
		parent := newTestDocumentAble(nil).abstractNode
		an := newAbstractNode(parent, context.Document)
		pr := NewPathResolver(0, "")
		wd := Posixfy(pr.WorkingDir())
		Convey("Empty parameters means working directory", func() {
			So(an.normalizeAssetPath("", "", false), ShouldEqual, wd)
		})
		Convey("target means working directory + target", func() {
			So(an.normalizeAssetPath("a/b", "", false), ShouldEqual, wd+"/a/b")
		})
	})
	Convey("An abstractNode can compute relative path", t, func() {
		parent := newTestDocumentAble(nil).abstractNode
		an := newAbstractNode(parent, context.Document)
		Convey("Empty parameters means empty relative path", func() {
			So(an.relativePath(""), ShouldEqual, "")
		})
		Convey("Non-empty filename means filename if document has no basedir", func() {
			So(an.relativePath("a/b.txt"), ShouldEqual, "a/b.txt")
		})
	})

	Convey("An abstractNode can access marker keyword", t, func() {
		parent := newAbstractNode(nil, context.Section)
		an := newAbstractNode(parent, context.Document)
		Convey("Empty parameters means empty style means empty rune", func() {
			So(an.listMarkerKeyword(""), ShouldEqual, 0)
		})
		Convey("Non-empty listType means rune", func() {
			So(an.listMarkerKeyword("upperalpha"), ShouldEqual, 'A')
		})
	})

}

func newTestDocumentAble(parent *abstractNode) *testDocumentAble {
	an := newAbstractNode(parent, context.Document)
	testDocumentAble := &testDocumentAble{an}
	an.MainDocumentable(testDocumentAble)
	return testDocumentAble
}

func (td *testDocumentAble) Safe() safemode.SafeMode {
	return safemode.UNSAFE
}

func (td *testDocumentAble) BaseDir() string {
	return ""
}

func (td *testDocumentAble) PlaybackAttributes(map[string]interface{}) {
	//
}

func (td *testDocumentAble) CounterIncrement(counterName string, block *abstractNode) string {
	return ""
}

func (td *testDocumentAble) Counter(name, seed string) int {
	return -1
}

func (td *testDocumentAble) DocType() string {
	return ""
}
