package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// TODO: make type users and funcs method of users

type user struct {
	Username string // Username for logging in.
	Password []byte // Password for user to log in.
	Admin    bool   // True if admin user.
}

var (
	fnameUsers = folderConfig + "users.json" // File where users are stored.
	dbUsers    = map[string]user{}           // username, user.
	adminUser  = "chef"                      // default username that can create, modify and delete users. TODO: make this a role that can be assigned.
)

/*loadUsers tries to load the users from fnameUsers. If it failes, it creates
a new file with the default user and password.*/
func loadUsers() {
	err := readJSON(&dbUsers, fnameUsers)
	if err != nil {
		log.Printf("Unable to load users from '%v': %v", fnameUsers, err)
		log.Print("Setting default user")
		addUpdateUser("chef", "koken", true)
	}
}

/* addUpdateUser takes a username, a converted password and an
indicator if it is an admin user. It stores in the filename.
If the username already exists, the password is updated.*/
func addUpdateUser(un, p string, b bool) {
	if un != "" {
		pwd, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.MinCost)
		if err != nil {
			log.Print(err)
			return
		}
		dbUsers[un] = user{un, pwd, b}
		SaveToJSON(dbUsers, fnameUsers)
	}
}

/* userExists takes a username. It returns true if the username already exists,
false if it doesn't.*/
func userExists(un string) bool {
	_, ok := dbUsers[un]
	if ok {
		return true
	}
	return false
}

func isAdmin(un string) bool {
	u, ok := dbUsers[un]
	if ok {
		return u.Admin
	}
	return false
}

/* removeUser takes a username and removes the user.*/
func removeUser(un string) {
	delete(dbUsers, un)
	SaveToJSON(dbUsers, fnameUsers)
}

/* checkPwd takes a username and a password. It compares this password
with the password stored for the user and returns an error if it does not
match.*/
func checkPwd(un, p string) error {
	err := fmt.Errorf("Username and/or password do not match")
	// lookup username
	u, ok := dbUsers[un]
	if !ok {
		return err
	}
	err2 := bcrypt.CompareHashAndPassword(u.Password, []byte(p))
	switch {
	case err2 == bcrypt.ErrMismatchedHashAndPassword:
		return err
	case err2 != nil:
		return err2
	default:
		return nil
	}
}

/* users returns the users as a slice of string.*/
func users() []string {
	xs := make([]string, len(dbUsers))
	i := 0
	for k, _ := range dbUsers {
		xs[i] = k
	}
	return xs
}
