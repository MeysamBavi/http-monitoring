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
)

func (d *DocGenerator) specifyUrlsCreateOperation() {
	op := openapi3.Operation{}
	op.
		WithSummary("Creates a new url for user").
		WithDescription("Creates a new url for user")

	d.handleError(d.reflector.SetRequest(&op, new(request.URL), http.MethodPost))
	d.handleError(d.reflector.SetJSONResponse(&op, new(model.URL), http.StatusCreated))
	d.handleError(d.reflector.SetJSONResponse(&op, echo.NewHTTPError(http.StatusUnauthorized), http.StatusUnauthorized))
	d.handleError(d.reflector.SetJSONResponse(&op, echo.NewHTTPError(http.StatusBadRequest), http.StatusBadRequest))

	d.handleError(d.reflector.SpecEns().AddOperation(http.MethodPost, urlGroup+"", op))
}

func (d *DocGenerator) specifyUrlsGetAllOperation() {
	op := openapi3.Operation{}
	op.
		WithSummary("Returns all urls of user").
		WithDescription("Returns all urls of user")

	d.handleError(d.reflector.SetJSONResponse(&op, new([]model.URL), http.StatusOK))
	d.handleError(d.reflector.SetJSONResponse(&op, echo.NewHTTPError(http.StatusUnauthorized), http.StatusUnauthorized))
	d.handleError(d.reflector.SetJSONResponse(&op, echo.NewHTTPError(http.StatusNotFound), http.StatusNotFound))

	d.handleError(d.reflector.SpecEns().AddOperation(http.MethodGet, urlGroup+"", op))
}

func (d *DocGenerator) specifyUrlsGetDayStatsOperation() {
	op := openapi3.Operation{}
	op.
		WithSummary("Returns stats for a day").
		WithDescription("Returns stats for a day")

	d.handleError(d.reflector.SetRequest(&op, new(request.DayStats), http.MethodGet))
	d.handleError(d.reflector.SetJSONResponse(&op, new([]model.DayStat), http.StatusOK))
	d.handleError(d.reflector.SetJSONResponse(&op, echo.NewHTTPError(http.StatusUnauthorized), http.StatusUnauthorized))
	d.handleError(d.reflector.SetJSONResponse(&op, echo.NewHTTPError(http.StatusBadRequest), http.StatusBadRequest))
	d.handleError(d.reflector.SetJSONResponse(&op, echo.NewHTTPError(http.StatusNotFound), http.StatusNotFound))

	d.handleError(d.reflector.SpecEns().AddOperation(http.MethodGet, urlGroup+"/{id}/stats", op))
}
