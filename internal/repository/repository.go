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
	GetUseById(id int) (models.User, error)
	UpdateUser(usr models.User) error
	Authenticate(email, testPassword string) (int, string, error)
	AllReservations() ([]models.Reservation, error)
	AllNewReservations() ([]models.Reservation, error)
	GetReservationById(id int) (models.Reservation, error)
	UpdateReservation(res models.Reservation) error
	DeleteReservation(id int) error
	UpdateProcessed(id, processed int) error
	AllRooms() ([]models.Room, error)
	GetRestrictionsByRoomForDate(roomId int, startDate, endDate time.Time) ([]models.RoomRestriction, error)
	InsertRoomRestrictionForRoom(roomId int, startDate time.Time) error
	DeleteRoomRestrictionById(id int) error
}
