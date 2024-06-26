package graph

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	opentracingLog "github.com/opentracing/opentracing-go/log"
)

// Originally taken from public domain with an Unlicense and modified: https://github.com/99designs/gqlgen-contrib/blob/master/gqlopentracing/tracer.go
// This is a gqlgen extension that enables tracing through opentracing for each overall operation that gqlgen encounters
// as well as each field that is accessed

type tracer struct{}

type GraphQLTracer interface {
	graphql.HandlerExtension
	graphql.FieldInterceptor
	graphql.ResponseInterceptor
}

// Creates a GraphQL tracing extension that can be used with a gqlgen server like this (you will need to prepend NewGraphQLTracer
// with whatever package prefix you assign it in the import statement)
//
//		graphqlSrv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
//				Resolvers: &graph.Resolver{},
//		}))
//	    graphqlSrv.Use(NewGraphQLTracer())
func NewGraphQLTracer() GraphQLTracer {
	return tracer{}
}

// HandlerExtension: sets up basic metadata for the extension

func (t tracer) ExtensionName() string {
	return "OpenTracing"
}

func (t tracer) Validate(schema graphql.ExecutableSchema) error {
	return nil
}

func (t tracer) InterceptResponse(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
	span, ctx := opentracing.StartSpanFromContext(ctx, "gqlgen-resolver")
	oc := graphql.GetOperationContext(ctx)
	ext.SpanKind.Set(span, "server")
	ext.Component.Set(span, "gqlgen")
	span.SetTag("query", oc.RawQuery)
	defer span.Finish()

	res := next(ctx)

	errList := graphql.GetErrors(ctx)
	if len(errList) > 0 {
		ext.Error.Set(span, true)
		span.LogFields(
			opentracingLog.String("event", "error"),
			opentracingLog.String("error.message", errList.Error()),
		)
	}

	return res
}

// FieldIntercepter: intercepts each individual GraphQL field before it is executed

func (t tracer) InterceptField(ctx context.Context, next graphql.Resolver) (interface{}, error) {
	fc := graphql.GetFieldContext(ctx)

	span, ctx := opentracing.StartSpanFromContext(ctx, fc.Object+"_"+fc.Field.Name)
	span.SetTag("resolver.object", fc.Object)
	span.SetTag("resolver.field", fc.Field.Name)
	defer span.Finish()

	res, err := next(ctx)

	errList := graphql.GetFieldErrors(ctx, fc)
	if len(errList) > 0 {
		ext.Error.Set(span, true)
		span.LogFields(
			opentracingLog.String("event", "error"),
		)

		for i, err := range errList {
			span.LogFields(
				opentracingLog.String(fmt.Sprintf("error.%d.message", i), err.Error()),
				opentracingLog.String(fmt.Sprintf("error.%d.kind", i), fmt.Sprintf("%T", err)),
			)
		}
	}

	if err != nil {
		ext.Error.Set(span, true)
		span.LogFields(
			opentracingLog.String("event", "error"),
			opentracingLog.String("error.message", err.Error()),
		)
	}

	return res, err
}
