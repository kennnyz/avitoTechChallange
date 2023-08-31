package http_delivery

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	mock_http_delivery "www.github.com/kennnyz/avitochallenge/internal/delivery/http/mocks"
	"www.github.com/kennnyz/avitochallenge/internal/models"
)

func TestSegmentValidation(t *testing.T) {
	handler := NewHandler(nil)
	testTable := []struct {
		name             string
		testId           int
		inputSegmentName string
		expectedError    error
	}{
		{
			name:             "success",
			testId:           1,
			inputSegmentName: "segment1",
			expectedError:    nil,
		},
		{
			name:             "empty segment name",
			testId:           2,
			inputSegmentName: "",
			expectedError:    models.SegmentNameEmptyErr,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			err := handler.validateSegment(&models.Segment{
				Name: testCase.inputSegmentName,
			})
			require.Equal(t, testCase.expectedError, err)
		})
	}
}

func TestCreateSegment(t *testing.T) {
	type mockBehavior func(s *mock_http_delivery.MockUserSegmentService)
	testTable := []struct {
		name               string
		testId             int
		mockBehavior       mockBehavior
		expectedStatusCode int
		expectedResponse   string
		expectedBody       string
		expectedMethod     string
	}{
		{
			name:   "success",
			testId: 1,
			mockBehavior: func(s *mock_http_delivery.MockUserSegmentService) {
				s.EXPECT().CreateSegment(context.Background(), "segment1").Return(nil)
			},
			expectedStatusCode: 200,
			expectedResponse:   `{"message":"success"}`,
			expectedBody:       `{"name":"segment1"}`,
			expectedMethod:     http.MethodPost,
		},
		{
			name:               "invalid json",
			testId:             2,
			mockBehavior:       func(s *mock_http_delivery.MockUserSegmentService) {},
			expectedBody:       `{"name":"segment1"`,
			expectedMethod:     http.MethodPost,
			expectedStatusCode: 400,
			expectedResponse:   `{"message":"invalid json"}`,
		},
		{
			name:               "empty segment name",
			testId:             3,
			mockBehavior:       func(s *mock_http_delivery.MockUserSegmentService) {},
			expectedStatusCode: 400,
			expectedResponse:   `{"message":"segment name is empty"}`,
			expectedBody:       `{"name":""}`,
			expectedMethod:     http.MethodPost,
		},
		{
			name:               "method not provided",
			testId:             4,
			mockBehavior:       func(s *mock_http_delivery.MockUserSegmentService) {},
			expectedStatusCode: 405,
			expectedResponse:   `{"message":"method not provided"}`,
			expectedBody:       `{"name":"segment1"}`,
			expectedMethod:     http.MethodGet,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			urlShorty := mock_http_delivery.NewMockUserSegmentService(ctrl)
			testCase.mockBehavior(urlShorty)

			handler := NewHandler(urlShorty)
			mux := http.NewServeMux()
			mux.HandleFunc("/create-segment", handler.createSegment)

			requestBody := bytes.NewBufferString(testCase.expectedBody)
			req := httptest.NewRequest(testCase.expectedMethod, "/create-segment", requestBody)
			recorder := httptest.NewRecorder()

			mux.ServeHTTP(recorder, req)
			assert.Equal(t, testCase.expectedStatusCode, recorder.Code)
			assert.Equal(t, testCase.expectedResponse+"\n", recorder.Body.String())
		})
	}
}

func TestDeleteSegment(t *testing.T) {
	type mockBehavior func(s *mock_http_delivery.MockUserSegmentService)
	testTable := []struct {
		name               string
		testId             int
		userId             int
		mockBehavior       mockBehavior
		expectedMethod     string
		expectedStatusCode int
		expectedResponse   string
		expectedBody       string
	}{
		{
			name:         "segment not found",
			testId:       1,
			expectedBody: `{"name":"segment2"}`,
			mockBehavior: func(s *mock_http_delivery.MockUserSegmentService) {
				s.EXPECT().DeleteSegment(context.Background(), "segment2").Return(models.SegmentNotFoundErr("segment2"))
			},
			expectedStatusCode: 500,
			expectedResponse:   `{"message":"segment segment2 not found"}`,
			expectedMethod:     http.MethodDelete,
		},
		{
			name:         "success",
			testId:       2,
			expectedBody: `{"name":"segment1"}`,
			mockBehavior: func(s *mock_http_delivery.MockUserSegmentService) {
				s.EXPECT().DeleteSegment(context.Background(), "segment1").Return(nil)
			},
			expectedStatusCode: 200,
			expectedResponse:   `{"message":"success"}`,
			expectedMethod:     http.MethodDelete,
		},
		{
			name:               "invalid json",
			testId:             3,
			expectedBody:       `{"name":"segment1"`,
			mockBehavior:       func(s *mock_http_delivery.MockUserSegmentService) {},
			expectedMethod:     http.MethodDelete,
			expectedStatusCode: 400,
			expectedResponse:   `{"message":"invalid json"}`,
		},
		{
			name:               "empty segment name",
			testId:             4,
			expectedBody:       `{"name":""}`,
			mockBehavior:       func(s *mock_http_delivery.MockUserSegmentService) {},
			expectedMethod:     http.MethodDelete,
			expectedStatusCode: 400,
			expectedResponse:   `{"message":"segment name is empty"}`,
		},
		{
			name:               "method not provided",
			testId:             5,
			expectedBody:       `{"name":"segment1"}`,
			mockBehavior:       func(s *mock_http_delivery.MockUserSegmentService) {},
			expectedMethod:     http.MethodGet,
			expectedStatusCode: 405,
			expectedResponse:   `{"message":"method not provided"}`,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			urlShorty := mock_http_delivery.NewMockUserSegmentService(ctrl)
			testCase.mockBehavior(urlShorty)

			handler := NewHandler(urlShorty)
			mux := http.NewServeMux()
			mux.HandleFunc("/delete-segment", handler.deleteSegment)

			requestBody := bytes.NewBufferString(testCase.expectedBody)
			req := httptest.NewRequest(testCase.expectedMethod, "/delete-segment", requestBody)
			recorder := httptest.NewRecorder()

			mux.ServeHTTP(recorder, req)
			assert.Equal(t, testCase.expectedStatusCode, recorder.Code)
			assert.Equal(t, testCase.expectedResponse+"\n", recorder.Body.String())
		})
	}
}
