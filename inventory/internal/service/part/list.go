package part

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/alexander-kartavtsev/starship/inventory/internal/model"
	"github.com/alexander-kartavtsev/starship/platform/pkg/tracing"
)

func (s *service) List(ctx context.Context, filter model.PartsFilter) (map[string]model.Part, error) {
	ctx, span := tracing.StartSpan(ctx, "inventory.repository.List",
		trace.WithAttributes(
			attribute.StringSlice("part_uuids", filter.Uuids),
			attribute.StringSlice("part_names", filter.Names),
			attribute.StringSlice("manufacturer_countries", filter.ManufacturerCountries),
			attribute.StringSlice("manufacturer_names", filter.ManufacturerNames),
			attribute.StringSlice("tags", filter.Tags),
		),
	)
	defer span.End()

	parts, errRep := s.inventoryRepository.List(ctx, filter)
	if errRep != nil {
		span.RecordError(errRep)
		span.End()
		return map[string]model.Part{}, errRep
	}
	span.SetAttributes(
		attribute.Int("count parts in d/b response ", len(parts)),
	)

	return parts, nil
}
