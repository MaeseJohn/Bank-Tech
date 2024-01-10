package tests

import (
	"API_Rest/models"
	"net/http"
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestCreateUser(t *testing.T) {
	CreateEnviroment()
	var createUsersTests = []struct {
		name   string
		user   models.User
		status int
	}{
		//Invaild uuid
		{"Invalid uuid", *models.NewUser("hola", "jonathan@gmail.com", "jonathan", "jonathan", "issuer", 1000), http.StatusUnprocessableEntity},
		//Invalid email
		{"Invalid email", *models.NewUser("4ddbb37b-efd5-4564-a2ba-c4ac80925b9f", "jonathangmail.com", "jonathan", "jonathan", "issuer", 1000), http.StatusUnprocessableEntity},
		//Invaild name
		{"Invalid user name", *models.NewUser("4ddbb37b-efd5-4564-a2ba-c4ac80925b9f", "jonathan@gmail.com", "jonathan", "1234", "issuer", 1000), http.StatusUnprocessableEntity},
		//Invalid role
		{"Invalid role", *models.NewUser("4ddbb37b-efd5-4564-a2ba-c4ac80925b9f", "jonathan@gmail.com", "jonathan", "jonathan", "ramon", 1000), http.StatusUnprocessableEntity},
		//201 Created
		{"201 Created", validUsers[0], http.StatusCreated},
		{"201 Created", validUsers[1], http.StatusCreated},
		{"201 Created", validUsers[2], http.StatusCreated},
		//Try to Create user with the same uuid
		{"Diplicate key value uuid", *models.NewUser("ae215592-5c65-11ee-8c99-0242ac120002", "loris@gmail.com", "loris", "loris", "issuer", 600), http.StatusInternalServerError},
		//Try to Create useer with the same email
		{"Duplicate key value email", *models.NewUser("eb84a444-8863-4c91-8a40-c801e45955b4", "jonathan@gmail.com", "jonathan", "jonathan", "issuer", 1000), http.StatusInternalServerError},
	}

	for _, u := range createUsersTests {
		t.Run(u.name, func(t *testing.T) {
			got := CreateUserRequest(&u.user)
			assert.Equal(t, got, u.status)
		})
	}
}

func TestGetUser(t *testing.T) {
	CreateEnviroment()
	CreateValidUsers(t)
	var emtyUser models.User
	var getUserData = []struct {
		name   string
		email  string
		status int
		user   models.User
	}{
		{"Valid user", validUsers[0].Email, http.StatusOK, getUsersData[0]},
		{"Valid user", validUsers[1].Email, http.StatusOK, getUsersData[1]},
		{"Valid user", validUsers[2].Email, http.StatusOK, getUsersData[2]},
		{"Unprocesable email", "jonathan", http.StatusUnprocessableEntity, emtyUser},
		{"Unprocesable email", "jose@lsakdjflksdf", http.StatusUnprocessableEntity, emtyUser},
		{"Unregister email", "mermelada@gmail.es", http.StatusNotFound, emtyUser},
		{"Unregister email", "pajaritos@gmail.com", http.StatusNotFound, emtyUser},
	}

	for _, req := range getUserData {
		t.Run(req.name, func(t *testing.T) {
			user, status := GetUserRequest(req.email)
			assert.Equal(t, req.status, status)
			assert.Equal(t, req.user, *user)
		})

	}
}

func TestGetUsers(t *testing.T) {
	CreateEnviroment()
	CreateValidUsers(t)
	got := GetUsersRequest()
	assert.Equal(t, got, getUsersData)
}

func TestLogin(t *testing.T) {
	CreateEnviroment()
	CreateValidUsers(t)
	var loginData = []struct {
		name     string
		Email    string
		Password string
		Status   int
		Token    string
	}{
		{"Unprocessable email", "12341234", "jonathan", http.StatusUnprocessableEntity, ""},
		{"Unregister email", "jonasfas@gmail.com", "jonathan", http.StatusUnauthorized, ""},
		{"Incorrect password", validUsers[0].Email, "Maese", http.StatusUnauthorized, ""},
		{"Valid parameters", validUsers[0].Email, validUsers[0].Password, http.StatusOK, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2Z1bmRzIjoxMDAwLCJ1c2VyX2lkIjoiNGRkYmIzN2ItZWZkNS00NTY0LWEyYmEtYzRhYzgwOTI1YjlmIiwidXNlcl9yb2xlIjoiaXNzdWVyIn0.H30BMZ4Eq8Ujj0_YJJUhC7IB49qxJxIfwdRy-GAsAI8"},
		{"Valid parameters", validUsers[1].Email, validUsers[1].Password, http.StatusOK, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2Z1bmRzIjo2MDAsInVzZXJfaWQiOiJhZTIxNTU5Mi01YzY1LTExZWUtOGM5OS0wMjQyYWMxMjAwMDIiLCJ1c2VyX3JvbGUiOiJpc3N1ZXIifQ.Oy882kyZjMbe8b8E5nd7HZiNLXQQYSbDWM66O3bPnjQ"},
		{"Valid parameters", validUsers[2].Email, validUsers[2].Password, http.StatusOK, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2Z1bmRzIjoxMDAwMCwidXNlcl9pZCI6IjY1MjBhMTQ4LTg4ZjUtNDVhNC1hMTg3LTg4M2NmZWJjZjk4NSIsInVzZXJfcm9sZSI6ImludmVzdG9yIn0.XxNTtL521_lVAgwioVwZpPUPbhJpYDOYt8-YNYgl-ZM"},
	}

	for _, u := range loginData {
		t.Run(u.name, func(t *testing.T) {
			status, token := LoginUserRequest(u.Email, u.Password)
			assert.Equal(t, status, u.Status)
			assert.Equal(t, token, u.Token)
		})
	}
}
