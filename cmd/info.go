package main

import (
	"fmt"
	"io"
	"os"

	"github.com/landru29/cnc-drilling/internal/information"
	"github.com/spf13/cobra"
)

func infoCommand(files *[]string, config *information.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "info <filename.dxf>",
		Short: "Display informations about DXF",
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

				if err := information.Process(fileDesc, cmd.OutOrStdout(), *config); err != nil {
					return err
				}
			}

			return nil
		},
	}
}
