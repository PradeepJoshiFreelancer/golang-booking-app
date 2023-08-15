package dbrepo

import (
	"context"
	"time"

	"github.com/pradeepj4u/bookings/cmd/models"
)

func (m *postgresDbRepo) AllUsers() bool {
	return true
}

// Insert the reservation in the table
func (m *postgresDbRepo) InsertReservation(res models.Reservation) (int, error) {
	cntx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newId int
	statement := `insert into reservations (first_name
											, last_name
											, email, phone
											, start_date
											, end_date
											, room_id
											, created_at
											, updated_at)
	values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`

	err := m.DB.QueryRowContext(cntx, statement,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomId,
		time.Now(),
		time.Now(),
	).Scan(&newId)
	if err != nil {
		return 0, err
	}
	return newId, nil
}

// Inserts the room restriciton into the database
func (m *postgresDbRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	cntx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	statement := `insert into room_restrictions (start_date
											, end_date
											, room_id
											, reservation_id
											, restriction_id
											, created_at
											, updated_at)
	values ($1, $2, $3, $4, $5, $6, $7)`

	_, err := m.DB.ExecContext(cntx, statement,
		r.StartDate,
		r.EndDate,
		r.RoomId,
		r.ReservationId,
		r.RestrictionId,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}
	return nil
}

// Search for a room availibility for a specific date range.
func (m *postgresDbRepo) SearchAvailibilityForDateRangeByRoomId(roomId int, startDate, endDate time.Time) (bool, error) {
	cntx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rowCount int

	statement := `
				select 
					count(*) 
				from 
					room_restrictions 
				where 
					room_id = $1 and
					$2 < end_date and 
					$3 > start_date`

	err := m.DB.QueryRowContext(cntx, statement, roomId, startDate, endDate).Scan(&rowCount)

	if err != nil {
		return false, err
	}
	if rowCount == 0 {
		return true, nil
	}

	return false, nil
}

// Search for a available rooms within a specific date range.
func (m *postgresDbRepo) SearchAvailibilityForAllRooms(startDate, endDate time.Time) ([]models.Room, error) {
	cntx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room

	stmt := `
			select 
			id, room_name  
		from
			rooms
		where 
			id not in
				(select rr.room_id from 
					room_restrictions rr 
				where 
					$1 < rr.end_date and 
					$2 > rr.start_date)`

	rows, err := m.DB.QueryContext(cntx, stmt, startDate, endDate)
	if err != nil {
		return rooms, err
	}

	for rows.Next() {
		var room models.Room
		rows.Scan(&room.ID, &room.RoomName)
		rooms = append(rooms, room)
	}

	if rows.Err() != nil {
		return rooms, err
	}
	return rooms, nil

}

// Search for a Room by Id

func (m *postgresDbRepo) GetRoomById(id int) (models.Room, error) {
	cntx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room models.Room
	stmt := `
		select 
			r.id
			, r.room_name 
			, r.created_at 
			, r.updated_at  
		from rooms r
		where id = $1 
	`
	err := m.DB.QueryRowContext(cntx, stmt, id).Scan(&room.ID, &room.RoomName, &room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		return room, err
	}

	return room, nil
}
