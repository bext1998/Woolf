package agents

import (
	"os"

	"gopkg.in/yaml.v3"
)

func LoadRole(path string) (Role, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Role{}, err
	}
	var role Role
	if err := yaml.Unmarshal(data, &role); err != nil {
		return Role{}, err
	}
	return role, nil
}
