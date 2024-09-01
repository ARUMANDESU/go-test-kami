package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/ARUMANDESU/go-test-kami/internal/api"
	"github.com/ARUMANDESU/go-test-kami/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type HTTPReservationTestSuite struct {
	httpServer       *httptest.Server
	reservationSuite *ReservationSuite
}

func NewHTTPReservationTestSuite(t *testing.T) *HTTPReservationTestSuite {
	t.Helper()

	reservationSuite := NewReservationSuite(t)

	httpAPI := api.NewAPI(logger.Plug(), reservationSuite.service)

	httpServer := httptest.NewServer(httpAPI.ChiRouter())

	t.Cleanup(func() {
		httpServer.Close()
	})

	return &HTTPReservationTestSuite{
		httpServer:       httpServer,
		reservationSuite: reservationSuite,
	}
}

func TestHTTPReservation_GetRoomReservations(t *testing.T) {
	suite := NewHTTPReservationTestSuite(t)

	tests := []struct {
		name     string
		roomID   string
		expected int
	}{
		{
			name:     "room has 3 reservations",
			roomID:   "018c0f7e-7b9b-7f4b-8e4b-4b5e4b5e4b5f",
			expected: 3,
		},
		{
			name:     "room has 1 reservation",
			roomID:   "018c0f7e-7b9b-7f4b-8e4b-4b5e4b5e4b6d",
			expected: 1,
		},
		{
			name:     "room has no reservations",
			roomID:   "018c0f7e-7b9b-7f4b-8e4b-4b5e4b5e4b5e",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", suite.httpServer.URL+"/v1/reservations/"+tt.roomID, nil)
			res := httptest.NewRecorder()

			suite.httpServer.Config.Handler.ServeHTTP(res, req)

			require.Equal(t, http.StatusOK, res.Code)
			assert.Equal(t, "application/json", res.Header().Get("Content-Type"))

			var body map[string]interface{}
			err := json.NewDecoder(res.Body).Decode(&body)
			require.NoError(t, err)
			if tt.expected == 0 {
				assert.Nil(t, body["reservations"])
				return
			}
			assert.Equal(t, tt.expected, len(body["reservations"].([]interface{})))
		})
	}
}

func TestHTTPReservation_ReserveRoom(t *testing.T) {
	suite := NewHTTPReservationTestSuite(t)

	tests := []struct {
		name     string
		roomID   string
		start    string
		end      string
		expected int
	}{
		{
			name:     "room is available",
			roomID:   "018c0f7e-7b9b-7f4b-8e4b-4b5e4b5e4b5f",
			start:    "2021-09-01T10:00:00Z",
			end:      "2021-09-01T11:00:00Z",
			expected: http.StatusCreated,
		},
		{
			name:     "room is not available",
			roomID:   "018c0f7e-7b9b-7f4b-8e4b-4b5e4b5e4b5f",
			start:    "2021-09-01T10:00:00Z",
			end:      "2021-09-01T11:00:00Z",
			expected: http.StatusConflict,
		},
		{
			name:     "room is not available",
			roomID:   "018c0f7e-7b9b-7f4b-8e4b-4b5e4b5e4b5f",
			start:    "2021-09-01T10:00:00Z",
			end:      "2021-09-01T11:00:00Z",
			expected: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", suite.httpServer.URL+"/v1/reservations", nil)
			req.Header.Set("Content-Type", "application/json")
			body := map[string]interface{}{
				"room_id":    tt.roomID,
				"start_time": tt.start,
				"end_time":   tt.end,
			}

			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(body)
			require.NoError(t, err)
			req.Body = io.NopCloser(&buf)
			require.NoError(t, err)

			res := httptest.NewRecorder()

			suite.httpServer.Config.Handler.ServeHTTP(res, req)

			require.Equal(t, tt.expected, res.Code)
			assert.Equal(t, "application/json", res.Header().Get("Content-Type"))
		})
	}
}

func TestHTTPReservation_ReserveRoom_Concurrently(t *testing.T) {
	suite := NewHTTPReservationTestSuite(t)

	var errCount int
	var mu sync.Mutex

	const concurrency = 10

	wg := sync.WaitGroup{}
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()

			req := httptest.NewRequest("POST", suite.httpServer.URL+"/v1/reservations", nil)
			req.Header.Set("Content-Type", "application/json")
			body := map[string]interface{}{
				"room_id":    "018c0f7e-7b9b-7f4b-8e4b-4b5e4b5e4b5f",
				"start_time": "2021-09-01T10:00:00Z",
				"end_time":   "2021-09-01T11:00:00Z",
			}

			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(body)
			require.NoError(t, err)
			req.Body = io.NopCloser(&buf)
			require.NoError(t, err)

			res := httptest.NewRecorder()

			suite.httpServer.Config.Handler.ServeHTTP(res, req)

			mu.Lock()
			if res.Code != http.StatusCreated {
				errCount++
			}
			mu.Unlock()
		}()
	}

	wg.Wait()

	assert.Equal(t, concurrency-1, errCount) // concurrency-1 because one request should succeed

	req := httptest.NewRequest("GET", suite.httpServer.URL+"/v1/reservations/018c0f7e-7b9b-7f4b-8e4b-4b5e4b5e4b5f", nil)
	res := httptest.NewRecorder()

	suite.httpServer.Config.Handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, "application/json", res.Header().Get("Content-Type"))

	var body map[string]interface{}
	err := json.NewDecoder(res.Body).Decode(&body)
	require.NoError(t, err)
	assert.Equal(t, 4, len(body["reservations"].([]interface{})))

}

func TestHTTPReservation_GetRoomReservations_InvalidRoomID(t *testing.T) {
	t.Run("invalid uuid", func(t *testing.T) {
		suite := NewHTTPReservationTestSuite(t)

		req := httptest.NewRequest("GET", suite.httpServer.URL+"/v1/reservations/invalid-room-id", nil)
		res := httptest.NewRecorder()

		suite.httpServer.Config.Handler.ServeHTTP(res, req)

		require.Equal(t, http.StatusBadRequest, res.Code)
		assert.Equal(t, "application/json", res.Header().Get("Content-Type"))

		var body map[string]interface{}
		err := json.NewDecoder(res.Body).Decode(&body)
		require.NoError(t, err)
		assert.True(t, strings.Contains(body["error"].(string), "must be a valid UUID"))
	})

}

func TestHTTPReservation_ReserveRoom_InvalidTimeFormat(t *testing.T) {
	suite := NewHTTPReservationTestSuite(t)

	req := httptest.NewRequest("POST", suite.httpServer.URL+"/v1/reservations", nil)
	req.Header.Set("Content-Type", "application/json")
	body := map[string]interface{}{
		"room_id":    "018c0f7e-7b9b-7f4b-8e4b-4b5e4b5e4b5f",
		"start_time": "invalid-time-format",
		"end_time":   "2021-09-01T11:00:00Z",
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	require.NoError(t, err)
	req.Body = io.NopCloser(&buf)
	require.NoError(t, err)

	res := httptest.NewRecorder()

	suite.httpServer.Config.Handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusBadRequest, res.Code)
	assert.Equal(t, "application/json", res.Header().Get("Content-Type"))

	var resBody map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&resBody)
	require.NoError(t, err)
	assert.True(t, strings.Contains(resBody["error"].(string), "invalid argument"))
}

func TestHTTPReservation_ReserveRoom_EndBeforeStart(t *testing.T) {
	suite := NewHTTPReservationTestSuite(t)

	req := httptest.NewRequest("POST", suite.httpServer.URL+"/v1/reservations", nil)
	req.Header.Set("Content-Type", "application/json")
	body := map[string]interface{}{
		"room_id":    "018c0f7e-7b9b-7f4b-8e4b-4b5e4b5e4b5f",
		"start_time": "2021-09-01T11:00:00Z",
		"end_time":   "2021-09-01T10:00:00Z",
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	require.NoError(t, err)
	req.Body = io.NopCloser(&buf)
	require.NoError(t, err)

	res := httptest.NewRecorder()

	suite.httpServer.Config.Handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusBadRequest, res.Code)
	assert.Equal(t, "application/json", res.Header().Get("Content-Type"))

	var resBody map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&resBody)
	require.NoError(t, err)
	assert.True(t, strings.Contains(resBody["error"].(string), "end time must be after start time"))
}

func TestHTTPReservation_ReserveRoom_MissingFields(t *testing.T) {
	suite := NewHTTPReservationTestSuite(t)

	req := httptest.NewRequest("POST", suite.httpServer.URL+"/v1/reservations", nil)
	req.Header.Set("Content-Type", "application/json")
	body := map[string]interface{}{
		"room_id": "018c0f7e-7b9b-7f4b-8e4b-4b5e4b5e4b5f",
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	require.NoError(t, err)
	req.Body = io.NopCloser(&buf)
	require.NoError(t, err)

	res := httptest.NewRecorder()

	suite.httpServer.Config.Handler.ServeHTTP(res, req)

	require.Equal(t, http.StatusBadRequest, res.Code)
	assert.Equal(t, "application/json", res.Header().Get("Content-Type"))

	var resBody map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&resBody)
	require.NoError(t, err)
	assert.True(t, strings.Contains(resBody["error"].(string), "invalid argument"))
}
