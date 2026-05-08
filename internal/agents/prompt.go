package agents

func BuildPrompt(role Role, draft string) string {
	return role.SystemPrompt + "\n\n" + draft
}
