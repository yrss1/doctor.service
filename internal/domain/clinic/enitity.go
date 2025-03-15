package clinic

type Entity struct {
	ID      string  `json:"id" db:"id"`
	Name    *string `json:"name" db:"name"`
	Address *string `json:"address" db:"address"`
	Phone   *string `json:"phone" db:"phone"`
}
