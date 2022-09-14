package api

import (
	"errors"
	"net/http"

	"github.com/MeysamBavi/http-monitoring/internal/auth"
	"github.com/MeysamBavi/http-monitoring/internal/model"
	"github.com/MeysamBavi/http-monitoring/internal/store"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

type AlertHandler struct {
	Logger     *zap.Logger
	AlertStore store.Alert
	JwtHandler *auth.JwtHandler
}

func (h *AlertHandler) Register(group *echo.Group) {
	group.Use(middleware.JWTWithConfig(h.JwtHandler.Config()))
	group.GET("/get", h.get)
}

func (h *AlertHandler) get(c echo.Context) error {
	claims := h.JwtHandler.ParseToUserClaims(c)

	urlId, err := model.ParseId(c.QueryParam("id"))
	if err != nil {
		h.Logger.Error("error parsing url id", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "could not parse url id from query parameter")
	}

	ctx := c.Request().Context()
	alerts, err := h.AlertStore.GetByUrlId(ctx, urlId)

	if err != nil {
		var notFound store.NotFoundError
		if errors.As(err, &notFound) {
			h.Logger.Error("url not found", zap.Error(err))
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
	}

	for i, a := range alerts {
		if a.UserId != *claims.UserId {
			if i > 0 {
				h.Logger.Warn("alerts with different user ids exist for url", zap.Any("url_id", urlId))
			}
			return echo.NewHTTPError(http.StatusForbidden, "not allowed to access this url")
		}
	}

	if alerts == nil {
		alerts = make([]*model.Alert, 0)
	}

	return c.JSON(http.StatusOK, alerts)
}
