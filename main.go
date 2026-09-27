package main

import (
	"bytes"
	"fmt"
	"text/template"

	// "time"
	"sync"
)

// Recipient is one row of the campaign CSV: who to email and the name
// used to personalise the template.
type Recipient struct{
	Name string
	Email string
}

// main wires up a producer/consumer pipeline: one goroutine reads recipients
// from the CSV into a channel, and a pool of workers drains that channel and
// sends one email per recipient.
func main(){
	fmt.Println("welcome to go-mailer")

	// Unbuffered channel shared by the producer and all workers.
	recipientCh := make(chan Recipient)

	// Tracks the workers so main waits until every email has been sent.
	var wg sync.WaitGroup


	// Producer: closes recipientCh when done, which ends the workers' range loops.
	go func(){
		loadRecipient("./emails.csv",recipientCh)

	}()

	fmt.Println("starting consumer")

	// Number of concurrent SMTP senders.
	workerCount := 5

	for i := range workerCount{
		fmt.Printf("starting worker no %v \n",i)
		wg.Add(1)
		go emailWorker(i,recipientCh, &wg)
	}
	wg.Wait()
}

// executeTemplate renders email.tmpl for r and returns the raw message
// (headers + body) ready to pass to smtp.SendMail.
func executeTemplate(r Recipient) (string,error){
	t,err := template.ParseFiles("email.tmpl")

	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	
	err1 := t.Execute(&tpl,r)

	if err1 != nil {
		return "",err1

	}

	return tpl.String(),nil

}