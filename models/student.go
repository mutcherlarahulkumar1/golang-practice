package models

type Student struct {
	ID        int    `json:"ID,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Email     string `json:"email,omitempty"`
	Class     string `json:"class,omitempty"`
}

type StudentPatch struct {
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Email     string `json:"email,omitempty"`
	Class     string `json:"class,omitempty"`
}
