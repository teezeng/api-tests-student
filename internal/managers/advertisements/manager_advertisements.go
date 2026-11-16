package advertisements

import (
	"net/http"
	"testing"

	adClient "api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/constants/fixtures"
	"api-tests-template/internal/helpers/response"
	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/utils"
)

func NewValidAdvertisement(photosCount int) models.CreateAdvertisementRequest {
	photos := make([][]byte, 0, photosCount)
	for range photosCount {
		photos = append(photos, fixtures.JPEGPhoto())
	}

	return models.CreateAdvertisementRequest{
		Title:       fixtures.TitlePrefix + utils.RandomString(fixtures.TitleRandomLength),
		Description: fixtures.Description,
		Price:       fixtures.Price,
		Quantity:    fixtures.Quantity,
		Photos:      photos,
	}
}

func CreateAdvertisement(t *testing.T, token string, req models.CreateAdvertisementRequest) models.AdvertisementResponse {
	t.Helper()

	raw := CreateAdvertisementRaw(t, token, req)
	raw.RequireStatus(t, http.StatusCreated)

	return response.Decode[models.AdvertisementResponse](t, raw, "создания объявления")
}

func CreateAdvertisementRaw(t *testing.T, token string, req models.CreateAdvertisementRequest) response.Raw {
	t.Helper()
	return response.Read(t, adClient.HttpPostAdvertisement(t, token, req))
}

func GetAdvertisementByID(t *testing.T, id string) models.AdvertisementResponse {
	t.Helper()

	raw := response.Read(t, adClient.HttpGetAdvertisementByID(t, id))
	raw.RequireStatus(t, http.StatusOK)

	return response.Decode[models.AdvertisementResponse](t, raw, "получения объявления")
}

func GetAdvertisementsList(t *testing.T, token, text string) models.AdvertisementsList {
	t.Helper()

	raw := response.Read(t, adClient.HttpGetAdvertisements(t, token, text))
	raw.RequireStatus(t, http.StatusOK)

	return response.Decode[models.AdvertisementsList](t, raw, "поиска объявлений")
}

func GetAdvertisementPhotos(t *testing.T, token, id string) []models.Photo {
	t.Helper()

	raw := response.Read(t, adClient.HttpGetAdvertisementPhotos(t, token, id))
	raw.RequireStatus(t, http.StatusOK)

	return response.Decode[[]models.Photo](t, raw, "получения фотографий объявления")
}

func DeleteAdvertisement(t *testing.T, token, id string) {
	t.Helper()

	raw := DeleteAdvertisementRaw(t, token, id)
	raw.RequireStatus(t, http.StatusNoContent)
}

func DeleteAdvertisementRaw(t *testing.T, token, id string) response.Raw {
	t.Helper()
	return response.Read(t, adClient.HttpDeleteAdvertisement(t, token, id))
}

func DeleteIfCreated(t *testing.T, ownerToken string, raw response.Raw) {
	t.Helper()

	if raw.StatusCode != http.StatusCreated {
		return
	}

	created := response.Decode[models.AdvertisementResponse](t, raw, "создания объявления")
	DeleteAdvertisementRaw(t, ownerToken, created.ID)
}
