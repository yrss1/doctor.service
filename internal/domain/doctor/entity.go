package doctor

type Entity struct {
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
