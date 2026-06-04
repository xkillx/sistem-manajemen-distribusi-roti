package dto

import (
	"time"

	sqlcdb "github.com/smdr/backend/internal/db"
)

type UserResponse struct {
	ID                 int64     `json:"id"`
	Username           string    `json:"username"`
	Name               string    `json:"name"`
	Phone              string    `json:"phone"`
	Role               string    `json:"role"`
	Active             bool      `json:"active"`
	MustChangePassword bool      `json:"must_change_password"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func UserToResponse(u sqlcdb.User) UserResponse {
	return UserResponse{
		ID:                 u.ID,
		Username:           u.Username,
		Name:               u.Name,
		Phone:              u.Phone,
		Role:               u.Role,
		Active:             u.Active,
		MustChangePassword: u.MustChangePassword,
		CreatedAt:          u.CreatedAt.Time,
		UpdatedAt:          u.UpdatedAt.Time,
	}
}

type ProductResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	SKU       string    `json:"sku"`
	Price     int64     `json:"price"`
	Active    bool      `json:"active"`
	CreatedBy int64     `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ProductToResponse(p sqlcdb.Product) ProductResponse {
	return ProductResponse{
		ID:        p.ID,
		Name:      p.Name,
		SKU:       p.Sku,
		Price:     p.Price,
		Active:    p.Active,
		CreatedBy: p.CreatedBy,
		CreatedAt: p.CreatedAt.Time,
		UpdatedAt: p.UpdatedAt.Time,
	}
}

type WarungResponse struct {
	ID                  int64      `json:"id"`
	Name                string     `json:"name"`
	OwnerName           string     `json:"owner_name"`
	Phone               string     `json:"phone"`
	Address             string     `json:"address"`
	Latitude            *float64   `json:"latitude"`
	Longitude           *float64   `json:"longitude"`
	Active              bool       `json:"active"`
	CreatedBy           int64      `json:"created_by"`
	AcquiredFromCheckin bool       `json:"acquired_from_checkin"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

func WarungToResponse(w sqlcdb.Warung) WarungResponse {
	var lat, lng *float64
	if w.Latitude.Valid {
		lat = &w.Latitude.Float64
	}
	if w.Longitude.Valid {
		lng = &w.Longitude.Float64
	}

	return WarungResponse{
		ID:                  w.ID,
		Name:                w.Name,
		OwnerName:           w.OwnerName,
		Phone:               w.Phone,
		Address:             w.Address,
		Latitude:            lat,
		Longitude:           lng,
		Active:              w.Active,
		CreatedBy:           w.CreatedBy,
		AcquiredFromCheckin: w.AcquiredFromCheckin,
		CreatedAt:           w.CreatedAt.Time,
		UpdatedAt:           w.UpdatedAt.Time,
	}
}

type VisitResponse struct {
	ID              int64      `json:"id"`
	SalesID         int64      `json:"sales_id"`
	WarungID        int64      `json:"warung_id"`
	Status          string     `json:"status"`
	CheckInLat      float64    `json:"check_in_lat"`
	CheckInLng      float64    `json:"check_in_lng"`
	CheckInAccuracy *float64   `json:"check_in_accuracy"`
	CheckInTime     time.Time  `json:"check_in_time"`
	DistanceMeters  *float64   `json:"distance_meters"`
	BusinessDate    string     `json:"business_date"`
	Notes           string     `json:"notes"`
	CancelledAt     *time.Time `json:"cancelled_at"`
	CancelledBy     *int64     `json:"cancelled_by"`
	CancelReason    string     `json:"cancel_reason"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

func VisitToResponse(v sqlcdb.Visit) VisitResponse {
	resp := VisitResponse{
		ID:           v.ID,
		SalesID:      v.SalesID,
		WarungID:     v.WarungID,
		Status:       v.Status,
		CheckInLat:   v.CheckInLat,
		CheckInLng:   v.CheckInLng,
		CheckInTime:  v.CheckInTime.Time,
		Notes:        v.Notes,
		CancelReason: v.CancelReason,
		CreatedAt:    v.CreatedAt.Time,
		UpdatedAt:    v.UpdatedAt.Time,
	}

	if v.CheckInAccuracy.Valid {
		resp.CheckInAccuracy = &v.CheckInAccuracy.Float64
	}
	if v.DistanceMeters.Valid {
		resp.DistanceMeters = &v.DistanceMeters.Float64
	}
	if v.CancelledAt.Valid {
		t := v.CancelledAt.Time
		resp.CancelledAt = &t
	}
	if v.CancelledBy.Valid {
		resp.CancelledBy = &v.CancelledBy.Int64
	}
	if v.BusinessDate.Valid {
		resp.BusinessDate = v.BusinessDate.Time.Format("2006-01-02")
	}

	return resp
}
