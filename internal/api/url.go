package api

import (
	"errors"
	"github.com/MeysamBavi/http-monitoring/internal/auth"
	"github.com/MeysamBavi/http-monitoring/internal/model"
	"github.com/MeysamBavi/http-monitoring/internal/request"
	"github.com/MeysamBavi/http-monitoring/internal/store"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"net/http"
)

type UrlHandler struct {
	Logger     *zap.Logger
	UrlStore   store.Url
	JwtHandler *auth.JwtHandler
}

func (h *UrlHandler) Register(group *echo.Group) {
	group.Use(middleware.JWTWithConfig(h.JwtHandler.Config()))
	group.GET("", h.getAll)
	group.POST("", h.create)
	group.GET("/:id/stats", h.getDayStats)
}

func (h *UrlHandler) create(c echo.Context) error {
	claims := h.JwtHandler.ParseToUserClaims(c)

	var req request.URL

	if err := c.Bind(&req); err != nil {
		h.Logger.Error("error binding request", zap.Error(err),
			zap.Any("user_id", claims.UserId),
			zap.String("request_id", c.Response().Header().Get(echo.HeaderXRequestID)))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := req.Validate(); err != nil {
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
		h.Logger.Error("error adding url", zap.Error(err),
			zap.Any("user_id", claims.UserId),
			zap.String("request_id", c.Response().Header().Get(echo.HeaderXRequestID)))
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusCreated, url)
}

func (h *UrlHandler) getAll(c echo.Context) error {
	claims := h.JwtHandler.ParseToUserClaims(c)

	ctx := c.Request().Context()
	urls, err := h.UrlStore.GetByUserId(ctx, *claims.UserId)

	if err != nil {
		var notFound store.NotFoundError
		if errors.As(err, &notFound) {
			h.Logger.Error("error getting user urls", zap.Error(notFound),
				zap.Any("user_id", claims.UserId),
				zap.String("request_id", c.Response().Header().Get(echo.HeaderXRequestID)))
			return echo.NewHTTPError(http.StatusNotFound, "no urls found")
		}

		h.Logger.Error("error getting user urls", zap.Error(err),
			zap.Any("user_id", claims.UserId),
			zap.String("request_id", c.Response().Header().Get(echo.HeaderXRequestID)))
		return echo.ErrInternalServerError
	}

	if urls == nil {
		urls = make([]*model.URL, 0)
	}

	return c.JSON(http.StatusOK, urls)
}

func (h *UrlHandler) getDayStats(c echo.Context) error {
	claims := h.JwtHandler.ParseToUserClaims(c)

	var dayStats request.DayStats
	if err := c.Bind(&dayStats); err != nil {
		h.Logger.Error("error binding the request", zap.Error(err),
			zap.Any("user_id", claims.UserId),
			zap.String("request_id", c.Response().Header().Get(echo.HeaderXRequestID)))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := dayStats.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	stats, err := h.UrlStore.GetDayStats(ctx, *claims.UserId, dayStats.ParseUrlId(), dayStats.DayFilter())

	if err != nil {
		var notFound store.NotFoundError
		if errors.As(err, &notFound) {
			h.Logger.Error("error getting url stats", zap.Error(notFound),
				zap.Any("user_id", claims.UserId),
				zap.String("request_id", c.Response().Header().Get(echo.HeaderXRequestID)))
			return echo.NewHTTPError(http.StatusNotFound, "no stats found for url")
		}

		h.Logger.Error("error getting day stats", zap.Error(err),
			zap.Any("user_id", claims.UserId),
			zap.String("request_id", c.Response().Header().Get(echo.HeaderXRequestID)))
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, stats)
}
