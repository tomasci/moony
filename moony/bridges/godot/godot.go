package godot

import (
	"github.com/google/uuid"
	"math/rand/v2"
	"moony/moony/bridges/godot/materials"
)

// godot bridge concept

type ID string

type XYZ struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type Size XYZ
type Position XYZ
type Rotation XYZ
type Scale XYZ

type Transform struct {
	Position *Position `json:"position"`
	Rotation *Rotation `json:"rotation"`
	Scale    *Scale    `json:"scale"`
}

type Object struct {
	// link to object
	ID uuid.UUID `json:"id"`
	// object size
	Size *Size `json:"size"`
	// object transform props
	Transform *Transform `json:"transform"`
}

// SpawnObject – this function like any other function here is to standardize server response for the same action called from different places
func SpawnObject(at Position, object Object, material materials.Material) []any {
	return []any{object, at, material}
}

func Move(to Position, object Object) {
	return
}

func RotateObject(to Rotation, object Object) {
	return
}

func GetObject(id ID) {
	return
}

func RemoveObject(id ID) {
	return
}

func RandomFloat(min float64, max float64) float64 {
	return min + rand.Float64()*(max-min)
}
