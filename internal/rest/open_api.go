package rest

import (
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ghodss/yaml"
	"github.com/gorilla/mux"
)

//go:generate go run ../../cmd/openapi-gen/main.go -path .
//go:generate oapi-codegen -package openapi3 -generate types  -o ../../pkg/openapi3/URL_types.gen.go openapi3.yaml
//go:generate oapi-codegen -package openapi3 -generate client -o ../../pkg/openapi3/client.gen.go     openapi3.yaml

// NewOpenAPI3 instantiates the OpenAPI specification for this service.
func NewOpenAPI3() openapi3.T {
	swagger := openapi3.T{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Title:       "URL API",
			Description: "REST APIs used for interacting with the URL Service",
			Version:     "0.0.0",
			License: &openapi3.License{
				Name: "MIT",
				URL:  "https://opensource.org/licenses/MIT",
			},
			Contact: &openapi3.Contact{
				URL: "https://github.com/Oguzyildirim/url-info",
			},
		},
		Servers: openapi3.Servers{
			&openapi3.Server{
				Description: "Local development",
				URL:         "http://127.0.0.1:9234",
			},
		},
	}

	swagger.Components.Schemas = openapi3.Schemas{
		"URL": openapi3.NewSchemaRef("",
			openapi3.NewObjectSchema().
				WithProperty("id", openapi3.NewUUIDSchema()).
				WithProperty("HTMLVersion", openapi3.NewStringSchema()).
				WithProperty("headingsCount", openapi3.NewStringSchema()).
				WithProperty("pageTitle", openapi3.NewStringSchema()).
				WithProperty("linksCount", openapi3.NewInt32Schema()).
				WithProperty("inaccessibleLinksCount", openapi3.NewInt32Schema()).
				WithProperty("HaveLoginForm", openapi3.NewBoolSchema())),
	}

	swagger.Components.RequestBodies = openapi3.RequestBodies{
		"SearchURLsRequest": &openapi3.RequestBodyRef{
			Value: openapi3.NewRequestBody().
				WithDescription("Request used for creating a URL info.").
				WithRequired(true).
				WithJSONSchema(openapi3.NewSchema().
					WithProperty("URL", openapi3.NewStringSchema().
						WithMinLength(10)),
				),
		},
	}

	swagger.Components.Responses = openapi3.Responses{
		"ErrorResponse": &openapi3.ResponseRef{
			Value: openapi3.NewResponse().
				WithDescription("Response when errors happen.").
				WithContent(openapi3.NewContentWithJSONSchema(openapi3.NewSchema().
					WithProperty("error", openapi3.NewStringSchema()))),
		},
		"SearchURLsResponse": &openapi3.ResponseRef{
			Value: openapi3.NewResponse().
				WithDescription("Response returned back after creating URLs.").
				WithContent(openapi3.NewContentWithJSONSchema(openapi3.NewSchema().
					WithPropertyRef("URL", &openapi3.SchemaRef{
						Ref: "#/components/schemas/URL",
					}))),
		},
		"ReadURLsResponse": &openapi3.ResponseRef{
			Value: openapi3.NewResponse().
				WithDescription("Response returned back after searching one URL.").
				WithContent(openapi3.NewContentWithJSONSchema(openapi3.NewSchema().
					WithPropertyRef("URL", &openapi3.SchemaRef{
						Ref: "#/components/schemas/URL",
					}))),
		},
		"ReadURLsByCountryResponse": &openapi3.ResponseRef{
			Value: openapi3.NewResponse().
				WithDescription("Response returned back after searching URLs by country.").
				WithContent(openapi3.NewContentWithJSONSchema(openapi3.NewSchema().
					WithPropertyRef("URL", &openapi3.SchemaRef{
						Ref: "#/components/schemas/URL",
					}))),
		},
	}

	swagger.Paths = openapi3.Paths{
		"/URLs": &openapi3.PathItem{
			Post: &openapi3.Operation{
				OperationID: "CreateURL",
				RequestBody: &openapi3.RequestBodyRef{
					Ref: "#/components/requestBodies/SearchURLsRequest",
				},
				Responses: openapi3.Responses{
					"400": &openapi3.ResponseRef{
						Ref: "#/components/responses/ErrorResponse",
					},
					"500": &openapi3.ResponseRef{
						Ref: "#/components/responses/ErrorResponse",
					},
					"201": &openapi3.ResponseRef{
						Ref: "#/components/responses/SearchURLsResponse",
					},
				},
			},
		},
		"/URLs/{URLId}": &openapi3.PathItem{
			Delete: &openapi3.Operation{
				OperationID: "DeleteURL",
				Parameters: []*openapi3.ParameterRef{
					{
						Value: openapi3.NewPathParameter("URLId").
							WithSchema(openapi3.NewUUIDSchema()),
					},
				},
				Responses: openapi3.Responses{
					"200": &openapi3.ResponseRef{
						Value: openapi3.NewResponse().WithDescription("URL updated"),
					},
					"404": &openapi3.ResponseRef{
						Value: openapi3.NewResponse().WithDescription("URL not found"),
					},
					"500": &openapi3.ResponseRef{
						Ref: "#/components/responses/ErrorResponse",
					},
				},
			},
			Get: &openapi3.Operation{
				OperationID: "ReadURL",
				Parameters: []*openapi3.ParameterRef{
					{
						Value: openapi3.NewPathParameter("URLId").
							WithSchema(openapi3.NewUUIDSchema()),
					},
				},
				Responses: openapi3.Responses{
					"200": &openapi3.ResponseRef{
						Ref: "#/components/responses/ReadURLsResponse",
					},
					"404": &openapi3.ResponseRef{
						Value: openapi3.NewResponse().WithDescription("URL not found"),
					},
					"500": &openapi3.ResponseRef{
						Ref: "#/components/responses/ErrorResponse",
					},
				},
			},
		},
	}

	return swagger
}

func RegisterOpenAPI(r *mux.Router) {
	swagger := NewOpenAPI3()

	r.HandleFunc("/openapi3.json", func(w http.ResponseWriter, r *http.Request) {
		renderResponse(w, &swagger, http.StatusOK)
	}).Methods(http.MethodGet)

	r.HandleFunc("/openapi3.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-yaml")

		data, _ := yaml.Marshal(&swagger)

		_, _ = w.Write(data)

		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)
}
