package main

import (
	"fmt"
	"io"
	"os"

	"github.com/landru29/cnc-drilling/internal/driller"
	"github.com/spf13/cobra"
)

func drillCommand(files *[]string, speedMillimeterPerMinute *float64, securityZ *float64) *cobra.Command {
	var (
		drillingDeep float64
	)

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

				if err := driller.Process(fileDesc, cmd.OutOrStdout(), *speedMillimeterPerMinute, drillingDeep, *securityZ); err != nil {
					return err
				}

				if err := footer(cmd.OutOrStdout(), file); err != nil {
					return err
				}
			}

			return nil
		},
	}

	output.Flags().Float64VarP(&drillingDeep, "deep", "d", 5, "drilling deep in millimeters")

	return output
}
