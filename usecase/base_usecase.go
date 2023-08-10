package usecase

type BaseUsecase[T any] interface {
	FindAll() ([]T, error)
	FindById(id string) (*T, error)
	SaveData(payload *T) error
	DeleteData(id string) error
}

type BaseSearchUsecase[T any] interface {
	SearchBy(by map[string]interface{}) ([]T, error)
}
