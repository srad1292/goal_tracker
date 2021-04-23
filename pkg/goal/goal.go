package goal

type Goal struct {
	Goal     int    `json:"goal"`
	GoalName string `json:"goalName"`
	Unit     string `json:"unit"`
	Active   bool   `json:"active"`
}

type GoalsResponse struct {
	Goals []Goal `json:"goals"`
}
