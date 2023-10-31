package configs

import (
	"encoding/base64"

	"gopkg.in/yaml.v3"
)

type Secret string

func (s *Secret) UnmarshalYAML(v *yaml.Node) error {
	value, err := base64.StdEncoding.DecodeString(v.Value)
	if err != nil {
		return err
	}

	*s = Secret(value)
	return nil
}
