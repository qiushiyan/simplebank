package limit

import (
	"net/http"
	"strconv"

	"github.com/qiushiyan/simplebank/foundation/validate"
)

// Limiter represents the limit and offset for a database query
type Limiter struct {
	Limit  int32 `json:"limit"  validate:"min=1,max=100"`
	Offset int32 `json:"offset" validate:"min=0"`
}

func NewLimiter() Limiter {
	return Limiter{
		Limit:  10,
		Offset: 0,
	}
}

func WithPage(pageId int32, pageSize int32) Limiter {
	return Limiter{
		Limit:  pageSize,
		Offset: (pageId - 1) * pageSize,
	}
}

func Parse(r *http.Request, defaultPageId int32, defaultPageSize int32) (Limiter, error) {
	pageId := r.URL.Query().Get("page_id")
	pageSize := r.URL.Query().Get("page_size")

	if pageId == "" {
		pageId = strconv.Itoa(int(defaultPageId))
	}

	if pageSize == "" {
		pageSize = strconv.Itoa(int(defaultPageSize))
	}

	id, err := strconv.Atoi(pageId)
	if err != nil {
		return Limiter{}, err
	}

	size, err := strconv.Atoi(pageSize)
	if err != nil {
		return Limiter{}, err
	}

	limiter := WithPage(int32(id), int32(size))
	if err := validate.Check(limiter); err != nil {
		return Limiter{}, err
	}

	return limiter, nil
}
