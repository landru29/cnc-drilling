package information

import (
	"fmt"
	"io"
	"math"

	"github.com/yofu/dxf"
)

type Config struct {
	SpeedMillimeterPerMinute float64
	SecurityZ                float64
	Deepness                 float64
	DeepPerTry               float64
	Layers                   []string
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

// TryDeeps is the set of deeps during all tries.
func (c Config) TryDeeps() []float64 {
	if c.DeepPerTry <= 0 {
		return []float64{c.Deepness}
	}

	output := make([]float64, int(math.Ceil(c.Deepness/c.DeepPerTry)))

	maxFullTry := int(math.Floor(c.Deepness / c.DeepPerTry))

	output[len(output)-1] = math.Mod(c.Deepness, c.DeepPerTry) + float64(maxFullTry)*c.DeepPerTry

	for index := range maxFullTry {
		output[index] = float64(index+1) * c.DeepPerTry
	}

	return output
}
