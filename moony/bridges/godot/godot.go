package godot

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
	ID ID `json:"id"`
	// object size
	Size *Size `json:"size"`
	// object transform props
	Transform *Transform `json:"transform"`
}

func SpawnObject(at Position, object Object) {
	// todo: convert it into transferable format
	return
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
