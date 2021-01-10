package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/UFOKN/nabu/internal/flows"
	"github.com/UFOKN/nabu/internal/objects"
	"github.com/UFOKN/nabu/internal/semsearch"
	"github.com/UFOKN/nabu/internal/tika"

	//	"../../internal/semsearch"

	"github.com/spf13/viper"
)

var prefixVal, viperVal, modeVal string

// example source s3://noaa-nwm-retro-v2.0-pds/full_physics/2017/201708010001.CHRTOUT_DOMAIN1.comp
func init() {
	log.SetFlags(log.Lshortfile)

	flag.StringVar(&viperVal, "cfg", "config.json", "Configuration file")
	flag.StringVar(&prefixVal, "prefix", "", "Prefix to override config file setting")
	flag.StringVar(&modeVal, "mode", "", "What Nabu should do: tika, txtai, object, prefix")

}

func main() {
	var v1 *viper.Viper
	var err error

	// Set up some logging approaches
	f, err := os.OpenFile("naburun.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.SetFlags(log.Lshortfile)
	// log.SetOutput(ioutil.Discard) // turn off all logging
	//wrt := io.MultiWriter(os.Stdout, f)
	//log.SetOutput(wrt)

	// Parse flags
	flag.Parse()

	if isFlagPassed("cfg") {
		v1, err = readConfig(viperVal, map[string]interface{}{})
		if err != nil {
			panic(fmt.Errorf("Error when reading config: %v", err))
		}
	}

	mc := objects.MinioConnection(v1)

	// Override prefix in config if flag set
	if isFlagPassed("prefix") {
		out := v1.GetStringMapString("objects")
		b := out["bucket"]
		p := prefixVal
		// r := out["region"]
		// v1.Set("objects", map[string]string{"bucket": b, "prefix": NEWPREFIX, "region": r})
		v1.Set("objects", map[string]string{"bucket": b, "prefix": p})
	}

	// Select run mod

	if !isFlagPassed("mode") {
		fmt.Println("Mode must be set -mode one of tika, txtai, object, prefix")
		os.Exit(0)
	}

	switch modeVal {

	case "tika":
		fmt.Println("Tika extract text from objects")
		err = tika.Build(v1, mc)
		if err != nil {
			log.Println(err)
		}

	case "txtai":
		fmt.Println("Index descriptions to txtai")
		err = semsearch.ObjectAssembly(v1, mc)
		if err != nil {
			log.Println(err)
		}

	case "object":
		fmt.Println("Load graph object to triplestore")
		spql := v1.GetStringMapString("sparql")
		s, err := flows.PipeLoad(v1, mc, "bucket", "object", spql["endpoint"])
		if err != nil {
			log.Println(err)
		}
		log.Println(string(s))

	case "prefix":
		fmt.Println("Load graphs in prefix to triplestore")
		err = flows.ObjectAssembly(v1, mc)
		if err != nil {
			log.Println(err)
		}

	}

	// txtai testing
	// err = semsearch.ObjectAssembly(v1, mc)
	// if err != nil {
	// 	log.Println(err)
	// }

	// Graph load an entire prefix
	//err = flows.ObjectAssembly(v1, mc)
	//if err != nil {
	//log.Println(err)
	//}

	//Graph Load a single object
	// spql := v1.GetStringMapString("sparql")
	// s, err := flows.PipeLoad(v1, mc, spql["endpoint"])
	// if err != nil {
	// 	log.Println(err)
	// }
	// log.Println(string(s))
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

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
