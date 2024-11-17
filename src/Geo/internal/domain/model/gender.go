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
