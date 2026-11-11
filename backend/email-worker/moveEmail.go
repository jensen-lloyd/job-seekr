package main

import (
    "github.com/emersion/go-imap"
    "github.com/emersion/go-imap/client"
)

func moveEmail(c *client.Client, uid uint32, source string, destination string) error {

    //get mailbox, now as Read/Write
    c.Select(source, false)

    seqSet := new(imap.SeqSet)
    seqSet.AddNum(uid)

    // Copy email to destination folder
    err := c.UidMove(
        seqSet,
        destination,
    )

    if err != nil {
        return err
    }

    return nil
}
