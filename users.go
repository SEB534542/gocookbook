package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type user struct {
	Username string // Username for logging in.
	Password []byte // Password for user to log in.
}

var (
	fnameUsers = folderConfig + "users.json" // File where users are stored.
	dbUsers    = map[string]user{}           // username, user
)

/*loadUsers tries to load the users from fnameUsers. If it failes, it creates
a new file with the default user and password.*/
func loadUsers() {
	err := readJSON(&dbUsers, fnameUsers)
	if err != nil {
		log.Printf("Unable to load users from '%v': %v", fnameUsers, err)
		log.Print("Setting default user")
		p, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.MinCost)
		if err != nil {
			log.Fatal(err)
		}
		addUpdateUser("admin", p)
	}
}

/* addUser takes a username and a converted password, and stores it in
the filename. If the username already exists, the password is updated.*/
func addUpdateUser(un string, pwd []byte) {
	dbUsers[un] = user{un, pwd}
	SaveToJSON(dbUsers, fnameUsers)
}

/* removeUser takes a username and removes the user.*/
func removeUser(un string) {
	delete(dbUsers, un)
	SaveToJSON(dbUsers, fnameUsers)
}

/* checkPwd takes a username and a password. It compares this password
with the password stored for the user and returns an error if it does not
match.*/
func checkPwd(un, pwd string) error {
	err := fmt.Errorf("Username and/or password do not match")
	// lookup username
	u, ok := dbUsers[un]
	if !ok {
		return err
	}
	err2 := bcrypt.CompareHashAndPassword(u.Password, []byte(pwd))
	switch {
	case err2 == bcrypt.ErrMismatchedHashAndPassword:
		return err
	case err2 != nil:
		return err2
	default:
		return nil
	}
}
