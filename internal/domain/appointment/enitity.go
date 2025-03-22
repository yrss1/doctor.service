package appointment

type Entity struct {
	ID         string  `json:"id" db:"id"`
	DoctorID   *string `json:"doctor_id" db:"doctor_id"`
	UserID     *string `json:"user_id" db:"user_id"`
	ScheduleID *string `json:"schedule_id" db:"schedule_id"`
	Status     *string `json:"status" db:"status"`
}
