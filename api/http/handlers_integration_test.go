package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"ws/model"
)

func TestSetUserIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	tt := []struct {
		name       string
		method     string
		login      string
		password   string
		input      model.User
		want       string
		statusCode int
	}{
		{
			name:     "with a valid user and valid credentials",
			method:   http.MethodPost,
			login:    USERNAME,
			password: PASSWORD,

			input: model.User{
				ID:         "ID1",
				Name:       "Name1",
				SignupTime: time.Now(),
			},

			want:       "",
			statusCode: http.StatusCreated,
		},
		{
			name:     "without Id and invalid credentials",
			method:   http.MethodPost,
			login:    "test",
			password: "test",
			input: model.User{
				ID:         "",
				Name:       "Name1",
				SignupTime: time.Now(),
			},

			want:       "",
			statusCode: http.StatusUnauthorized,
		},
		{
			name:     "without Id and invalid credentials and invalid method",
			method:   http.MethodDelete,
			login:    "test",
			password: "test",
			input: model.User{
				ID:         "",
				Name:       "Name1",
				SignupTime: time.Now(),
			},

			want:       "",
			statusCode: http.StatusMethodNotAllowed,
		},
	}

	ah := NewAppHandlers()
	ts := httptest.NewServer(ah.SetupRoutes())
	defer ts.Close()
	client := ts.Client()

	for _, tc := range tt {

		t.Run(tc.name, func(t *testing.T) {
			userJson, _ := json.Marshal(&tc.input)
			req, err := http.NewRequest(tc.method, fmt.Sprintf("%s/users", ts.URL), bytes.NewBuffer(userJson))
			if err != nil {
				t.Fatal(err)
			}

			req.SetBasicAuth(tc.login, tc.password)

			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			if res.StatusCode != tc.statusCode {
				t.Errorf("expected status %v, received status %v", tc.statusCode, res.StatusCode)
			}
		})
	}
}

func TestGetUserIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	tt := []struct {
		name       string
		method     string
		login      string
		password   string
		ID         string
		want       model.User
		statusCode int
	}{
		{
			name:     "with valid id",
			method:   http.MethodGet,
			login:    USERNAME,
			password: PASSWORD,
			ID:       "ID1",

			want: model.User{
				ID:         "ID1",
				Name:       "ID1Name",
				SignupTime: time.Now(),
			},
			statusCode: http.StatusOK,
		},
		{
			name:     "with inexisting Id",
			method:   http.MethodGet,
			login:    USERNAME,
			password: PASSWORD,
			ID:       "ID4",

			want:       model.User{},
			statusCode: http.StatusNotFound,
		},
		{
			name:     "with invalid Id",
			method:   http.MethodGet,
			login:    USERNAME,
			password: PASSWORD,
			ID:       "",

			want:       model.User{},
			statusCode: http.StatusBadRequest,
		},
	}

	ah := NewAppHandlers()
	ts := httptest.NewServer(ah.SetupRoutes())
	defer ts.Close()
	client := ts.Client()

	for _, tc := range tt {

		t.Run(tc.name, func(t *testing.T) {

			/* prepare data for test */
			if tc.ID != "" && tc.statusCode != http.StatusNotFound {
				userJson, _ := json.Marshal(&tc.want)
				req, err := http.NewRequest("POST", fmt.Sprintf("%s/users", ts.URL), bytes.NewBuffer(userJson))
				if err != nil {
					t.Fatal(err)
				}

				req.SetBasicAuth(tc.login, tc.password)

				_, err = client.Do(req)
				if err != nil {
					t.Fatal(err)
				}
			}

			req, err := http.NewRequest(tc.method, fmt.Sprintf("%s/users/%s", ts.URL, tc.ID), nil)
			if err != nil {
				t.Fatal(err)
			}

			req.SetBasicAuth(tc.login, tc.password)

			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}

			if res.StatusCode != tc.statusCode {
				t.Errorf("expected status %v, received status %v", tc.statusCode, res.StatusCode)
			}

			// compare response only for valid result
			if res.StatusCode == http.StatusOK {
				want, _ := json.Marshal(&tc.want)
				body, _ := ioutil.ReadAll(res.Body)
				if bytes.Compare(body, want) != 0 {
					t.Errorf("expected response  %s, received response %s", want, body)
				}
			}
			ah.users = nil
		})
	}
}

func TestGetUsersIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	tt := []struct {
		name       string
		login      string
		password   string
		method     string
		want       []model.User
		statusCode int
	}{
		{
			name:     "with valid users",
			login:    USERNAME,
			password: PASSWORD,
			method:   http.MethodGet,

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
			name:     "with empty users",
			method:   http.MethodGet,
			login:    USERNAME,
			password: PASSWORD,

			want:       []model.User{},
			statusCode: http.StatusOK,
		},
	}

	ah := NewAppHandlers()
	ts := httptest.NewServer(ah.SetupRoutes())
	defer ts.Close()
	client := ts.Client()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ah := NewAppHandlers()

			if len(tc.want) > 0 {
				/* prepare data for test */
				for _, u := range tc.want {
					userJson, _ := json.Marshal(&u)
					req, err := http.NewRequest("POST", fmt.Sprintf("%s/users", ts.URL), bytes.NewBuffer(userJson))
					if err != nil {
						t.Fatal(err)
					}

					req.SetBasicAuth(tc.login, tc.password)

					_, err = client.Do(req)
					if err != nil {
						t.Fatal(err)
					}
				}
			}

			req, err := http.NewRequest(tc.method, fmt.Sprintf("%s/users", ts.URL), nil)
			if err != nil {
				t.Fatal(err)
			}

			req.SetBasicAuth(tc.login, tc.password)

			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}

			if res.StatusCode != tc.statusCode {
				t.Errorf("expected status %v, received status %v", tc.statusCode, res.StatusCode)
			}

			// compare response only for valid result
			if res.StatusCode == http.StatusOK && len(tc.want) > 0 {
				want, _ := json.Marshal(&tc.want)
				body, _ := ioutil.ReadAll(res.Body)
				if bytes.Compare(body, want) != 0 {
					t.Errorf("expected response  %s, received response %s", want, body)
				}
			}
			ah.users = nil
		})
	}
}
