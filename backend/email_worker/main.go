package main

import (
    "fmt"
    "log"
    "time"

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

        _ = emails


        // loop sorted emails to filter LinkedIn/Seek

            // filter out by sender addresses

            // filter out by subjects

            // add email (sender, data, subject, body, status) to list

        // loop job emails to filter relevant

            // email already read?
                //remove from list

            // email older than 90 days?
                // moveToDelete
                // remove from list

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
