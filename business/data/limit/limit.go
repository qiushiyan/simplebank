package limit

import (
	"net/http"

	"github.com/qiushiyan/simplebank/foundation/validate"
	"github.com/spf13/cast"
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

func WithPaging(pageId int32, pageSize int32) Limiter {
	return Limiter{
		Limit:  pageSize,
		Offset: (pageId - 1) * pageSize,
	}
}

func Parse(r *http.Request, defaultPageId int32, defaultPageSize int32) (Limiter, error) {
	id := defaultPageId
	size := defaultPageSize

	values := r.URL.Query()

	if pageId := values.Get("page_id"); pageId != "" {
		pid, err := cast.ToInt32E(pageId)
		if err != nil {
			return Limiter{}, err
		}
		id = pid
	}

	if pageSize := values.Get("page_size"); pageSize != "" {
		ps, err := cast.ToInt32E(pageSize)
		if err != nil {
			return Limiter{}, err
		}
		size = ps
	}

	l := WithPaging(id, size)

	if err := validate.Check(l); err != nil {
		return Limiter{}, err
	}

	return l, nil
}
