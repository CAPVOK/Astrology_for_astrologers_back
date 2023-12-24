package api

import (
	"bytes"
	"crypto/rand"
	"errors"
	"net/http"
	"net/mail"
	"space/internal/app/ds"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
)

type UserSingleton struct {
	userID  int
	isAdmin bool
}

func loadUserData() UserSingleton {
	userData := UserSingleton{
		userID:  1,
		isAdmin: false,
	}
	return userData
}

func singleton() (int, bool, error) {
	userData := loadUserData()
	return userData.userID, userData.isAdmin, nil
}

// Login godoc
//
//	@Summary		login user
//	@Description	create user session and put it into cookie
//	@Tags			auth
//	@Accept			json
//	@Param			body	body		ds.Credentials	true	"user credentials"
//	@Success		200		{object}	object{body=object{id=int}}
//	@Failure		400		{object} object{status=string, message=string}
//	@Failure		401		{object} object{status=string, message=string}
//	@Failure		404		{object} object{status=string, message=string}
//	@Failure		409		{object} object{status=string, message=string}
//	@Failure		500		{object} object{status=string, message=string}
//	@Router			/auth/login [post]
func (h *Handler) Login(ctx *gin.Context) {
	auth, err := h.auth(ctx)
	if auth == true {
		ctx.JSON(ds.GetHttpStatusCode(err), gin.H{
			"status":  "fail",
			"message": "Вы должны быть авторизованы\"",
		})
		return
	}

	//if err != nil {
	//	ds.WriteError(w, err.Error(), http.StatusBadRequest)
	//	return
	//}
	//defer h.CloseAndAlert(r.Body)

	var c ds.Credentials
	if err := ctx.ShouldBindJSON(&c); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Неверное тело запроса",
		})
		return
	}

	if err = checkCredentials(c); err != nil {
		ctx.JSON(ds.GetHttpStatusCode(err), gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}
	c.Email = strings.TrimSpace(c.Email)

	expectedUser, err := h.Repo.GetByEmail(c.Email)
	if err != nil {
		ctx.JSON(ds.GetHttpStatusCode(err), gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	if !checkPasswords(expectedUser.Password, c.Password) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Логин или пароль неверны",
		})
		return
	}

	session := ds.Session{
		Token:     uuid.NewString(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
		UserID:    int(expectedUser.ID),
		Role:      expectedUser.Role,
	}
	if err = h.RedisRepo.Add(session); err != nil {
		ctx.JSON(ds.GetHttpStatusCode(err), gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	//ctx.SetCookie("session_token", session.Token, int(session.ExpiresAt.Unix()), "/", "localhost", false, true)
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "session_token",
		Value:    session.Token,
		Expires:  session.ExpiresAt,
		Path:     "/",
		HttpOnly: true,
	})

	ctx.JSON(http.StatusOK, gin.H{
		"id": expectedUser.ID,
	})
}

// Logout godoc
//
//	@Summary		logout user
//	@Description	delete current session and nullify cookie
//	@Tags			auth
//	@Success		204
//	@Failure		400	{object}	object{err=string}
//	@Failure		401	{object}	object{err=string}
//	@Failure		404	{object}	object{err=string}
//	@Failure		409	{object}	object{err=string}
//	@Failure		500	{object}	object{err=string}
//	@Router			/auth/logout [post]
func (h *Handler) Logout(ctx *gin.Context) {
	auth, err := h.auth(ctx)
	if auth != true {
		ctx.JSON(ds.GetHttpStatusCode(err), gin.H{
			"status":  "fail",
			"message": "Вы должны быть авторихованы",
		})
		return
	}

	t, err := ctx.Cookie("session_token")

	if err = h.RedisRepo.DeleteByToken(t); err != nil {
		ctx.JSON(ds.GetHttpStatusCode(err), gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return

	}

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now(),
		Path:     "/",
		HttpOnly: true,
	})

	ctx.Status(http.StatusNoContent)
}

// Register godoc
//
//	@Summary		register user
//	@Description	add new user to db and return it id
//	@Tags			auth
//	@Produce		json
//	@Accept			json
//	@Param			body	body		ds.Credentials	true	"user credentials"
//	@Success		200		{object}	object{body=object{id=int}}
//	@Failure		400		{object} object{status=string, message=string}
//	@Failure		401		{object} object{status=string, message=string}
//	@Failure		409		{object} object{status=string, message=string}
//	@Failure		500		{object} object{status=string, message=string}
//	@Router			/auth/register [post]
func (h *Handler) Register(ctx *gin.Context) {
	auth, err := h.auth(ctx)
	if auth == true {
		ctx.JSON(ds.GetHttpStatusCode(err), gin.H{
			"status":  "fail",
			"message": "Вы должны быть неавторизованы",
		})
		return
	}

	var user ds.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Неверное тело запроса",
		})
		return
	}
	user.Email = strings.TrimSpace(user.Email)
	if err = checkCredentials(ds.Credentials{Email: user.Email, Password: user.Password}); err != nil {
		ctx.JSON(ds.GetHttpStatusCode(err), gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}
	user.Role = ds.Usr

	salt := make([]byte, 8)
	rand.Read(salt)
	user.Password = HashPassword(salt, user.Password)
	var id int
	if id, err = h.Repo.AddUser(user); err != nil {
		ctx.JSON(ds.GetHttpStatusCode(err), gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	session := ds.Session{
		Token:     uuid.NewString(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
		UserID:    id,
		Role:      ds.Usr,
	}
	if err = h.RedisRepo.Add(session); err != nil {
		ctx.JSON(ds.GetHttpStatusCode(err), gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	//ctx.SetCookie("session_token", session.Token, int(session.ExpiresAt.Unix()), "/", "localhost", false, true)
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "session_token",
		Value:    session.Token,
		Expires:  session.ExpiresAt,
		Path:     "/",
		HttpOnly: true,
	})

	ctx.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

func (h *Handler) auth(ctx *gin.Context) (bool, error) {
	c, err := ctx.Request.Cookie("session_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return false, ds.ErrUnauthorized
		}

		return false, ds.ErrBadRequest
	}
	if c.Expires.After(time.Now()) {
		return false, ds.ErrUnauthorized
	}
	sessionToken := c.Value
	sc, err := h.RedisRepo.SessionExists(sessionToken)
	if err != nil {
		return false, err
	}
	if sc.UserID == 0 {
		return false, ds.ErrUnauthorized
	}

	return true, ds.ErrAlreadyExists
}

func HashPassword(salt []byte, password []byte) []byte {
	hashedPass := argon2.IDKey(password, salt, 1, 64*1024, 4, 32)
	return append(salt, hashedPass...)
}

func checkPasswords(passHash []byte, plainPassword []byte) bool {
	salt := passHash[0:8]
	userPassHash := HashPassword(salt, plainPassword)
	return bytes.Equal(userPassHash, passHash)
}

func valid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func checkCredentials(cred ds.Credentials) error {
	if cred.Email == "" || len(cred.Password) == 0 {
		return ds.ErrWrongCredentials
	}

	if !valid(cred.Email) {
		return ds.ErrWrongCredentials
	}

	return nil
}
