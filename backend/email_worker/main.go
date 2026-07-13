package main

import (
    "log"
    "time"
    "fmt"

    "github.com/emersion/go-imap"
)


func main() {
	//server := "10.0.0.9:1143"
    server := "127.0.0.1:1143"
	username := "jl.110@protonmail.com"
	password := "SSF9Tigm7mIFk4iKhD18VQ"

    fmt.Println("Connecting to %u", server)

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

		log.Println("Connected successfully!")



        // getInbox
        emails := make([]Email, 0)
        mbox, err := c.Select("INBOX", true)

        criteria := imap.NewSearchCriteria()
        criteria.WithoutFlags = []string{imap.SeenFlag}

        ids, err := c.Search(criteria)
        if err != nil {
            log.Fatal(err)
        }


        if err != nil {
            log.Fatal(err)
        }

        _ = mbox //to make the compiler shut up!




        seqset := new(imap.SeqSet)
        seqset.AddNum(ids...)

        messages := make(chan *imap.Message, 10)
        done := make(chan error, 1)

        go func() {
            section := &imap.BodySectionName{}
            done <- c.Fetch(seqset, []imap.FetchItem{
                imap.FetchUid,
                imap.FetchEnvelope,
                section.FetchItem(),
            }, messages)
        }()

        for msg := range messages {
            if msg.Envelope == nil {
                continue
            }

            if len(msg.Envelope.From) > 0 {
                from := msg.Envelope.From[0]

                emails = append(emails, Email{
                ID:      msg.Uid,
                Sender:  from.MailboxName + from.HostName,
                Subject: msg.Envelope.Subject,
                Date:    msg.Envelope.Date,
                Body:    "body",
                })


            }

        }

        if err := <-done; err != nil {
            log.Fatal(err)
        }





        // loop email envelopes

            // filter out by sender addresses

            // filter out by subjects

            // add email (sender, data, subject, body, status) to list

        // loop emails

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
