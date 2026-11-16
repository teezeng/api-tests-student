package response

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

type Raw struct {
	StatusCode int
	Body       string
}

func Read(t *testing.T, resp *http.Response) Raw {
	t.Helper()
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "не удалось прочитать тело ответа")

	return Raw{StatusCode: resp.StatusCode, Body: string(body)}
}

func (r Raw) RequireStatus(t *testing.T, want int) {
	t.Helper()
	require.Equalf(t, want, r.StatusCode, "неожиданный HTTP status code. Тело ответа: %s", r.Body)
}

func Decode[T any](t *testing.T, r Raw, entity string) T {
	t.Helper()

	var out T
	require.NoErrorf(t, json.Unmarshal([]byte(r.Body), &out),
		"не удалось разобрать ответ (%s). Тело: %s", entity, r.Body)

	return out
}
