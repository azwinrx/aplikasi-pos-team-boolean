package dto

import (
	"time"
)

// ReservationListItem response for reservation list endpoint
type ReservationListItem struct {
	ID              int64     `json:"id"`
	CustomerName    string    `json:"customer_name"`
	CustomerPhone   string    `json:"customer_phone"`
	ReservationDate string    `json:"reservation_date"` // format: YYYY-MM-DD
	CheckIn         time.Time `json:"check_in"`
	CheckOut        time.Time `json:"check_out"`
	TotalPrice      float64   `json:"total_price"`
}

// ReservationDetail response for reservation detail endpoint
type ReservationDetail struct {
	ID              int64     `json:"id"`
	CustomerName    string    `json:"customer_name"`
	CustomerPhone   string    `json:"customer_phone"`
	TableNumber     string    `json:"table_number"`
	PaxNumber       int       `json:"pax_number,omitempty"`
	ReserveDate     string    `json:"reserve_date"`
	ReservationTime time.Time `json:"reservation_time"`
	ReservedEnd     time.Time `json:"reserved_end"`
	DepositFee      float64   `json:"deposit_fee,omitempty"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	TotalPrice      float64   `json:"total_price,omitempty"`
}

// ReservationUpdateRequest for update/add reservation
type ReservationUpdateRequest struct {
	TableNumber     string  `json:"table_number"`
	PaxNumber       int     `json:"pax_number"`
	ReserveDate     string  `json:"reserve_date"`     // format: YYYY-MM-DD
	ReservationTime string  `json:"reservation_time"` // format: HH:MM (24h)
	DurationMinutes int     `json:"duration_minutes"` // for reserved_end calculation
	DepositFee      float64 `json:"deposit_fee"`
	Status          string  `json:"status"`
	CustomerName    string  `json:"customer_name"`
	CustomerPhone   string  `json:"customer_phone"`
}

// ReservationCreateRequest for create reservation
type ReservationCreateRequest struct {
	TableNumber     string  `json:"table_number"`
	PaxNumber       int     `json:"pax_number"`
	ReserveDate     string  `json:"reserve_date"`     // format: YYYY-MM-DD
	ReservationTime string  `json:"reservation_time"` // format: HH:MM (24h)
	DurationMinutes int     `json:"duration_minutes"` // for reserved_end calculation
	DepositFee      float64 `json:"deposit_fee"`
	Status          string  `json:"status"`
	CustomerName    string  `json:"customer_name"`
	CustomerPhone   string  `json:"customer_phone"`
}
