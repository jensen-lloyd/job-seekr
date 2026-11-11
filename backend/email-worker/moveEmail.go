package main

import (
    "log"

    "github.com/emersion/go-imap"
    "github.com/emersion/go-imap/client"
)

func moveEmail(c *client.Client, uid uint32, source string, destination string) error {

    // Select source mailbox in read/write mode
    _, err := c.Select(source, false)
    if err != nil {
        log.Printf("Failed to select mailbox %s: %v", source)
        return err
    }

    seqSet := new(imap.SeqSet)
    seqSet.AddNum(uid)

    // Move email
    err = c.UidMove(seqSet, destination)
    if err != nil {
        return err
    }

    return nil

}
