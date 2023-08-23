package dbrepo

import (
	"context"
	"errors"
	"time"

	"github.com/pradeepj4u/bookings/cmd/models"
	"golang.org/x/crypto/bcrypt"
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

// Get the user for a given id
func (m *postgresDbRepo) GetUseById(id int) (models.User, error) {
	cntx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
			select u.id 
				, u.first_name
				, u.last_name 
				, u.email 
				, u.password
				, u.access_level 
				, u.created_at 
				, u.updated_at 
			from users u
			where u.id = $1 `

	row := m.DB.QueryRowContext(cntx, stmt, id)

	var usr models.User

	err := row.Scan(&usr.ID, &usr.LastName, &usr.Email, &usr.Password, &usr.AccessLevel, &usr.CreatedAt, &usr.UpdatedAt)

	if err != nil {
		return usr, err
	}
	return usr, nil
}

// Update the user information
func (m *postgresDbRepo) UpdateUser(usr models.User) error {
	cntx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
	update users set  u.first_name = $2 
					, u.last_name = $3
					, u.email = $4
					, u.access_level = $5
					,u.updated_at = $6
				where u.id = $1`

	_, err := m.DB.ExecContext(cntx, stmt,
		usr.ID,
		usr.FirstName,
		usr.LastName,
		usr.Email,
		usr.AccessLevel,
		time.Now())
	if err != nil {
		return err
	}
	return nil
}

// Authenticate if the given password is correct or now, it compares the hased values for password
func (m *postgresDbRepo) Authenticate(email, testPassword string) (int, string, error) {
	cntx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	row := m.DB.QueryRowContext(cntx, "select u.id , u.password from users u where u.email = $1", email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("Password mismatch")
	} else if err != nil {
		return 0, "", err
	}
	return id, hashedPassword, nil

}

// Reutuns List of all repositories
func (m *postgresDbRepo) AllReservations() ([]models.Reservation, error) {
	cntx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation
	stmt := `
		select 
		r.id 
		, r.first_name 
		, r.last_name 
		, r.email 
		, r.phone 
		, r.start_date 
		, r.end_date 
		, r.room_id 
		, r.created_at 
		, r.updated_at 
		, r.processed
		, rm.room_name 
	from reservations r,
		rooms rm
	where r.room_id = rm.id 
	order by r.start_date 
	`
	rows, err := m.DB.QueryContext(cntx, stmt)
	if err != nil {
		return reservations, err
	}
	defer rows.Close()
	for rows.Next() {
		var i models.Reservation

		err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomId,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Processed,
			&i.Room.RoomName,
		)
		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, i)
	}
	if rows.Err() != nil {
		return reservations, rows.Err()
	}
	return reservations, nil
}

// Reutuns List of all New repositories
func (m *postgresDbRepo) AllNewReservations() ([]models.Reservation, error) {
	cntx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation
	stmt := `
	select 
		r.id 
	, r.first_name 
	, r.last_name 
	, r.email 
	, r.phone 
	, r.start_date 
	, r.end_date 
	, r.room_id 
	, r.created_at 
	, r.updated_at 
	, r.processed 
	, rm.room_name 
	from reservations r,
	rooms rm
	where r.room_id = rm.id
	and 	r.processed = 0
	order by r.start_date 
	`
	rows, err := m.DB.QueryContext(cntx, stmt)
	if err != nil {
		return reservations, err
	}
	defer rows.Close()
	for rows.Next() {
		var i models.Reservation

		err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomId,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Processed,
			&i.Room.RoomName,
		)
		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, i)
	}
	if rows.Err() != nil {
		return reservations, rows.Err()
	}
	return reservations, nil
}

// Reutuns one reservation for an Id of all New repositories
func (m *postgresDbRepo) GetReservationById(id int) (models.Reservation, error) {
	cntx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservation models.Reservation

	stmt := `
	select 
		r.id 
		, r.first_name 
		, r.last_name 
		, r.email 
		, r.phone 
		, r.start_date 
		, r.end_date 
		, r.room_id 
		, r.created_at 
		, r.updated_at 
		, r.processed 
		, rm.room_name 
	from reservations r,
	rooms rm
	where r.room_id = rm.id
	and   r.id = $1
	`
	row := m.DB.QueryRowContext(cntx, stmt, id)

	err := row.Scan(
		&reservation.ID,
		&reservation.FirstName,
		&reservation.LastName,
		&reservation.Email,
		&reservation.Phone,
		&reservation.StartDate,
		&reservation.EndDate,
		&reservation.RoomId,
		&reservation.CreatedAt,
		&reservation.UpdatedAt,
		&reservation.Processed,
		&reservation.Room.RoomName,
	)

	if err != nil {
		return reservation, err
	}

	return reservation, nil
}

// Update the reservation information
func (m *postgresDbRepo) UpdateReservation(res models.Reservation) error {
	cntx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
	update reservations set  first_name = $2 
					, last_name = $3
					, email = $4
					, phone = $5
					,updated_at = $6
				where id = $1`

	_, err := m.DB.ExecContext(cntx, stmt,
		res.ID,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		time.Now())

	if err != nil {
		return err
	}
	return nil
}

// Delete the reservation information
func (m *postgresDbRepo) DeleteReservation(id int) error {
	cntx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `delete from reservations where id=$1`

	_, err := m.DB.ExecContext(cntx, stmt, id)

	if err != nil {
		return err
	}
	return nil
}

// Update the processed status in reservations
func (m *postgresDbRepo) UpdateProcessed(id, processed int) error {
	cntx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
	update reservations set  processed = $1
				where id = $2`

	_, err := m.DB.ExecContext(cntx, stmt,
		processed, id)

	if err != nil {
		return err
	}
	return nil
}

// AllRooms returns all Rooms
func (m *postgresDbRepo) AllRooms() ([]models.Room, error) {
	cntx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var rooms []models.Room
	stmt := `
		select id , room_name , created_at ,updated_at from rooms order by room_name
	`
	rows, err := m.DB.QueryContext(cntx, stmt)
	if err != nil {
		return rooms, err
	}
	defer rows.Close()

	for rows.Next() {
		var room models.Room

		err := rows.Scan(
			&room.ID,
			&room.RoomName,
			&room.CreatedAt,
			&room.UpdatedAt,
		)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}
	if rows.Err() != nil {
		return rooms, rows.Err()
	}
	return rooms, nil
}

// Get Room Restrictions for a given room for specific date range
func (m *postgresDbRepo) GetRestrictionsByRoomForDate(roomId int, startDate, endDate time.Time) ([]models.RoomRestriction, error) {
	cntx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var restricitons []models.RoomRestriction

	stmt := `
	select 
		  rr.id 
		, coalesce (rr.reservation_id ,0)
		, rr.restriction_id 
		, rr.room_id 
		, rr.start_date 
		, rr.end_date 
	from room_restrictions rr
	where rr.room_id =$1
	and $2 < rr.end_date 
	and $3 >= rr.start_date  
	`
	rows, err := m.DB.QueryContext(cntx, stmt, roomId, startDate, endDate)
	if err != nil {
		return restricitons, err
	}
	defer rows.Close()

	for rows.Next() {
		var restriciton models.RoomRestriction
		err := rows.Scan(
			&restriciton.ID,
			&restriciton.ReservationId,
			&restriciton.RestrictionId,
			&restriciton.RoomId,
			&restriciton.StartDate,
			&restriciton.EndDate,
		)
		if err != nil {
			return restricitons, err
		}
		restricitons = append(restricitons, restriciton)
	}
	if err = rows.Err(); err != nil {
		return restricitons, err
	}
	return restricitons, nil

}

// InsertRoomRestrictionForRoomId inserts a new Block for room restrictions
func (m *postgresDbRepo) InsertRoomRestrictionForRoom(roomId int, startDate time.Time) error {
	cntx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		insert into room_restrictions (
			start_date
		, end_date
		, room_id 
		, restriction_id
		, created_at
		, updated_at
		)
	values ($1,$2,$3,$4,$5,$6)
	`
	_, err := m.DB.ExecContext(cntx, stmt, startDate, startDate.AddDate(0, 0, 1), roomId, 1, time.Now(), time.Now())
	if err != nil {
		return err
	}
	return nil
}

// DeleteRoomRestrictionById Deletes the room restrictions
func (m *postgresDbRepo) DeleteRoomRestrictionById(id int) error {
	cntx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		delete from room_restrictions where id =$1
	`
	_, err := m.DB.ExecContext(cntx, stmt, id)
	if err != nil {
		return err
	}
	return nil
}
