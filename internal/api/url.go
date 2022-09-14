package api

import (
	"errors"
	"net/http"

	"github.com/MeysamBavi/http-monitoring/internal/auth"
	"github.com/MeysamBavi/http-monitoring/internal/model"
	"github.com/MeysamBavi/http-monitoring/internal/request"
	"github.com/MeysamBavi/http-monitoring/internal/store"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

type UrlHandler struct {
	Logger     *zap.Logger
	UrlStore   store.Url
	JwtHandler *auth.JwtHandler
}

func (h *UrlHandler) Register(group *echo.Group) {
	group.Use(middleware.JWTWithConfig(h.JwtHandler.Config()))
	group.POST("/create", h.create)
	group.GET("/get", h.getAll)
	group.GET("/stat", h.getDayStat)
}

func (h *UrlHandler) create(c echo.Context) error {
	claims := h.JwtHandler.ParseToUserClaims(c)

	var req request.URL

	if err := c.Bind(&req); err != nil {
		h.Logger.Error("error binding request", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "could not parse request")
	}

	if err := req.Validate(); err != nil {
		h.Logger.Error("validation error", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	url := &model.URL{
		UserId:    *claims.UserId,
		Url:       req.Url,
		Threshold: req.Threshold,
		Interval:  req.Interval,
	}

	err := h.UrlStore.Add(ctx, url)

	if err != nil {
		h.Logger.Error("error adding url", zap.Error(err))
		return echo.ErrInternalServerError
	}

	h.Logger.Info("url created", zap.Any("url", url))

	return c.JSON(http.StatusCreated, url)
}

func (h *UrlHandler) getAll(c echo.Context) error {
	claims := h.JwtHandler.ParseToUserClaims(c)

	ctx := c.Request().Context()
	urls, err := h.UrlStore.GetByUserId(ctx, *claims.UserId)

	if err != nil {
		var notFound store.NotFoundError
		if errors.As(err, &notFound) {
			h.Logger.Error("error getting user urls", zap.Error(notFound))
			return echo.NewHTTPError(http.StatusNotFound, notFound)
		}

		h.Logger.Error("error getting user urls", zap.Error(err))
		return echo.ErrInternalServerError
	}

	if urls == nil {
		urls = make([]*model.URL, 0)
	}

	return c.JSON(http.StatusOK, urls)
}

func (h *UrlHandler) getDayStat(c echo.Context) error {
	claims := h.JwtHandler.ParseToUserClaims(c)

	urlId, err := model.ParseId(c.QueryParam("id"))
	if err != nil {
		h.Logger.Error("error parsing url id", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "could not parse url id from query parameter")
	}

	date, err := model.ParseDate(c.QueryParam("date"))
	if err != nil {
		h.Logger.Error("error parsing date", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "could not parse date from query parameter")
	}

	ctx := c.Request().Context()
	stat, err := h.UrlStore.GetDayStat(ctx, *claims.UserId, urlId, date)

	if err != nil {
		var notFound store.NotFoundError
		if errors.As(err, &notFound) {
			h.Logger.Error("error getting url stat", zap.Error(notFound))
			return echo.NewHTTPError(http.StatusNotFound, notFound)
		}

		return echo.ErrBadGateway.Internal
	}

	return c.JSON(http.StatusOK, stat)
}
