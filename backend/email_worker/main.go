package main

import (
    "fmt"
    "log"
    "time"
    "strconv"

    //"github.com/emersion/go-imap"
)


func main() {
	//server := "10.0.0.9:1143"
    server := "127.0.0.1:1143"
	username := "jl.110@protonmail.com"
	password := "SSF9Tigm7mIFk4iKhD18VQ"

    fmt.Println("Connecting to ", server)

    Outer:
	for {
        fmt.Println("starting for loop")
		c, err := connectIMAP(server, username, password)
		if err != nil {
			log.Printf("Failed to connect: %v", err)
			log.Println("Retrying in 15 minutes...")
			time.Sleep(15 * time.Minute)
			continue Outer
		}
		log.Println("Connected successfully")



        emails, err := getMail(c)
        if err != nil {
            log.Fatal(err)
        } else {
            log.Println("Pulled all unread emails from server")
        }


        emails, old_emails, err := filterEmails(c, emails)
        if err != nil {
            log.Fatal(err)
        } else {
            log.Println("Filtered emails for sites (sender addr & subj) and age")
        }
        log.Println("Old emails to move: " + strconv.Itoa(len(old_emails)))



        // move old job emails to Job Hunting/To Delete
        for _, email := range old_emails {

            err := moveToDelete(c, email.ID)
            if err != nil {
                log.Printf("Failed to move email %d: %v", email.ID, err)
                continue
            }

            log.Printf("Moved old email: %s", email.Subject)
        }





        // loop filtered emails

            // extractJobs from body

            // loop jobs in email

                // generate jobID

                // check for unique jobID

                // create job record

                // publish to correct queue

            // moveToDelete


        c.Logout()

        fmt.Println("Processing complete. Waiting 15mins\n\n")
        time.Sleep(15 * time.Minute)

    }



}
