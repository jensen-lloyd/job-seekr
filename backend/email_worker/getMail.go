package main

import (
    "io"
    "sort"

    "github.com/emersion/go-imap"
    "github.com/emersion/go-imap/client"
)

func getMail(c *client.Client) ([]Email, error) {
    // getInbox
    mbox, err := c.Select("INBOX", true)
    _ = mbox //to make the compiler shut up!

    criteria := imap.NewSearchCriteria()
    criteria.WithoutFlags = []string{imap.SeenFlag}

    ids, err := c.Search(criteria)
    if err != nil {
        return nil, err
    }


    if err != nil {
        return nil, err
    }



    seqset := new(imap.SeqSet)
    seqset.AddNum(ids...)

    messages := make(chan *imap.Message, 10)
    done := make(chan error, 1)

    go func() {
        section := &imap.BodySectionName{}
        done <- c.Fetch(seqset, []imap.FetchItem{
            imap.FetchUid,
            imap.FetchEnvelope,
            imap.FetchFlags,
            section.FetchItem(),
        }, messages)
    }()

    emails := make([]Email, 0, len(ids))

    for msg := range messages {
        if msg.Envelope == nil {
            continue
        }

        // get body
        //section := &imap.BodySectionName{}
        
        section := &imap.BodySectionName{
            Peek: true,
        }


        body := ""

        if r := msg.GetBody(section); r != nil {
            b, err := io.ReadAll(r)
            if err != nil {
                return nil, err
            }

            body = string(b)
        }



        from := msg.Envelope.From[0]

        emails = append(emails, Email{
        ID:         msg.Uid,
        SenderName: from.PersonalName,
        SenderAddr: from.MailboxName + "@" + from.HostName,
        Subject:    msg.Envelope.Subject,
        Date:       msg.Envelope.Date,
        Body:       body,
        })




    }

    // check whether fetch failed
    if err := <-done; err != nil {
        return nil, err
    }


    // sort emails oldest at the beginning to newest at end

    sort.Slice(emails, func(i, j int) bool {
        return emails[i].Date.Before(emails[j].Date)
    })



    return emails, nil
}
