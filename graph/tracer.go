package graph

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	opentracingLog "github.com/opentracing/opentracing-go/log"
)

// Originally taken from public domain with an Unlicense: https://github.com/99designs/gqlgen-contrib/blob/master/gqlopentracing/tracer.go

type Tracer struct{}

var _ interface {
	graphql.HandlerExtension
	graphql.OperationInterceptor
	graphql.FieldInterceptor
} = &Tracer{}

// HandlerExtension: sets up basic metadata for the extension

func (t Tracer) ExtensionName() string {
	return "OpenTracing"
}

func (t Tracer) Validate(schema graphql.ExecutableSchema) error {
	return nil
}

// OperationIntercepter: intercepts the entire query operation before it hits gqlgen

func (t Tracer) InterceptOperation(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
	oc := graphql.GetOperationContext(ctx)

	span, ctx := opentracing.StartSpanFromContext(ctx, oc.RawQuery)
	ext.SpanKind.Set(span, "server")
	ext.Component.Set(span, "gqlgen")
	defer span.Finish()

	return next(ctx)
}

// FieldIntercepter: intercepts each individual GraphQL field before it is executed

func (t Tracer) InterceptField(ctx context.Context, next graphql.Resolver) (interface{}, error) {
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
