package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "strings"
    "time"

    amqp "github.com/rabbitmq/amqp091-go"
)

const (
    seekQueue     = "job-seek-scrape"
    linkedinQueue = "job-linkedin-scrape"
)



func main() {
	server := "100.64.74.102:1143"
	username := "jl.110@protonmail.com"
	password := "SSF9Tigm7mIFk4iKhD18VQ"

    db_dump := false
    delete_when_done := true
    delete_old := true
    email_age_cutoff := 7 //age of email in days before it is classed 'old' and ignored


    Outer:
	for {
        fmt.Println("Connecting to IMAP @", server)
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
            err_connect := connectMongo()
            if err_connect != nil {
                log.Printf("MongoDB connection failed: %v", err_connect)
                mongoReady <- false
                return
            }
            mongoReady <- true
        }()


        // Connect to rabbitMQ queue
        conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
        if err != nil {
            log.Fatal(err)
        }
        defer conn.Close()

        ch, err := conn.Channel()
        if err != nil {
            log.Fatal(err)
        }
        defer ch.Close()

        // Declare queues (no change if they exist already)
        _, err = ch.QueueDeclare(
            seekQueue,
            true,  // durable
            false, // delete when unused
            false, // exclusive
            false, // no-wait
            nil,
        )
        if err != nil {
            log.Fatal(err)
        }

        _, err = ch.QueueDeclare(
            linkedinQueue,
            true,
            false,
            false,
            false,
            nil,
        )
        if err != nil {
            log.Fatal(err)
        }


        // Enable publisher confirms
        if err := ch.Confirm(false); err != nil {
            log.Fatalf("Failed to enable publisher confirms: %v", err)
        }

        confirms := ch.NotifyPublish(make(chan amqp.Confirmation, 1))


        // Pull unread emails from inbox
        emails, err := getMail(c, "INBOX", true)
        if err != nil {
            log.Fatal(err)
        } else {
            log.Printf("Pulled %d unread emails from server", len(emails))
        }


        emails, old_emails, err := filterEmails(emails, email_age_cutoff)
        if err != nil {
            log.Fatal(err)
        } else {
            log.Println("Filtered emails for sites (sender addr & subj) and age")
        }
        log.Printf("%d old emails to be moved\n", len(old_emails))



        // move old job emails to Job Hunting/toDelete
        // performed in a goroutine asynchonously
        if delete_old == true {
            not_moved := 0
            for i, email := range old_emails {

                err := moveEmail(c, email.ID, "INBOX", "Folders/Job Hunting/toDelete")
                if err != nil {
                    log.Printf("Failed to move email %d: %v", email.ID, err)
                    not_moved += 1
                    continue
                }
                        log.Printf("Moved old email %d/%d: %s", i+1, len(old_emails), email.Subject)
            }
            log.Printf("Successfully moved %d of %d old emails", (len(old_emails)-not_moved), len(old_emails))
        }



        // move job emails to Job Hunting/toBeProcessed
        // performed in a goroutine asynchonously
        not_moved := 0
        for i, email := range emails {

            err := moveEmail(c, email.ID, "INBOX", "Folders/Job Hunting/toBeProcessed")
            if err != nil {
                log.Printf("Failed to move email %d: %v", email.ID, err)
                not_moved += 1
                continue
            }
                    log.Printf("Moved detected email %d/%d: %s", i+1, len(emails), email.Subject)
        }
        log.Printf("Successfully moved %d of %d detected emails -> toBeProcessed", (len(emails)-not_moved), len(emails))



        // Check if DB successfully connected and initialised
        dbReady := <- mongoReady
        if !dbReady {
            log.Fatal("Unable to connect to MongoDB. Cannot proceed with operations")
        }


        // Pull emails from toBeProcessed
        emails, err = getMail(c, "Folders/Job Hunting/toBeProcessed", false)
        if err != nil {
            log.Fatal(err)
        } else {
            log.Printf("Pulled all %d emails from toBeProcessed folder", len(emails))
        }



        // loop filtered emails
        log.Printf("Extracting jobs from %d emails\n", len(emails))
        jobs, err := extractJobs(emails)
        if err != nil {
            log.Fatal(err)
        } else {
            log.Printf("%d jobs extracted from %d emails", len(jobs), len(emails))
        }


        // Process each job, add to DB and queue
        log.Printf("Adding jobs to DB and queue")
        for i, job := range jobs {
            
            // enter job record in DB
            job.DateAdded = time.Now()

            db_result, err := jobsCollection.InsertOne(context.Background(), job)

            if err != nil {
                log.Printf("Failed to add job to MongoDB: %v", err)
                continue
            }

            jobID := fmt.Sprintf("%v", db_result.InsertedID)

            log.Printf(
                "Added to DB job %d/%d: %s %s %s",
                i+1,
                len(jobs),
                job.Platform,
                job.JobURL,
                jobID,
            )

            // publish to correct queue
            queue := ""
            switch job.Platform {
            case "seek":
                queue = seekQueue

            case "linkedin":
                queue = linkedinQueue

            default:
                log.Printf("Unknown platform %q, not adding job to queue", job.Platform)
                continue
            }

            // Convert MongoDB ObjectID into hex string
            jobID = strings.TrimPrefix(jobID, "ObjectID(\"")
            jobID = strings.TrimSuffix(jobID, "\")")

            // Create queue message
            queueJob := QueueJob{
                JobID: jobID,
            }

            body, err := json.Marshal(queueJob)
            if err != nil {
                log.Printf("Failed to encode queue message: %v", err)
                continue
            }

            // Publish message
            err = ch.Publish(
                "",    // default exchange
                queue, // routing key = queue name
                false, // mandatory
                false, // immediate (unused)
                amqp.Publishing{
                    ContentType:  "application/json",
                    DeliveryMode: amqp.Persistent,
                    Body:         body,
                },
            )

            if err != nil {
                log.Printf("Failed to publish job %s: %v", jobID, err)
                continue
            }

            // Wait for RabbitMQ confirmation
            confirmation := <-confirms

            if !confirmation.Ack {
                log.Printf("RabbitMQ rejected job %s", jobID)
                continue
            }

            log.Printf("Queued job %s to %s", jobID, queue)


        }
        


        // move processed job emails from toBeProcessed -> toDelete
        // performed in a goroutine asynchonously
        if delete_when_done == true {
            not_moved := 0
            if len(emails) == 0 {
                log.Printf("No processed emails to delete")
            }
            for i, email := range emails {

                err := moveEmail(c, email.ID, "Folders/Job Hunting/toBeProcessed", "Folders/Job Hunting/toDelete")
                if err != nil {
                    log.Printf("Failed to move email %d: %v", email.ID, err)
                    not_moved += 1
                    continue
                }

                log.Printf("Moved completed email %d/%d: %s", i+1, len(emails), email.Subject)
            }
            log.Printf("Successfully moved %d of %d completed emails -> toDelete", (len(emails)-not_moved), len(emails))
        }


        // Close connection to SMTP server
        c.Logout()
        log.Printf("IMAP connection closed")

        // Dump DB contents oldest to newest
        if db_dump == true {
            log.Printf("MongoDB contents (newest is last):")
            dumpJobs()
        }


        log.Println("Processing complete. Waiting 5mins\n\n")
        time.Sleep(5 * time.Minute)



    }

}
