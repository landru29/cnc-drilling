package main

import (
	"github.com/landru29/cnc-drilling/internal/configuration"
	"github.com/landru29/cnc-drilling/internal/geometry"
	"github.com/landru29/cnc-drilling/internal/operations/boxer"
	"github.com/landru29/cnc-drilling/internal/operations/surfacer"
	"github.com/spf13/cobra"
)

func boxCommand(config *configuration.Config) *cobra.Command {
	var (
		surface   geometry.Box
		method    surfacer.Method
		edgeDeepZ float64 = 51.0
	)

	output := &cobra.Command{
		Use:   "box",
		Short: "Generate gcode to machine a box",
		RunE: func(cmd *cobra.Command, args []string) error {
			return boxer.Process(surface, cmd.OutOrStdout(), cmd.OutOrStderr(), *config, method, edgeDeepZ)
		},
	}

	output.Flags().VarP(&surface, "box", "b", "box projection on XY plane [(minX, minY), (maxX, maxY)]")
	output.Flags().VarP(&method, "method", "m", "surfacing method (zigzag, spiral, spiral-from-center)")
	output.Flags().Float64VarP(&config.DeepXYPerTry, "deep-xy-per-try", "s", 5.0, "step between each pass (XY plane) in mm. to be calculated depending on the tool radius")
	output.Flags().Float64VarP(&config.DeepZ, "deep-z", "d", config.DeepZ, "surface deep in millimeters")
	output.Flags().Float64VarP(&edgeDeepZ, "edge-deep-z", "e", edgeDeepZ, "edge deep in millimeters")
	output.Flags().Float64VarP(&config.DeepZPerTry, "deep-z-per-try", "", config.DeepZPerTry, "max deep in millimeters during one try")

	return output
}
