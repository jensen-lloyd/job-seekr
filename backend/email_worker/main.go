package main

import (
    "fmt"
    "log"
    "time"
    //"strconv"
)


func main() {
	//server := "10.0.0.9:1143"
    server := "127.0.0.1:1143"
	username := "jl.110@protonmail.com"
	password := "SSF9Tigm7mIFk4iKhD18VQ"


    Outer:
	for {
        fmt.Println("Connecting to IMAP @ ", server)
		c, err := connectIMAP(server, username, password)
		if err != nil {
			log.Printf("Failed to connect: %v", err)
			log.Println("Retrying...")
			time.Sleep(5 * time.Second)
			continue Outer
		}
		log.Println("Connected to IMAP server successfully")



        emails, err := getMail(c)
        if err != nil {
            log.Fatal(err)
        } else {
            log.Printf("Pulled all %d unread emails from server", len(emails))
        }


        emails, old_emails, err := filterEmails(emails)
        if err != nil {
            log.Fatal(err)
        } else {
            log.Println("Filtered emails for sites (sender addr & subj) and age")
        }
        log.Printf("Old emails to move: %d", len(old_emails))



        // move old job emails to Job Hunting/To Delete
        for i, email := range old_emails {

            err := moveToDelete(c, email.ID)
            if err != nil {
                log.Printf("Failed to move email %d: %v", email.ID, err)
                continue
            }

            log.Printf("Moved old email %d/%d: %s", i, len(old_emails), email.Subject)
        }





        // loop filtered emails
        log.Println("Extracting jobs from emails")
        jobs, err := extractJobs(emails)
        if err != nil {
            log.Fatal(err)
        } else {
            log.Printf("%d jobs extracted from %d emails", len(jobs), len(emails))
        }
        fmt.Print(jobs)



        // check for unique jobID

        // create job record

        // publish to correct queue

        // moveToDelete
        




        c.Logout()

        fmt.Println("Processing complete. Waiting 5mins\n\n")
        time.Sleep(5 * time.Minute)

    }



}
