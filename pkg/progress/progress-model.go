package progress

type Progress struct {
	Progress    int    `json:"progress"`
	Amount      int    `json:"amount"`
	SessionDate string `json:"sessionDate"`
	Goal        int    `json:"goal"`
}
