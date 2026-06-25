package main

func main() {

    for {

        // connectIMAP
            // if fails, break out

        // getInbox

        // loop email envelopes

            // filter out by sender addresses

            // filter out by subjects

            // add email (sender, data, subject, body, status) to list

        // loop emails

            // email already read?
                //remove from list

            // email older than 90 days?
                // moveToDelete
                // remove from list

        // loop filtered emails

            // extractJobs from body

            // loop jobs in email

                // generate jobID

                // check for unique jobID

                // create job record

                // publish to correct queue

            // moveToDelete


    }
    // once loop ends or is broken out of
    // sleep 5mins

}
