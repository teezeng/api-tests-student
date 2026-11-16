package advertisements

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/managers/advertisements/models"
)

const photosFieldName = "photos"

func buildCreateAdvertisementForm(t *testing.T, req models.CreateAdvertisementRequest) (body, contentType string) {
	t.Helper()

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	require.NoError(t, w.WriteField("title", req.Title))
	require.NoError(t, w.WriteField("description", req.Description))
	if req.Price != 0 {
		require.NoError(t, w.WriteField("price", fmt.Sprintf("%g", req.Price)))
	}
	if req.Quantity != 0 {
		require.NoError(t, w.WriteField("quantity", strconv.Itoa(req.Quantity)))
	}
	for i, photo := range req.Photos {
		fw, err := w.CreateFormFile(photosFieldName, fmt.Sprintf("photo_%d.jpg", i))
		require.NoError(t, err)
		_, err = fw.Write(photo)
		require.NoError(t, err)
	}
	require.NoError(t, w.Close())

	return buf.String(), w.FormDataContentType()
}

func HttpPostAdvertisement(t *testing.T, token string, req models.CreateAdvertisementRequest) *http.Response {
	body, contentType := buildCreateAdvertisementForm(t, req)

	return apiRunner.GetRunner().
		Auth(token).
		Create().
		Post(path.AdvertisementPath).
		ContentType(contentType).
		Body(body).
		Expect(t).
		End().Response
}

func HttpGetAdvertisementByID(t *testing.T, id string) *http.Response {
	return apiRunner.GetRunner().
		Create().
		Get(path.AdvertisementPath).
		Query("id", id).
		Expect(t).
		End().Response
}

func HttpGetAdvertisements(t *testing.T, token, textForSearch string) *http.Response {
	return apiRunner.GetRunner().
		Auth(token).
		Create().
		Get(path.AdvertisementsPath).
		Query("search", textForSearch).
		Expect(t).
		End().Response
}

func HttpGetAdvertisementPhotos(t *testing.T, token, id string) *http.Response {
	return apiRunner.GetRunner().
		Auth(token).
		Create().
		Get(fmt.Sprintf(path.AdvertisementPhotosPath, id)).
		Expect(t).
		End().Response
}

func HttpDeleteAdvertisement(t *testing.T, token, id string) *http.Response {
	return apiRunner.GetRunner().
		Auth(token).
		Create().
		Delete(path.AdvertisementPath).
		Query("id", id).
		Expect(t).
		End().Response
}
