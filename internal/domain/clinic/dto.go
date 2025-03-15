package clinic

type Request struct {
	ID      string  `json:"id"`
	Name    *string `json:"name" `
	Address *string `json:"address" `
	Phone   *string `json:"phone""`
}

type Response struct {
	ID      string `json:"id"`
	Name    string `json:"name" `
	Address string `json:"address" `
	Phone   string `json:"phone""`
}

func ParseFromEntity(data Entity) (res Response) {
	res = Response{
		ID:      data.ID,
		Name:    *data.Name,
		Address: *data.Address,
		Phone:   *data.Phone,
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
