/*
Copyright © 2024 Jozu.com
*/
package build

import (
	"context"
	"fmt"
	"strings"

	"kitops/pkg/lib/constants"
	"kitops/pkg/lib/filesystem"
	"kitops/pkg/lib/repo"
	"kitops/pkg/output"

	"github.com/spf13/cobra"
	"oras.land/oras-go/v2/registry"
)

var (
	shortDesc = `Build a model`
	longDesc  = `Build a model TODO`
)

type buildOptions struct {
	modelFile   string
	contextDir  string
	configHome  string
	storageHome string
	fullTagRef  string
	modelRef    *registry.Reference
	extraRefs   []string
}

func BuildCommand() *cobra.Command {
	opts := &buildOptions{}

	cmd := &cobra.Command{
		Use:   "build",
		Short: shortDesc,
		Long:  longDesc,
		Run:   runCommand(opts),
	}

	cmd.Flags().StringVarP(&opts.modelFile, "file", "f", "", "Path to the model file")
	cmd.Flags().StringVarP(&opts.fullTagRef, "tag", "t", "", "Tag for the model. Example: -t registry/repository:tag1,tag2")
	cmd.Args = cobra.ExactArgs(1)
	return cmd
}

func runCommand(opts *buildOptions) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		err := opts.complete(cmd.Context(), args)
		if err != nil {
			output.Fatalf("Failed to process configuration: %s", err)
			return
		}
		err = RunBuild(cmd.Context(), opts)
		if err != nil {
			output.Fatalf("Failed to build model kit: %s", err)
			return
		}
	}
}

func (opts *buildOptions) complete(ctx context.Context, args []string) error {
	opts.contextDir = args[0]

	if opts.modelFile == "" {
		opts.modelFile = filesystem.FindKitfileInPath(opts.contextDir)
	}

	configHome, ok := ctx.Value(constants.ConfigKey{}).(string)
	if !ok {
		return fmt.Errorf("default config path not set on command context")
	}
	opts.configHome = configHome
	opts.storageHome = constants.StoragePath(opts.configHome)

	if opts.fullTagRef != "" {
		modelRef, extraRefs, err := repo.ParseReference(opts.fullTagRef)
		if err != nil {
			return fmt.Errorf("failed to parse reference %s: %w", opts.fullTagRef, err)
		}
		opts.modelRef = modelRef
		opts.extraRefs = extraRefs
	}
	printConfig(opts)
	return nil
}

func printConfig(opts *buildOptions) {
	output.Debugf("Using storage path: %s", opts.storageHome)
	output.Debugf("Context dir: %s", opts.contextDir)
	output.Debugf("Model file: %s", opts.modelFile)
	if opts.modelRef != nil {
		output.Debugf("Building %s", opts.modelRef.String())
	} else {
		output.Debugln("No tag or reference specified")
	}
	if len(opts.extraRefs) > 0 {
		output.Debugf("Additional tags: %s", strings.Join(opts.extraRefs, ", "))
	}
}