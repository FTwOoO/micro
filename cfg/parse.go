package cfg

import (
	"encoding/json"
	"flag"
	"os"
)

func NewConfiguration(path string, obj interface{}) (err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}

	body := make([]byte, 1024*16)
	n, err := f.Read(body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body[:n], obj)
	return
}

func ParseConfigFileFromFlags(configPointer interface{}) {
	var err error

	flagSet := flag.NewFlagSet("goservice", flag.ExitOnError)
	configFilePath := flagSet.String("conf", "", "config file path")
	err = flagSet.Parse(os.Args[1:])
	if err != nil {
		panic(err)
	}

	err = NewConfiguration(*configFilePath, configPointer)
	if err != nil {
		panic(err)
	}
}