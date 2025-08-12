package domain

type Repository interface {
	Create(user User) (User, error)
	Delete(userID uint) error
}
