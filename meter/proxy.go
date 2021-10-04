package meter

import (
	"fmt"
	"strings"

	"github.com/evcc-io/evcc/api"
	"github.com/evcc-io/evcc/templates"
	"gopkg.in/yaml.v3"
)

func init() {
	for _, tmpl := range templates.ByClass(templates.Meter) {
		println(strings.ToUpper(tmpl.Type))
		println("")

		// render the proxy
		sample, err := tmpl.RenderProxy()
		if err != nil {
			panic(err)
		}

		println("-- proxy --")
		println(string(sample))
		println("")

		instantiateFunc := instantiateFunction(tmpl)
		registry.Add(tmpl.Type, instantiateFunc)

		// render all usages
		for _, usage := range tmpl.Usages() {
			println("--", usage, "--")

			b, err := tmpl.RenderResult(map[string]interface{}{
				"usage": usage,
			})
			if err != nil {
				panic(err)
			}

			println(string(b))
			println("")
		}
	}
}

func instantiateFunction(tmpl templates.Template) func(map[string]interface{}) (api.Meter, error) {
	return func(other map[string]interface{}) (api.Meter, error) {
		b, err := tmpl.RenderResult(other)
		if err != nil {
			return nil, err
		}

		fmt.Println("-- instantiated --")
		println(string(b))
		println("")

		var instance struct {
			Type  string
			Other map[string]interface{} `yaml:",inline"`
		}

		if err := yaml.Unmarshal(b, &instance); err != nil {
			return nil, err
		}

		return NewFromConfig(instance.Type, instance.Other)
	}
}