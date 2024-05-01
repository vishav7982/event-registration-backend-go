package models

import (
	"errors"
	"restapp/db"
	"restapp/utils"
)

type User struct {
	ID       int64
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

func (u *User) Save() error {
	query := "INSERT INTO users(email,password) VALUES(?,?)"
	statement, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer statement.Close()
	HashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}
	result, err := statement.Exec(u.Email, HashedPassword)
	if err != nil {
		return err
	}
	userId, err := result.LastInsertId()
	u.ID = userId
	return err
}

func (u *User) ValidateCredentials() error {
	query := "SELECT id,email,password FROM users WHERE email= ?"
	row := db.DB.QueryRow(query, u.Email)
	var retrievedEmail, retrievedPassword string
	err := row.Scan(&u.ID, &retrievedEmail, &retrievedPassword)
	if err != nil {
		return err
	}
	if !utils.CheckPasswordHash(u.Password, retrievedPassword) {
		return errors.New("credentials invalid")
	}
	//User validated succesfully
	return nil

}
