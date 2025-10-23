package models

type Teacher struct {
	ID        int    `json:"ID,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Email     string `json:"email,omitempty"`
	Class     string `json:"class,omitempty"`
	Subject   string `json:"subject,omitempty"`
}
