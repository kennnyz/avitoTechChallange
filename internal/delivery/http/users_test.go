package http_delivery

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	mock_http_delivery "www.github.com/kennnyz/avitochallenge/internal/delivery/http/mocks"
	"www.github.com/kennnyz/avitochallenge/internal/models"
)

func TestAddUserToSegment(t *testing.T) {
	type mockBehavior func(s *mock_http_delivery.MockUserSegmentService)
	testTable := []struct {
		name               string
		testId             int
		mockBehavior       mockBehavior
		expectedStatusCode int
		expectedResponse   string
		expectedBody       string
		expectedMethod     string
		path               string
		handle             func(w http.ResponseWriter, r *http.Request)
	}{
		{
			testId: 1,
			name:   "success",
			mockBehavior: func(s *mock_http_delivery.MockUserSegmentService) {
				s.EXPECT().AddUserToSegments(gomock.Any(), gomock.Any()).Return(&models.AddUserToSegmentResponse{
					UserID:           1000,
					AddedSegments:    []string{"segment1", "segment2"},
					DeletedSegments:  []string{"segment3"},
					NotExistSegments: nil,
				}, nil)
			},
			expectedStatusCode: 200,
			expectedResponse:   `{"id":1000,"added_segments":["segment1","segment2"],"deleted_segments":["segment3"],"not_exist_segments":null}`,
			expectedBody:       `{"id":1000,"segments_to_add":["segment1","segment2"],"segments_to_delete":["segment3"]}`,
			expectedMethod:     http.MethodPost,
			path:               "/add-user-to-segment",
		},
		{
			testId:             2,
			name:               "wrong method",
			mockBehavior:       func(s *mock_http_delivery.MockUserSegmentService) {},
			expectedStatusCode: 405,
			path:               "/add-user-to-segment",
			expectedResponse:   `{"message":"method not provided"}`,
			expectedBody:       `{"id":1000,"segments_to_add":["segment1","segment2"],"segments_to_delete":["segment3"]}`,
			expectedMethod:     http.MethodGet,
		},
		{
			testId:             3,
			name:               "invalid json",
			mockBehavior:       func(s *mock_http_delivery.MockUserSegmentService) {},
			expectedStatusCode: 400,
			expectedResponse:   `{"message":"invalid json"}`,
			expectedBody:       `{"id":1000,"segments_to_adding":["segment1","segment2"],"segments_to_delete":["segment3"]`,
			expectedMethod:     http.MethodPost,
			path:               "/add-user-to-segment",
		},
		{
			testId: 4,
			name:   "user not found add user to segment",
			mockBehavior: func(s *mock_http_delivery.MockUserSegmentService) {
				s.EXPECT().AddUserToSegments(gomock.Any(), gomock.Any()).Return(nil, models.UserNotFoundErr)
			},
			expectedStatusCode: 500,
			expectedResponse:   `{"message":"user not found"}`,
			expectedBody:       `{"id":1000,"segments_to_add":["segment1","segment2"],"segments_to_delete":["segment3"]}`,
			expectedMethod:     http.MethodPost,
			path:               "/add-user-to-segment",
		},
		{
			testId:             5,
			name:               "wrong json",
			mockBehavior:       func(s *mock_http_delivery.MockUserSegmentService) {},
			expectedStatusCode: 400,
			expectedResponse:   `{"message":"invalid json"}`,
			expectedBody:       `{"id":asd,"segments_to_add":["segment1","segment2"],"segments_to_delete":["segment3"]}`,
			expectedMethod:     http.MethodPost,
			path:               "/add-user-to-segment",
		},
		{
			testId: 6,
			name:   "get all active segments",
			mockBehavior: func(s *mock_http_delivery.MockUserSegmentService) {
				s.EXPECT().GetActiveUserSegments(context.Background(), 1000).Return([]string{"segment1", "segment2"}, nil)
			},
			expectedStatusCode: 200,
			expectedResponse:   `["segment1","segment2"]`,
			expectedBody:       `{"userid":1000}`,
			expectedMethod:     http.MethodGet,
			path:               "/active-user-segments",
		},
		{
			testId: 7,
			name:   "user not found get all active segments",
			mockBehavior: func(s *mock_http_delivery.MockUserSegmentService) {
				s.EXPECT().GetActiveUserSegments(context.Background(), 1000).Return(nil, models.UserNotFoundErr)
			},
			expectedStatusCode: 500,
			expectedResponse:   `{"message":"user not found"}`,
			expectedBody:       `{"userid":1000}`,
			expectedMethod:     http.MethodGet,
			path:               "/active-user-segments",
		},
		{
			testId:             8,
			name:               "wrong method get all active segments",
			mockBehavior:       func(s *mock_http_delivery.MockUserSegmentService) {},
			expectedStatusCode: 405,
			path:               "/active-user-segments",
			expectedResponse:   `{"message":"method not provided"}`,
			expectedBody:       `{"userid":1000}`,
			expectedMethod:     http.MethodPost,
		},
		{
			testId:             9,
			name:               "wrong json get all active segments",
			mockBehavior:       func(s *mock_http_delivery.MockUserSegmentService) {},
			expectedStatusCode: 400,
			expectedResponse:   `{"message":"invalid json"}`,
			expectedBody:       `{"userid":asd}`,
			expectedMethod:     http.MethodGet,
			path:               "/active-user-segments",
		},
		{
			testId:             10,
			name:               "bad user id get all active segments",
			mockBehavior:       func(s *mock_http_delivery.MockUserSegmentService) {},
			expectedStatusCode: 400,
			expectedResponse:   `{"message":"bad user id"}`,
			expectedBody:       `{"userid":-1}`,
			expectedMethod:     http.MethodGet,
			path:               "/active-user-segments",
		},
		{
			testId:             11,
			name:               "bad user id add user to segment",
			mockBehavior:       func(s *mock_http_delivery.MockUserSegmentService) {},
			expectedStatusCode: 400,
			expectedResponse:   `{"message":"bad user id"}`,
			expectedBody:       `{"id":-1,"segments_to_add":["segment1","segment2"],"segments_to_delete":["segment3"]}`,
			expectedMethod:     http.MethodPost,
			path:               "/add-user-to-segment",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mock_http_delivery.NewMockUserSegmentService(ctrl)
			testCase.mockBehavior(service)

			handler := NewHandler(service)
			mux := http.NewServeMux()
			mux.HandleFunc("/add-user-to-segment", handler.addUserToSegment)
			mux.HandleFunc("/active-user-segments", handler.getActiveUserSegments)

			requestBody := bytes.NewBufferString(testCase.expectedBody)
			req := httptest.NewRequest(testCase.expectedMethod, testCase.path, requestBody)
			recorder := httptest.NewRecorder()

			mux.ServeHTTP(recorder, req)
			assert.Equal(t, testCase.expectedStatusCode, recorder.Code)
			assert.Equal(t, testCase.expectedResponse+"\n", recorder.Body.String())
		})
	}
}
