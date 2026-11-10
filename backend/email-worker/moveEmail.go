package main

import (
    "github.com/emersion/go-imap"
    "github.com/emersion/go-imap/client"
)

func moveEmail(c *client.Client, uid uint32, destination string) error {

    //get mailbox, now as Read/Write
    c.Select("INBOX", false)

    seqSet := new(imap.SeqSet)
    seqSet.AddNum(uid)

    // Copy email to destination folder
    err := c.UidCopy(
        seqSet,
        destination,
    )

    if err != nil {
        return err
    }

    return nil
}
