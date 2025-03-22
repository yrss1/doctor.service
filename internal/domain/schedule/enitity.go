package schedule

import "time"

type Entity struct {
	ID          string     `json:"id" db:"id"`
	DoctorID    *string    `json:"doctor_id" db:"doctor_id"`
	SlotStart   *time.Time `json:"slot_start" db:"slot_start"`
	SlotEnd     *time.Time `json:"slot_end" db:"slot_end"`
	IsAvailable *bool      `json:"is_available" db:"is_available"`
}
