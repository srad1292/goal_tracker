package goal

import (
	"context"
	"fmt"
	"log"

	"github.com/srad1292/goal_tracker/pkg/database"
)

func GetGoalsFromPersistence(onlyActive bool) (GoalsResponse, error) {
	db := database.GetDatabase()

	var query string
	if onlyActive {
		query = `
			select * 
			from goal 
			where active=true
			order by goal_name;
		`
	} else {
		query = `
			select * 
			from goal 
			order by goal_name;
		`
	}

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

func AddGoalToPersistence(newGoal Goal) (Goal, error) {
	db := database.GetDatabase()

	var query string = `
		insert into goal (goal_name, unit, active) 
		values
		($1, $2, $3)
		returning goal;
	`
	dbGoalId, err := db.Query(context.Background(), query, newGoal.GoalName, newGoal.Unit, newGoal.Active)

	if err != nil {
		log.Printf("Error creating goal: %v", err)
		return newGoal, err
	}

	var goal int = 0

	for dbGoalId.Next() {
		err = dbGoalId.Scan(&goal)
		if err != nil {
			log.Printf("Error scanning new goal records: %v", err)
			return newGoal, err
		} else {
			newGoal.Goal = goal
		}
	}

	return newGoal, nil
}

func UpdateGoalInPersistence(goal Goal, goalId int) (Goal, error) {
	db := database.GetDatabase()

	var query string = `
		update goal 
		set goal_name = $1,
		unit = $2,
		active = $3
		where goal = $4;
	`
	_, err := db.Query(context.Background(), query, goal.GoalName, goal.Unit, goal.Active, goalId)

	if err != nil {
		log.Printf("Error updating goal: %v", err)
		return goal, err
	}

	goal.Goal = goalId

	return goal, nil
}
