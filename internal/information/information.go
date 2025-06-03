package information

import (
	"fmt"
	"io"

	"github.com/yofu/dxf"
)

type Config struct {
	SpeedMillimeterPerMinute float64
	SecurityZ                float64
	Deepness                 float64
	Layer                    []string
}

func Process(in io.Reader, out io.Writer) error {
	drawing, err := dxf.FromReader(in)
	if err != nil {
		return err
	}

	defer func(closer io.Closer) {
		_ = closer.Close()
	}(drawing)

	if _, err := fmt.Fprintf(out, "%d layer(s) found:\n", len(drawing.Layers)); err != nil {
		return err
	}

	for _, layer := range drawing.Layers {
		isDefault := ""

		if drawing.CurrentLayer.Name() == layer.Name() {
			isDefault = " [default]"
		}

		if _, err := fmt.Fprintf(out, "\t* %s%s\n", layer.Name(), isDefault); err != nil {
			return err
		}
	}

	return nil
}
