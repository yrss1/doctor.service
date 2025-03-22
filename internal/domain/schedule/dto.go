package schedule

import "time"

type Request struct {
	ID          string     `json:"id"`
	DoctorID    *string    `json:"doctor_id"`
	SlotStart   *time.Time `json:"slot_start"`
	SlotEnd     *time.Time `json:"slot_end"`
	IsAvailable *bool      `json:"is_available"`
}

type Response struct {
	ID          string    `json:"id"`
	DoctorID    string    `json:"doctor_id"`
	SlotStart   time.Time `json:"slot_start"`
	SlotEnd     time.Time `json:"slot_end"`
	IsAvailable bool      `json:"is_available"`
}

func ParseFromEntity(data Entity) (res Response) {
	res = Response{
		ID:          data.ID,
		DoctorID:    *data.DoctorID,
		SlotStart:   *data.SlotStart,
		SlotEnd:     *data.SlotEnd,
		IsAvailable: *data.IsAvailable,
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
