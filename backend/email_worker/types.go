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
