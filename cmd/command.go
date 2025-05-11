package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/landru29/cnc-drilling/internal/driller"

	"github.com/spf13/cobra"
)

func mainCommand() *cobra.Command {
	var (
		speedMillimeterPerMinute float64
		drillingDeep             float64
		securityZ                float64
	)

	output := &cobra.Command{
		Use:   "cnc-drilling <filename.dxf>",
		Short: "Generate gcode to drill from dxf",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("no file to process")
			}

			for _, file := range args {
				filename := filepath.Base(file)
				fileDesc, err := os.Open(file)
				if err != nil {
					return fmt.Errorf("%s: %w", file, err)
				}

				defer func(closer io.Closer) {
					_ = closer.Close()
				}(fileDesc)

				if _, err := fmt.Fprintf(cmd.OutOrStdout(), "; File: %s\n", filename); err != nil {
					return err
				}

				if err := driller.Process(fileDesc, cmd.OutOrStdout(), speedMillimeterPerMinute, drillingDeep, securityZ); err != nil {
					return err
				}

				if _, err := fmt.Fprintf(cmd.OutOrStdout(), "\n; End of file: %s\n\n", filename); err != nil {
					return err
				}
			}

			return nil
		},
	}

	output.Flags().Float64VarP(&speedMillimeterPerMinute, "speed", "s", 60, "speed in millimeters per minute")
	output.Flags().Float64VarP(&drillingDeep, "deep", "d", 5, "drilling deep in millimeters")
	output.Flags().Float64VarP(&securityZ, "security-z", "z", 10, "Z security in millimeters")

	return output
}
