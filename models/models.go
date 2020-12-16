package models

// Job schema of the jobs table
type Job struct {
	ID       int64  `json:"jobid"`
	Name     string `json:"jobname"`
	Openings int64  `json:"openings"`
	Location string `json:"location"`
}
