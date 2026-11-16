package auth

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/auth"
	"api-tests-template/internal/constants/fixtures"
	"api-tests-template/internal/helpers/response"
	"api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"
)

// Login выполняет авторизацию пользователя и возвращает токен с данными пользователя, ожидает 200
func Login(t *testing.T, username, password string) models.LoginOkResponse {
	require.NotEmpty(t, username, "Отсутствует логин пользователя")
	require.NotEmpty(t, password, "Отсутствует пароль пользователя")

	request := models.LoginRequest{
		Email:    username,
		Password: password,
	}

	loginRawResponse := auth.HttpPostAuthLogin(t, request)
	defer loginRawResponse.Body.Close()

	respArray, errResponse := io.ReadAll(loginRawResponse.Body)
	require.NoError(t, errResponse)

	loginResponse := models.LoginOkResponse{}
	err := json.Unmarshal(respArray, &loginResponse)
	require.NoError(t, err)

	return loginResponse
}

func RegisterRandomUser(t *testing.T) models.LoginOkResponse {
	t.Helper()

	email := strings.ToLower(fixtures.UserEmailPrefix+utils.RandomString(fixtures.UserRandomLength)) +
		fixtures.UserEmailDomain

	request := models.RegisterRequest{
		Email:     email,
		Password:  fixtures.UserPasswordPrefix + utils.RandomString(fixtures.UserRandomLength),
		FirstName: fixtures.UserFirstName,
	}

	raw := response.Read(t, auth.HttpPostAuthRegister(t, request))
	raw.RequireStatus(t, http.StatusCreated)

	registered := response.Decode[models.LoginOkResponse](t, raw, "регистрации пользователя")
	if registered.Token == "" {
		return Login(t, request.Email, request.Password)
	}

	return registered
}
