package messaging

import "tinder-geo/internal/domain/model"

type ProfileDto struct {
	ID          int64        `json:"id"`
	Gender      model.Gender `json:"gender"`
	Age         int16        `json:"age"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Photos      []string     `json:"photos"`
}
