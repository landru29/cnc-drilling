package main

import (
	"fmt"
	"io"
	"os"

	"github.com/landru29/cnc-drilling/internal/configuration"
	"github.com/landru29/cnc-drilling/internal/information"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func infoCommand(files *[]string, config *configuration.Config) *cobra.Command {
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

func configFileCommand(config *configuration.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "save-config",
		Short: "save config gile",
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := yaml.Marshal(&config)
			if err != nil {
				return err
			}

			// Sauvegarder dans un fichier
			return os.WriteFile("config.yaml", data, 0644)
		},
	}
}
