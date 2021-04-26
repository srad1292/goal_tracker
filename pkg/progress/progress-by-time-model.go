package progress

type ProgressByTime struct {
	ProgressDate  string `json:"progressDate"`
	TotalProgress int    `json:"totalProgress"`
	Goal          int    `json:"goal"`
	GoalName      string `json:"goalName"`
	Unit          string `json:"unit"`
}

type ProgressByTimeResponse struct {
	GroupedSessions []ProgressByTime `json:"groupedSessions"`
}
