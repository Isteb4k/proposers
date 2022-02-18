package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Slash struct {
	m       [SignedBlocksWindow]bool
	index   int64
	offline bool
}

func checkValidatorSlashes(heightFrom int64) {
	file, err := os.Open(statPath)

	if err != nil {
		log.Fatalf("Error when opening file: %s", err)
	}

	validatorsState := make(map[string]*Slash)

	fileScanner := bufio.NewScanner(file)

	for fileScanner.Scan() {
		columns := strings.Split(fileScanner.Text(), ";")

		height, _ := strconv.ParseInt(columns[0], 10, 64)
		if height < heightFrom {
			continue
		}

		missedValidators := make(map[string]bool)

		if len(columns[2]) > 3 {
			validators := strings.Split(columns[2], ",")

			for _, id := range validators {
				if _, ok := validatorsState[id]; !ok {
					validatorsState[id] = &Slash{
						m:     [24]bool{},
						index: 0,
					}
				}

				missedValidators[id] = true
			}
		}

		for validator, info := range validatorsState {
			if info.offline {
				continue
			}

			info.index = (info.index + 1) % SignedBlocksWindow
			index := info.index

			// Update signed block bit array & counter
			// This counter just tracks the sum of the bit array
			// That way we avoid needing to read/write the whole array each time
			missedInWindow := info.m[index]

			switch missedValidators[validator] {
			case true:
				// If missed < 24 then missed = missed + 1
				if countMissedBlocks(info.m) < SignedBlocksWindow && !missedInWindow {
					info.m[index] = true
				}
			case false:
				if (countMissedBlocks(info.m) > 0) && (missedInWindow) {
					info.m[index] = false
				}
			}

			if countMissedBlocks(info.m) > 11 {
				fmt.Println(validator, validatorNames[validator], height)
				info.offline = true
			}
		}
	}
}

func countMissedBlocks(m [SignedBlocksWindow]bool) int64 {
	var count int64

	for _, missed := range m {
		if missed {
			count++
		}
	}

	return count
}

const SignedBlocksWindow int64 = 24
