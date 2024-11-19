package model

import "encoding/json"

type Profile struct {
	ID          int64    `json:"id"`
	Age         int16    `json:"age"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Photos      []string `json:"photos"`
}

func (p Profile) MarshalBinary() (data []byte, err error) {
	return json.Marshal(p)
}
