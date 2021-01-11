package models

const (
	TaskQueued  = "Queued"
	TaskSuccess = "Success"
	TaskFailed  = "Failed"
)

type Task struct {
	TaskID int64      `json:"task_id"`
	Status string     `json:"status"`
	Error  string     `json:"error,omitempty"`
	Stats  *TaskStats `json:"stats,omitempty"`
	TaskInfo
}

type TaskInfo struct {
	SellerID int64  `json:"-"`
	FileURL  string `json:"-"`
}

type TaskStats struct {
	Created int `json:"created"`
	Updated int `json:"updated"`
	Deleted int `json:"deleted"`
	Invalid int `json:"invalid"`
}
