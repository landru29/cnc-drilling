package main

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"

	"github.com/landru29/cnc-drilling/internal/information"
	"github.com/spf13/cobra"
)

func mainCommand() *cobra.Command {
	var (
		files  []string
		config information.Config
	)

	output := &cobra.Command{
		Use:   "cnc-router <filename.dxf>",
		Short: "Generate gcode from dxf",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("no file to process")
			}

			files = make([]string, len(args))

			copy(files, args)

			return nil
		},
	}

	output.PersistentFlags().Float64VarP(&config.SpeedMillimeterPerMinute, "feed", "f", 60, "speed in millimeters per minute")
	output.PersistentFlags().Float64VarP(&config.SecurityZ, "security-z", "z", 10, "Z security in millimeters")
	output.PersistentFlags().StringArrayVarP(&config.Layers, "layer", "l", nil, "layer to filter")

	output.AddCommand(
		drillCommand(&files, &config),
		engraveCommand(&files, &config),
		infoCommand(&files),
	)

	return output
}

func header(writer io.Writer, filename string) error {
	if _, err := fmt.Fprintf(
		writer,
		"; File: %s\n",
		filepath.Base(filename),
	); err != nil {
		return err
	}

	return nil
}

func footer(writer io.Writer, filename string) error {
	if _, err := fmt.Fprintf(writer, ";\n; End of file: %s\n\n", filepath.Base(filename)); err != nil {
		return err
	}

	return nil
}
