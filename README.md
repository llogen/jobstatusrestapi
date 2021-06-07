# jobstatusrestapi
restapi that can be used to manage the jobstatuses

using like this: \
Reading a job status: GET /readjobstatus/ID \
Updating a job status: POST /updatejobstatus/?{ID:ID,Status:false/true} \
Adding a job status: POST /addjobstatus/?{ID:ID,Status:false/true} \
Removing a job status: POST /removejobstatus/?{ID:ID,Status:false/true}
