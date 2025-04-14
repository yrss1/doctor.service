package room

type Entity struct {
	AppointmentID string `json:"appointment_id"`
	UserID        string `json:"user_id"`
	DoctorID      string `json:"doctor_id"`
}
