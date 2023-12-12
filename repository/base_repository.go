package repository

type BaseRepository[T any] interface {
	Save(payload *T) error
	Get(id string) (*T, error)
	List() ([]T, error)
	Delete(id string) error
}

type BaseSearch[T any] interface {
	Search(by map[string]any) ([]T, error)
}
