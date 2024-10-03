package tests

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/fnacarellidev/microsservices/auth/handlers"
	"github.com/fnacarellidev/microsservices/tests/testutil"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthTestSuite struct {
	suite.Suite
	pgContainer *testutil.PostgresContainer
	router      *httprouter.Router
	jwt         string
}

func (suite *AuthTestSuite) SetupSuite() {
	var err error

	suite.router = httprouter.New()
	suite.router.POST("/auth/register", handlers.RegisterHandler)
	suite.router.POST("/auth/login", handlers.LoginHandler)
	suite.pgContainer, err = testutil.CreatePostgresContainer()
	if err != nil {
		log.Fatal("Failed:", err)
	}

	os.Setenv("DB_URL", suite.pgContainer.ConnectionString)
}

func (suite *AuthTestSuite) Test000Register() {
	t := suite.T()
	rr := httptest.NewRecorder()
	data := []byte(`{
		"username": "fabin",
		"password": "fabin123"
	}`)
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(data))
	suite.router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func (suite *AuthTestSuite) Test001RegisterUserThatAlreadyExists() {
	t := suite.T()
	rr := httptest.NewRecorder()
	data := []byte(`{
		"username": "fabin",
		"password": "fabin123"
	}`)
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(data))
	suite.router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func (suite *AuthTestSuite) Test002Login() {
	t := suite.T()
	rr := httptest.NewRecorder()
	data := []byte(`{
		"username": "fabin",
		"password": "fabin123"
	}`)
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(data))
	suite.router.ServeHTTP(rr, req)
	for _, cookie := range rr.Result().Cookies() {
		if cookie.Name == "jwt" {
			suite.jwt = cookie.Value
		}
	}

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.NotEmpty(t, suite.jwt)
}

func TestAuthSuite(t *testing.T) {
	os.Setenv("GO_TESTING", "")
	suite.Run(t, new(AuthTestSuite))
}
