package http

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"ws/model"
)

func TestSetUser(t *testing.T) {

	tt := []struct {
		name       string
		method     string
		input      model.User
		want       string
		statusCode int
	}{
		{
			name:   "with a valid user",
			method: http.MethodPost,
			input: model.User{
				ID:         "ID1",
				Name:       "Name1",
				SignupTime: time.Now(),
			},

			want:       "",
			statusCode: http.StatusCreated,
		},
		{
			name:   "with a user already existing",
			method: http.MethodPost,
			input: model.User{
				ID:         "ID1",
				Name:       "Name1",
				SignupTime: time.Now(),
			},

			want:       "",
			statusCode: http.StatusBadRequest,
		},

		{
			name:   "without Id",
			method: http.MethodPost,
			input: model.User{
				ID:         "",
				Name:       "Name1",
				SignupTime: time.Now(),
			},

			want:       "",
			statusCode: http.StatusBadRequest,
		},
		{
			name:   "without Name",
			method: http.MethodPost,
			input: model.User{
				ID:         "ID1",
				Name:       "",
				SignupTime: time.Now(),
			},

			want:       "",
			statusCode: http.StatusBadRequest,
		},
		{
			name:   "without signup time",
			method: http.MethodPost,
			input: model.User{
				ID:   "ID1",
				Name: "Name",
			},

			want:       "",
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range tt {

		t.Run(tc.name, func(t *testing.T) {
			ah := NewAppHandlers()
			if tc.name == "with a user already existing" {
				ah.users[tc.input.ID] = tc.input
			}
			rec := httptest.NewRecorder()
			userJson, _ := json.Marshal(&tc.input)
			req, err := http.NewRequest(
				tc.method,
				"/users",
				bytes.NewBuffer(userJson),
			)

			if err != nil {
				t.Fatalf("Could not create a request %v", err)
			}

			ah.SetUser(rec, req)

			if rec.Code != tc.statusCode {
				b, _ := ioutil.ReadAll(rec.Body)
				t.Logf("error message: %v\n", b)
				t.Errorf("expected status %v, received status %v", tc.statusCode, rec.Code)
			}

		})
	}
}

func TestGetUser(t *testing.T) {

	tt := []struct {
		name       string
		method     string
		ID         string
		want       model.User
		statusCode int
	}{
		{
			name:   "with valid id",
			method: http.MethodGet,
			ID:     "ID1",

			want: model.User{
				ID:         "ID1",
				Name:       "ID1Name",
				SignupTime: time.Now(),
			},
			statusCode: http.StatusOK,
		},
		{
			name:   "with inexisting Id",
			method: http.MethodGet,
			ID:     "ID4",

			want:       model.User{},
			statusCode: http.StatusNotFound,
		},
		{
			name:   "with invalid Id",
			method: http.MethodGet,
			ID:     "",

			want:       model.User{},
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range tt {

		t.Run(tc.name, func(t *testing.T) {
			ah := NewAppHandlers()

			if tc.ID != "" && tc.statusCode != http.StatusNotFound {
				ah.users[tc.ID] = tc.want
			}

			rec := httptest.NewRecorder()
			req, err := http.NewRequest(
				tc.method,
				"/users"+"/"+tc.ID,
				nil,
			)

			if err != nil {
				t.Fatalf("Could not create a request %v", err)
			}

			ah.GetUser(rec, req)

			if rec.Code != tc.statusCode {
				t.Errorf("expected status %v, received status %v", tc.statusCode, rec.Code)
				return
			}

			// compare response only for valid result
			if rec.Code == http.StatusOK {
				want, _ := json.Marshal(&tc.want)
				resp := rec.Result()
				body, _ := ioutil.ReadAll(resp.Body)
				if bytes.Compare(body, want) != 0 {
					t.Errorf("expected response  %s, received response %s", want, body)
				}
			}
			ah.users = nil
		})
	}
}
func TestGetUsers(t *testing.T) {

	tt := []struct {
		name       string
		method     string
		want       []model.User
		statusCode int
	}{
		{
			name:   "with valid users",
			method: http.MethodGet,

			want: []model.User{
				{
					ID:         "ID1",
					Name:       "ID1Name",
					SignupTime: time.Now(),
				},
				{
					ID:         "ID2",
					Name:       "ID2Name",
					SignupTime: time.Now(),
				},
			},
			statusCode: http.StatusOK,
		},
		{
			name:   "with empty users",
			method: http.MethodGet,

			want:       []model.User{},
			statusCode: http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ah := NewAppHandlers()

			if len(tc.want) > 0 {
				for _, u := range tc.want {
					ah.users[u.ID] = u
				}
			}

			rec := httptest.NewRecorder()
			req, err := http.NewRequest(
				tc.method,
				"/users",
				nil,
			)

			if err != nil {
				t.Fatalf("Could not create a request %v", err)
			}

			ah.GetUsers(rec, req)

			if rec.Code != tc.statusCode {
				t.Errorf("expected status %v, received status %v", tc.statusCode, rec.Code)
				return
			}

			// compare response only for valid result
			if rec.Code == http.StatusOK && len(tc.want) > 0 {
				want, _ := json.Marshal(&tc.want)
				resp := rec.Result()
				body, _ := ioutil.ReadAll(resp.Body)
				if bytes.Compare(body, want) != 0 {
					t.Errorf("expected response  %s, received response %s", want, body)
				}
			}
			ah.users = nil
		})
	}
}
