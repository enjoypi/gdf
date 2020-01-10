package main

import (
	"github.com/enjoypi/gdf/generators"
"k8s.io/gengo/args"

"github.com/spf13/pflag"
"k8s.io/klog"
)

func main() {
	klog.InitFlags(nil)
	arguments := args.Default()

	// Override defaults.
	arguments.OutputFileBaseName = "dirtyflag_generated"

	// Custom args.
	customArgs := &generators.CustomArgs{}
	pflag.CommandLine.StringSliceVar(&customArgs.BoundingDirs, "bounding-dirs", customArgs.BoundingDirs,
		"Comma-separated list of import paths which bound the types for which deep-copies will be generated.")
	arguments.CustomArgs = customArgs

	// Run it.
	if err := arguments.Execute(
		generators.NameSystems(),
		generators.DefaultNameSystem(),
		generators.Packages,
	); err != nil {
		klog.Fatalf("Error: %v", err)
	}
	klog.V(2).Info("Completed successfully.")
}
