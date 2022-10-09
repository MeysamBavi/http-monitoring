package apidoc

import (
	"github.com/MeysamBavi/http-monitoring/internal/model"
	"github.com/MeysamBavi/http-monitoring/internal/request"
	"github.com/labstack/echo/v4"
	"github.com/swaggest/openapi-go/openapi3"
	"net/http"
)

const (
	userGroup = "/users"
)

func (d *DocGenerator) specifyUsersCreateOperation() {
	op := openapi3.Operation{}
	op.
		WithSummary("Creates a new user").
		WithDescription("Creates a new user with the given username and password")

	d.handleError(d.reflector.SetRequest(&op, new(request.User), http.MethodPost))
	d.handleError(d.reflector.SetJSONResponse(&op, new(model.User), http.StatusCreated))
	d.handleError(d.reflector.SetJSONResponse(&op, echo.NewHTTPError(http.StatusBadRequest), http.StatusBadRequest))

	d.handleError(d.reflector.SpecEns().AddOperation(http.MethodPost, userGroup+"", op))
}

func (d *DocGenerator) specifyUsersLoginOperation() {
	op := openapi3.Operation{}
	op.
		WithSummary("Authenticates user and generates JWT token").
		WithDescription("Authenticates user and generates JWT token")

	d.handleError(d.reflector.SetRequest(&op, new(request.User), http.MethodPost))
	d.handleError(d.reflector.SetStringResponse(&op, http.StatusOK, "JWT token"))
	d.handleError(d.reflector.SetJSONResponse(&op, echo.NewHTTPError(http.StatusBadRequest), http.StatusBadRequest))
	d.handleError(d.reflector.SetJSONResponse(&op, echo.NewHTTPError(http.StatusNotFound), http.StatusNotFound))
	d.handleError(d.reflector.SetJSONResponse(&op, echo.NewHTTPError(http.StatusUnauthorized), http.StatusUnauthorized))

	d.handleError(d.reflector.SpecEns().AddOperation(http.MethodPost, userGroup+"/login", op))
}
