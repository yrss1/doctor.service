package doctor

type Entity struct {
	ID                 string         `json:"id" db:"id"`
	Name               *string        `json:"name" db:"name"`
	Specialization     *string        `json:"specialization" db:"specialization"`
	Experience         *string        `json:"experience" db:"experience"`
	Price              *string        `json:"price" db:"price"`
	Rating             *float64       `json:"rating" db:"rating"`
	Address            *string        `json:"address" db:"address"`
	Phone              *string        `json:"phone" db:"phone"`
	PhotoURL           *string        `json:"photo_url" db:"photo_url"`
	ClinicName         *string        `json:"clinic_name" db:"clinic_name"`
	AvailableSchedules []ScheduleSlot `json:"available_schedules" db:"available_schedules"`
}

type ScheduleSlot struct {
	ScheduleID string `json:"schedule_id"`
	SlotStart  string `json:"slot_start"`
	SlotEnd    string `json:"slot_end"`
}

type EntityWithRaw struct {
	Entity
	AvailableSchedulesRaw []byte `db:"available_schedules"`
}
