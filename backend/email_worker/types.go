package main

import(
    "time"
)

type Email struct {
    ID      uint32
    Sender  string
    Subject string
    Date    time.Time
    Body    string
}
