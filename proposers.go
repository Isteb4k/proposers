package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
)

func checkMassPassesForProposers(heightFrom int64) {
	file, err := os.Open(statPath)
	if err != nil {
		log.Fatalf("Error when opening file: %s", err)
	}

	fileScanner := bufio.NewScanner(file)

	// passes map stores the number of validators that skipped a block for proposer
	passes := make(map[string][]int)
	proposers := make(map[string]int)

	for fileScanner.Scan() {
		columns := strings.Split(fileScanner.Text(), ";")
		height, _ := strconv.ParseInt(columns[0], 10, 64)
		if height < heightFrom {
			continue
		}

		proposer := columns[1]
		missedValidators := columns[2]

		proposers[proposer] = proposers[proposer] + 1

		validatorsCount := len(strings.Split(missedValidators, ","))

		if validatorsCount > 1 {
			passes[proposer] = append(passes[proposer], validatorsCount)
		}
	}

	// log table header
	w := tabwriter.NewWriter(os.Stdout, 15, 0, 1, ' ', 0)
	fmt.Fprintln(w, "\t       id\t\t     proposer:\t    miss\t    all\t     %\t   avg \t   min \t   max \t")
	w.Flush()

	// log info about mass passes for every proposer
	for proposer, massPasses := range passes {
		var min = math.MaxInt32
		var max int
		var avg int
		for _, missedValidatorsCount := range massPasses {
			if missedValidatorsCount > max {
				max = missedValidatorsCount
			}
			if missedValidatorsCount < min {
				min = missedValidatorsCount
			}
			avg += missedValidatorsCount
		}

		str := fmt.Sprintf(
			"%s\t %s:\t%d\t%d\t%0.2f\t%d\t%d\t%d\t",
			proposer,
			validatorNames[proposer],
			len(massPasses),
			proposers[proposer],
			float64(len(massPasses))/float64(proposers[proposer])*100,
			avg/len(massPasses),
			min,
			max,
		)

		// log table row
		w := tabwriter.NewWriter(os.Stdout, 15, 0, 1, ' ', 0)
		fmt.Fprintln(w, str)
		w.Flush()
	}

	if err := fileScanner.Err(); err != nil {
		log.Fatalf("Error while reading file: %s", err)
	}

	file.Close()
}
