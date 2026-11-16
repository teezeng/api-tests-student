package favorites

import (
	"fmt"
	"net/http"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/managers/favorites/models"
)

func HttpGetFavorites(t *testing.T, token string) *http.Response {
	return apiRunner.GetRunner().
		Auth(token).
		Create().
		Get(path.FavoritesPath).
		Expect(t).
		End().Response
}

func HttpPostFavorite(t *testing.T, token, adID string) *http.Response {
	return apiRunner.GetRunner().
		Auth(token).
		Create().
		Post(path.FavoritesPath).
		ContentType("application/json").
		JSON(models.AddFavoriteRequest{AdvertisementID: adID}).
		Expect(t).
		End().Response
}

func HttpDeleteFavorite(t *testing.T, token, adID string) *http.Response {
	return apiRunner.GetRunner().
		Auth(token).
		Create().
		Delete(fmt.Sprintf(path.FavoriteByIDPath, adID)).
		Expect(t).
		End().Response
}
