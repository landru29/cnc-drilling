package main

import (
	"fmt"
	"io"
	"os"

	"github.com/landru29/cnc-drilling/internal/configuration"
	"github.com/landru29/cnc-drilling/internal/driller"
	"github.com/spf13/cobra"
)

func drillCommand(files *[]string, config *configuration.Config) *cobra.Command {
	output := &cobra.Command{
		Use:   "drill <filename.dxf>",
		Short: "Generate gcode to drill from dxf",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, file := range *files {
				fileDesc, err := os.Open(file)
				if err != nil {
					return fmt.Errorf("%s: %w", file, err)
				}

				defer func(closer io.Closer) {
					_ = closer.Close()
				}(fileDesc)

				if err := header(cmd.OutOrStdout(), file); err != nil {
					return err
				}

				if err := driller.Process(
					fileDesc,
					cmd.OutOrStdout(),
					*config,
				); err != nil {
					return err
				}

				if err := footer(cmd.OutOrStdout(), file); err != nil {
					return err
				}
			}

			return nil
		},
	}

	output.Flags().Float64VarP(&config.Deepness, "deep", "d", config.Deepness, "drilling deep in millimeters")
	output.Flags().Float64VarP(&config.DeepPerTry, "deep-per-try", "", config.DeepPerTry, "max deep in millimeters during one try")
	output.Flags().VarP(&config.Origin, "origin", "o", "shift origin")

	return output
}
