package progress

type ProgressByTime struct {
	ProgressDate  int    `json:"progressDate"`
	TotalProgress int    `json:"totalProgress"`
	SessionDate   string `json:"sessionDate"`
	Goal          int    `json:"goal"`
	GoalName      string `json:"goalName"`
	Unit          string `json:"unit"`
}

type ProgressByTimeResponse struct {
	Progress []ProgressByTime `json:"progress"`
}
