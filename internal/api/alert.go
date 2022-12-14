package api

import (
	"errors"
	"github.com/MeysamBavi/http-monitoring/internal/request"
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
	group.GET("/:id", h.get)
}

func (h *AlertHandler) get(c echo.Context) error {
	claims := h.JwtHandler.ParseToUserClaims(c)

	var alert request.Alert
	if err := c.Bind(&alert); err != nil {
		h.Logger.Error("error binding request",
			zap.Error(err),
			zap.Any("user_id", claims.UserId),
			zap.String("request_id", c.Response().Header().Get(echo.HeaderXRequestID)))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := alert.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	alerts, err := h.AlertStore.GetByUrlId(ctx, alert.ParseUrlId())

	if err != nil {
		var notFound store.NotFoundError
		if errors.As(err, &notFound) {
			h.Logger.Error("url not found",
				zap.Error(notFound),
				zap.Any("user_id", claims.UserId),
				zap.String("request_id", c.Response().Header().Get(echo.HeaderXRequestID)))
			return echo.NewHTTPError(http.StatusNotFound, "url not found")
		}

		h.Logger.Error("error getting alert", zap.Error(err))
		return echo.ErrInternalServerError
	}

	for i, a := range alerts {
		if a.UserId != *claims.UserId {
			if i > 0 {
				h.Logger.Warn("alerts with different user ids exist for url", zap.Any("url_id", alert.UrlId))
			}
			return echo.NewHTTPError(http.StatusForbidden, "not allowed to access this url")
		}
	}

	if alerts == nil {
		alerts = make([]*model.Alert, 0)
	}

	return c.JSON(http.StatusOK, alerts)
}
