package main

import (
	"fmt"
)

type Recipient struct{
	Name string
	Email string
}

func main(){
	fmt.Println("welcome to go-mailer") 

	loadRecipient("./emails.csv")

}