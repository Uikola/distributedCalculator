package entity

// ComputingResource структура для работы с вычислительными ресурсами.
type ComputingResource struct {
	ID       int64
	Name     string
	Task     *string
	Occupied bool
}
