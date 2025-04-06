package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

const (
	postgresDSN = "your-dsn-here" // Replace with your real DSN
)

func main() {
	db, err := sql.Open("postgres", postgresDSN)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	doctorIDs, err := getAllDoctorIDs(db)
	if err != nil {
		log.Fatalf("Failed to get doctor IDs: %v", err)
	}

	err = generateSchedules(db, doctorIDs)
	if err != nil {
		log.Fatalf("Failed to generate schedules: %v", err)
	}

	fmt.Println("✅ Schedules generated successfully!")
}

// 📥 Load all doctor UUIDs from DB
func getAllDoctorIDs(db *sql.DB) ([]string, error) {
	rows, err := db.Query(`SELECT id FROM doctors`)
	if err != nil {
		return nil, fmt.Errorf("query doctor ids failed: %w", err)
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan doctor id failed: %w", err)
		}
		ids = append(ids, id)
	}

	return ids, nil
}

// 📅 Generate 2 weeks of schedule per doctor
func generateSchedules(db *sql.DB, doctorIDs []string) error {
	now := time.Now()
	startDate := now.Truncate(24 * time.Hour)
	endDate := startDate.AddDate(0, 0, 14) // 2 weeks

	for _, doctorID := range doctorIDs {
		for day := startDate; day.Before(endDate); day = day.AddDate(0, 0, 1) {
			weekday := day.Weekday()

			// Skip Sundays
			if weekday == time.Sunday {
				continue
			}

			slots := []struct {
				startHour int
				endHour   int
			}{
				{9, 13},
				{14, 18},
			}

			if weekday == time.Saturday {
				slots = []struct {
					startHour int
					endHour   int
				}{
					{9, 13},
					{14, 15},
				}
			}

			for _, slot := range slots {
				for hour := slot.startHour; hour < slot.endHour; hour++ {
					start := time.Date(day.Year(), day.Month(), day.Day(), hour, 0, 0, 0, time.Local)
					end := start.Add(time.Hour)

					_, err := db.Exec(`
						INSERT INTO schedule (doctor_id, slot_start, slot_end)
						VALUES ($1, $2, $3)
						ON CONFLICT (doctor_id, slot_start, slot_end) DO NOTHING
					`, doctorID, start, end)

					if err != nil {
						return fmt.Errorf("insert failed for doctor %s at %s: %w", doctorID, start, err)
					}
				}
			}
		}
	}
	return nil
}
