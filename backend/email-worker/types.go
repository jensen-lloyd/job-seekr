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
	ID       string    `bson:"id"` //hash of JobURL 
	Platform string    `bson:"platform"`
	JobURL   string    `bson:"job_url"` //is just the job number on respective platform
	DateAdded time.Time `bson:"date_added"`
}
