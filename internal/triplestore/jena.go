package triplestore

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/UFOKN/nabu/internal/graph"
)

//JenaUpdateNQ TODO   rename to jenaUpdateNQ(s []bytes, graph string)  need a jenaUpateNT(s []bytes)
func JenaUpdateNQ(s []byte, sue string) ([]byte, error) {

	nt, g, err := graph.NQToNTCtx(string(s))
	if err != nil {
		log.Println("Error in nqToNTCtx")
		log.Println(err)
	}

	// INSERT DATA { GRAPH <http://example/bookStore> { <http://example/book1>  ns:price  42  }  }

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
	// fmt.Println("response Body:", string(body))
	// fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)

	return body, err
}

// JenaUpateNT TODO   rename to jenaUpdateNQ(s []bytes, graph string)  need a jenaUpateNT(s []bytes)
func JenaUpateNT(s []byte, sue string) ([]byte, error) {

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
