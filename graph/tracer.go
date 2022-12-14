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
	graphql.OperationInterceptor
	graphql.FieldInterceptor
}

// Creates a GraphQL tracing extension that can be used with a gqlgen server like this:
//
//		graphqlSrv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
//				Resolvers: &graph.Resolver{},
//		}))
//	    graphqlSrv.Use(graph.NewGraphQLTracer())
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

// OperationIntercepter: intercepts the entire query operation before it hits gqlgen

func (t tracer) InterceptOperation(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
	oc := graphql.GetOperationContext(ctx)

	span, ctx := opentracing.StartSpanFromContext(ctx, oc.RawQuery)
	ext.SpanKind.Set(span, "server")
	ext.Component.Set(span, "gqlgen")
	defer span.Finish()

	return next(ctx)
}

// FieldIntercepter: intercepts each individual GraphQL field before it is executed

func (t tracer) InterceptField(ctx context.Context, next graphql.Resolver) (interface{}, error) {
	fc := graphql.GetFieldContext(ctx)

	span := opentracing.SpanFromContext(ctx)
	span.SetOperationName(fc.Object + "_" + fc.Field.Name)
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

	return res, err
}
