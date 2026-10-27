package httpx

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type ListParams struct {
	Page  int
	Limit int
	Sort  string
	Desc  bool
}

func ParseListParams(c *fiber.Ctx, allowedSort map[string]bool, defSort string) ListParams {
	page := 1
	limit := 20
	if v := c.Query("page"); v != "" {
		if p, err := strconv.Atoi(v); err == nil && p > 0 {
			page = p
		} else {
			page = -1 // indicate invalid? caller can handle; but we just parse simply
		}
	}
	if v := c.Query("limit"); v != "" {
		if l, err := strconv.Atoi(v); err == nil && l > 0 {
			if l > 100 {
				l = 100
			}
			limit = l
		} else {
			limit = -1
		}
	}
	sort := defSort
	if v := c.Query("sort"); v != "" {
		s := v
		desc := false
		if strings.HasPrefix(s, "-") {
			desc = true
			s = strings.TrimPrefix(s, "-")
		}
		if allowedSort != nil {
			if _, ok := allowedSort[s]; ok {
				sort = s
			} else {
				sort = defSort
			}
		} else {
			sort = s
		}
		return ListParams{Page: page, Limit: limit, Sort: sort, Desc: desc}
	}
	return ListParams{Page: page, Limit: limit, Sort: sort, Desc: false}
}

func (p ListParams) Offset() int {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit < 1 {
		p.Limit = 20
	}
	return (p.Page - 1) * p.Limit
}

func TotalPages(total int64, limit int) int {
	if limit <= 0 {
		limit = 20
	}
	if total <= 0 {
		return 0
	}
	pages := int(total) / limit
	if int(total)%limit != 0 {
		pages++
	}
	return pages
}
