package main

import (
    "strings"
    "time"
)

func filterEmails(input_emails []Email) ([]Email, []Email, error) {
    cutoff := time.Now().AddDate(0, -1, 0) //set cutoff as 1 month before current time

    var job_emails []Email //makes slice we will store positive results in
    var old_emails []Email //slice to store results in for emails to be deleted


    //loop through each email
    for _, email := range input_emails {

        isJobEmail := false

        // SEEK Job Alerts
        if email.SenderName == "SEEK Job Alerts" &&
           email.SenderAddr == "jobmail@s.seek.com.au" {

            job_emails = append(job_emails, email)
            isJobEmail = true
        }

        // LinkedIn Job Alerts
        if email.SenderName == "LinkedIn Job Alerts" &&
           email.SenderAddr == "jobalerts-noreply@linkedin.com" {

            job_emails = append(job_emails, email)
            isJobEmail = true
        }

        // LinkedIn company hiring emails
        if email.SenderAddr == "jobs-noreply@linkedin.com" &&
           strings.Contains(strings.ToLower(email.Subject), "is hiring") {

            job_emails = append(job_emails, email)
            isJobEmail = true
        }

        // SEEK Recommendations
        if email.SenderName == "Seek Recommendations" &&
           email.SenderAddr == "noreply@s.seek.com.au" {

            job_emails = append(job_emails, email)
            isJobEmail = true
        }



        // Not a job email, ignore it
        if !isJobEmail {
            continue
        }


        // Now check age ONLY for job emails
        if email.Date.Before(cutoff) {
            old_emails = append(old_emails, email)
            continue
        }


        // Otherwise it's a valid current job email
        job_emails = append(job_emails, email)

                
    } //end of loop




    return job_emails, old_emails, nil
}
