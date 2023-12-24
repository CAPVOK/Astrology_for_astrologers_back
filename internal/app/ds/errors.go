package ds

import (
	"errors"
	"net/http"
)

var (
	ErrInternalServerError = errors.New("Внутренняя ошибка сервера")
	ErrNotFound            = errors.New("Не найдено")
	ErrBadRequest          = errors.New("Ошибка при отправке запроса")
	ErrUnauthorized        = errors.New("Неавторизован")
	ErrWrongCredentials    = errors.New("Имя пользователя или пароль неверны")
	ErrInvalidToken        = errors.New("Неверные токены сессии")
	ErrAlreadyExists       = errors.New("Уже существует")
	ErrOutOfRange          = errors.New("Id неверен")
	ErrWrongUser           = errors.New("У вас недостаточно прав")
)

func GetHttpStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	switch err {
	case ErrWrongCredentials:
		return http.StatusBadRequest
	case ErrUnauthorized:
		return http.StatusForbidden
	case ErrWrongUser:
		return http.StatusForbidden

	case ErrInvalidToken:
		return http.StatusBadRequest
	case ErrBadRequest:
		return http.StatusBadRequest

	case ErrNotFound:
		return http.StatusNotFound
	case ErrOutOfRange:
		return http.StatusNotFound

	case ErrAlreadyExists:
		return http.StatusConflict

	default:
		return http.StatusInternalServerError
	}
}
