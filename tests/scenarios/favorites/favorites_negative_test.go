package favorites

import (
	"fmt"
	"net/http"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/constants/fixtures"
	"api-tests-template/internal/helpers/response"
	"api-tests-template/internal/managers/favorites"

	base "api-tests-template/tests"
)

func (s *TestSuite) TestNegativeUnauthorized() {
	endpoints := []struct {
		name string
		call func() response.Raw
	}{
		{
			name: "GET /favorites",
			call: func() response.Raw { return favorites.GetFavoritesRaw(s.T(), "") },
		},
		{
			name: "POST /favorites",
			call: func() response.Raw { return favorites.AddFavoriteRaw(s.T(), "", fixtures.MissingAdvertisementID) },
		},
		{
			name: "DELETE /favorites/{ad_id}",
			call: func() response.Raw { return favorites.DeleteFavoriteRaw(s.T(), "", fixtures.MissingAdvertisementID) },
		},
	}

	for _, endpoint := range endpoints {
		s.Run(endpoint.name+" без токена → 401", func() {
			base.Step(endpoint.name + " без токена")
			raw := endpoint.call()

			base.Check("Статус 401 и сообщение об ошибке")
			assertions.UnauthorizedResponse(s.T(), raw)
		})
	}
}

func (s *TestSuite) TestNegativeNotFound() {
	s.Run("POST /favorites с несуществующим ad_id → 404", func() {
		base.Step("Добавление в избранное несуществующего объявления")
		raw := favorites.AddFavoriteRaw(s.T(), s.buyer.Token, fixtures.MissingAdvertisementID)

		base.Check("Статус 404 и сообщение об ошибке")
		assertions.ErrorResponse(s.T(), raw, http.StatusNotFound)
	})

	s.Run("DELETE /favorites/{ad_id} с несуществующим ad_id → 404", func() {
		base.Step("Удаление из избранного несуществующего объявления")
		raw := favorites.DeleteFavoriteRaw(s.T(), s.buyer.Token, fixtures.MissingAdvertisementID)

		base.Check("Статус 404 и сообщение об ошибке")
		assertions.ErrorResponse(s.T(), raw, http.StatusNotFound)
	})
}

func (s *TestSuite) TestNegativeAddFavoriteBusinessRules() {
	s.Run("добавление своего объявления → 403", func() {
		own := s.createAdvertisementAs(s.buyer.Token)

		base.Step(fmt.Sprintf("Владелец добавляет свое объявление id=%s в избранное", own.ID))
		raw := favorites.AddFavoriteRaw(s.T(), s.buyer.Token, own.ID)

		base.Check("Статус 403 и сообщение об ошибке")
		assertions.ErrorResponse(s.T(), raw, http.StatusForbidden)
	})

	s.Run("повторное добавление того же объявления → 409", func() {
		ad := s.createAdvertisementAs(s.seller.Token)

		base.Step(fmt.Sprintf("Первое добавление объявления id=%s в избранное", ad.ID))
		favorites.AddFavorite(s.T(), s.buyer.Token, ad.ID)

		base.Step("Повторное добавление того же объявления")
		raw := favorites.AddFavoriteRaw(s.T(), s.buyer.Token, ad.ID)

		base.Check("Статус 409 и сообщение об ошибке")
		assertions.ErrorResponse(s.T(), raw, http.StatusConflict)
	})
}
