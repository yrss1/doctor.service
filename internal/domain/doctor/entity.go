package doctor

type Entity struct {
	ID             string   `json:"id" db:"id"`
	Name           *string  `json:"name" db:"name"`
	Specialization *string  `json:"specialization" db:"specialization"`
	Experience     *string  `json:"experience" db:"experience"`
	Price          *string  `json:"price" db:"price"`
	Rating         *float64 `json:"rating" db:"rating"`
	Address        *string  `json:"address" db:"address"`
	Phone          *string  `json:"phone" db:"phone"`
	ClinicID       *string  `json:"clinic_id" db:"clinic_id"`
}
