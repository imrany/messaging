package v1

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/imrany/smart_spore_hub/server/database/processes/course"
)

// GetAllCourses - GET /api/v1/courses
func GetAllCourses(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	w.Header().Set("content-type", "application/json")
	courses, err := course.List(ctx, nil, 50, 0)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(Response{
		Success: true,
		Data:    courses,
	})
}
