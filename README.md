First table creation in ElephantSQL -> browser -> execute SQL Query

CREATE TABLE jobs (
    jobid SERIAL PRIMARY KEY,
    jobname TEXT,
    openings INT,
    location TEXT
);

URL, DB details in .env file

Refer router/router.go and models/models.go before playing around POSTMAN

To add a new job: /api/newjob  (corresponding json, refer struct in models.go)
To search a job: /api/searchjob/{name}
To delete a job: /api/deletejob/{id}
To view all jobs: /api/getalljobs