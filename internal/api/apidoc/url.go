package apidoc

import (
	"github.com/MeysamBavi/http-monitoring/internal/model"
	"github.com/MeysamBavi/http-monitoring/internal/request"
	"github.com/labstack/echo/v4"
	"github.com/swaggest/openapi-go/openapi3"
	"net/http"
)

const (
	urlGroup = "/urls"
	urlTag   = "Urls"
)

func (d *DocGenerator) specifyUrlsCreateOperation() {
	op := openapi3.Operation{}
	op.
		WithSecurity(map[string][]string{securityName: {}}).
		WithSummary("Creates a new url for user").
		WithDescription("Creates a new url for user").
		WithID("createUrl").
		WithTags(urlTag)

	d.handleError(d.reflector.SetRequest(&op, new(request.URL), http.MethodPost))
	d.handleError(d.reflector.SetJSONResponse(&op, new(model.URL), http.StatusCreated))
	d.handleError(d.reflector.SetJSONResponse(&op, echo.NewHTTPError(http.StatusUnauthorized), http.StatusUnauthorized))
	d.handleError(d.reflector.SetJSONResponse(&op, echo.NewHTTPError(http.StatusBadRequest), http.StatusBadRequest))

	d.handleError(d.reflector.SpecEns().AddOperation(http.MethodPost, urlGroup+"", op))
}

func (d *DocGenerator) specifyUrlsGetAllOperation() {
	op := openapi3.Operation{}
	op.
		WithSecurity(map[string][]string{securityName: {}}).
		WithSummary("Returns all urls of user").
		WithDescription("Returns all urls of user in a list").
		WithID("getAllUrls").
		WithTags(urlTag)

	d.handleError(d.reflector.SetJSONResponse(&op, new([]model.URL), http.StatusOK))
	d.handleError(d.reflector.SetJSONResponse(&op, echo.NewHTTPError(http.StatusUnauthorized), http.StatusUnauthorized))
	d.handleError(d.reflector.SetJSONResponse(&op, echo.NewHTTPError(http.StatusNotFound), http.StatusNotFound))

	d.handleError(d.reflector.SpecEns().AddOperation(http.MethodGet, urlGroup+"", op))
}

func (d *DocGenerator) specifyUrlsGetDayStatsOperation() {
	op := openapi3.Operation{}
	op.
		WithSecurity(map[string][]string{securityName: {}}).
		WithSummary("Returns url monitoring stats").
		WithDescription("Returns monitoring stats for a specific url. Stats can be filtered using query parameters").
		WithID("getDayStats").
		WithTags(urlTag)

	d.handleError(d.reflector.SetRequest(&op, new(request.DayStats), http.MethodGet))
	d.handleError(d.reflector.SetJSONResponse(&op, new([]model.DayStat), http.StatusOK))
	d.handleError(d.reflector.SetJSONResponse(&op, echo.NewHTTPError(http.StatusUnauthorized), http.StatusUnauthorized))
	d.handleError(d.reflector.SetJSONResponse(&op, echo.NewHTTPError(http.StatusBadRequest), http.StatusBadRequest))
	d.handleError(d.reflector.SetJSONResponse(&op, echo.NewHTTPError(http.StatusNotFound), http.StatusNotFound))

	d.handleError(d.reflector.SpecEns().AddOperation(http.MethodGet, urlGroup+"/{id}/stats", op))
}
