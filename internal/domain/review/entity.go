package review

type Entity struct {
	ID       string  `json:"id" db:"id"`
	DoctorID *string `json:"doctor_id" db:"doctor_id"`
	UserID   *string `json:"user_id" db:"user_id"`
	Rating   *string `json:"rating" db:"rating"`
	Comment  *string `json:"comment" db:"comment"`
}
