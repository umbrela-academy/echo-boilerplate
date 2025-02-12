package users

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/alexferl/echo-boilerplate/util"
)

type ListUsersResponse struct {
	Users []*PublicUser `json:"users"`
}

func (h *Handler) ListUsers(c echo.Context) error {
	page, perPage, limit, skip := util.ParsePaginationParams(c)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	count, err := h.Mapper.Count(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed counting users: %v", err)
	}

	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(skip))
	result, err := h.Mapper.Find(ctx, nil, []*PublicUser{}, opts)
	if err != nil {
		return fmt.Errorf("failed getting users: %v", err)
	}

	uri := fmt.Sprintf("http://%s%s", c.Request().Host, c.Request().URL.Path)
	util.SetPaginationHeaders(c.Response().Header(), int(count), page, perPage, uri)

	resp := &ListUsersResponse{Users: result.([]*PublicUser)}

	return h.Validate(c, http.StatusOK, resp)
}
