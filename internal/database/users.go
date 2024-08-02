package database

import (
	"errors"
)

type User struct {
	ID             int    `json:"id"`
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password"`
}

func (db *DB) CreateUser(email, hashedPassword string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	id := len(dbStructure.Users) + 1

	for _, user := range dbStructure.Users {
		if user.Email == email {
			return User{}, errors.New("email already taken")
		}
	}

	user := User{
		ID:             id,
		Email:          email,
		HashedPassword: hashedPassword,
	}

	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, user := range dbStructure.Users {
		if user.Email == email {
			return user, nil
		}
	}

	return User{}, errors.New("user doesn't exist")
}

func (db *DB) GetUserByID(id int) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		return User{}, errors.New("user doesn't exist")
	}

	return user, nil
}

func (db *DB) UpdateUser(userID int, email, hashedPassword string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbStructure.Users[userID]
	if !ok {
		return User{}, errors.New("user doesn't exist")
	}

	user.Email = email
	user.HashedPassword = hashedPassword
	dbStructure.Users[userID] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
