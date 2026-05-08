package agents

type Preset struct {
	Name        string   `yaml:"name"`
	DisplayName string   `yaml:"display_name"`
	Roles       []string `yaml:"roles"`
}
