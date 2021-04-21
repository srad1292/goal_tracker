package goal

import (
	"context"
	"fmt"
	"log"

	"github.com/srad1292/goal_tracker/pkg/database"
)

func GetGoalsFromPersistence() GoalsResponse {
	db := database.GetDatabase()

	goals, error := db.Query(context.Background(), "select * from goal where active=true")

	if error != nil {
		log.Printf("Error getting goals: %v", error)
		return GoalsResponse{
			Goals: []Goal{},
		}
	}

	// goals.
	for goals.Next() {
		var goal int
		var goal_name string
		var unit string
		var active bool

		error = goals.Scan(&goal, &goal_name, &unit, &active)
		if error != nil {
			// handle this error
			log.Printf("Error reading goals records: %v", error)
		} else {
			fmt.Printf("Goal: %d, Name: %s, Unit: %s", goal, goal_name, unit)
		}
	}

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
