package models

import "gorm.io/gorm"

type SeatStatus string

const (
	Available SeatStatus = "available"
	Reserved  SeatStatus = "reserved"
)

type Seat struct {
	gorm.Model
	X      int        `json:"x" gorm:"uniqueIndex:idx_seat_xy"` // 加上 unique index 名稱
	Y      int        `json:"y" gorm:"uniqueIndex:idx_seat_xy"` // 同一個 index 名稱表示組合唯一
	UserID *int       `json:"user_id"`
	Status SeatStatus `json:"status"`
}
