package model

import (
	"encoding/json"
	"fmt"
)

type Gender int

const (
	M Gender = iota
	F
)

func (g Gender) String() string {
	if g == M {
		return "M"
	}
	return "F"
}

// todo: why pointer?
func (g *Gender) UnmarshalJSON(data []byte) error {
	var s string

	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	if s == "M" {
		*g = M
		return nil
	}

	if s == "F" {
		*g = F
		return nil
	}

	return fmt.Errorf("undefined value: %s", s)
}

func (g Gender) MarshalJSON() ([]byte, error) {
	return json.Marshal(g.String())
}

func (g Gender) Invert() Gender {
	if g == M {
		return F
	}
	return M
}
