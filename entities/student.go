package entities

type Student struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	DOB     string `json:"DOB"`
	Branch  string `json:"branch"`
	Phone   string `json:"phone"`
	Company string `json:"company"`
	Status  string `json:"status"`
}
