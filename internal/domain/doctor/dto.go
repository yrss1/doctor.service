package doctor

type Request struct {
	ID             string   `json:"id"`
	Name           *string  `json:"name"`
	Specialization *string  `json:"specialization"`
	Experience     *string  `json:"experience"`
	Price          *string  `json:"price"`
	Rating         *float64 `json:"rating"`
	Address        *string  `json:"address"`
	Phone          *string  `json:"phone"`
	Gender         *string  `json:"gender"`
	VisitType      *string  `json:"visit_type"`
	ClinicName     *string  `json:"clinic_name"`
}

type Response struct {
	ID                 string         `json:"id"`
	Name               string         `json:"name"`
	Specialization     string         `json:"specialization"`
	Experience         string         `json:"experience"`
	Price              string         `json:"price"`
	Rating             float64        `json:"rating"`
	Address            string         `json:"address"`
	Phone              string         `json:"phone"`
	Gender             string         `json:"gender"`
	VisitType          string         `json:"visit_type"`
	ClinicName         string         `json:"clinic_name"`
	AvailableSchedules []ScheduleSlot `json:"available_schedules"`
}

func ParseFromEntity(data Entity) (res Response) {
	res = Response{
		ID:                 data.ID,
		Name:               *data.Name,
		Specialization:     *data.Specialization,
		Experience:         *data.Experience,
		Price:              *data.Price,
		Rating:             *data.Rating,
		Address:            *data.Address,
		Phone:              *data.Phone,
		Gender:             *data.Gender,
		VisitType:          *data.VisitType,
		ClinicName:         *data.ClinicName,
		AvailableSchedules: data.AvailableSchedules,
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
