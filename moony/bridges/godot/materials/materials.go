package materials

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

type Material int

const (
	Red Material = iota
	Green
	Blue
	MaterialCount
)

func (m Material) String() string {
	names := []string{Red: "red", Green: "green", Blue: "blue"}
	if m < Red || m > MaterialCount {
		return "unknown"
	}

	return names[m]
}

func (m Material) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.String())
}

func (m *Material) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	switch s {
	case "red":
		*m = Red
	case "green":
		*m = Green
	case "blue":
		*m = Blue
	default:
		err = fmt.Errorf("invalid material: %s", s)
	}

	return nil
}

func Random() Material {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return Material(r.Intn(int(MaterialCount)))
}
