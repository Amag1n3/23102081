STAGE 1
REST API endpoints:

DB schema:
Every student should have a unique ID like we do, say, the ENROLLMENT NUMBER.



/placements
The placement details will be hosted on the "/placement" endpoint as a GET response (for user=student).
The response would be an "array" of JSON of this format:
{
 student_name string
 enrollment_no int64
 batch int64
 package_offered int64
 offered_by string
}

for user=ADMIN
the hosting will be /placement but for a POST request
the header should have appropriate authorization token which can be used to distinguish between user and admin



/events
The event details will be hosted on the "/event" endpoint as a GET response.
The response would be an "array" of JSON of this format:
{
 event_name string
 rsvp string
 batches_invited string
 start_date_time time --
                        |---->Unix-time data type
 end_date_time time ----     
}

same logic for the admin access to post new event as /placements endpoint
must have and authorization header to distinguish admin and normal user/student


/results
Again, hosted on /results
but this time, 3 types of users can have varying levels of access
students --> can only see
teacher --> add marks on existing/already happened exams
admin --> can create exam events, see any students marks and even change them.


Stage 2
I'd use PostgreSQL for it's light-weight, fast and feature-rich
The problem of finding the correct student which takes O(N) time would be the sole bottleneck of the DB when records grow from 50,000 to 5,000,000



Stage 3
Yes, it's logically correct — it fetches unread notifications for a student sorted by newest first.
Why is it slow?
With 5,000,000 rows, there's no index on studentID or isRead, so the database does a full table scan on every request. At this scale that's expensive.
Fix — Add an index
sql:
CREATE INDEX idx_notifications_student_unread
ON notifications (studentID, isRead, createdAt DESC);
This turns the lookup from O(n) full scan to O(log n) index seek.


No. Each index adds overhead on every INSERT/UPDATE
Only index what you actually query on.
Query — Students who got a placement notification in the last 7 days
SELECT DISTINCT studentID
FROM notifications
WHERE notificationType = 'Placement'
  AND createdAt >= NOW() - INTERVAL '7 days';


Stage 4
Cache the Fetched data with REDIS on the user's device.
when refreshed, serve the cached data
dont make a new request until either sufficient time passes or the server's database updates

alternatively, we can use websockets as well
yes, they'll be harder to implement but the losses due to a persistent connection which keeps updating by itself might be worth considering the gains by not sending the server a request every single time



Stage 5
No the sending and saving of the emails in the server shouldnt happen simultaneously
save the email to DB first then retry if the sending of emails fails in-between

function notify_all(student_ids, message):
    insert_to_db(student_ids, message)
    
    for student_id in student_ids:
        enqueue(job={ student_id, message })

worker():
    job = dequeue()
    retry(send_email(job.student_id, job.message), attempts=3)
    push_to_app(job.student_id, job.message)



Stage 6
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	authURL = "http://4.224.186.213/evaluation-service/auth"
	apiURL  = "http://4.224.186.213/evaluation-service/notifications"
	topN    = 10
)

var creds = map[string]string{
	"email":        "amoghtyagi22092005@gmail.com",
	"name":         "amogh tyagi",
	"rollNo":       "23102081",
	"accessCode":   "MdprhE",
	"clientID":     "90e96a35-ea8a-4823-8125-4b2a9170a57d",
	"clientSecret": "UybTNQAChyzwYuxB",
}

func getToken() (string, error) {
	b, _ := json.Marshal(creds)
	resp, err := http.Post(authURL, "application/json", bytes.NewReader(b))
	if err != nil {
		return "", fmt.Errorf("auth request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("auth parse failed: %w", err)
	}
	if result.AccessToken == "" {
		return "", fmt.Errorf("empty token received")
	}
	return result.AccessToken, nil
}
