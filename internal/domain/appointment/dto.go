package appointment

type Request struct {
	ID         string  `json:"id"`
	DoctorID   *string `json:"doctor_id"`
	UserID     *string `json:"user_id"`
	ScheduleID *string `json:"schedule_id"`
	Status     *string `json:"status"`
	MeetingURL *string `json:"meeting_url"`
}

type Response struct {
	ID         string  `json:"id"`
	DoctorID   *string `json:"doctor_id"`
	UserID     *string `json:"user_id"`
	ScheduleID *string `json:"schedule_id"`
	Status     *string `json:"status"`
	MeetingURL *string `json:"meeting_url"`
}

func ParseFromEntity(data Entity) (res Response) {
	res = Response{
		ID:         data.ID,
		DoctorID:   data.DoctorID,
		UserID:     data.UserID,
		ScheduleID: data.ScheduleID,
		Status:     data.Status,
		MeetingURL: data.MeetingURL,
	}

	return
}

func ParseFromEntities(data []Entity) (res []Response) {
	res = make([]Response, 0)
	for _, obj := range data {
		res = append(res, ParseFromEntity(obj))
	}
	return
}
