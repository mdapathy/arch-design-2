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
			  name: "package-out",
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

var fileSystemsResults = [][]bool{
	{
		true,
		true,
		true,
		true,
		false,
	},
	{
		true,
		true,
		false,
		false,
		false,
	},
	{
		true,
		true,
		true,
		true,
		true,
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
				t.Logf("Gennerated ninja build file:\n%s", text)

				//build rule
				if strings.Contains(text, "out/bin/package-out:") != fileSystemsResults[index][0] {
					t.Errorf("Generated ninja file does not have build of the test module")
				}
				if strings.Contains(text, " g.gomodule.binaryBuild | main.go\n") != fileSystemsResults[index][1] {
					t.Errorf("Generated ninja file's build depends on test files too")
				}

				//test rule
				if strings.Contains(text, "out/reports/test.txt") != fileSystemsResults[index][2] {
					t.Errorf("Generated ninja file does not create test result's file")
				}

				if strings.Contains(text, "g.gomodule.gotest | main_test.go main.go") != fileSystemsResults[index][3] {
					t.Errorf("Generated ninja file's test rule does not depend on all files")
				}

				//vendor rule
				if strings.Contains(text, "build vendor: g.gomodule.vendor | go.mod") != fileSystemsResults[index][4] {
					t.Errorf("Generated ninja file has a vendor build rule")
				}

			}

		})
	}
}
