package main

type NamespaceWatcherConfig struct {
	Defaults struct {
		Mapping struct {
			Cluster string `json:"cluster,omitempty"`
			Project string `json:"project,omitempty"`
		} `json:"mapping"`
		Creation string `json:"creation,omitempty"` // default allow
	} `json:"defaults"`
	Rules struct {
		Mapping  []NamespaceMatchConfig
		Creation []NamespaceCreationConfig
	} `json:"rules"`
}

type NamespaceMatchConfig struct {
	Regex   string `json:"regex"`
	Cluster string `json:"cluster"`
	Project string `json:"project"`
}

type NamespaceCreationConfig struct {
	Regex  string `json:"regex"`
	Action string `json:"action"`
}
