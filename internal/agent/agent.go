package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jonwraymond/tooldiscovery/discovery"
	"github.com/jonwraymond/tooldiscovery/tooldoc"
	"github.com/jonwraymond/toolexec/run"
	"github.com/jonwraymond/toolfoundation/adapter"
	"github.com/jonwraymond/toolfoundation/model"
	"github.com/jonwraymond/toolprotocol/a2a"
	"github.com/jonwraymond/toolprotocol/wire"
)

// Agent provides A2A agent operations backed by ApertureStack discovery and execution.
type Agent struct {
	Name             string
	Description      string
	Version          string
	DocumentationURL string
	IconURL          string
	BaseURL          string

	Discovery *discovery.Discovery
	Runner    run.Runner
	MaxSkills int
}

// AgentCard returns the A2A agent card representation.
func (a *Agent) AgentCard(ctx context.Context) (any, error) {
	if a.Discovery == nil {
		return nil, fmt.Errorf("discovery not configured")
	}
	if a.Name == "" || a.Description == "" || a.Version == "" {
		return nil, fmt.Errorf("agent name, description, and version are required")
	}

	tools, err := a.canonicalTools(ctx)
	if err != nil {
		return nil, err
	}

	provider := adapter.CanonicalProvider{
		Name:        a.Name,
		Description: a.Description,
		Version:     a.Version,
		Capabilities: map[string]any{
			"streaming": true,
			"tasks":     true,
		},
		Skills: tools,
		SourceMeta: map[string]any{
			"supportedInterfaces": []adapter.A2AAgentInterface{
				{
					URL:             a.BaseURL,
					ProtocolBinding: "jsonrpc",
					ProtocolVersion: wire.A2AVersion,
				},
			},
			"documentationUrl": a.DocumentationURL,
			"iconUrl":          a.IconURL,
		},
	}

	return adapter.NewA2AAdapter().FromCanonicalProvider(&provider)
}

// ListSkills returns the agent's skills in A2A tool format.
func (a *Agent) ListSkills(ctx context.Context) ([]wire.Tool, error) {
	if a.Discovery == nil {
		return nil, fmt.Errorf("discovery not configured")
	}
	results, err := a.Discovery.Search(ctx, "", a.skillLimit())
	if err != nil {
		return nil, err
	}

	tools := make([]wire.Tool, 0, len(results))
	for _, res := range results {
		doc, _ := a.Discovery.DescribeTool(res.Summary.ID, tooldoc.DetailFull)
		tools = append(tools, wire.Tool{
			Name:        res.Summary.ID,
			Description: chooseDescription(res.Summary.Summary, res.Summary.ShortDescription),
			InputSchema: normalizeSchema(doc.Tool),
		})
	}

	return tools, nil
}

// Invoke runs a skill via the runner.
func (a *Agent) Invoke(ctx context.Context, skillID string, args map[string]any) (a2a.InvokeResult, error) {
	if a.Runner == nil {
		return a2a.InvokeResult{}, fmt.Errorf("runner not configured")
	}
	result, err := a.Runner.Run(ctx, skillID, args)
	if err != nil {
		return a2a.InvokeResult{}, err
	}

	content := []wire.Content{
		{
			Type: wire.ContentTypeText,
			Text: stringify(result.Structured),
		},
	}
	return a2a.InvokeResult{Content: content}, nil
}

func (a *Agent) canonicalTools(ctx context.Context) ([]adapter.CanonicalTool, error) {
	results, err := a.Discovery.Search(ctx, "", a.skillLimit())
	if err != nil {
		return nil, err
	}
	tools := make([]adapter.CanonicalTool, 0, len(results))
	for _, res := range results {
		desc := chooseDescription(res.Summary.Summary, res.Summary.ShortDescription)
		tools = append(tools, adapter.CanonicalTool{
			Namespace:   res.Summary.Namespace,
			Name:        res.Summary.Name,
			DisplayName: res.Summary.Name,
			Summary:     res.Summary.Summary,
			Description: desc,
			Category:    res.Summary.Category,
			Tags:        res.Summary.Tags,
			InputModes:  res.Summary.InputModes,
			OutputModes: res.Summary.OutputModes,
		})
	}
	return tools, nil
}

func (a *Agent) skillLimit() int {
	if a.MaxSkills <= 0 {
		return 500
	}
	return a.MaxSkills
}

func normalizeSchema(tool *model.Tool) map[string]any {
	if tool == nil {
		return nil
	}
	if tool.InputSchema == nil {
		return nil
	}
	if m, ok := tool.InputSchema.(map[string]any); ok {
		return m
	}
	data, err := json.Marshal(tool.InputSchema)
	if err != nil {
		return nil
	}
	var out map[string]any
	if err := json.Unmarshal(data, &out); err != nil {
		return nil
	}
	return out
}

func chooseDescription(summary, short string) string {
	if summary != "" {
		return summary
	}
	return short
}

func stringify(v any) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("%v", v)
	}
	return string(data)
}
