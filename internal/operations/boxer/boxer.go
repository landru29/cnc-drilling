package boxer

import (
	"io"

	"github.com/landru29/cnc-drilling/internal/configuration"
	"github.com/landru29/cnc-drilling/internal/geometry"
	"github.com/landru29/cnc-drilling/internal/operations/edger"
	"github.com/landru29/cnc-drilling/internal/operations/surfacer"
)

func Process(box geometry.Box, out io.Writer, info io.Writer, config configuration.Config, method surfacer.Method, edgeDeepZ float64) error {
	if err := surfacer.Process(box, out, info, config, method); err != nil {
		return err
	}

	config.DeepZ += edgeDeepZ
	if err := edger.Process(box, out, info, config, true); err != nil {
		return err
	}

	return nil
}
