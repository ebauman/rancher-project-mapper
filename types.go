package main

type NamespaceWatcherConfig struct {
	Rules []NamespaceMatchConfig `json:"rules"`
}

type NamespaceMatchConfig struct {
	Regex   string `json:"regex"`
	Cluster string `json:"cluster"`
	Project string `json:"project"`
	Default bool   `json:"default,omitempty"`
}
