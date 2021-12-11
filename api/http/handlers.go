package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"ws/logger"
	"ws/model"

	"github.com/julienschmidt/httprouter"
)

const (
	USERNAME = "idt"
	PASSWORD = "idt123"
)

type appHandlers struct {
	users map[string]model.User
}

func NewAppHandlers() *appHandlers {
	return &appHandlers{
		users: make(map[string]model.User),
	}
}

func (ah *appHandlers) SetupRoutes() *httprouter.Router {
	router := httprouter.New()
	router.HandlerFunc("GET", "/", logger.Log(ah.IndexHandler))

	router.HandlerFunc("POST", "/users", logger.Log(validator(ah.SetUser)))
	router.HandlerFunc("GET", "/users/:id", logger.Log(validator(ah.GetUser)))
	router.HandlerFunc("GET", "/users", logger.Log(validator(ah.GetUsers)))

	return router
}

func validator(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if errMsg := validateBasicAuth(r); errMsg == "" {
			f(w, r)
		} else {
			http.Error(w, errMsg, http.StatusUnauthorized)
		}
	}
}

func validateBasicAuth(r *http.Request) (errMsg string) {
	u, p, ok := r.BasicAuth()
	if !ok {
		errMsg = "check if BasicAuthroization header passed"
		return
	}
	if u != USERNAME {
		errMsg = "check if you've set right username"
		return
	}
	if p != PASSWORD {
		errMsg = fmt.Sprintf("check if you've set right password for: %s", u)
		return
	}

	return ""
}

func validateUser(u model.User) (errorMessage string) {
	if u.ID == "" {
		errorMessage = "id field can not be empty"
		return
	}

	if u.Name == "" {
		errorMessage = "name field can not be empty"
		return
	}
	if time.Time.IsZero(u.SignupTime) == true {
		errorMessage = "signupTime field can not be empty"
		return
	}

	return
}

func (ah *appHandlers) SetUser(w http.ResponseWriter, r *http.Request) {
	user := model.User{}

	//decode body
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return

	}

	// validate user input
	if errMsg := validateUser(user); errMsg != "" {
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	// if user already existed
	if _, ok := ah.users[user.ID]; ok {
		http.Error(w, fmt.Sprintf("user with id: %s is existing", user.ID), http.StatusBadRequest)
		return
	}

	ah.users[user.ID] = user

	w.WriteHeader(http.StatusCreated)
}

func (ah *appHandlers) GetUser(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/users/")
	if id == "" {
		http.Error(w, "id parameter can not be empty", http.StatusBadRequest)
		return
	}

	if user, ok := ah.users[id]; ok {
		uJ, err := json.Marshal(&user)
		if err != nil {
			http.Error(w, fmt.Sprintf("internal server error in getting the user with id: %s", id), http.StatusInternalServerError)
			return
		}

		w.Write(uJ)
		return
	} else {
		http.Error(w, fmt.Sprintf(" user with id: %s not found", id), http.StatusNotFound)
		return
	}
}

func (ah *appHandlers) GetUsers(w http.ResponseWriter, r *http.Request) {

	if len(ah.users) == 0 {
		usJ, _ := json.Marshal([]model.User{})
		w.Write(usJ)
		return
	}

	users := []model.User{}

	// convert map of users to the users slice
	for _, u := range ah.users {
		users = append(users, u)
	}

	if len(users) > 0 {
		usJ, err := json.Marshal(&users)
		if err != nil {
			http.Error(w, "internal server error in getting users", http.StatusInternalServerError)
			return
		}

		w.Write(usJ)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (ah *appHandlers) IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.ServeFile(w, r, "public/index.html")
	}
}
