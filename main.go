package main

import (
	"io/ioutil"
	"os"

	"github.com/GoogleContainerTools/kpt-functions-catalog/functions/go/starlark/third_party/sigs.k8s.io/kustomize/kyaml/fn/runtime/starlark"
	"sigs.k8s.io/kustomize/kyaml/fn/framework"
	"sigs.k8s.io/kustomize/kyaml/fn/framework/command"
)

type StarlarkProcessor struct {
	sourceFile string
}

func (sp *StarlarkProcessor) Process(rl *framework.ResourceList) error {
	source, err := ioutil.ReadFile(sp.sourceFile)
	if err != nil {
		return err
	}

	starFltr := &starlark.SimpleFilter{
		Name:           "starlark",
		Program:        string(source),
		FunctionConfig: rl.FunctionConfig,
	}
	rl.Items, err = starFltr.Filter(rl.Items)
	if err != nil {
		rl.Result = &framework.Result{
			Name: "starlark",
			Items: []framework.ResultItem{
				{
					Message:  err.Error(),
					Severity: framework.Error,
				},
			},
		}
		return rl.Result
	}
	return nil
}

func main() {
	sp := StarlarkProcessor{}
	cmd := command.Build(&sp, command.StandaloneEnabled, false)
	cmd.Flags().StringVar(&sp.sourceFile, "source", "", "the filename of the starlark source.")
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
