package gomodule

import (
	"bytes"
	"github.com/google/blueprint"
	"github.com/roman-mazur/bood"
	"strings"
	"testing"
)

func TestTestedBinFactory(t *testing.T) {
	ctx := blueprint.NewContext()

	ctx.MockFileSystem(map[string][]byte{
		"Blueprints": []byte(`
			go_binary {
			  name: "package-out",
			  pkg: ".",
              testPkg: ".",
			  srcs: [ "main_test.go", "main.go",],
			}
		`),
		"main.go":      nil,
		"main_test.go": nil,
	})

	ctx.RegisterModuleType("go_binary", TestedBinFactory)

	cfg := bood.NewConfig()

	_, errs := ctx.ParseBlueprintsFiles(".", cfg)
	if len(errs) != 0 {
		t.Fatalf("Syntax errors in the test blueprint file: %s", errs)
	}

	_, errs = ctx.PrepareBuildActions(cfg)
	if len(errs) != 0 {
		t.Errorf("Unexpected errors while preparing build actions: %s", errs)
	}
	buffer := new(bytes.Buffer)
	if err := ctx.WriteBuildFile(buffer); err != nil {
		t.Errorf("Error writing ninja file: %s", err)
	} else {
		text := buffer.String()
		t.Logf("Gennerated ninja build file:\n%s", text)

		if !strings.Contains(text, "out/bin/package-out: ") {
			t.Errorf("Generated ninja file does not have build of the test module")
		}

		if !strings.Contains(text, "build out/bin/package-out: g.gomodule.binaryBuild | main.go") {
			t.Errorf("Generated ninja file does not support testing of all files of the test module")
		}

		if !strings.Contains(text, "build out/bin/package-out: g.gomodule.binaryBuild | main.go\n") {
			t.Errorf("Generated ninja file's build has dependencies  of the test module")
		}

		if strings.Contains(text, "build vendor: g.gomodule.vendor | go.mod") {
			t.Errorf("Generated ninja file has a vendor build rule")
		}

	}
}
