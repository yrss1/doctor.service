package doctor

type Request struct {
	ID           int      `json:"id"`
	Name         *string  `json:"name"`
	Specialty    *string  `json:"specialty"`
	Experience   *int     `json:"experience"`
	Price        *float64 `json:"price"`
	Address      *string  `json:"address"`
	ClinicName   *string  `json:"clinic_name"`
	Phone        *string  `json:"phone"`
	Email        *string  `json:"email"`
	PhotoURL     *string  `json:"photo_url"`
	Education    *string  `json:"education"`
	Rating       *float64 `json:"rating"`
	ReviewsCount *int     `json:"reviews_count"`
	IsActive     *bool    `json:"is_active"`
}

type Response struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Specialty    string  `json:"specialty"`
	Experience   int     `json:"experience"`
	Price        float64 `json:"price"`
	Address      string  `json:"address"`
	ClinicName   string  `json:"clinic_name"`
	Phone        string  `json:"phone"`
	Email        string  `json:"email"`
	PhotoURL     string  `json:"photo_url"`
	Education    string  `json:"education"`
	Rating       float64 `json:"rating"`
	ReviewsCount int     `json:"reviews_count"`
	IsActive     bool    `json:"is_active"`
}

func ParseFromEntity(data Entity) (res Response) {
	res = Response{
		ID:           data.ID,
		Name:         *data.Name,
		Specialty:    *data.Specialty,
		Experience:   *data.Experience,
		Price:        *data.Price,
		Address:      *data.Address,
		ClinicName:   *data.ClinicName,
		Phone:        *data.Phone,
		Email:        *data.Email,
		PhotoURL:     *data.PhotoURL,
		Education:    *data.Education,
		Rating:       *data.Rating,
		ReviewsCount: *data.ReviewsCount,
		IsActive:     *data.IsActive,
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
