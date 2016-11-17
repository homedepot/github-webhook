package configuration

type TriggerRule map[string]string

type TriggerRun struct {
	Path string   `json:"path," yaml:"path" mapstructure:"Path"`
	Args []string `json:"args,omitempty" yaml:"args,omitempty" mapstructure:"Args"`
	Dir  string   `json:"dir,omitempty" yaml:"dir" mapstructure:"Dir"`
}

type Trigger struct {
	Name string `json:"name," yaml:"name" mapstructure:"Name"`
	Event string `json:"event,omitempty" yaml:"event" mapstructure:"Event"`
	Rules TriggerRule `json:"rules,omitempty" yaml:"rules,omitempty" mapstructure:"Rules"`
	Run TriggerRun `json:"run,omitempty" yaml:"run,omitempty" mapstructure:"Run"`
}

type Configuration struct {
	Triggers []Trigger `json:"triggers,omitempty" yaml:"triggers,omitempty" mapstructure:"Triggers"`
}
