package config

import (
	"fmt"
	"github.com/spf13/viper"
)

// frig frig... do not use lowercase... those are private variables
type ImplNetwork struct {
	Orgname string
}

var implNetworkTemplate = map[string]interface{}{
	"implementation_network": map[string]string{
		"orgname": "indexer",
	},
}

func GetImplNetworkConfig(viperConfig *viper.Viper) (ImplNetwork, error) {
	sub := viperConfig.Sub("implementation_network")
	return readImpleNetworkConfig(sub)
}

// use config.Sub("gleaner)
func readImpleNetworkConfig(implementation_networkSubtress *viper.Viper) (ImplNetwork, error) {
	// Check if the configuration subtree exists
	if implementation_networkSubtress == nil {
		return ImplNetwork{}, fmt.Errorf("implementation_network configuration section is missing")
	}

	var gleanerCfg ImplNetwork

	// Set defaults from template
	for key, value := range implNetworkTemplate {
		implementation_networkSubtress.SetDefault(key, value)
	}

	// Check for required keys before binding
	requiredKeys := []string{"orgname"}
	for _, key := range requiredKeys {
		if !implementation_networkSubtress.IsSet(key) {
			return ImplNetwork{}, fmt.Errorf("required key '%s' is missing in implementation_network configuration", key)
		}
	}

	implementation_networkSubtress.BindEnv("orgname", "IMPLEMENTATION_NETWORK_ORGNAME")
	implementation_networkSubtress.AutomaticEnv()

	err := implementation_networkSubtress.Unmarshal(&gleanerCfg)
	if err != nil {
		return ImplNetwork{}, fmt.Errorf("error when parsing implementation_network config: %v", err)
	}

	// Validate the unmarshaled configuration
	if gleanerCfg.Orgname == "" {
		return ImplNetwork{}, fmt.Errorf("orgname cannot be empty in implementation_network configuration")
	}

	return gleanerCfg, nil
}
