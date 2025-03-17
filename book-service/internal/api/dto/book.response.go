package dto

import "time"

type CreateResponse struct {
	ID          uint      `json:"idBook"`
	Title       string    `json:"tittle"`
	Author      string    `json:"author"`
	Description string    `json:"description"`
	Price       int       `json:"price"`
	Stock       int       `json:"stock"`
	CreatedAt   time.Time `json:"createdAt"`
}
