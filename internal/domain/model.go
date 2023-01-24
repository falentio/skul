package domain

import (
	"net/url"
	"strconv"
	"time"

	"gorm.io/gorm"

	"github.com/falentio/raid-go"
	"github.com/falentio/skul/internal/pkg/response"
)

const FileSlugPrefix = "fil"

type Model struct {
	ID        raid.Raid      `json:"id" gorm:"type:varchar(32);primaryKey" validate:"-"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

type PaginateOptions struct {
	Count  int `json:"count"`
	Offset int `json:"offset"`
	Page   int `json:"page"`
}

func (o *PaginateOptions) PageFromQuery(q url.Values) error {
	if q.Has("page") {
		v := q.Get("page")
		i, err := strconv.Atoi(v)
		if err != nil {
			return response.NewBadRequest(nil, "invalid url value for \"page\" query value, received %q", v)
		}

		o.Count = 10
		o.Offset = (i - 1) * 10
		o.Page = i
		return nil
	}
	if q.Has("offset") {
		v := q.Get("offset")
		i, err := strconv.Atoi(v)
		if err != nil {
			return response.NewBadRequest(nil, "invalid url value for \"offset\" query value, received %q", v)
		}
		o.Offset = i
	}
	if q.Has("count") {
		v := q.Get("count")
		i, err := strconv.Atoi(v)
		if err != nil {
			return response.NewBadRequest(nil, "invalid url value for \"count\" query value, received %q", v)
		}
		o.Count = i
	}
	if o.Count == 0 {
		o.Offset = 0
		o.Page = 1
	}
	o.Page = o.Offset/o.Count + 1
	return nil
}
