package net

import (
	"github-events-microservices/api/config"
	"github-events-microservices/stores"
	mockstores "github-events-microservices/stores/mocks"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestRequestsHandler_List(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	type fields struct {
		storesMap          map[string]stores.ReadStore
		supportedDataTypes []string
	}
	type args struct {
		writer  httptest.ResponseRecorder
		request *http.Request
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "invalid data type",
			fields: fields{
				storesMap:          map[string]stores.ReadStore{},
				supportedDataTypes: []string{"events"},
			},
			args: args{
				writer: *httptest.NewRecorder(),
				request: &http.Request{URL: &url.URL{
					RawQuery: "dataType=invalidType",
				}},
			},
			wantErr: true,
		},
		{
			name: "invalid limit - not an int",
			fields: fields{
				storesMap:          map[string]stores.ReadStore{},
				supportedDataTypes: []string{"events"},
			},
			args: args{
				writer: *httptest.NewRecorder(),
				request: &http.Request{URL: &url.URL{
					RawQuery: "limit=invalidLimit",
				}},
			},
			wantErr: true,
		},
		{
			name: "invalid limit - negative int",
			fields: fields{
				storesMap:          map[string]stores.ReadStore{},
				supportedDataTypes: []string{"events"},
			},
			args: args{
				writer: *httptest.NewRecorder(),
				request: &http.Request{URL: &url.URL{
					RawQuery: "limit=-1",
				}},
			},
			wantErr: true,
		},
		{
			name: "invalid order type",
			fields: fields{
				storesMap:          map[string]stores.ReadStore{},
				supportedDataTypes: []string{"events"},
			},
			args: args{
				writer: *httptest.NewRecorder(),
				request: &http.Request{URL: &url.URL{
					RawQuery: "orderType=invalidOrderType",
				}},
			},
			wantErr: true,
		},
		{
			name: "empty data",
			fields: fields{
				storesMap:          mockStoresMap(mockCtrl),
				supportedDataTypes: []string{"events"},
			},
			args: args{
				writer: *httptest.NewRecorder(),
				request: &http.Request{URL: &url.URL{
					RawQuery: "dataType=events",
				}},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := RequestsHandler{
				storesMap:          tt.fields.storesMap,
				supportedDataTypes: tt.fields.supportedDataTypes,
			}
			receiver.List(&tt.args.writer, tt.args.request)
			if tt.wantErr {
				if tt.args.writer.Code != 400 {
					t.Errorf("expected: 400, got: %d", tt.args.writer.Code)
				}
			} else if !tt.wantErr && tt.args.writer.Code != 200 {
				t.Errorf("expected: 200, got: %d", tt.args.writer.Code)
			} else {
				result := string(tt.args.writer.Body.Bytes())
				if result != "[]" {
					t.Errorf("expected: [], got: %s", result)
				}
			}

		})
	}
}

func mockStoresMap(mockCtrl *gomock.Controller) map[string]stores.ReadStore {
	storesMap := make(map[string]stores.ReadStore)
	eventsStoreMock := mockstores.NewMockReadStore(mockCtrl)
	eventsStoreMock.EXPECT().
		Get(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)
	storesMap[config.ApiConfiguration.EventsCollection] = eventsStoreMock
	return storesMap
}
