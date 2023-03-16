package booking

import (
	"net/http"
	"se/jwt-api/orm"
	"time"

	"github.com/gin-gonic/gin"
)

type BookingBody struct {
	UserID string
	CarID  string
	Start  time.Time
	End    time.Time
}

func BookingCar(c *gin.Context) {
	var json BookingBody
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	layout := "2006-01-02"
	start, err := time.Parse(layout, json.Start.Format(layout))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start time format"})
		return
	}
	end, err := time.Parse(layout, json.End.Format(layout))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end time format"})
		return
	}
	// before คือ ตัวแปร end นั้นมาก่อน start ไหม ส่วน Equal เช็คว่าเวลามันเท่ากันหรือไม่
	if end.Before(start) || end.Equal(start) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "End time must be greater than start time"})
		return
	}
	// Query the database using Gorm
	var results []orm.Booking
	orm.Db.Where("car_id = ? AND ((start <= ? AND end >= ?) OR (start >= ? AND start <= ?) OR (end >= ? AND end <= ?))", json.CarID, end, start, start, end, start, end).Find(&results)

	if len(results) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Booking conflict"})
		return
	}
	// Create the booking
	booking := orm.Booking{UserID: json.UserID, CarID: json.CarID, Start: start, End: end}
	if err := orm.Db.Create(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": booking})
}
