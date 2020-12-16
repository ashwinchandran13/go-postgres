First table creation in ElephantSQL -> browser -> execute SQL Query

CREATE TABLE jobs (
    jobid SERIAL PRIMARY KEY,
    jobname TEXT,
    openings INT,
    location TEXT
);
