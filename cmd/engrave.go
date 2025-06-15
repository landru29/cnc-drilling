package main

import (
	"fmt"
	"io"
	"os"

	"github.com/landru29/cnc-drilling/internal/engraver"
	"github.com/landru29/cnc-drilling/internal/information"
	"github.com/spf13/cobra"
)

func engraveCommand(files *[]string, config *information.Config) *cobra.Command {
	output := &cobra.Command{
		Use:   "engrave <filename.dxf>",
		Short: "Generate gcode to engrave from dxf",
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

				if err := engraver.Process(
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

	output.Flags().Float64VarP(&config.Deepness, "deep", "d", 1, "engrave deep in millimeters")
	output.Flags().Float64VarP(&config.DeepPerTry, "deep-per-try", "", 5, "max deep in millimeters during one try")
	output.Flags().VarP(&config.Origin, "origin", "o", "shift origin")

	return output
}
