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
	var gleanerCfg ImplNetwork
	for key, value := range implNetworkTemplate {
		implementation_networkSubtress.SetDefault(key, value)
	}
	implementation_networkSubtress.BindEnv("orgname", "IMPLEMENTATION_NETWORK_ORGNAME")

	implementation_networkSubtress.AutomaticEnv()
	// config already read. substree passed
	err := implementation_networkSubtress.Unmarshal(&gleanerCfg)
	if err != nil {
		panic(fmt.Errorf("error when parsing gleaner config: %v", err))
	}
	return gleanerCfg, err
}
