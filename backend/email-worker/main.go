package main

import (
    "fmt"
    "log"
    "time"
)


func main() {
	server := "100.64.74.102:1143"
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


        // Connect to and initialise DB in background
        mongoReady := make(chan bool)

        go func() {
            err := connectMongo()
            if err != nil {
                log.Printf("MongoDB connection failed: %v", err)
                mongoReady <- false
                return
            }

            err := initialiseMongo()
            if err != nil {
                log.Printf("MongoDB initialisation failed: %v", err)
                mongoReady <- false
                return
            }

            mongoReady <- true

        }()




        // Pull unread emails from server
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
        log.Printf("%d old emails to be moved\n", len(old_emails))



        // move old job emails to Job Hunting/To Delete
        // performed in a goroutine asynchonously
        go func() {
            if len(old_emails) == 0 {
                return
            }
            for i, email := range old_emails {

                err := moveToDelete(c, email.ID)
                if err != nil {
                    log.Printf("Failed to move email %d: %v", email.ID, err)
                    continue
                }

                log.Printf("Moved old email %d/%d: %s", i+1, len(old_emails), email.Subject)
            }
            log.Printf("Successfully moved %d old emails", len(old_emails))
        }()





        // loop filtered emails
        log.Printf("Extracting jobs from %d emails\n", len(emails))
        jobs, err := extractJobs(emails)
        if err != nil {
            log.Fatal(err)
        } else {
            log.Printf("%d jobs extracted from %d emails", len(jobs)+1, len(emails))
        }
        fmt.Println(jobs)



        // Check if DB successfully connected and initialised
        ready := <- mongoReady
        if !ready {
            log.Fatal("Unable to connect to MongoDB. Cannot proceed with operations")
        } else if ready {
            log.Println("MongoDB connected successfully")
        }



        // Process each job, add to DB and queue
        for _, job := range jobs {
            // check for unique jobID
            exists, err := jobExists(job.ID)

            if err != nil {
                log.Printf("MongoDB lookup failed: %v", err)
                continue
            }

            if exists {
                log.Printf("Job already exists: %s", job.ID)
                continue
            }

            // Job doesn't exist
            // Continue processing it



            // create job record

            // publish to correct queue

            // moveToDelete
        }
        




        c.Logout()

        log.Println("Processing complete. Waiting 5mins\n\n")
        time.Sleep(5 * time.Minute)

    }

}
