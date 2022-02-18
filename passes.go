package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func checkValidatorPasses(heightFrom int64) {
	file, err := os.Open(statPath)
	if err != nil {
		log.Fatalf("Error when opening file: %s", err)
	}

	fileScanner := bufio.NewScanner(file)

	missedValidators := make(map[string][]string)

	for fileScanner.Scan() {
		columns := strings.Split(fileScanner.Text(), ";")
		height, _ := strconv.ParseInt(columns[0], 10, 64)
		if height < heightFrom {
			continue
		}

		if len(columns[2]) < 2 {
			continue
		}

		proposer := columns[1]
		for _, validator := range strings.Split(columns[2], ",") {
			missedValidators[validator] = append(missedValidators[validator], proposer)
		}
	}

	for validator, proposers := range missedValidators {
		fmt.Println(validator, validatorNames[validator])

		p := make(map[string]int)
		for _, proposer := range proposers {
			p[proposer] += 1
		}

		for k, v := range p {
			fmt.Println(v, ": ", k, validatorNames[k])
		}

		fmt.Println(len(proposers), "\n")
	}

	if err = fileScanner.Err(); err != nil {
		log.Fatalf("Error while reading file: %s", err)
	}

	file.Close()
}
