package util

import (
	"encoding/json"
	"fmt"
	"log"
)

func PrintObjInJson(in any) {
	byte, err := json.MarshalIndent(in, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(byte))
}
