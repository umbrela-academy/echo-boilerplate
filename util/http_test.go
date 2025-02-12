package util

import (
	"fmt"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaginate(t *testing.T) {
	testCases := []struct {
		query       string
		link        string
		xPage       string
		xPerPage    string
		xTotal      string
		xTotalPages string
		xNextPage   string
		xPrevPage   string
	}{
		{
			query: "",
			link: `<http://example.com/users?per_page=10&page=1>; rel=last, ` +
				`<http://example.com/users?per_page=10&page=1>; rel=first`,
			xPage:       "1",
			xPerPage:    "10",
			xTotal:      "10",
			xTotalPages: "1",
			xNextPage:   "",
			xPrevPage:   "",
		},
		{
			query: "per_page=10&page=1",
			link: `<http://example.com/users?per_page=10&page=1>; rel=last, ` +
				`<http://example.com/users?per_page=10&page=1>; rel=first`,
			xPage:       "1",
			xPerPage:    "10",
			xTotal:      "10",
			xTotalPages: "1",
			xNextPage:   "",
			xPrevPage:   "",
		},
		{
			query: "per_page=1&page=2",
			link: `<http://example.com/users?per_page=1&page=3>; rel=next, ` +
				`<http://example.com/users?per_page=1&page=10>; rel=last, ` +
				`<http://example.com/users?per_page=1&page=1>; rel=first, ` +
				`<http://example.com/users?per_page=1&page=1>; rel=prev`,
			xPage:       "2",
			xPerPage:    "1",
			xTotal:      "10",
			xTotalPages: "10",
			xNextPage:   "3",
			xPrevPage:   "1",
		},
		{
			query: "per_page=1&page=10",
			link: `<http://example.com/users?per_page=1&page=10>; rel=last, ` +
				`<http://example.com/users?per_page=1&page=1>; rel=first, ` +
				`<http://example.com/users?per_page=1&page=9>; rel=prev`,
			xPage:       "10",
			xPerPage:    "1",
			xTotal:      "10",
			xTotalPages: "10",
			xNextPage:   "",
			xPrevPage:   "9",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.query, func(t *testing.T) {
			resp := httptest.NewRecorder()

			total, _ := strconv.Atoi(tc.xTotal)
			page, _ := strconv.Atoi(tc.xPage)
			perPage, _ := strconv.Atoi(tc.xPerPage)

			uri := fmt.Sprintf("http://%s%s", "example.com", "/users")
			SetPaginationHeaders(resp.Header(), total, page, perPage, uri)

			h := resp.Header()

			assert.Equal(t, tc.xPage, h.Get("X-Page"))
			assert.Equal(t, tc.xPerPage, h.Get("X-Per-Page"))
			assert.Equal(t, tc.xTotal, h.Get("X-Total"))
			assert.Equal(t, tc.xTotalPages, h.Get("X-Total-Pages"))
			assert.Equal(t, tc.xNextPage, h.Get("X-Next-Page"))
			assert.Equal(t, tc.xPrevPage, h.Get("X-Prev-Page"))
			assert.Equal(t, tc.link, h.Get("Link"))
		})
	}
}
