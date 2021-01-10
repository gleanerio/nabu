package triplestore

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/UFOKN/nabu/internal/graph"
)

//BlazeUpdateNQ updates the blaze triple store
func BlazeUpdateNQ(s []byte, sue string) ([]byte, error) {
	nt, g, err := graph.NQToNTCtx(string(s))
	if err != nil {
		log.Printf("nqToNTCtx err: %s triple: %s", err, string(s))
	}

	p := "INSERT DATA { "
	pab := []byte(p)
	gab := []byte(fmt.Sprintf(" graph <%s>  { ", g))
	u := " } }"
	uab := []byte(u)
	pab = append(pab, gab...)
	pab = append(pab, []byte(nt)...)
	pab = append(pab, uab...)

	// fmt.Println(string(pab))
	//su := "INSERT DATA {" + s + "}"

	req, err := http.NewRequest("POST", sue, bytes.NewBuffer(pab))
	req.Header.Set("Content-Type", "application/sparql-update")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("response Body:", string(body))
		log.Println("response Status:", resp.Status)
		log.Println("response Headers:", resp.Header)
	}

	return body, err
}

// BlazeUpateNT TODO   rename to jenaUpdateNQ(s []bytes, graph string)  need a jenaUpateNT(s []bytes)
func BlazeUpateNT(s []byte, sue string) ([]byte, error) {

	p := "INSERT DATA { "
	u := " }"
	pab := []byte(p)
	uab := []byte(u)
	pab = append(pab, s...)
	pab = append(pab, uab...)

	// fmt.Println(string(pab))
	//su := "INSERT DATA {" + s + "}"

	req, err := http.NewRequest("POST", sue, bytes.NewBuffer(pab))
	req.Header.Set("Content-Type", "application/sparql-update")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("response Body:", string(body))
		log.Println("response Status:", resp.Status)
		log.Println("response Headers:", resp.Header)
	}

	return body, err
}
