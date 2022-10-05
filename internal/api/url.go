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
	"net/url"
	"strconv"
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
			return echo.NewHTTPError(http.StatusNotFound, "no urls found")
		}

		h.Logger.Error("error getting user urls", zap.Error(err))
		return echo.ErrInternalServerError
	}

	if urls == nil {
		urls = make([]*model.URL, 0)
	}

	return c.JSON(http.StatusOK, urls)
}

func (h *UrlHandler) getDayStats(c echo.Context) error {
	claims := h.JwtHandler.ParseToUserClaims(c)

	urlId, err := model.ParseId(c.Param("id"))
	if err != nil {
		h.Logger.Error("error parsing url id", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "invalid url id")
	}

	dateFilter, err := generateDayFilter(c.QueryParams())

	if err != nil {
		h.Logger.Error("error parsing date", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "invalid query params")
	}

	ctx := c.Request().Context()
	stats, err := h.UrlStore.GetDayStats(ctx, *claims.UserId, urlId, dateFilter)

	if err != nil {
		var notFound store.NotFoundError
		if errors.As(err, &notFound) {
			h.Logger.Error("error getting url stats", zap.Error(notFound))
			return echo.NewHTTPError(http.StatusNotFound, "no stats found for url")
		}

		return echo.ErrBadGateway.Internal
	}

	return c.JSON(http.StatusOK, stats)
}

func generateDayFilter(values url.Values) (func(date model.Date) bool, error) {
	hasDay := values.Has("day")
	hasMonth := values.Has("month")
	hasYear := values.Has("year")

	var day, month, year int
	var err error
	if hasDay {
		day, err = strconv.Atoi(values.Get("day"))
		if err != nil {
			return nil, errors.New("invalid day")
		}
	}

	if hasMonth {
		month, err = strconv.Atoi(values.Get("month"))
		if err != nil {
			return nil, errors.New("invalid month")
		}
	}

	if hasYear {
		year, err = strconv.Atoi(values.Get("year"))
		if err != nil {
			return nil, errors.New("invalid year")
		}
	}

	return func(date model.Date) bool {
		if hasDay && date.Day != day {
			return false
		}

		if hasMonth && date.Month != month {
			return false
		}

		if hasYear && date.Year != year {
			return false
		}
		return true
	}, nil
}
