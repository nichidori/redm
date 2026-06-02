package rest

type PaginationMeta struct {
	TotalCount int `json:"total_count"`
	Offset     int `json:"offset"`
	Limit      int `json:"limit"`
}

type IDName struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CustomField struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Value    any    `json:"value"`
	Multiple bool   `json:"multiple,omitempty"`
}

type APIError struct {
	Errors []string `json:"errors"`
}

func (e APIError) Error() string {
	if len(e.Errors) == 0 {
		return "unknown API error"
	}
	return e.Errors[0]
}
