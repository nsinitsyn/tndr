package model

type Profile struct {
	ID int64
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
