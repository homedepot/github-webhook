package configuration

import (
	"github.com/mitchellh/mapstructure"
	"fmt"
	"strings"
	"gopkg.in/yaml.v2"
	"github.com/concourse/atc"
	"github.com/homedepot/github-webhook/events"
	"io/ioutil"
)

/*

triggers:
  - name: master
  	event: push
    rules:
   		ref: refs/heads/master
      	head_commit.author.username: username
      	pusher.name: username
      	commits.modified.README.md
  	run:
		path: sh
		args:
		- -exc
		- |
		  env
		  ls -l

 */

func LoadFile(filePath string) (Configuration, error) {

	if len(filePath) == 0 {
		return Configuration{}, nil
	}

	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return Configuration{}, err
	}

	config, err := LoadBytes(bytes)
	if err != nil {
		return Configuration{}, fmt.Errorf("failed to load %s: %s", filePath, err)
	}

	return config, nil
}

func LoadBytes(configBytes []byte) (Configuration, error) {
	var untypedInput map[string]interface{}

	if err := yaml.Unmarshal(configBytes, &untypedInput); err != nil {
		return Configuration{}, err
	}

	var config Configuration
	var metadata mapstructure.Metadata

	msConfig := &mapstructure.DecoderConfig{
		Metadata:         &metadata,
		Result:           &config,
		WeaklyTypedInput: true,
		DecodeHook:       atc.SanitizeDecodeHook,
	}

	decoder, err := mapstructure.NewDecoder(msConfig)
	if err != nil {
		return Configuration{}, err
	}

	if err := decoder.Decode(untypedInput); err != nil {
		return Configuration{}, err
	}

	if len(metadata.Unused) > 0 {
		keys := strings.Join(metadata.Unused, ", ")
		return Configuration{}, fmt.Errorf("extra keys in the task configuration: %s", keys)
	}

	err = Validate(config)
	if err != nil {
		return Configuration{}, err
	}

	return config, nil
}


func Validate(config Configuration) error {
	messages := []string{}

	messages = append(messages, validateTriggers(config)...)

	if len(messages) > 0 {
		return fmt.Errorf("invalid configuration:\n%s", strings.Join(messages, "\n"))
	}

	return nil
}

func validateTriggers(config Configuration) []string {
	messages := []string{}

	if len(config.Triggers) > 0 {
		for _, trigger := range config.Triggers {
			messages = append(messages, validateTrigger(trigger)...)
		}
	}

	return messages
}

func validateTrigger(trigger Trigger) []string {
	messages := []string{}

	if len(trigger.Name) == 0 {
		messages = append(messages, "  missing trigger name")
	}

	if len(trigger.Event) == 0 {
		messages = append(messages, "  missing github event type")
	}

	_, found := events.EventTypes[trigger.Event]
	if !found {
		messages = append(messages, "  unsupported github event type")
	}


	if len(trigger.Run.Path) == 0 {
		messages = append(messages, "  missing trigger run path")
	}

	messages = append(messages, validateRule(trigger.Rules)...)

	return messages
}

func validateRule(rule TriggerRule) []string {
	messages := []string{}
	for key, match := range rule {
		if len(key) == 0 {
			messages = append(messages, "  missing rule key")
		}
		if len(match) == 0 {
			messages = append(messages, "  missing rule match")
		}
	}

	return messages
}

