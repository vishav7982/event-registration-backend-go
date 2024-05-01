package models

import (
	"errors"
	"restapp/db"
	"time"
)

type Event struct {
	ID          int64
	Name        string    `binding:"required"`
	Description string    `binding:"required"`
	Location    string    `binding:"required"`
	DateTime    time.Time `binding:"required"`
	UserID      int64
}

func (e *Event) Save() error {
	//Later : add it to a database
	query := `
	INSERT INTO events(name,description,location,dateTime,user_id)
	VALUES(?,?,?,?,?)
	`
	preparedQuery, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer preparedQuery.Close() // will close the connection after running this wrapped function which is save
	result, err := preparedQuery.Exec(e.Name, e.Description, e.Location, e.DateTime, e.UserID)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	e.ID = id
	return err

}

func GetAllEvents() ([]Event, error) {
	var events []Event

	rows, err := db.DB.Query("SELECT * FROM events")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var event Event
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func GetEventById(id int64) (*Event, error) {
	query := "SELECT * FROM events WHERE id= ?"
	row := db.DB.QueryRow(query, id)

	var event Event
	err := row.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)

	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (e *Event) Update() error {
	query := `
	UPDATE events
	SET name = ?, description=?, location=?, dateTime=?
	WHERE id = ?
	`
	statement, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(e.Name, e.Description, e.Location, e.DateTime, e.ID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteEvent(id int64) error {
	query := `
	DELETE FROM events WHERE id = ?
	`
	statement, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer statement.Close()
	_, err = statement.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

func (e *Event) Register(userId int64) error {
	query := `
	INSERT INTO registrations(event_id,user_id) values(?,?)
	`
	statement, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(e.ID, userId)
	return err
}
func (e *Event) CheckRegistration(userId int64) (bool, error) {
	var count int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM registrations WHERE event_id = ? AND user_id = ?", e.ID, userId).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func (e *Event) CancelRegistration(userId int64) error {
	// Check if registration exists before attempting cancellation
	exists, err := e.CheckRegistration(userId)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("registration does not exist")
	}
	query := `
	DELETE FROM registrations WHERE event_id = ? AND user_id = ?
	`
	statement, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(e.ID, userId)
	return err
}
