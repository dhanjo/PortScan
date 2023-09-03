// scanner.go

package scanner

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/t94j0/nmap"
)

type CustomResult struct {
	Hosts map[string]nmap.Host `json:"Hosts"`
}

func Scan(domain string) (string, error) {
	var wg sync.WaitGroup

	scan := nmap.Init()
	scan = scan.AddHosts(domain)
	scan = scan.AddFlags("-sV", "-sS")

	wg.Add(1)

	result, err := scan.Run()
	if err != nil {
		log.Fatal(err)
	}

	customResult := CustomResult{
		Hosts: result.Hosts,
	}

	jsonData, err := json.MarshalIndent(customResult, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	if string(jsonData) == "{}" {
		return "", fmt.Errorf("No meaningful data for domain: %s", domain)
	}

	defer wg.Done()
	return string(jsonData), nil
}
