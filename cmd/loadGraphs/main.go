package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/UFOKN/nabu/internal/flows"
	"github.com/UFOKN/nabu/internal/objects"
	"github.com/gosuri/uiprogress"
	"github.com/spf13/viper"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var dgraphBkts arrayFlags

func init() {
	log.SetFlags(log.Lshortfile)
}

// TODO Make a "save to file" version too?
func main() {
	flag.Var(&dgraphBkts, "bckt", "List of buckets to load into the triple store")
	flag.Parse()

	v1, err := readConfig("config", map[string]interface{}{
		"sqlfile": "",
		"bucket":  "",
		"sparql": map[string]string{
			"address":  "",
			"port":     "", // 3030 for jena,
			"database": "", // doa/update for jena
		},
		"minio": map[string]string{
			"address":   "localhost",
			"port":      "9000",
			"accesskey": "",
			"secretkey": "",
		},
	})
	if err != nil {
		panic(fmt.Errorf("Error when reading config: %v\n", err))
	}

	sc := v1.GetStringMapString("sparql")
	sue := fmt.Sprintf("http://%s:%s/%s", sc["address"], sc["port"], sc["database"])
	log.Println(sue)

	mc := objects.MinioConnection(v1)

	uiprogress.Start() // start bar rendering

	// Load metadata graphs
	for b := range dgraphBkts {
		log.Printf("Getting objects to load %s \n", dgraphBkts[b])
		err := flows.LoadGraph(mc, dgraphBkts[b], sue)
		if err != nil {
			log.Printf("Error in loadgraph flow for %s \n", dgraphBkts[b])
		}
	}

	uiprogress.Stop() // stop rendering
}

func readConfig(filename string, defaults map[string]interface{}) (*viper.Viper, error) {
	v := viper.New()
	for key, value := range defaults {
		v.SetDefault(key, value)
	}
	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.AutomaticEnv()
	err := v.ReadInConfig()
	return v, err
}
