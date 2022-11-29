package entities

type Student struct {
	ID      string  `json:"id,omitempty"`
	Name    string  `json:"name,omitempty"`
	DOB     string  `json:"DOB,omitempty"`
	Branch  string  `json:"branch,omitempty"`
	Phone   string  `json:"phone,omitempty"`
	Company Company `json:"company"`
	Status  string  `json:"status,omitempty"`
}
