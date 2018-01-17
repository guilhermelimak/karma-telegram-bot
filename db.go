package main

import (
	"encoding/json"
	"fmt"

	scribble "github.com/nanobox-io/golang-scribble"
)

const dir = "./data"
const collection = "karma"

// Connection bla bla bl
type Connection struct {
	driver *scribble.Driver
}

func (c *Connection) getAllRecords() []User {
	var users []User

	records, err := c.driver.ReadAll(collection)

	if err != nil {
		fmt.Println("Error", err)
	}

	for _, u := range records {
		userFound := User{}

		if err := json.Unmarshal([]byte(u), &userFound); err != nil {
			fmt.Println("Error", err)
		}

		users = append(users, userFound)
	}

	fmt.Printf("%+v", users)

	return users
}

func (c *Connection) get(id string) User {
	user := User{}

	if err := c.driver.Read(collection, id, &user); err != nil {
		c.save(User{ID: id})
	}

	return user
}

func (c *Connection) save(u User) {
	fmt.Printf("Saving: %+v\n", u)
	c.driver.Write(collection, u.ID, u)
}

func (c *Connection) updateKarma(userID string, increment int) {
	user := c.get(userID)
	user.Karma = user.Karma + increment
	user.ID = userID
	c.save(user)
}

func initDB() Connection {
	db, err := scribble.New(dir, nil)
	c := Connection{driver: db}

	if err != nil {
		fmt.Println("Error initDB: ", err)
	}

	return c
}
