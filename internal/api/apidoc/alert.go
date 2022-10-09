package apidoc

import (
	"github.com/MeysamBavi/http-monitoring/internal/model"
	"github.com/MeysamBavi/http-monitoring/internal/request"
	"github.com/labstack/echo/v4"
	"github.com/swaggest/openapi-go/openapi3"
	"net/http"
)

const (
	alertGroup = "/alerts"
	alertTag   = "Alerts"
)

func (d *DocGenerator) specifyAlertsGetOperation() {
	op := openapi3.Operation{}
	op.
		WithSecurity(map[string][]string{securityName: {}}).
		WithSummary("Gets all alerts").
		WithDescription("Gets all alerts").
		WithID("getAlerts").
		WithTags(alertTag)

	d.handleError(d.reflector.SetRequest(&op, new(request.Alert), http.MethodGet))
	d.handleError(d.reflector.SetJSONResponse(&op, new([]model.Alert), http.StatusOK))
	d.handleError(d.reflector.SetJSONResponse(&op, echo.NewHTTPError(http.StatusUnauthorized), http.StatusUnauthorized))
	d.handleError(d.reflector.SetJSONResponse(&op, echo.NewHTTPError(http.StatusBadRequest), http.StatusBadRequest))
	d.handleError(d.reflector.SetJSONResponse(&op, echo.NewHTTPError(http.StatusNotFound), http.StatusNotFound))
	d.handleError(d.reflector.SetJSONResponse(&op, echo.NewHTTPError(http.StatusForbidden), http.StatusForbidden))

	d.handleError(d.reflector.SpecEns().AddOperation(http.MethodGet, alertGroup+"/{id}", op))
}
