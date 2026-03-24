package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOKResponse(t *testing.T) {

	type sampleResponse struct {
		Product struct {
			Name string `json:"name"`
		}
	}

	sample := ApiResponse[sampleResponse]{
		Success: true,
		Message: "Success",
		Data: sampleResponse{
			Product: struct {
				Name string `json:"name"`
			}{
				Name: "test",
			},
		},
	}

	t.Run("succesful http200 json response", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		OKResponse(recorder, http.StatusOK, sample)

		assert.Equal(t, http.StatusOK, recorder.Code, "Expected status code 200 OK")
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"), "Expected Content-Type to be application/json")

		//deserialize body
		var response ApiResponse[sampleResponse]
		json.Unmarshal(recorder.Body.Bytes(), &response)

		assert.Equal(t, true, response.Success, "Expected success to be true")
		assert.Equal(t, "Success", response.Message, "Expected message to be Success")
		assert.Equal(t, "test", response.Data.Product.Name, "Expected product name to be test")
	})
}

func TestErrorResponse(t *testing.T) {
	t.Run("json response for a given http status code", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ErrorResponse(recorder, http.StatusInternalServerError, "Some error occurred")

		assert.Equal(t, http.StatusInternalServerError, recorder.Code, "Expected status code 500 Internal Server Error")
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"), "Expected Content-Type to be application/json")

		//deserialize body
		var response ApiResponse[any]
		json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.Equal(t, "Some error occurred", response.Message, "Some error occurred")
	})
}
