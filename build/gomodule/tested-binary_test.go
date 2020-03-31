package gomodule

import (
	"bytes"
	"github.com/google/blueprint"
	"github.com/roman-mazur/bood"
	"strings"
	"testing"
)

var fileSystems = []map[string][]byte{
	{
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
	},
	{
		"Blueprints": []byte(`
			go_binary {
			  name: "package-out",
			  pkg: ".",
			  srcs: [ "main_test.go", "main.go",],
			}
		`),
		"main.go":      nil,
		"main_test.go": nil,
	},

	{
		"Blueprints": []byte(`
			go_binary {
			  name: "new-package",
			  pkg: ".",
              testPkg: ".",
			  srcs: [ "main_test.go", "main.go",],
			  vendorFirst: true
			}
		`),
		"main.go":      nil,
		"main_test.go": nil,
	},
}

var expectedOutput = [][]string{
	{
		"out/bin/package-out:",
		"g.gomodule.binaryBuild | main.go\n",
		"out/reports/test.txt",
		"g.gomodule.gotest | main_test.go main.go",
	},
	{
		"out/bin/package-out:",
		"g.gomodule.binaryBuild | main.go\n",
	},
	{
		"out/bin/new-package",
		"g.gomodule.binaryBuild | main.go vendor\n",
		"build vendor: g.gomodule.vendor | go.mod\n",
		"out/reports/test.txt",
		"g.gomodule.gotest | main_test.go main.go",
	},
}

func TestTestedBinFactory(t *testing.T) {
	for index, fileSystem := range fileSystems {
		t.Run(string(index), func(t *testing.T) {
			ctx := blueprint.NewContext()

			ctx.MockFileSystem(fileSystem)

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
		//		t.Logf("Gennerated ninja build file:\n%s", text)
				for _, expectedStr := range expectedOutput[index] {
					//build rule
					if strings.Contains(text, expectedStr) != true {
						t.Errorf("Generated ninja file does not have expected string `%s`", expectedStr)
					}
				}

			}

		})
	}
}
