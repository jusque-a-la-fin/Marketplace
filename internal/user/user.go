package user

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserRepo interface {
	// GetUserID получает идентификатор пользователя
	GetUserID(username string) (string, error)
	// SignIn авторизует уже зарегистрированного пользователя
	SignIn(usr *User) (*User, int, error)
	// SignUp регистрирует нового пользователя
	SignUp(usr *User) (*User, int, error)
}
