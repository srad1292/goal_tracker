package progress

type ProgressSession struct {
	Progress    int    `json:"progress"`
	Amount      int    `json:"amount"`
	SessionDate string `json:"sessionDate"`
	Goal        int    `json:"goal"`
	GoalName    string `json:"goalName"`
	Unit        string `json:"unit"`
}

type ProgressResponse struct {
	Sessions []ProgressSession `json:"sessions"`
}
