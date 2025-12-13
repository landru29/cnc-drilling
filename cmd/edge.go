package main

import (
	"github.com/landru29/cnc-drilling/internal/configuration"
	"github.com/landru29/cnc-drilling/internal/geometry"
	"github.com/landru29/cnc-drilling/internal/operations/edger"
	"github.com/spf13/cobra"
)

func edgeCommand(config *configuration.Config) *cobra.Command {
	var (
		surface          geometry.Box
		counterclockwise bool = false
	)

	output := &cobra.Command{
		Use:   "edge",
		Short: "Generate gcode to machine the edge of a box",
		RunE: func(cmd *cobra.Command, args []string) error {
			surface.Offset(config.Origin.Value)

			return edger.Process(surface, cmd.OutOrStdout(), cmd.OutOrStderr(), *config, !counterclockwise)
		},
	}

	output.Flags().VarP(&surface, "box", "b", "box projection on XY plane [(minX, minY), (maxX, maxY)]")
	output.Flags().Float64VarP(&config.DeepXYPerTry, "deep-xy-per-try", "s", 5.0, "step between each pass (XY plane) in mm. to be calculated depending on the tool radius")
	output.Flags().Float64VarP(&config.DeepZ, "deep-z", "d", config.DeepZ, "surface deep in millimeters")
	output.Flags().Float64VarP(&config.DeepZPerTry, "deep-z-per-try", "", config.DeepZPerTry, "max deep in millimeters during one try")
	output.Flags().BoolVarP(&counterclockwise, "counterclockwise", "c", counterclockwise, "counterclockwise edge machining")

	return output
}
