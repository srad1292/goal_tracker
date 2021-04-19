package goal

func GetGoalsFromPersistence() GoalsResponse {
	return GoalsResponse{
		Goals: []Goal{
			{
				Goal:     1,
				GoalName: "Push Ups",
				Unit:     "",
			},
			{
				Goal:     3,
				GoalName: "Drawing",
				Unit:     "Minutes",
			},
		},
	}
}
