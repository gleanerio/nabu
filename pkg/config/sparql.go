package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Sparql struct {
	Endpoint     string
	EndpointBulk string
	Authenticate bool
	Username     string
	Password     string
}

var sparqlTemplate = map[string]interface{}{
	"sparql": map[string]string{
		"endpoint":     "http://coreos.lan:3030/testing/sparql",
		"endpointBulk": "http://coreos.lan:3030/testing/data",
		"authenticate": "False",
		"username":     "",
		"password":     "",
	},
}

func GetSparqlConfig(viperConfig *viper.Viper) (Sparql, error) {
	sub := viperConfig.Sub("sparql")
	return ReadSparqlConfig(sub)
}

func ReadSparqlConfig(viperSubtree *viper.Viper) (Sparql, error) {
	var sparql Sparql
	for key, value := range sparqlTemplate {
		viperSubtree.SetDefault(key, value)
	}
	viperSubtree.BindEnv("endpoint", "SPARQL_ENDPOINT")
	viperSubtree.BindEnv("endpointBulk", "SPARQL_ENDPOINTBULK")
	viperSubtree.BindEnv("authenticate", "SPARQL_AUTHENTICATE")
	viperSubtree.BindEnv("username", "SPARQL_USERNAME")
	viperSubtree.BindEnv("password", "SPARQL_PASSWORD")

	viperSubtree.AutomaticEnv()
	// config already read. substree passed
	err := viperSubtree.Unmarshal(&sparql)
	if err != nil {
		panic(fmt.Errorf("error when parsing sparql endpoint config: %v", err))
	}
	return sparql, err
}
