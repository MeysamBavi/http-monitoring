package apidoc

import (
	"github.com/swaggest/openapi-go/openapi3"
	"go.uber.org/zap"
)

const (
	securityName = "bearer_token"
)

type DocGenerator struct {
	logger    *zap.Logger
	reflector *openapi3.Reflector
}

func NewDocGenerator(logger *zap.Logger) *DocGenerator {
	return &DocGenerator{logger: logger, reflector: newReflector()}
}

func newReflector() *openapi3.Reflector {
	reflector := openapi3.Reflector{}
	reflector.Spec = &openapi3.Spec{Openapi: "3.0.3"}
	reflector.Spec.Info.
		WithTitle("http-monitoring").
		WithDescription("http-monitoring is a simple http monitoring service")
	return &reflector
}

func (d *DocGenerator) OpenAPISpecAsYaml() ([]byte, error) {
	d.specifySecurity()
	d.specifyOperations()
	return d.reflector.SpecEns().MarshalYAML()
}

func (d *DocGenerator) specifyOperations() {
	d.specifyUsersCreateOperation()
	d.specifyUsersLoginOperation()

	d.specifyUrlsCreateOperation()
	d.specifyUrlsGetAllOperation()
	d.specifyUrlsGetDayStatsOperation()

	d.specifyAlertsGetOperation()
}

func (d *DocGenerator) handleError(err error) {
	if err != nil {
		d.logger.Fatal("error while generating openapi spec", zap.Error(err))
	}
}

func (d *DocGenerator) specifySecurity() {
	d.reflector.SpecEns().ComponentsEns().SecuritySchemesEns().WithMapOfSecuritySchemeOrRefValuesItem(
		securityName,
		openapi3.SecuritySchemeOrRef{
			SecurityScheme: &openapi3.SecurityScheme{
				HTTPSecurityScheme: (&openapi3.HTTPSecurityScheme{}).
					WithScheme("Bearer").
					WithBearerFormat("JWT").
					WithDescription("JWT token for user authentication"),
			},
		},
	)
}
