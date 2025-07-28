package user

import (
	"encoding/json"
	"log"
	hdr "marketplace/internal/handlers"
	"marketplace/internal/token"
	"marketplace/internal/user"
	"marketplace/internal/utils"
	"net/http"
	"regexp"
	"strings"
)

const (
	// minLoginLen — минимальная длина логина
	minLoginLen int = 3
	// minLoginLen — максимальная длина логина
	maxLoginLen int = 20
	// minPasswordLen — минимальная длина пароля
	minPasswordLen int = 8
	// maxPasswordLen — максимальная длина пароля
	maxPasswordLen int = 30
)

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// SignIn авторизует уже зарегистрированного пользователя
func (hnd *UserHandler) SignIn(wrt http.ResponseWriter, rqt *http.Request) {
	usr, _ := ProcessRequest(wrt, rqt)
	if usr == nil {
		return
	}

	user, code, err := hnd.UserRepo.SignIn(usr)
	if err != nil {
		log.Println(err)
	}

	switch code {
	case hdr.UnauthorizedCode:
		errSend := hdr.SendUnauthorized(wrt, err.Error())
		if errSend != nil {
			log.Printf("error while sending the unauthorized error message: %v\n", errSend)
		}
		return

	case hdr.InternalServerErrorCode:
		errSend := hdr.SendInternalServerError(wrt, err.Error())
		if errSend != nil {
			log.Printf("error while sending the internal server error message: %v\n", errSend)
		}
		return
	}

	ProcessToken(wrt, rqt, user)
}

// SignUp регистрирует нового пользователя
func (hnd *UserHandler) SignUp(wrt http.ResponseWriter, rqt *http.Request) {
	usr, arq := ProcessRequest(wrt, rqt)
	if usr == nil || arq == nil {
		return
	}

	isErr := validateUsername(wrt, arq.Username)
	if !isErr {
		return
	}

	isErr = validatePassword(wrt, arq.Password)
	if !isErr {
		return
	}

	user, code, err := hnd.UserRepo.SignUp(usr)
	if err != nil {
		log.Println(err)
	}

	switch code {
	case hdr.BadRequestCode:
		errStr := "ошибка: такой логин уже занят"
		errSend := hdr.SendBadReq(wrt, errStr)
		if errSend != nil {
			log.Printf("error while sending the bad request error message: %v\n", errSend)
		}
		return

	case hdr.InternalServerErrorCode:
		errSend := hdr.SendInternalServerError(wrt, err.Error())
		if errSend != nil {
			log.Printf("error while sending the internal server error message: %v\n", errSend)
		}
		return
	}

	ProcessToken(wrt, rqt, user)
	errJSON := json.NewEncoder(wrt).Encode(user)
	if errJSON != nil {
		log.Printf("error while sending response body: %v\n", errJSON)
	}
}

func ProcessRequest(wrt http.ResponseWriter, rqt *http.Request) (*user.User, *AuthRequest) {
	var arq AuthRequest
	err := json.NewDecoder(rqt.Body).Decode(&arq)
	if err != nil {
		errSend := hdr.SendBadReq(wrt, "error while deserialization")
		if errSend != nil {
			log.Printf("error while sending the bad request error message: %v\n", errSend)
		}
		return nil, nil
	}

	if arq.Username == "" && arq.Password == "" {
		errSend := hdr.SendBadReq(wrt, "ошибка: пользователь не отправил логин и пароль")
		if errSend != nil {
			log.Printf("error while sending the bad request error message: %v\n", errSend)
		}
		return nil, nil
	}

	if arq.Username == "" {
		errSend := hdr.SendBadReq(wrt, "ошибка: пользователь не отправил логин")
		if errSend != nil {
			log.Printf("error while sending the bad request error message: %v\n", errSend)
		}
		return nil, nil
	}

	if arq.Password == "" {
		errSend := hdr.SendBadReq(wrt, "ошибка: пользователь не отправил пароль")
		if errSend != nil {
			log.Printf("error while sending the bad request error message: %v\n", errSend)
		}
		return nil, nil
	}

	usr := &user.User{
		Username: arq.Username,
		Password: arq.Password,
	}
	return usr, &arq
}

func ProcessToken(wrt http.ResponseWriter, rqt *http.Request, thisUser *user.User) {
	tokenString, errToken := token.CreateJWTtoken(thisUser.Username)
	if errToken != nil {
		errSend := hdr.SendInternalServerError(wrt, errToken.Error())
		if errSend != nil {
			log.Printf("error while sending the internal server error message: %v\n", errSend)
		}

		log.Println("error while creating the JWT token: ", errToken)
	}

	wrt.Header().Set("Authorization", tokenString)
	wrt.Header().Set("Content-Type", "application/json")
	wrt.WriteHeader(http.StatusOK)
}

// validateUsername валидирует логин
func validateUsername(wrt http.ResponseWriter, username string) bool {
	check := utils.CheckLen(username, "недостаточная", "превышена допустимая", "логина", "логин", minLoginLen, maxLoginLen)
	if check != "" {
		errSend := hdr.SendBadReq(wrt, check)
		if errSend != nil {
			log.Printf("error while sending the bad request error message: %v\n", errSend)
		}
		return false
	}

	re := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !re.MatchString(username) {
		errSend := hdr.SendBadReq(wrt, "ошибка: недопустимый/ые символ/ы в логине -> логин может содержать только буквы, цифры, '-' и '_'")
		if errSend != nil {
			log.Printf("error while sending the bad request error message: %v\n", errSend)
		}
		return false
	}
	return true
}

// validatePassword валидирует пароль
func validatePassword(wrt http.ResponseWriter, password string) bool {
	check := utils.CheckLen(password, "недостаточная", "превышена допустимая", "пароля", "пароль", minPasswordLen, maxPasswordLen)
	if check != "" {
		errSend := hdr.SendBadReq(wrt, check)
		if errSend != nil {
			log.Printf("error while sending the bad request error message: %v\n", errSend)
		}
		return false
	}

	hasUpper := false
	hasLower := false
	hasNumber := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasNumber = true
		case strings.ContainsRune("!@#$%^&*()-_=+[]{}|;:,.<>?/", char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		errSend := hdr.SendBadReq(wrt, "ошибка: в пароле отсутствует хотя бя одна заглавная буква")
		if errSend != nil {
			log.Printf("error while sending the bad request error message: %v\n", errSend)
		}
		return false
	}

	if !hasLower {
		errSend := hdr.SendBadReq(wrt, "ошибка: в пароле отсутствует хотя бя одна строчная буква")
		if errSend != nil {
			log.Printf("error while sending the bad request error message: %v\n", errSend)
		}
		return false
	}

	if !hasNumber {
		errSend := hdr.SendBadReq(wrt, "ошибка: в пароле отсутствует хотя бя одна цифра")
		if errSend != nil {
			log.Printf("error while sending the bad request error message: %v\n", errSend)
		}
		return false
	}

	if !hasSpecial {
		errSend := hdr.SendBadReq(wrt, "ошибка: в пароле отсутствует хотя бя один специальный символ")
		if errSend != nil {
			log.Printf("error while sending the bad request error message: %v\n", errSend)
		}
		return false
	}
	return true
}
