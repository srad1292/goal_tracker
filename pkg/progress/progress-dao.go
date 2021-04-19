package progress

func GetProgressFromPersistence() ProgressResponse {
	return ProgressResponse{
		Progress: []ProgressPersistence{
			{
				Progress:    1,
				Amount:      15,
				SessionDate: "2021-04-18",
				Goal:        1,
				GoalName:    "Push Ups",
				Unit:        "",
			},
			{
				Progress:    2,
				Amount:      85,
				SessionDate: "2021-04-18",
				Goal:        3,
				GoalName:    "Drawing",
				Unit:        "Minutes",
			},
		},
	}
}
