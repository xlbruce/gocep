package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type CepInfo map[string]string

func (c *CepInfo) Json(writer io.Writer) {
	bytes, err := json.MarshalIndent(*c, "", "  ")
	if err != nil {
		panic(err)
	}
	writer.Write(bytes)
}

func (c *CepInfo) Piped(writer io.Writer) {
	var values []string
	for key, value := range *c {
		values = append(values, fmt.Sprintf("%s:%s", key, value))
	}
	bytes := []byte(strings.Join(values, "|"))
	writer.Write(bytes)
}

const (
	urlTemplate = "https://viacep.com.br/ws/%s/json/"
)

var (
	cep     = flag.String("cep", "", "CEP number. For instance: 01001000")
	format  = flag.String("format", "json", "Output format. Possible values: json, piped")
	cepInfo *CepInfo
)

func usage() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s: \n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Usage()
	os.Exit(1)
}

func main() {
	flag.Parse()
	if *cep == "" {
		usage()
	}

	url := fmt.Sprintf(urlTemplate, *cep)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	cepInfo = newCepInfo(resp.Body)

	switch *format {
	case "json":
		cepInfo.Json(os.Stdout)
	case "piped":
		cepInfo.Piped(os.Stdout)
	}
}

func newCepInfo(reader io.Reader) *CepInfo {
	info := &CepInfo{}
	json.NewDecoder(reader).Decode(info)
	return info
}
