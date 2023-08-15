package repository

import (
	"time"

	"github.com/pradeepj4u/bookings/cmd/models"
)

type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(r models.RoomRestriction) error
	SearchAvailibilityForDateRangeByRoomId(roomId int, startDate, endDate time.Time) (bool, error)
	SearchAvailibilityForAllRooms(startDate, endDate time.Time) ([]models.Room, error)
	GetRoomById(id int) (models.Room, error)
}
