package main

import (
    "crypto/sha256"
    "fmt"
    "log"
    "net/http"
    "regexp"
    "strings"
    "time"
)


func getSeekJobID(trackingURL string) (string, error) {

    // Clean the URL
    trackingURL = strings.ReplaceAll(trackingURL, "=", "")
	trackingURL = strings.ReplaceAll(trackingURL, `\r`, "")
	trackingURL = strings.ReplaceAll(trackingURL, `\n`, "")
	trackingURL = strings.ReplaceAll(trackingURL, "\r", "")
	trackingURL = strings.ReplaceAll(trackingURL, "\n", "")


    // Request the SEEK tracking URL
    response, err := http.Get(trackingURL)
    if err != nil {
        return "", err
    }

    defer response.Body.Close()

    // Get the final URL after redirects
    realURL := response.Request.URL.String()


    // Extract the SEEK job number
    seekJobRegex := regexp.MustCompile(
        `https://au\.seek\.com/job/([0-9]+)`,
    )

    matches := seekJobRegex.FindStringSubmatch(realURL)


    if len(matches) < 2 {
        return "", fmt.Errorf("could not find SEEK job ID in URL: %s", realURL)
    }


    // Return the job number
    return matches[1], nil
}



func extractJobs(emails []Email) ([]Job, error) {
    var jobs []Job //makes slice to store job objects in


    linkedinRegex := regexp.MustCompile(`https://www\.linkedin\.com/comm/jobs/view/([0-9]+)`)


    seekTrackingRegex := regexp.MustCompile(
	`https://email\.s\.seek\.com\.au/uni/ss/c/[^"]*`,
)


    //loop to find all relevant URLs in email

    for i, email := range emails { 
        log.Printf("Processing URLs from email %d/%d\n", i+1, len(emails))

        //process Linkedin URLs
        linkedinMatches := linkedinRegex.FindAllStringSubmatch(email.Body, -1)

        for _, match := range linkedinMatches {
            jobID := match[1]

            if len(jobID) < 9 {
                continue
            }

            // Set platform
            platform := "linkedin"

            // Create job object and add to slice
            hash := sha256.Sum256([]byte(jobID + platform))

            jobs = append(jobs, Job{
                ID: fmt.Sprintf("%x", hash),
                Platform: platform,
                JobURL: jobID,
                DateAdded: time.Now(),
            })

        }


        //process Seek URLs
        seekMatches := seekTrackingRegex.FindAllString(email.Body, -1)

        for _, match := range seekMatches {

            //Get jobID from email URL
            jobID, err := getSeekJobID(match)
            if err != nil {
                continue
            }

            if len(jobID) < 7 {
                continue
            }

            // Set platform
            platform := "seek"

            // Create job object and add to slice
            hash := sha256.Sum256([]byte(jobID + platform))

            jobs = append(jobs, Job{
                ID: fmt.Sprintf("%x", hash),
                Platform: platform,
                JobURL: jobID,
                DateAdded: time.Now(),
            })


        }



        if i > 10 {
            return jobs, nil
        }

    }


    // Return list of jobs
    return jobs, nil
}
