package main

import (
    "io"
    "sort"

    "github.com/emersion/go-imap"
    "github.com/emersion/go-imap/client"
)

// mailbox: IMAP folder to retrieve emails from
// unreadOnly: true = unread only, false = all emails
func getMail(c *client.Client, mailbox string, unreadOnly bool) ([]Email, error) {

    // Select mailbox (read-only so fetching doesn't mark emails as read)
    _, err := c.Select(mailbox, true)
    if err != nil {
        return nil, err
    }

    // Build search criteria
    criteria := imap.NewSearchCriteria()

    if unreadOnly {
        criteria.WithoutFlags = []string{imap.SeenFlag}
    }

    // Search using UIDs (not message sequence numbers)
    ids, err := c.UidSearch(criteria)
    if err != nil {
        return nil, err
    }

    // Nothing to retrieve
    if len(ids) == 0 {
        return []Email{}, nil
    }

    // Create UID sequence set
    seqset := new(imap.SeqSet)
    seqset.AddNum(ids...)

    // Fetch email data
    messages := make(chan *imap.Message, 10)
    done := make(chan error, 1)

    section := &imap.BodySectionName{
        Peek: true,
    }

    go func() {
        done <- c.UidFetch(seqset, []imap.FetchItem{
            imap.FetchUid,
            imap.FetchEnvelope,
            imap.FetchFlags,
            section.FetchItem(),
        }, messages)
    }()

    emails := make([]Email, 0, len(ids))

    var readErr error

    for msg := range messages {

        if msg.Envelope == nil {
            continue
        }

        body := ""

        if r := msg.GetBody(section); r != nil {
            b, err := io.ReadAll(r)
            if err != nil {
                readErr = err
                continue
            }

            body = string(b)
        }

        // Safely handle emails without a sender
        senderName := ""
        senderAddr := ""

        if len(msg.Envelope.From) > 0 {
            from := msg.Envelope.From[0]
            senderName = from.PersonalName
            senderAddr = from.MailboxName + "@" + from.HostName
        }

        emails = append(emails, Email{
            ID:         msg.Uid,
            SenderName: senderName,
            SenderAddr: senderAddr,
            Subject:    msg.Envelope.Subject,
            Date:       msg.Envelope.Date,
            Body:       body,
        })
    }

    // Check whether fetch failed
    if err := <-done; err != nil {
        return nil, err
    }

    if readErr != nil {
        return nil, readErr
    }

    // Sort oldest to newest
    sort.Slice(emails, func(i, j int) bool {
        return emails[i].Date.Before(emails[j].Date)
    })

    return emails, nil
}
