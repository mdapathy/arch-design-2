package gomodule

import (
	"fmt"
	"github.com/google/blueprint"
	"github.com/roman-mazur/bood"
	"path"
	"regexp"
)

var (
	// Package context used to define Ninja build rules.
	pctx = blueprint.NewPackageContext("github.com/mdapathy/build/gomodule")

	// Ninja rule to execute go build.
	goBuild = pctx.StaticRule("binaryBuild", blueprint.RuleParams{
		Command:     "cd $workDir && go build -o $outputPath $pkg",
		Description: "build go command $pkg",
	}, "workDir", "outputPath", "pkg")

	// Ninja rule to execute go mod vendor.
	goVendor = pctx.StaticRule("vendor", blueprint.RuleParams{
		Command:     "cd $workDir && go mod vendor",
		Description: "vendor dependencies of $name",
	}, "workDir", "name")

	goTest = pctx.StaticRule("gotest", blueprint.RuleParams{
		Command:     "cd ${workDir}  && go test -v ${benchmark} ${pkg} > ${outputfile}",
		Description: "test ${pkg}",
	}, "workDir", "pkg", "benchmark", "outputfile")
)

type testedBinaryModule struct {
	blueprint.SimpleName

	properties struct {
		Name           string
		Pkg            string
		TestPkg        string
		Srcs           []string
		SrcsExclude    []string
		BenchmarkTimes string
		VendorFirst    bool
	}
}

func (tb *testedBinaryModule) GenerateBuildActions(ctx blueprint.ModuleContext) {
	name := ctx.ModuleName()
	config := bood.ExtractConfig(ctx)
	config.Debug.Printf("Adding build actions for go binary module '%s'", name)

	outputPath := path.Join(config.BaseOutputDir, "bin", name)
	testOutput := path.Join(config.BaseOutputDir, "reports", "test.txt")

	var testInputs []string
	var buildInputs []string
	inputErrors := false

	for _, src := range tb.properties.Srcs {
		if matches, err := ctx.GlobWithDeps(src, tb.properties.SrcsExclude); err == nil {
			testInputs = append(testInputs, matches...)
			for _, i := range matches {
				if val, _ := regexp.Match(".*_test\\.go$", []byte(i)); val == false {
					buildInputs = append(buildInputs, i)
				}
			}

		} else {
			ctx.PropertyErrorf("srcs", "Cannot resolve files that match pattern %s", src)
			inputErrors = true
		}
	}
	if inputErrors {
		return
	}

	if tb.properties.VendorFirst {
		vendorDirPath := path.Join(ctx.ModuleDir(), "vendor")
		ctx.Build(pctx, blueprint.BuildParams{
			Description: fmt.Sprintf("Vendor dependencies of %s", name),
			Rule:        goVendor,
			Outputs:     []string{vendorDirPath},
			Implicits:   []string{path.Join(ctx.ModuleDir(), "go.mod")},
			Optional:    true,
			Args: map[string]string{
				"workDir": ctx.ModuleDir(),
				"name":    name,
			},
		})

		buildInputs = append(buildInputs, vendorDirPath)
		testInputs = append(testInputs, vendorDirPath)
	}

	if len(tb.properties.Srcs) > 0 && len(tb.properties.Pkg) > 0 {
		ctx.Build(pctx, blueprint.BuildParams{
			Description: fmt.Sprintf("Build %s as Go binary", name),
			Rule:        goBuild,
			Outputs:     []string{outputPath},
			Implicits:   buildInputs,

			Args: map[string]string{
				"outputPath": outputPath,
				"workDir":    ctx.ModuleDir(),
				"pkg":        tb.properties.Pkg,
			},
		})
	}

	if len(tb.properties.TestPkg) > 0 {
		bench := "-bench=."

		if len(tb.properties.BenchmarkTimes) > 0 {
			bench += fmt.Sprintf(" -benchtime=%sx", tb.properties.BenchmarkTimes)
		}

		ctx.Build(pctx, blueprint.BuildParams{
			Description: fmt.Sprintf("Test module %s", tb.properties.TestPkg),
			Rule:        goTest,
			Outputs:     []string{testOutput},
			Implicits:   testInputs,
			Args: map[string]string{
				"benchmark":  bench,
				"outputfile": testOutput,
				"workDir":    ctx.ModuleDir(),
				"pkg":        tb.properties.TestPkg,
			},
		})
	}

}

func TestedBinFactory() (blueprint.Module, []interface{}) {
	mType := &testedBinaryModule{}
	return mType, []interface{}{&mType.SimpleName.Properties, &mType.properties}
}
