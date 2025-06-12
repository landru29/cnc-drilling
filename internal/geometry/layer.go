package geometry

import (
	"slices"

	"github.com/yofu/dxf/entity"
)

// FilterEntities filters the entities by layers (if specified).
func FilterEntities(entities entity.Entities, layers ...string) entity.Entities {
	output := entity.Entities{}

	for _, entity := range entities {
		if len(layers) > 0 && !slices.Contains(layers, entity.Layer().Name()) {
			continue
		}

		output = append(output, entity)
	}

	return output
}
