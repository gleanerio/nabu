package config

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"log"
)

//type EndPoints struct {
//	EndPoints map[string][]EndPoint
//	//EndPoints [][]EndPoints
//}

type EndPoint struct {
	Server       string
	Type         string
	URL          string
	Method       string
	ContentType  string
	Authenticate bool
	Username     string
	Password     string
}

//var EndPointsTemplate = map[string]interface{}{
//	"endpoints": map[string]string{
//		"name":         "",
//		"url":          "",
//		"method":       "",
//		"contentType":  "",
//		"authenticate": "",
//		"username":     "",
//		"password":     "",
//	},
//}

func GetEndPointsConfig(v1 *viper.Viper) ([]EndPoint, error) {
	var subtreeKey = "endpoints"
	var endpointsCfg []EndPoint

	if v1 == nil {
		return nil, fmt.Errorf("GetEndPointsConfig: viperConfig is nil")
	}

	err := v1.UnmarshalKey(subtreeKey, &endpointsCfg)
	if err != nil {
		log.Fatal("error when parsing ", subtreeKey, " config: ", err)
		//No sources, so nothing to run
	}

	return endpointsCfg, err
}

func GetEndpoint(v1 *viper.Viper, set, servertype string) (*EndPoint, error) {

	epcfg, err := GetEndPointsConfig(v1)
	if err != nil {
		log.Fatalf("error getting endpoint node in config %v", err)
	}

	for _, item := range epcfg {
		if item.Server == set && item.Type == servertype {
			return &item, nil // return the item if found
		}
	}

	return nil, errors.New("unable to find the set and or servertype you requested in the config")
}
