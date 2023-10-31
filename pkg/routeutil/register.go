package routeutil

import (
	"fmt"
	"net/http"
	"path"
	"reflect"
	"sync"

	"github.com/gin-gonic/gin"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/apperror"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/swag"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/swag/endpoint"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/swag/swagger"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/transportutil"
	"mmlabel.gitlab.com/mm-printing-backend/version"
)

var api = swag.New(
	swag.Title("mm-printing backend"),
	swag.Version(version.Version),
	swag.ContactEmail("support@mm-printing.com"),
	swag.BasePath("/"),
	swag.Schemes("http", "https"),
	swag.SecurityScheme("Authorization", swagger.APIKeySecurity("Authorization", "header")),
	swag.SecurityScheme("DeviceID", swagger.APIKeySecurity("DeviceID", "header")),
	swag.Description(fmt.Sprintf("GoVersion: %s <br/> GitHash: %s", version.GoVersion, version.GitHash)+apperror.ExposeDocs()),
)

var once sync.Once

func ServingDocs(c *gin.Context) {
	api.Host = c.Request.Host
	api.Tags = []swagger.Tag{}

	once.Do(func() {
		for k, d := range api.Definitions {
			for k2, p := range d.Properties {
				pv := reflect.New(p.GoType)
				switch v := pv.Interface().(type) {
				case enum.Enum:
					p.Enum = v.EnumDescriptions()
					p.Format = "string"

					api.Definitions[k].Properties[k2] = p
				}
			}
		}
	})

	c.JSON(http.StatusOK, api)
}

type RegisterOption int

const (
	RegisterOptionSkipAuth RegisterOption = iota + 1
	RegisterOptionDeprecated
)

func AddEndpoint(
	g *gin.RouterGroup,
	relativePath string,
	handler gin.HandlerFunc,
	req, resp interface{},
	description string,
	opts ...RegisterOption,
) {
	path := joinPaths(g.BasePath(), relativePath)

	swaggerOpts := []endpoint.Option{
		endpoint.Tags(g.BasePath()),
		endpoint.Security("Authorization", []string{}...),
		endpoint.Body(req, "Request Body", true),
		endpoint.Response(http.StatusOK, resp, "Response"),
	}

	for _, o := range opts {
		switch o {
		case RegisterOptionSkipAuth:
			transportutil.RegisterPublicEndpoint(path)
		case RegisterOptionDeprecated:
			swaggerOpts = append(swaggerOpts, endpoint.Deprecated(true))
		}
	}

	api.AddEndpoint(endpoint.New(
		http.MethodPost,
		path,
		description,
		swaggerOpts...,
	))

	g.POST(relativePath, handler)
}

func AddEndpointGet(
	g *gin.RouterGroup,
	relativePath string,
	handler gin.HandlerFunc,
	req, resp interface{},
	description string,
	opts ...RegisterOption,
) {
	path := joinPaths(g.BasePath(), relativePath)

	swaggerOpts := []endpoint.Option{
		endpoint.Tags(g.BasePath()),
		endpoint.Security("Authorization", []string{}...),
		endpoint.Response(http.StatusOK, resp, "Response"),
	}

	for _, o := range opts {
		switch o {
		case RegisterOptionSkipAuth:
			transportutil.RegisterPublicEndpoint(path)
		case RegisterOptionDeprecated:
			swaggerOpts = append(swaggerOpts, endpoint.Deprecated(true))
		}
	}

	api.AddEndpoint(endpoint.New(
		http.MethodGet,
		path,
		description,
		swaggerOpts...,
	))

	g.GET(relativePath, handler)
}

func AddCustomEndpoint(
	g *gin.RouterGroup,
	relativePath string,
	handler gin.HandlerFunc,
	description string,
	opts ...endpoint.Option,
) {
	opts = append(opts, endpoint.Security("Authorization", []string{}...))
	path := joinPaths(g.BasePath(), relativePath)
	api.AddEndpoint(endpoint.New(
		http.MethodPost,
		path,
		description,
		opts...,
	))

	g.POST(relativePath, handler)
}

func joinPaths(absolutePath, relativePath string) string {
	if relativePath == "" {
		return absolutePath
	}

	finalPath := path.Join(absolutePath, relativePath)
	appendSlash := lastChar(relativePath) == '/' && lastChar(finalPath) != '/'
	if appendSlash {
		return finalPath + "/"
	}
	return finalPath
}

func lastChar(str string) uint8 {
	if str == "" {
		panic("The length of the string can't be 0")
	}
	return str[len(str)-1]
}

func Parameter(p swagger.Parameter) endpoint.Option {
	return func(b *endpoint.Builder) {
		if b.Endpoint.Parameters == nil {
			b.Endpoint.Parameters = []swagger.Parameter{}
		}

		b.Endpoint.Parameters = append(b.Endpoint.Parameters, p)
	}
}
