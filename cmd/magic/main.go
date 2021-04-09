package main

import (
	"fmt"
	"github.com/scottcagno/go-blog/tools"
	"log"
)

// example model
type User struct {
	ID       int    `html:"number"`
	Name     string `html:"text"`
	Email    string `html:"email"`
	Password string `html:"password"`
	IsActive bool   `html:"checkbox"`
}

func main() {

	user := &User{
		ID:       5,
		Name:     "John Doe",
		Email:    "jdoe@example.com",
		Password: "jdoe123",
		IsActive: true,
	}

	modl, err := tools.MakeModel(user, "html")
	if err != nil {
		log.Printf("got error: %v\n", err)
	}
	fmt.Printf("Model(%s)\n", modl.Name)
	for i := 0; i < modl.Count; i++ {
		fld := modl.Fields[i]
		fmt.Printf("Field %d: name=%s, type=%s, value=%v, tag=%s\n", i, fld.Name, fld.Type, fld.Value, fld.Tag)
	}

}
