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
