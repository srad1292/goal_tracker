package goal

import (
	"context"
	"fmt"
	"log"

	"github.com/srad1292/goal_tracker/pkg/database"
)

func GetGoalsFromPersistence() (GoalsResponse, error) {
	db := database.GetDatabase()

	var query string = `
		select * 
		from goal 
		where active=true
		order by goal_name
	`
	dbGoals, err := db.Query(context.Background(), query)

	if err != nil {
		log.Printf("Error getting goals: %v", err)
		return GoalsResponse{
			Goals: []Goal{},
		}, err
	}

	goals := make([]Goal, 0)
	for dbGoals.Next() {
		var goal int
		var goal_name string
		var unit string
		var active bool

		err = dbGoals.Scan(&goal, &goal_name, &unit, &active)
		if err != nil {
			// handle this error
			log.Printf("Error reading goals records: %v", err)
			return GoalsResponse{
				Goals: []Goal{},
			}, err
		} else {
			fmt.Printf("Goal: %d, Name: %s, Unit: %s\n", goal, goal_name, unit)
			goals = append(goals, Goal{
				Goal:     goal,
				GoalName: goal_name,
				Unit:     unit,
				Active:   active,
			})
		}
	}

	return GoalsResponse{
		Goals: goals,
	}, nil
}
