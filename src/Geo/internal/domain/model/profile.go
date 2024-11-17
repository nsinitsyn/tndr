package model

import "encoding/json"

type Profile struct {
	ID          int64    `json:"id"`
	Age         int16    `json:"age"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Photos      []string `json:"photos"`
}

type ProfilesSortable []Profile

func (p ProfilesSortable) Len() int {
	return len(p)
}

func (p ProfilesSortable) Less(i, j int) bool {
	return p[i].ID < p[j].ID
}

func (p ProfilesSortable) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p Profile) MarshalBinary() (data []byte, err error) {
	return json.Marshal(p)
}
