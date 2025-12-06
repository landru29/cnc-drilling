package main

import (
	"github.com/landru29/cnc-drilling/internal/configuration"
	"github.com/landru29/cnc-drilling/internal/geometry"
	"github.com/landru29/cnc-drilling/internal/operations/surfacer"
	"github.com/spf13/cobra"
)

func surfaceCommand(config *configuration.Config) *cobra.Command {
	var (
		surface geometry.Box
		method  surfacer.Method
	)

	output := &cobra.Command{
		Use:   "surface",
		Short: "Generate gcode to surface a rectangle area",
		RunE: func(cmd *cobra.Command, args []string) error {
			return surfacer.Process(surface, cmd.OutOrStdout(), cmd.OutOrStderr(), *config, method)
		},
	}

	output.Flags().VarP(&surface, "area", "a", "surface area [(minX, minY), (maxX, maxY)]")
	output.Flags().VarP(&method, "method", "m", "surfacing method (zigzag, spiral, spiral-from-center)")
	output.Flags().Float64VarP(&config.DeepXYPerTry, "deep-xy-per-try", "s", 5.0, "step between each pass (XY plane) in mm. to be calculated depending on the tool radius")
	output.Flags().Float64VarP(&config.DeepZ, "deep-z", "d", config.DeepZ, "engrave deep in millimeters")
	output.Flags().Float64VarP(&config.DeepZPerTry, "deep-z-per-try", "", config.DeepZPerTry, "max deep in millimeters during one try")

	return output
}
