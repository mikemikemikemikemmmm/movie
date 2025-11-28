package sql

import (
	"backend/internal/config"
	"backend/internal/models"
	"backend/internal/structs"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var SqlDB *gorm.DB

func InitSQL() error {
	db, err := gorm.Open(postgres.Open(config.GetConfig().SqlUrl), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}
	SqlDB = db
	return nil
}
func CheckSqlReady() error {
	sqlDB, err := SqlDB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
func ReserveSeats(reserveData *structs.ReservePostData) error {
	if len(reserveData.SeatIds) == 0 {
		return fmt.Errorf("seatIDs is empty")
	}
	result := SqlDB.Model(&models.Seat{}).
		Where("id IN ?", reserveData.SeatIds).
		Updates(map[string]interface{}{
			"user_id": reserveData.UserId,
			"status":  models.Reserved,
		})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func RollbackReserveSeats(reserveData *structs.ReservePostData) error {
	if len(reserveData.SeatIds) == 0 {
		return fmt.Errorf("seatIDs is empty")
	}
	result := SqlDB.Model(&models.Seat{}).
		Where("id IN ?", reserveData.SeatIds).
		Updates(map[string]interface{}{
			"user_id": nil,
			"status":  models.Available,
		})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func GetAllSeats() ([]models.Seat, error) {
	var seats []models.Seat
	result := SqlDB.Find(&seats)
	if result.Error != nil {
		log.Print(result.Error)
		return nil, result.Error
	}
	return seats, nil
}
