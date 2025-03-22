package review

type Request struct {
	ID       string  `json:"id"`
	DoctorID *string `json:"doctor_id"`
	UserID   *string `json:"user_id"`
	Rating   *string `json:"rating"`
	Comment  *string `json:"comment"`
}

type Response struct {
	ID       string  `json:"id"`
	DoctorID *string `json:"doctor_id"`
	UserID   *string `json:"user_id"`
	Rating   *string `json:"rating"`
	Comment  *string `json:"comment"`
}

func ParseFromEntity(data Entity) (res Response) {
	res = Response{
		ID:       data.ID,
		DoctorID: data.DoctorID,
		UserID:   data.UserID,
		Rating:   data.Rating,
		Comment:  data.Comment,
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
