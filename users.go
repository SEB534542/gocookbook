package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// Users represents a file location and a map containing all the users (of type user).
type Users struct {
	uns   map[string]user // username, user.
	fname string          // location of json file.
}

// user represents a username, with a password and an indicator if the user is an admin.
type user struct {
	Username string // Username for logging in.
	Password []byte // Password for user to log in.
	Admin    bool   // True if admin user.
}

// CreateUsers takes a file name, loads the Users from the JSON and returns it.
func CreateUsers(fname string) Users {
	dbUsers := Users{
		uns:   map[string]user{},
		fname: fname,
	}
	dbUsers.Load()
	return dbUsers
}

/*Load tries to load the Users from the filename stored in Users. If it failes, it
a new Users is created with the default user and password as specified in this
method.*/
func (dbUsers Users) Load() {
	err := readJSON(&dbUsers.uns, dbUsers.fname)
	if err != nil {
		log.Printf("Unable to load users from '%v': %v", dbUsers.fname, err)
		log.Print("Setting default user")
		dbUsers.AddUpdate("chef", "koken", true)
	}
}

/* AddUpdate takes a username, a password and an
indicator if it is an admin user. If the username already exists, the password is
updated, else a new user is added, after which the updated Users is stored.*/
func (dbUsers Users) AddUpdate(un, p string, b bool) {
	if un != "" {
		pwd, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.MinCost)
		if err != nil {
			log.Print(err)
			return
		}
		dbUsers.uns[un] = user{un, pwd, b}
		SaveToJSON(dbUsers.uns, dbUsers.fname)
	}
}

/* Exists takes a username. It returns true if the username already exists,
false if it doesn't.*/
func (dbUsers Users) Exists(un string) bool {
	_, ok := dbUsers.uns[un]
	if ok {
		return true
	}
	return false
}

/* IsAdmin takes a username and returns triue if the user is and admin.
It returns false if the it is not an admin, or user doesn't exists.*/
func (dbUsers Users) IsAdmin(un string) bool {
	u, ok := dbUsers.uns[un]
	if ok {
		return u.Admin
	}
	return false
}

/* Remove takes a username and removes the user.*/
func (dbUsers Users) Remove(un string) {
	delete(dbUsers.uns, un)
	SaveToJSON(dbUsers.uns, dbUsers.fname)
}

/* CheckPwd takes a username and a password. It compares this password
with the password stored for the user and returns an error if it does not
match.*/
func (dbUsers Users) CheckPwd(un, p string) error {
	err := fmt.Errorf("Username and/or password do not match")
	// lookup username
	u, ok := dbUsers.uns[un]
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

/* Users returns all users as a slice of string.*/
func (dbUsers Users) Users() []string {
	xs := make([]string, len(dbUsers.uns))
	i := 0
	for k, _ := range dbUsers.uns {
		xs[i] = k
	}
	return xs
}
