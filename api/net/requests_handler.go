package net

import (
	"encoding/json"
	"errors"
	"fmt"
	"github-events-microservices/api/config"
	"github-events-microservices/model"
	"github-events-microservices/stores"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type RequestsHandler struct {
	storesMap          map[string]stores.ReadStore
	supportedDataTypes []string
}

type ApiError struct {
	Error string `json:"error"`
}

type DataCount struct {
	Count int64 `json:"count"`
}

type ListParams struct {
	DataType string
	Limit    int64
	OrderBy  stores.OrderBy
}

func (receiver RequestsHandler) List(writer http.ResponseWriter, request *http.Request) {
	listParams, err := parseListParams(request)
	if err != nil {
		writeError(writer, http.StatusBadRequest, err.Error())
		return
	}
	store := receiver.storesMap[listParams.DataType]
	if store == nil {
		supportedDataTypes := []string{config.ApiConfiguration.EventsCollection, config.ApiConfiguration.ReposCollection, config.ApiConfiguration.UsersCollection}
		errorMessage := fmt.Sprintf("unknown data type: '%s'. Supported data types are: %s.", listParams.DataType, strings.Join(supportedDataTypes, ", "))
		writeError(writer, http.StatusBadRequest, errorMessage)
	} else {
		var results, _ = createResults(listParams.DataType)
		err := store.Get(listParams.Limit, listParams.OrderBy, &results)
		if err != nil {
			errorMessage := fmt.Sprintf("failed to list items: %s", err.Error())
			slog.Error(errorMessage)
			writeError(writer, http.StatusBadRequest, errorMessage)
		} else {
			writeJsonResponse(writer, results, listParams.DataType)
		}
	}
}

func (receiver RequestsHandler) Count(writer http.ResponseWriter, request *http.Request) {
	storeKey := getStoreKey(request)
	store := receiver.storesMap[storeKey]
	if store == nil {
		errorMessage := fmt.Sprintf("unknown data type: '%s'", storeKey)
		writeError(writer, http.StatusBadRequest, errorMessage)
	} else {
		count, err := store.Count()
		if err != nil {
			errorMessage := fmt.Sprintf("Failed to count '%s'", storeKey)
			writeError(writer, http.StatusInternalServerError, errorMessage)
		} else {
			writeJsonResponse(writer, DataCount{Count: count}, "count")
		}
	}
}

func parseListParams(request *http.Request) (*ListParams, error) {
	limit, err := getLimit(request)
	if err != nil {
		return nil, err
	}

	orderBy, err := getOrderBy(request)
	if err != nil {
		return nil, err
	}

	return &ListParams{
		DataType: getParam(request.URL.Query(), config.DataType, config.DefaultDataType),
		Limit:    *limit,
		OrderBy:  *orderBy,
	}, nil
}

func getStoreKey(request *http.Request) string {
	return getParam(request.URL.Query(), config.DataType, config.DefaultDataType)
}

func getLimit(request *http.Request) (*int64, error) {
	val, err := parseIntParam(config.LimitParamKey, request.URL.Query().Get(config.LimitParamKey), config.DefaultLimit)
	if err != nil {
		return nil, err
	} else if *val < 0 {
		return nil, errors.New("invalid limit. Limit must be a non-negative integer")
	}

	limit := int64(*val)
	return &limit, nil
}

func getOrderBy(request *http.Request) (*stores.OrderBy, error) {
	orderByColumn := getParam(request.URL.Query(), config.OrderByColumnQueryParam, config.DefaultOrderByColumn)
	orderTypeString := getParam(request.URL.Query(), config.OrderTypeQueryParam, config.Ascending)
	var orderType int
	if orderTypeString == config.Ascending {
		orderType = 1
	} else if orderTypeString == config.Descending {
		orderType = -1
	} else {
		return nil, errors.New(fmt.Sprintf("Invalid order type: '%s'. Please specify one of the following order types: '%s', '%s'", orderTypeString, config.Ascending, config.Descending))
	}

	return &stores.OrderBy{
		Column: orderByColumn,
		Order:  orderType,
	}, nil
}

func getParam(params url.Values, paramKey string, defaultValue string) string {
	value := params.Get(paramKey)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func parseIntParam(paramKey string, paramValue string, defaultValue int) (*int, error) {
	if len(paramValue) == 0 {
		return &defaultValue, nil
	} else {
		limit, err := strconv.Atoi(paramValue)
		if err != nil {
			errorMessage := fmt.Sprintf("invalid '%s': %s", paramKey, paramValue)
			slog.Warn(errorMessage)
			return nil, errors.New(errorMessage)
		}
		return &limit, nil
	}
}

func createResults(key string) (interface{}, error) {
	if key == config.ApiConfiguration.EventsCollection {
		return make([]model.Event, 0), nil
	} else if key == config.ApiConfiguration.ReposCollection {
		return make([]model.Repo, 0), nil
	} else if key == config.ApiConfiguration.UsersCollection {
		return make([]model.User, 0), nil
	} else {
		errorMessage := fmt.Sprintf("unknwon data type: %s", key)
		slog.Error(errorMessage)
		return nil, fmt.Errorf(errorMessage)
	}
}

func writeError(writer http.ResponseWriter, errorCode int, errorMessage string) {
	writer.WriteHeader(errorCode)

	bytes, err := json.Marshal(ApiError{Error: errorMessage})
	if err != nil {
		slog.Error("failed to serialize API error")
	} else {
		_, err := writer.Write(bytes)
		if err != nil {
			slog.Error("Failed to write error response")
		}
	}
}

func writeJsonResponse(writer http.ResponseWriter, response interface{}, responseKey string) {
	writer.WriteHeader(http.StatusOK)
	bytes, serializationError := json.Marshal(response)
	if serializationError != nil {
		slog.Error(fmt.Sprintf("failed to serialize response '%s' to json. reason: %s", responseKey, serializationError.Error()))
	}
	_, writeError := writer.Write(bytes)
	if writeError != nil {
		slog.Error(fmt.Sprintf("failed to write json response '%s' to json. reason: %s", string(bytes), writeError.Error()))
	}
}

func NewRequestsHandler(url string, port string) *RequestsHandler {
	storesMap := make(map[string]stores.ReadStore)
	storesMap[config.ApiConfiguration.EventsCollection] = createMongoStore(url, port, config.ApiConfiguration.EventsDb, config.ApiConfiguration.EventsCollection)
	storesMap[config.ApiConfiguration.ReposCollection] = createMongoStore(url, port, config.ApiConfiguration.ReposDb, config.ApiConfiguration.ReposCollection)
	storesMap[config.ApiConfiguration.UsersCollection] = createMongoStore(url, port, config.ApiConfiguration.UsersDb, config.ApiConfiguration.UsersCollection)

	return &RequestsHandler{
		storesMap:          storesMap,
		supportedDataTypes: []string{config.ApiConfiguration.EventsCollection, config.ApiConfiguration.ReposCollection, config.ApiConfiguration.UsersCollection},
	}
}

func createMongoStore(url string, port string, database string, collection string) *stores.MongoDbCollectionStore {
	mongoDbStore, err := stores.NewMongoDbStore(url, port, database, collection)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	return mongoDbStore
}
