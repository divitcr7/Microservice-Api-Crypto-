package domain

type Session struct {
	Id       int64  `json:"id" db:"_id"`
	TaskName string `json:"task_name" db:"task_name"`
	Interval int64  `json:"interval" db:"interval"`
}
