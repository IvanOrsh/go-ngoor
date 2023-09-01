package main

import (
	"html/template"
	"os"
)

type UserMeta struct {
	Visits int
}

type User struct {
	Name string
	Age int
	Bio string
	Meta UserMeta
}

func main() {
	t, err := template.ParseFiles("hello.gohtml")
	if err != nil {
		panic(err)
	}

	// user := struct{
	// 	Name string
	// }{
	// 	Name: "John Doe",
	// }

	user := User{
		Name: "John Doe",
		Age: 36,
		Bio: `<script>alert("Should never allow such things!");</script>`,
		Meta: UserMeta{
			 Visits: 15,
		},
	}

	err = t.Execute(os.Stdout, user)
	if err != nil {
		panic(err)
	}
}