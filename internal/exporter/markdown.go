package exporter

import (
	"fmt"
	"strings"

	"woolf/internal/session"
)

type MarkdownExporter struct{}

func (MarkdownExporter) Export(sess session.Session) ([]byte, error) {
	var b strings.Builder
	title := sess.Title
	if title == "" {
		title = sess.SessionID
	}
	fmt.Fprintf(&b, "# %s\n\n", title)
	fmt.Fprintf(&b, "- session_id: %s\n", sess.SessionID)
	fmt.Fprintf(&b, "- status: %s\n", sess.Status)
	fmt.Fprintf(&b, "- rounds_completed: %d\n", sess.Totals.RoundsCompleted)
	fmt.Fprintf(&b, "- total_tokens: %d\n", sess.Totals.TotalTokens)
	fmt.Fprintf(&b, "- total_cost_usd: %.6f\n\n", sess.Totals.TotalCostUSD)
	if sess.Source != nil {
		fmt.Fprintf(&b, "## Source\n\n")
		fmt.Fprintf(&b, "- type: %s\n", sess.Source.Type)
		if sess.Source.Path != "" {
			fmt.Fprintf(&b, "- path: %s\n", sess.Source.Path)
		}
		if sess.Source.ContentHash != "" {
			fmt.Fprintf(&b, "- hash: %s\n", sess.Source.ContentHash)
		}
		b.WriteString("\n")
	}
	if len(sess.AgentsConfig) > 0 {
		b.WriteString("## Agents\n\n")
		for _, agent := range sess.AgentsConfig {
			fmt.Fprintf(&b, "- %s (%s): %s\n", agent.DisplayName, agent.Name, agent.Model)
		}
		b.WriteString("\n")
	}
	for _, round := range sess.Rounds {
		fmt.Fprintf(&b, "## Round %d\n\n", round.RoundIndex)
		for _, response := range round.Responses {
			fmt.Fprintf(&b, "### %s\n\n", response.AgentName)
			if response.StanceTag != nil && *response.StanceTag != "" {
				fmt.Fprintf(&b, "- stance: %s\n", *response.StanceTag)
			}
			fmt.Fprintf(&b, "- status: %s\n", response.Status)
			fmt.Fprintf(&b, "- model: %s\n", response.Model)
			fmt.Fprintf(&b, "- tokens: prompt=%d completion=%d\n\n", response.Tokens.Prompt, response.Tokens.Completion)
			b.WriteString(strings.TrimSpace(response.Content))
			b.WriteString("\n\n")
		}
	}
	if len(sess.Summaries) > 0 {
		b.WriteString("## Summaries\n\n")
		for key, value := range sess.Summaries {
			fmt.Fprintf(&b, "### %s\n\n%s\n\n", key, strings.TrimSpace(value))
		}
	}
	return []byte(b.String()), nil
}
