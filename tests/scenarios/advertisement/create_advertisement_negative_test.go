package advertisement

import (
	"fmt"
	"net/http"
	"strings"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/constants/fixtures"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/advertisements/models"

	base "api-tests-template/tests"
)

func noMutation(*models.CreateAdvertisementRequest) {}

func tooLongTitle(r *models.CreateAdvertisementRequest) {
	r.Title = strings.Repeat("a", fixtures.MaxTitleLength+1)
}

func (s *TestSuite) TestCreateAdvertisementNegative() {
	cases := []struct {
		name       string
		token      string
		mutate     func(*models.CreateAdvertisementRequest)
		wantStatus int
	}{
		{
			name:       "пустой title",
			token:      s.login.Token,
			mutate:     func(r *models.CreateAdvertisementRequest) { r.Title = "" },
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "без фотографий",
			token:      s.login.Token,
			mutate:     func(r *models.CreateAdvertisementRequest) { r.Photos = nil },
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "отрицательная цена",
			token:      s.login.Token,
			mutate:     func(r *models.CreateAdvertisementRequest) { r.Price = fixtures.NegativePrice },
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "слишком длинный title",
			token:      s.login.Token,
			mutate:     tooLongTitle,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "без авторизации",
			token:      "",
			mutate:     noMutation,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "неверный токен",
			token:      fixtures.InvalidToken,
			mutate:     noMutation,
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			req := advertisements.NewValidAdvertisement(1)
			tc.mutate(&req)

			base.Step(fmt.Sprintf("POST /advertisement, кейс %q, ожидаем %d", tc.name, tc.wantStatus))
			raw := advertisements.CreateAdvertisementRaw(s.T(), tc.token, req)
			advertisements.DeleteIfCreated(s.T(), s.login.Token, raw)

			base.Check("Статус-код и сообщение об ошибке соответствуют ожидаемым")
			if tc.wantStatus == http.StatusUnauthorized {
				assertions.UnauthorizedResponse(s.T(), raw)
				return
			}
			assertions.ErrorResponse(s.T(), raw, tc.wantStatus)
		})
	}
}
