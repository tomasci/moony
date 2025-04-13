package utils

import (
	"fmt"
	"moony/moony/core/types"
)

func VecToStr[VectorType types.Vector2 | types.Vector3](vector VectorType, separator string) string {
	switch v := any(vector).(type) {
	case types.Vector2:
		return fmt.Sprintf("%d%s%d", v.X, separator, v.Y)
	case types.Vector3:
		return fmt.Sprintf("%d%s%d%s%d", v.X, separator, v.Y, separator, v.Z)
	default:
		return ""
	}
}
