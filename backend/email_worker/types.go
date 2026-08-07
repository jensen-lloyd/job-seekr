package main

import(
    "time"
)

type Email struct {
    ID         uint32
    SenderName string
    SenderAddr string
    Subject    string
    Date       time.Time
    Body       string
}


type Job struct {
    ID uint32 //hash of URL
    Platform string
    JobURL string
    DateAdded time.Time
}
