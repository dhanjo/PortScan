// main.go

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	scanner "Portscan/Scanner"
)

func main() {
	var wg sync.WaitGroup

	// Define flags
	fileName := flag.String("f", "", "Enter file name to scan ports")
	outputFileName := flag.String("o", "", "Enter output file name")
	flag.Parse()

	if flag.NFlag() == 0 {
		flag.Usage()
		return
	}

	data, err := ioutil.ReadFile(*fileName)

	if err != nil {
		fmt.Println("Error in reading file", err)
	}

	lines := strings.Split(string(data), "\n")

	// Create an output file if provided
	var outputFile *os.File
	if *outputFileName != "" {
		outputFile, err = os.Create(*outputFileName)
		if err != nil {
			fmt.Println("Error creating output file:", err)
		}
		defer outputFile.Close()
	}

	for _, line := range lines {
		wg.Add(1)
		go func(line string) {
			defer wg.Done()
			result, err := scanner.Scan(line)
			if err != nil {
				fmt.Println(err)
				return
			}

			// Print the result to the console
			fmt.Printf("%s\n", result)

			// Write the result to the output file if provided
			if outputFile != nil {
				_, err := outputFile.WriteString(result + "\n")
				if err != nil {
					fmt.Println("Error writing to output file:", err)
				}
			}
		}(line)
	}

	wg.Wait()
}
