package plugins

import (
	"gopkg.in/yaml.v3"
)

type Plugin interface {
	Setup(yaml.Node) error

}