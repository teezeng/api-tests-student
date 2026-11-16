package advertisement

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/constants/fixtures"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"

	base "api-tests-template/tests"
)

type TestSuite struct {
	suite.Suite
	login authModels.LoginOkResponse
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()
	base.Precondition("Авторизация пользователя с кредами из переменных окружения")
	s.login = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

func (s *TestSuite) TestPositiveCreateAndGet() {
	t := s.T()
	req := advertisements.NewValidAdvertisement(fixtures.MaxPhotosCount)

	base.Precondition(fmt.Sprintf("Создание объявления %q: цена=%g, количество=%d, фото=%d",
		req.Title, req.Price, req.Quantity, len(req.Photos)))
	created := advertisements.CreateAdvertisement(t, s.login.Token, req)

	t.Cleanup(func() {
		base.Postcondition(fmt.Sprintf("Удаление объявления id=%s", created.ID))
		advertisements.DeleteAdvertisementRaw(t, s.login.Token, created.ID)
	})

	base.Check("Ответ на создание содержит те же параметры, что были переданы")
	s.Require().NotEmpty(created.ID, "сервер вернул пустой id созданного объявления")
	s.Require().Equal(req.Title, created.Title, "title не совпал в ответе создания")
	s.Require().Equal(req.Description, created.Description, "description не совпал в ответе создания")
	s.Require().Equal(req.Price, created.Price, "price не совпал в ответе создания")
	s.Require().Equal(req.Quantity, created.Quantity, "quantity не совпал в ответе создания")
	s.Require().Len(created.Photos, len(req.Photos), "количество фото в ответе создания не совпадает")

	s.Run("GET /advertisement?id= возвращает переданные параметры", func() {
		base.Step(fmt.Sprintf("Запрос объявления по id=%s", created.ID))
		got := advertisements.GetAdvertisementByID(s.T(), created.ID)

		base.Check("Все параметры совпадают с переданными при создании")
		s.Require().Equal(created.ID, got.ID, "id не совпадает")
		s.Require().Equal(req.Title, got.Title, "title не совпадает")
		s.Require().Equal(req.Description, got.Description, "description не совпадает")
		s.Require().Equal(req.Price, got.Price, "price не совпадает")
		s.Require().Equal(req.Quantity, got.Quantity, "quantity не совпадает")
		s.Require().Equal(s.login.User.Id, got.UserID, "объявление принадлежит другому пользователю")
		s.Require().Len(got.Photos, len(req.Photos), "количество фото в объявлении не совпадает")
	})

	s.Run("GET /advertisements?search={title} находит созданное объявление", func() {
		base.Step(fmt.Sprintf("Поиск по полному названию %q", req.Title))
		list := advertisements.GetAdvertisementsList(s.T(), s.login.Token, req.Title)

		base.Check("Созданное объявление есть в результатах поиска")
		s.Require().NotEmpty(list.Items, "поиск вернул пустой список")
		assertions.ContainsAdvertisementByID(s.T(), list.Items, created.ID)
	})

	s.Run("GET /advertisements/{id}/photos возвращает загруженные фото", func() {
		base.Step(fmt.Sprintf("Запрос фотографий объявления id=%s", created.ID))
		photos := advertisements.GetAdvertisementPhotos(s.T(), s.login.Token, created.ID)

		base.Check("Количество, идентификаторы и URL фотографий корректны")
		s.Require().Len(photos, len(req.Photos), "количество фото в ручке фотографий не совпадает")
		for _, photo := range photos {
			s.Require().NotEmpty(photo.ID, "у фотографии пустой id")
			s.Require().NotEmpty(photo.URL, "у фотографии id=%s пустой URL", photo.ID)
		}
		s.Require().ElementsMatch(models.PhotoIDs(created.Photos), models.PhotoIDs(photos),
			"идентификаторы фото из создания и из ручки фотографий не совпадают")
	})
}
