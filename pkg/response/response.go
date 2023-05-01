package response

import (
	"net/http"
)

// JSON HTTP response
func JSON(w http.ResponseWriter, response Responser) error {
	return dr.JSON(w, response)
}
