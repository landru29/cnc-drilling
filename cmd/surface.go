package main

import (
	"github.com/landru29/cnc-drilling/internal/configuration"
	"github.com/landru29/cnc-drilling/internal/geometry"
	"github.com/landru29/cnc-drilling/internal/surfacer"
	"github.com/spf13/cobra"
)

func surfaceCommand(config *configuration.Config) *cobra.Command {
	var (
		surface geometry.Box
		step    float64
	)

	output := &cobra.Command{
		Use:   "surface",
		Short: "Generate gcode to surface a rectangle area",
		RunE: func(cmd *cobra.Command, args []string) error {
			return surfacer.Process(surface, step, cmd.OutOrStdout(), *config)
		},
	}

	output.Flags().VarP(&surface, "area", "a", "surface area [(minX, minY), (maxX, maxY)]")
	output.Flags().Float64VarP(&step, "step", "s", 5.0, "step between each pass (XY plane) in mm. to be calculated depending on the tool radius")
	output.Flags().Float64VarP(&config.Deepness, "deep", "d", config.Deepness, "engrave deep in millimeters")
	output.Flags().Float64VarP(&config.DeepStart, "deep-start", "", config.DeepStart, "initial deep in millimeters")
	output.Flags().Float64VarP(&config.DeepPerTry, "deep-per-try", "", config.DeepPerTry, "max deep in millimeters during one try")

	return output
}
