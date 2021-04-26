package progress

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/srad1292/goal_tracker/pkg/database"
)

func AddProgressToPersistence(newProgress Progress) (Progress, error) {
	db := database.GetDatabase()

	var query string = `
		insert into progress (amount, session_date, goal) 
		values
		($1, $2, $3)
		returning progress;
	`
	dbProgressId, err := db.Query(context.Background(), query, newProgress.Amount, newProgress.SessionDate, newProgress.Goal)

	if err != nil {
		log.Printf("Error creating progress: %v", err)
		return newProgress, err
	}

	var progress int = 0

	for dbProgressId.Next() {
		err = dbProgressId.Scan(&progress)
		if err != nil {
			log.Printf("Error scanning new progress records: %v", err)
			return newProgress, err
		} else {
			newProgress.Progress = progress
		}
	}

	return newProgress, nil
}

func UpdateProgressInPersistence(progress Progress, progressId int) (Progress, error) {
	db := database.GetDatabase()

	var query string = `
		update progress 
		set amount = $1,
		session_date = $2,
		goal = $3
		where progress = $4;
	`
	_, err := db.Query(context.Background(), query, progress.Amount, progress.SessionDate, progress.Goal, progressId)

	if err != nil {
		log.Printf("Error updating progress: %v", err)
		return progress, err
	}

	progress.Progress = progressId

	return progress, nil
}

func DeleteProgressFromPersistence(progressId int) error {
	db := database.GetDatabase()

	var query string = `
		delete
		from progress
		where progress = $1;
	`
	_, err := db.Query(context.Background(), query, progressId)

	if err != nil {
		log.Printf("Error deleting progress: %v", err)
		return err
	}

	return nil
}

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

func GetProgressByTimeFromPersistence(goalId int, year int, period string) (ProgressByTimeResponse, error) {
	db := database.GetDatabase()

	startDate := fmt.Sprintf(`'%d-01-01'`, year)
	endDate := fmt.Sprintf(`'%d-01-01'`, year+1)

	query := `
		select date_trunc($1, session_date) as progress_date, sum(amount) as total_progress, 
			progress.goal, goal.goal_name, goal.unit
		from progress
		left join goal on goal.goal = progress.goal
		where progress.goal = $2
		and session_date between $3 and $4
		group by progress_date, progress.goal, goal.goal_name, goal.unit 
		order by progress_date;
	`

	dbProgress, err := db.Query(context.Background(), query, period, goalId, startDate, endDate)

	if err != nil {
		log.Printf("Error getting progress by time for goal: %d, year: %d, period: %s, err: %v", goalId, year, period, err)
		return ProgressByTimeResponse{
			GroupedSessions: []ProgressByTime{},
		}, err
	}

	sessions := make([]ProgressByTime, 0)
	for dbProgress.Next() {
		var progressDate time.Time
		var totalProgress int
		var goal int
		var goal_name string
		var unit string

		err = dbProgress.Scan(&progressDate, &totalProgress, &goal, &goal_name, &unit)
		if err != nil {
			log.Printf("Error reading progress by time for goal: %d, year: %d, period: %s, err: %v", goalId, year, period, err)
			return ProgressByTimeResponse{
				GroupedSessions: []ProgressByTime{},
			}, err
		} else {
			sessions = append(sessions, ProgressByTime{
				ProgressDate:  progressDate.Format("2006-01-02"),
				TotalProgress: totalProgress,
				Goal:          goal,
				GoalName:      goal_name,
				Unit:          unit,
			})
		}
	}

	return ProgressByTimeResponse{
		GroupedSessions: sessions,
	}, nil
}

func GetBestProgressByTimeFromPersistence(goalId int, year int, period string, edge string) (ProgressByTimeResponse, error) {
	db := database.GetDatabase()

	startDate := fmt.Sprintf(`'%d-01-01'`, year)
	endDate := fmt.Sprintf(`'%d-01-01'`, year+1)

	query := `
		select date_trunc($1, session_date) as progress_date, sum(amount) as total_progress, 
			progress.goal, goal.goal_name, goal.unit
		from progress
		left join goal on goal.goal = progress.goal
		where progress.goal = $2
		and session_date between $3 and $4
		group by progress_date, progress.goal, goal.goal_name, goal.unit 
		having sum(amount) = (`

	if edge == "low" {
		query = query + `
			select min(grouped_sessions.total_amount)`
	} else {
		query = query + `
			select max(grouped_sessions.total_amount)`
	}

	query = query + `
		from(
			select sum(amount) as total_amount
				from progress
				where progress.goal = $2
				and session_date between $3 and $4
				group by date_trunc($1, session_date), progress.goal
			) as grouped_sessions
		)
		order by progress_date;
	`

	dbProgress, err := db.Query(context.Background(), query, period, goalId, startDate, endDate)

	if err != nil {
		log.Printf("Error getting progress by time for goal: %d, year: %d, period: %s, err: %v", goalId, year, period, err)
		return ProgressByTimeResponse{
			GroupedSessions: []ProgressByTime{},
		}, err
	}

	sessions := make([]ProgressByTime, 0)
	for dbProgress.Next() {
		var progressDate time.Time
		var totalProgress int
		var goal int
		var goal_name string
		var unit string

		err = dbProgress.Scan(&progressDate, &totalProgress, &goal, &goal_name, &unit)
		if err != nil {
			log.Printf("Error reading progress by time for goal: %d, year: %d, period: %s, err: %v", goalId, year, period, err)
			return ProgressByTimeResponse{
				GroupedSessions: []ProgressByTime{},
			}, err
		} else {
			sessions = append(sessions, ProgressByTime{
				ProgressDate:  progressDate.Format("2006-01-02"),
				TotalProgress: totalProgress,
				Goal:          goal,
				GoalName:      goal_name,
				Unit:          unit,
			})
		}
	}

	return ProgressByTimeResponse{
		GroupedSessions: sessions,
	}, nil
}

func GetBestSessionsFromPersistence(goalId int, year int, useLow bool) (ProgressResponse, error) {
	db := database.GetDatabase()

	startDate := fmt.Sprintf(`'%d-01-01'`, year)
	endDate := fmt.Sprintf(`'%d-01-01'`, year+1)

	query := `
		select progress.*, goal.goal_name, goal.unit 
		from progress
		left join goal on goal.goal = progress.goal
		where progress.goal = $1
		and session_date between $2 and $3
		and amount = (
			select `

	if useLow {
		query = query + `min(amount)`
	} else {
		query = query + "max(amount)"
	}

	query = query + `
		from progress 
			where goal = $1
			and session_date between $2 and $3
		)
		order by session_date;
	`

	dbProgress, err := db.Query(context.Background(), query, goalId, startDate, endDate)

	if err != nil {
		log.Printf("Error getting best sessions for goal: %d, year: %d, err: %v", goalId, year, err)
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
			log.Printf("Error reading records of best sessions for goal: %d, year: %d, : %v", goalId, year, err)
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
