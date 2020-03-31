package main

import (
	"flag"
	"github.com/google/blueprint"
	"github.com/mdapathy/arch-design-2/build/gomodule"
	"github.com/mdapathy/arch-design-2/build/jsmodule"
	boodmain "github.com/roman-mazur/bood"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

var (
	dryRun  = flag.Bool("dry-run", false, "Generate ninja build file but don't start the build")
	verbose = flag.Bool("v", false, "Display debugging logs")
)

func NewContext() *blueprint.Context {
	ctx := boodmain.PrepareContext()
	ctx.RegisterModuleType("go_binary", gomodule.TestedBinFactory)
	ctx.RegisterModuleType("js_bundle", jsmodule.JsMinimizedScriptFactory)
	return ctx
}

func main() {

	flag.Parse()

	config := boodmain.NewConfig()
	if !*verbose {
		config.Debug = log.New(ioutil.Discard, "", 0)
	}
	ctx := NewContext()

	ninjaBuildPath := boodmain.GenerateBuildFile(config, ctx)

	if !*dryRun {
		config.Info.Println("Starting the build now")
		cmd := exec.Command("ninja", append([]string{"-f", ninjaBuildPath}, flag.Args()...)...)
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			config.Info.Print(err)
			config.Info.Fatal("Error invoking ninja build. See logs above.")

		}
	}
}
