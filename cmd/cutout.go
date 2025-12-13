package main

import (
	"github.com/landru29/cnc-drilling/internal/configuration"
	"github.com/landru29/cnc-drilling/internal/geometry"
	"github.com/landru29/cnc-drilling/internal/operations/outcutter"
	"github.com/spf13/cobra"
)

func cutOutCommand(config *configuration.Config) *cobra.Command {
	var (
		box geometry.Box
	)

	output := &cobra.Command{
		Use:   "cutout",
		Short: "Generate gcode to cut out a rectangle",
		RunE: func(cmd *cobra.Command, args []string) error {
			box.Offset(config.Origin.Value)

			return outcutter.Process(box, cmd.OutOrStdout(), cmd.OutOrStderr(), *config)
		},
	}

	output.Flags().VarP(&box, "box", "b", "starting box area [(minX, minY), (maxX, maxY)]. The final part will be lowered by deep-xy on each side")
	output.Flags().Float64VarP(&config.DeepXYPerTry, "deep-xy-per-try", "s", 0.5, "step between each pass (XY plane) in mm. to be calculated depending on the tool radius")
	output.Flags().Float64VarP(&config.DeepXY, "deep-xy", "d", config.DeepXY, "engrave deep in millimeters")

	return output
}
