package grok

import (
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

//Settings ...
type Settings struct {
	API   *APISettings   `yaml:"api"`
	Mongo *MongoSettings `yaml:"mongo"`
	GCP   *GCPSettings   `yaml:"gcp"`
}

// APISettings ...
type APISettings struct {
	Host    string   `yaml:"host"`
	Swagger string   `yaml:"swagger"`
	Auth    *APIAuth `yaml:"auth"`
}

// MongoSettings ...
type MongoSettings struct {
	ConnectionString string `yaml:"connection_string"`
	Database         string `yaml:"database"`
}

// GCPSettings ...
type GCPSettings struct {
	ProjectID string `yaml:"project_id"`
	PubSub    struct {
		Fake     bool   `yaml:"fake"`
		Endpoint string `yaml:"endpoint"`
	} `yaml:"pubsub"`
}

// APIAuth ...
type APIAuth struct {
	Fake       bool         `yaml:"fake"`
	FakeConfig *FakeAPIAuth `yaml:"fake_config"`
	Tenant     string       `yaml:"tenant"`
	JWKS       string       `yaml:"jwks"`
	Audience   []string     `yaml:"audience"`
}

// FakeAPIAuth ...
type FakeAPIAuth struct {
	Claims        map[string]interface{} `yaml:"claims"`
	Authenticated bool                   `yaml:"authenticated"`
}

// FromYAML ...
func FromYAML(file string, dist interface{}) error {
	filename, _ := filepath.Abs(file)

	data, err := ioutil.ReadFile(filename)

	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, dist)
}
