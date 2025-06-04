package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

const (
	postgresDSN = "postgres://shipager_user:udD4TT3yraudlmXidrNIS6tovbDRdeV6@dpg-d0lf7j3uibrs73a90j8g-a.frankfurt-postgres.render.com/shipager"
	batchSize   = 1000
)

func main() {
	db, err := sql.Open("postgres", postgresDSN)
	if err != nil {
		log.Fatalf("❌ Failed to connect to DB: %v", err)
	}
	defer db.Close()

	doctorIDs, err := getAllDoctorIDs(db)
	if err != nil {
		log.Fatalf("❌ Failed to get doctor IDs: %v", err)
	}

	err = generateSchedules(db, doctorIDs)
	if err != nil {
		log.Fatalf("❌ Failed to generate schedules: %v", err)
	}

	// Обновляем уже прошедшие слоты — делаем их недоступными
	err = markPastSlotsUnavailable(db)
	if err != nil {
		log.Fatalf("❌ Failed to mark past slots unavailable: %v", err)
	}

	fmt.Println("✅ Schedules generated and past slots updated successfully!")
}

// Тип слота
type Slot struct {
	DoctorID  string
	StartTime time.Time
	EndTime   time.Time
}

// Получение всех UUID врачей
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

// Генерация расписаний и вставка в одной транзакции
func generateSchedules(db *sql.DB, doctorIDs []string) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	now := time.Now()
	startDate := now.Truncate(24 * time.Hour)
	endDate := startDate.AddDate(0, 0, 14)

	var allSlots []Slot

	for _, doctorID := range doctorIDs {
		for day := startDate; day.Before(endDate); day = day.AddDate(0, 0, 1) {
			weekday := day.Weekday()
			if weekday == time.Sunday {
				continue
			}

			slots := []struct{ startHour, endHour int }{
				{9, 13}, {14, 18},
			}

			if weekday == time.Saturday {
				slots = []struct{ startHour, endHour int }{
					{9, 13}, {14, 15},
				}
			}

			for _, slot := range slots {
				for hour := slot.startHour; hour < slot.endHour; hour++ {
					start := time.Date(day.Year(), day.Month(), day.Day(), hour, 0, 0, 0, time.Local)
					end := start.Add(time.Hour)
					allSlots = append(allSlots, Slot{
						DoctorID:  doctorID,
						StartTime: start,
						EndTime:   end,
					})
				}
			}
		}
	}

	// Вставка батчами
	for i := 0; i < len(allSlots); i += batchSize {
		end := i + batchSize
		if end > len(allSlots) {
			end = len(allSlots)
		}
		batch := allSlots[i:end]
		if err := insertSlotBatchTx(tx, batch); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit failed: %w", err)
	}

	return nil
}

// Вставка батча слотов в рамках транзакции
func insertSlotBatchTx(tx *sql.Tx, slots []Slot) error {
	query := "INSERT INTO schedule (doctor_id, slot_start, slot_end) VALUES "
	args := make([]interface{}, 0, len(slots)*3)
	valueStrings := make([]string, 0, len(slots))

	for i, s := range slots {
		start := i*3 + 1
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d)", start, start+1, start+2))
		args = append(args, s.DoctorID, s.StartTime, s.EndTime)
	}

	query += fmt.Sprintf("%s ON CONFLICT (doctor_id, slot_start, slot_end) DO NOTHING", joinWithComma(valueStrings))

	_, err := tx.Exec(query, args...)
	return err
}

// Вспомогательная функция для склейки строк
func joinWithComma(parts []string) string {
	result := ""
	for i, p := range parts {
		if i > 0 {
			result += ", "
		}
		result += p
	}
	return result
}

func markPastSlotsUnavailable(db *sql.DB) error {
	query := `
		UPDATE schedule
		SET is_available = FALSE
		WHERE slot_end < NOW()
	`
	_, err := db.Exec(query)
	return err
}
