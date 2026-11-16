package assertions

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"api-tests-template/internal/constants/fixtures"
	"api-tests-template/internal/helpers/response"
	advModels "api-tests-template/internal/managers/advertisements/models"
	favModels "api-tests-template/internal/managers/favorites/models"
)

func ContainsAdvertisementByTitle(t *testing.T, items []advModels.AdvertisementResponse, title string) {
	t.Helper()
	for _, item := range items {
		if item.Title == title {
			return
		}
	}
	t.Fatalf("объявление с названием %q не найдено среди %d элементов", title, len(items))
}

func ContainsAdvertisementByID(t *testing.T, items []advModels.AdvertisementResponse, id string) {
	t.Helper()
	for _, item := range items {
		if item.ID == id {
			return
		}
	}
	t.Fatalf("объявление с id=%s не найдено среди %d элементов", id, len(items))
}

func FavoritesContainAdvertisement(t *testing.T, items []favModels.FavoriteItem, adID string) {
	t.Helper()
	require.Truef(t, favoritesHaveAdvertisement(items, adID),
		"объявление id=%s не найдено в избранном (элементов: %d)", adID, len(items))
}

func FavoritesNotContainAdvertisement(t *testing.T, items []favModels.FavoriteItem, adID string) {
	t.Helper()
	require.Falsef(t, favoritesHaveAdvertisement(items, adID),
		"объявление id=%s все еще в избранном (элементов: %d)", adID, len(items))
}

func ErrorResponse(t *testing.T, raw response.Raw, wantStatus int) {
	t.Helper()
	raw.RequireStatus(t, wantStatus)

	require.Truef(t, gjson.Valid(raw.Body), "тело ошибки не является валидным JSON: %s", raw.Body)
	require.NotEmptyf(t, gjson.Get(raw.Body, "error").String(),
		"в теле ошибки не заполнено поле error: %s", raw.Body)

	description := gjson.Get(raw.Body, "message").String()
	if description == "" {
		description = gjson.Get(raw.Body, "details").String()
	}
	require.NotEmptyf(t, description,
		"в теле ошибки нет описания (message/details): %s", raw.Body)
}

func UnauthorizedResponse(t *testing.T, raw response.Raw) {
	t.Helper()
	ErrorResponse(t, raw, http.StatusUnauthorized)
	require.Equalf(t, fixtures.ErrorCodeUnauthorized, gjson.Get(raw.Body, "error").String(),
		"некорректное значение поля error. Тело: %s", raw.Body)
}

func favoritesHaveAdvertisement(items []favModels.FavoriteItem, adID string) bool {
	for _, item := range items {
		if item.Advertisement.ID == adID {
			return true
		}
	}
	return false
}
