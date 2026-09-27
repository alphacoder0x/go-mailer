package main

import (
	"fmt"
	"log"
	"net/smtp"
	"sync"
	"time"
)

// emailWorker receives recipients from ch until it is closed, renders the
// template for each one and sends it via SMTP. id is only used for logging.
func emailWorker(id int, ch chan Recipient,wg *sync.WaitGroup ) {
	defer wg.Done()

	for recipient := range ch {

		// Local Mailpit SMTP server (see readme.md).
		smtpHost := "localhost"
		smtpPort := "1025"
		

		// formattedMsg := fmt.Sprintf("To: %s\r\nSubject: Test Email\r\n\r\n%s\r\n",recipient.Email,"just testing our campaign")
		// message := []byte(formattedMsg)

		message,err := executeTemplate(recipient)

		if err != nil{
			fmt.Printf("worker: %d error parsing template for %s",id,recipient.Email)
			continue
		}

		// No auth: Mailpit accepts anything. Replace the sender address before real use.
		err1 := smtp.SendMail(smtpHost+":"+smtpPort,nil,"<youremail>@gmail.com",[]string{recipient.Email},[]byte(message))

		if err1 != nil{
			log.Fatal(err1)
		}

		// Small delay per worker to throttle the send rate.
		time.Sleep(50 * time.Millisecond)

		fmt.Printf("Worker %d: Sending email to %s\n",id,recipient.Email)

	}

	

}