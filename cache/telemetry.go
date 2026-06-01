package cache

import (
	"context"
	"encoding/json"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type TrackedTaskEnvelope struct {
	Payload  []byte            `json:"payload"`
	Metadata map[string]string `json:"metadata"` // Holds tracing context fields
}

// PrepareTaskEnvelope bundles business data alongside propagation carriers
func PrepareTaskEnvelope(ctx context.Context, payload any) ([]byte, error) {
	bytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	carrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, carrier)

	envelope := TrackedTaskEnvelope{
		Payload:  bytes,
		Metadata: carrier,
	}

	return json.Marshal(envelope)
}

// ExtractTaskContext parses an inbound redis task envelope back into a valid traced context
func ExtractTaskContext(parentCtx context.Context, envelopeBytes []byte) (context.Context, []byte, error) {
	var envelope TrackedTaskEnvelope
	if err := json.Unmarshal(envelopeBytes, &envelope); err != nil {
		return parentCtx, nil, err
	}

	ctx := otel.GetTextMapPropagator().Extract(parentCtx, propagation.MapCarrier(envelope.Metadata))
	return ctx, envelope.Payload, nil
}
