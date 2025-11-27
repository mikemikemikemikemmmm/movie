package structs

type ReservePostData struct {
	UserId  int   `json:"user_id" binding:"required"`
	SeatIds []int `json:"seat_ids" binding:"required,dive,gt=0"` // 每個數字必須大於0
}
