package domain

type User struct {
	ID         string
	Role       string
	BusinessID string
}

type Batch struct {
	ID         string
	BusinessID string
	Status     string
}
