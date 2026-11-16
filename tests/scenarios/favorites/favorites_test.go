package favorites

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	advModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/managers/favorites"

	base "api-tests-template/tests"
)

type TestSuite struct {
	suite.Suite
	buyer  authModels.LoginOkResponse
	seller authModels.LoginOkResponse
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	base.Precondition("Регистрация покупателя, который добавляет объявления в избранное")
	s.buyer = auth.RegisterRandomUser(s.T())

	base.Precondition("Регистрация продавца, которому принадлежат объявления")
	s.seller = auth.RegisterRandomUser(s.T())
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

func (s *TestSuite) createAdvertisementAs(ownerToken string) advModels.AdvertisementResponse {
	t := s.T()
	req := advertisements.NewValidAdvertisement(1)

	base.Precondition(fmt.Sprintf("Создание объявления %q", req.Title))
	ad := advertisements.CreateAdvertisement(t, ownerToken, req)

	t.Cleanup(func() {
		base.Postcondition(fmt.Sprintf("Удаление объявления id=%s", ad.ID))
		advertisements.DeleteAdvertisementRaw(t, ownerToken, ad.ID)
	})

	return ad
}

func (s *TestSuite) TestPositiveFavoritesLifecycle() {
	ad := s.createAdvertisementAs(s.seller.Token)

	s.Run("добавленное объявление появляется в избранном", func() {
		base.Step(fmt.Sprintf("Покупатель добавляет объявление id=%s в избранное", ad.ID))
		added := favorites.AddFavorite(s.T(), s.buyer.Token, ad.ID)

		base.Check("Ответ добавления и GET /favorites содержат это объявление")
		s.Require().Equal(ad.ID, added.Advertisement.ID, "в ответе добавления другое объявление")
		list := favorites.GetFavorites(s.T(), s.buyer.Token)
		assertions.FavoritesContainAdvertisement(s.T(), list.Items, ad.ID)
	})

	s.Run("удаленное из избранного объявление пропадает из списка", func() {
		base.Step(fmt.Sprintf("Покупатель удаляет объявление id=%s из избранного", ad.ID))
		favorites.DeleteFavorite(s.T(), s.buyer.Token, ad.ID)

		base.Check("GET /favorites больше не содержит объявление")
		list := favorites.GetFavorites(s.T(), s.buyer.Token)
		assertions.FavoritesNotContainAdvertisement(s.T(), list.Items, ad.ID)
	})

	s.Run("удаление самого объявления убирает его из избранного", func() {
		base.Step("Покупатель снова добавляет объявление в избранное")
		favorites.AddFavorite(s.T(), s.buyer.Token, ad.ID)
		listBefore := favorites.GetFavorites(s.T(), s.buyer.Token)
		assertions.FavoritesContainAdvertisement(s.T(), listBefore.Items, ad.ID)

		base.Step(fmt.Sprintf("Продавец удаляет объявление id=%s через DELETE /advertisement", ad.ID))
		advertisements.DeleteAdvertisement(s.T(), s.seller.Token, ad.ID)

		base.Check("GET /favorites покупателя больше не содержит удаленное объявление")
		listAfter := favorites.GetFavorites(s.T(), s.buyer.Token)
		assertions.FavoritesNotContainAdvertisement(s.T(), listAfter.Items, ad.ID)
	})
}
