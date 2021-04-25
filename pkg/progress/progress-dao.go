package progress

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/srad1292/goal_tracker/pkg/database"
)

func GetProgressFromPersistence(goalId int, year int) (ProgressResponse, error) {
	db := database.GetDatabase()

	startDate := fmt.Sprintf(`'%d-01-01'`, year)
	endDate := fmt.Sprintf(`'%d-01-01'`, year+1)

	query := `
		select progress.*, goal.goal_name, goal.unit 
		from progress
		left join goal on progress.goal = goal.goal 
		where progress.goal = $1
		and session_date between $2 and $3
		order by session_date;
	`

	dbProgress, err := db.Query(context.Background(), query, goalId, startDate, endDate)

	if err != nil {
		log.Printf("Error getting progress for goal: %d, year: %d, err: %v", goalId, year, err)
		return ProgressResponse{
			Sessions: []ProgressSession{},
		}, err
	}

	sessions := make([]ProgressSession, 0)
	for dbProgress.Next() {
		var progress int
		var amount int
		var sessionDate time.Time
		var goal int
		var goal_name string
		var unit string

		err = dbProgress.Scan(&progress, &amount, &sessionDate, &goal, &goal_name, &unit)
		if err != nil {
			log.Printf("Error reading progress records for goal: %d, year: %d, : %v", goalId, year, err)
			return ProgressResponse{
				Sessions: []ProgressSession{},
			}, err
		} else {
			sessions = append(sessions, ProgressSession{
				Progress:    progress,
				Amount:      amount,
				SessionDate: sessionDate.Format("2006-01-02"),
				Goal:        goal,
				GoalName:    goal_name,
				Unit:        unit,
			})
		}
	}

	return ProgressResponse{
		Sessions: sessions,
	}, nil
}
