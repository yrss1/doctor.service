package appointment

import "time"

type Entity struct {
	ID         string  `json:"id" db:"id"`
	DoctorID   *string `json:"doctor_id" db:"doctor_id"`
	UserID     *string `json:"user_id" db:"user_id"`
	ScheduleID *string `json:"schedule_id" db:"schedule_id"`
	Status     *string `json:"status" db:"status"`
	MeetingURL *string `json:"meeting_url" db:"meeting_url"`
}

type EntityView struct {
	ID        string    `json:"id" db:"appointment_id"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	SlotStart time.Time `json:"slot_start" db:"slot_start"`
	SlotEnd   time.Time `json:"slot_end" db:"slot_end"`

	DoctorID       string  `json:"doctor_id" db:"doctor_id"`
	DoctorName     string  `json:"doctor_name" db:"doctor_name"`
	Specialization string  `json:"specialization" db:"specialization"`
	DoctorPhone    string  `json:"doctor_phone" db:"doctor_phone"`
	DoctorGender   string  `json:"doctor_gender" db:"doctor_gender"`
	MeetingURL     *string `json:"meeting_url" db:"meeting_url"`
}
