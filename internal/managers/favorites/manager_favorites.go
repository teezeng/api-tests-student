package favorites

import (
	"net/http"
	"testing"

	favClient "api-tests-template/internal/client/http/favorites"
	"api-tests-template/internal/helpers/response"
	"api-tests-template/internal/managers/favorites/models"
)

func AddFavorite(t *testing.T, token, adID string) models.FavoriteResponse {
	t.Helper()

	raw := AddFavoriteRaw(t, token, adID)
	raw.RequireStatus(t, http.StatusCreated)

	return response.Decode[models.FavoriteResponse](t, raw, "добавления в избранное")
}

func AddFavoriteRaw(t *testing.T, token, adID string) response.Raw {
	t.Helper()
	return response.Read(t, favClient.HttpPostFavorite(t, token, adID))
}

func GetFavorites(t *testing.T, token string) models.FavoriteListResponse {
	t.Helper()

	raw := GetFavoritesRaw(t, token)
	raw.RequireStatus(t, http.StatusOK)

	return response.Decode[models.FavoriteListResponse](t, raw, "списка избранного")
}

func GetFavoritesRaw(t *testing.T, token string) response.Raw {
	t.Helper()
	return response.Read(t, favClient.HttpGetFavorites(t, token))
}

func DeleteFavorite(t *testing.T, token, adID string) {
	t.Helper()

	raw := DeleteFavoriteRaw(t, token, adID)
	raw.RequireStatus(t, http.StatusNoContent)
}

func DeleteFavoriteRaw(t *testing.T, token, adID string) response.Raw {
	t.Helper()
	return response.Read(t, favClient.HttpDeleteFavorite(t, token, adID))
}
