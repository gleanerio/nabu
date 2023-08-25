package graph

import (
	"bytes"
	"fmt"
	"github.com/gleanerio/nabu/pkg/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
)

// CLEAR removes ALL graphs
func Clear(v1 *viper.Viper) ([]byte, error) {
	//spql := v1.GetStringMapString("sparql")
	spql, _ := config.GetSparqlConfig(v1)
	d := fmt.Sprint("CLEAR ALL")

	pab := []byte(d)

	//req, err := http.NewRequest("POST", spql["endpoint"], bytes.NewBuffer(pab))
	req, err := http.NewRequest("POST", spql.EndpointBulk, bytes.NewBuffer(pab))
	if err != nil {
		log.Error(err)
	}
	req.Header.Set("Content-Type", "application/sparql-update")
	// req.Header.Set("Content-Type", "application/sparql-results+xml")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("response Body:", string(body))
		log.Error("response Status:", resp.Status)
		log.Error("response Headers:", resp.Header)
	}

	log.Trace(string(body))

	log.Println("ALERT graph was cleared")

	return body, err
}
