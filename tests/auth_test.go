package tests

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/fnacarellidev/microsservices/auth/handlers"
	"github.com/fnacarellidev/microsservices/diary/jwtaux"
	"github.com/fnacarellidev/microsservices/tests/testutil"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type AuthTestSuite struct {
	suite.Suite
	pgContainer   *testutil.PostgresContainer
	router        *httprouter.Router
	jwtCookie     *http.Cookie
	tempHS256File *os.File
}

func (suite *AuthTestSuite) SetupSuite() {
	var err error

	os.Setenv("GO_TESTING", "")
	suite.tempHS256File, err = os.Create("hs256secret.txt")
	if err != nil {
		log.Fatal(err)
	}

	suite.jwtCookie = nil
	suite.router = httprouter.New()
	suite.tempHS256File.Write([]byte("mockedhs256"))
	suite.router.POST("/auth/login", handlers.LoginHandler)
	suite.router.POST("/auth/register", handlers.RegisterHandler)
	suite.pgContainer, err = testutil.CreatePostgresContainer()
	if err != nil {
		log.Fatal(err)
	}

	os.Setenv("DB_URL", suite.pgContainer.ConnectionString)
}

func (suite *AuthTestSuite) TearDownSuite() {
	if err := suite.tempHS256File.Close(); err != nil {
		log.Fatal(err)
	}

	if err := os.Remove(suite.tempHS256File.Name()); err != nil {
		log.Fatal(err)
	}
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
			suite.jwtCookie = cookie
		}
	}

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.NotNil(t, suite.jwtCookie)
}

func (suite *AuthTestSuite) Test003DecodeJwtUsername() {
	t := suite.T()
	decodedJwt, err := jwtaux.GetDecodedJwtFromCookieHeader(*suite.jwtCookie)

	require.NoError(t, err)
	assert.Equal(t, "fabin", decodedJwt["username"])
}

func (suite *AuthTestSuite) Test004LoginWithInvalidPassword() {
	t := suite.T()
	rr := httptest.NewRecorder()
	data := []byte(`{
		"username": "fabin",
		"password": "fabin12"
	}`)
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(data))
	suite.router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}
