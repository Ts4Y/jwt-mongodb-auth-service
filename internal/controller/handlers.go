package controller

import (
	"fmt"
	"jwt-mongo-auth/internal/service"
	"net/http"

	"github.com/sirupsen/logrus"
)

type Controller struct {
	log               *logrus.Logger
	controllerService *service.Service
	mux               *http.ServeMux
}

func NewController(log *logrus.Logger, serv *service.Service) *Controller {
	return &Controller{
		log:               log,
		controllerService: serv,
		mux:               http.NewServeMux(),
	}
}

func (c *Controller) Run() {
	c.mux.HandleFunc("/login", c.GenerateTokens)
	c.mux.HandleFunc("/update",c.UpdateTokens)

	c.log.Infoln("Сервер успешно запущен на порту 9000")
	if err := http.ListenAndServe(":9000", c.mux); err != nil {
		c.log.Fatalln("Не удалось начать прослушивание,ошибка", err)
	}
}

func (c *Controller) GenerateTokens(w http.ResponseWriter, r *http.Request) {
	guid := r.URL.Query().Get("guid")

	tokens, err := c.controllerService.GenerateTokens(guid)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "AccessToken",
		Value: tokens.AccessToken,
	})

	http.SetCookie(w, &http.Cookie{
		Name:  "RefreshToken",
		Value: tokens.RefreshToken,
	})

	fmt.Fprint(w, "Success")
}

func (c *Controller) UpdateTokens(w http.ResponseWriter, r *http.Request) {
	guid := r.URL.Query().Get("guid")

	cookie, err := r.Cookie("RefreshToken")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, "Cookie с именем 'RefreshToken' не найден.", http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	refreshToken := cookie.Value

	tokens, err := c.controllerService.UpdateTokens(r.Context(), guid, refreshToken)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "AccessToken",
		Value: tokens.AccessToken,
	})

	http.SetCookie(w, &http.Cookie{
		Name:  "RefreshToken",
		Value: tokens.RefreshToken,
	})

	fmt.Fprint(w, "Success")
}
