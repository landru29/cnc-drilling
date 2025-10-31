package main

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/go-viper/mapstructure/v2"
	"github.com/kelseyhightower/envconfig"
	"github.com/landru29/cnc-drilling/internal/configuration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func mainCommand() (*cobra.Command, error) {
	var (
		files  []string
		config configuration.Config
	)

	viperConfiguration := viper.New()

	viperConfiguration.SetConfigName("drill")
	viperConfiguration.SetConfigType("yaml")
	viperConfiguration.AddConfigPath(".")
	viperConfiguration.AddConfigPath("$HOME/.cnc")
	viperConfiguration.AddConfigPath("/etc/cnc/")
	viperConfiguration.SetEnvPrefix("CNC")
	viperConfiguration.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viperConfiguration.AutomaticEnv()
	err := viperConfiguration.ReadInConfig()
	if err != nil {
		var notFoundErr viper.ConfigFileNotFoundError

		if !errors.As(err, &notFoundErr) {
			return nil, err
		}
	}

	if err := envconfig.Process("cnc", &config); err != nil {
		return nil, err
	}

	if err := viperConfiguration.Unmarshal(&config, func(decoderConfig *mapstructure.DecoderConfig) {
		decoderConfig.DecodeHook = mapstructure.ComposeDecodeHookFunc(
			configuration.DecodeOrigin,
		)
	}); err != nil {
		return nil, err
	}

	output := &cobra.Command{
		Use:   "cnc-router <filename.dxf>",
		Short: "Generate gcode from dxf",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			files = make([]string, len(args))

			copy(files, args)

			return nil
		},
	}

	output.PersistentFlags().Float64VarP(&config.Feed, "feed", "f", config.Feed, "speed in millimeters per minute")

	output.PersistentFlags().Float64VarP(&config.SecurityZ, "security-z", "z", config.SecurityZ, "Z security in millimeters")
	output.PersistentFlags().StringArrayVarP(&config.Layers, "layer", "l", config.Layers, "layer to filter")

	output.AddCommand(
		drillCommand(&files, &config),
		engraveCommand(&files, &config),
		infoCommand(&files, &config),
		configFileCommand(&config),
		surfaceCommand(&config),
	)

	return output, nil
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
