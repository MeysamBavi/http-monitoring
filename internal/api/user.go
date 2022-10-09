package api

import (
	"errors"
	"net/http"

	"github.com/MeysamBavi/http-monitoring/internal/auth"
	"github.com/MeysamBavi/http-monitoring/internal/model"
	"github.com/MeysamBavi/http-monitoring/internal/request"
	"github.com/MeysamBavi/http-monitoring/internal/store"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type UserHandler struct {
	Logger     *zap.Logger
	UserStore  store.User
	JwtHandler *auth.JwtHandler
}

func (h *UserHandler) Register(group *echo.Group) {
	group.POST("", h.create)
	group.POST("/login", h.login)
}

func (h *UserHandler) create(c echo.Context) error {

	var req request.User

	if err := c.Bind(&req); err != nil {
		h.Logger.Error("error binding request", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := req.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	user := model.User{
		Username: req.Username,
		Password: req.Password,
	}

	err := h.UserStore.Add(ctx, &user)

	if err != nil {
		var duplicate store.DuplicateError
		if errors.As(err, &duplicate) {
			h.Logger.Error("duplicate username", zap.Error(duplicate))
			return echo.NewHTTPError(http.StatusBadRequest, "this username is already taken")
		}

		h.Logger.Error("error adding user", zap.Error(err))
		return echo.ErrInternalServerError
	}

	h.Logger.Info("user created", zap.Any("user", user))

	return c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) login(c echo.Context) error {
	var req request.User

	if err := c.Bind(&req); err != nil {
		h.Logger.Error("error binding request", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := req.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	user, err := h.UserStore.GetByUsername(ctx, req.Username)

	if err != nil {
		var notFound store.NotFoundError
		if errors.As(err, &notFound) {
			h.Logger.Error("user not found", zap.Error(notFound))
			return echo.NewHTTPError(http.StatusNotFound, "user not found")
		}

		h.Logger.Error("error getting user", zap.Error(err))
		return echo.ErrInternalServerError
	}

	if user.Password != req.Password {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid password")
	}

	h.Logger.Info("user logged in", zap.Any("user", user))

	token, err := h.JwtHandler.GenerateFromUser(user)

	if err != nil {
		h.Logger.Error("error generating token", zap.Error(err))
		return echo.ErrInternalServerError
	}

	return c.String(http.StatusOK, token)
}
