package authorizer

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"

	"github.com/open-policy-agent/opa/v1/rego"
	"github.com/open-policy-agent/opa/v1/storage/inmem"
)

//go:embed policy.rego data.json
var fs embed.FS

type Authorizer interface {
	Evaluate(ctx context.Context, input map[string]any) (bool, error)
}

type Engine struct {
	query rego.PreparedEvalQuery
}

func NewEmbedded() (Authorizer, error) {
	ctx := context.Background()

	policySrc, err := fs.ReadFile("policy.rego")
	if err != nil {
		return nil, fmt.Errorf("read embedded policy: %w", err)
	}

	dataBytes, err := fs.ReadFile("data.json")
	if err != nil {
		return nil, fmt.Errorf("read embedded data.json: %w", err)
	}

	var data map[string]any
	if err := json.Unmarshal(dataBytes, &data); err != nil {
		return nil, fmt.Errorf("parse data.json: %w", err)
	}

	r := rego.New(
		rego.Query("data.authorizer.allow"),
		rego.Module("policy.rego", string(policySrc)),
		rego.Store(inmem.NewFromObject(data)),
	)

	query, err := r.PrepareForEval(ctx)
	if err != nil {
		return nil, fmt.Errorf("compile policy: %w", err)
	}

	return &Engine{query: query}, nil
}

func (e *Engine) Evaluate(ctx context.Context, input map[string]any) (bool, error) {
	results, err := e.query.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		return false, fmt.Errorf("eval: %w", err)
	}

	if len(results) == 0 || len(results[0].Expressions) == 0 {
		return false, nil
	}

	allowed, ok := results[0].Expressions[0].Value.(bool)
	return ok && allowed, nil
}
