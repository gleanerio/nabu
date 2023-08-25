package config

import (
	"fmt"
	"github.com/spf13/viper"
)

// frig frig... do not use lowercase... those are private variables
type Gleaner struct {
	Runid  string
	Summon string
	Mill   string
}

var GleanerTemplate = map[string]interface{}{
	"gleaner": map[string]string{
		"runid":  "indexer",
		"summon": "true",
		"mill":   "flase",
	},
}

func GetGleanerConfig(viperConfig *viper.Viper) (Gleaner, error) {
	sub := viperConfig.Sub("gleaner")
	return ReadGleanerConfig(sub)
}

// use config.Sub("gleaner)
func ReadGleanerConfig(gleanerSubtress *viper.Viper) (Gleaner, error) {
	var gleanerCfg Gleaner
	for key, value := range GleanerTemplate {
		gleanerSubtress.SetDefault(key, value)
	}
	gleanerSubtress.BindEnv("runid", "GLEANER_RUNID")
	gleanerSubtress.BindEnv("summon", "GLEANER_SUMMON")
	gleanerSubtress.BindEnv("mill", "GLEANER_MILL")

	gleanerSubtress.AutomaticEnv()
	// config already read. substree passed
	err := gleanerSubtress.Unmarshal(&gleanerCfg)
	if err != nil {
		panic(fmt.Errorf("error when parsing gleaner config: %v", err))
	}
	return gleanerCfg, err
}
