package main

import (
	"fmt"
	"io"
	"os"

	"github.com/landru29/cnc-drilling/internal/configuration"
	"github.com/landru29/cnc-drilling/internal/operations/pather"
	"github.com/spf13/cobra"
)

func pathCommand(files *[]string, config *configuration.Config) *cobra.Command {
	output := &cobra.Command{
		Use:   "path <filename.dxf>",
		Short: "Generate gcode to follow a path from dxf",
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

				if err := pather.Process(
					fileDesc, cmd.OutOrStdout(),
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

	output.Flags().Float64VarP(&config.DeepZ, "deep-z", "d", config.DeepZ, "engrave deep in millimeters")
	output.Flags().Float64VarP(&config.DeepZPerTry, "deep-z-per-try", "", config.DeepZPerTry, "max deep in millimeters during one try")
	output.Flags().VarP(&config.Origin, "origin", "o", "shift origin")

	return output
}
