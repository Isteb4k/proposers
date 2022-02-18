package main

import (
	"flag"
	"log"
)

const statPath = "/home/centos/proposers.txt"

type Usecase = string

const (
	CheckProposersMassPassesCase = "proposers"
	CheckValidatorPassesCase     = "passes"
	CheckValidatorSlashesCase    = "slashes"
)

func main() {
	var usecase Usecase
	flag.StringVar(&usecase, "case", "", "Case")

	var heightFrom int64
	flag.Int64Var(&heightFrom, "from", 0, "From")
	flag.Parse()

	switch usecase {
	case CheckProposersMassPassesCase:
		checkMassPassesForProposers(heightFrom)
	case CheckValidatorPassesCase:
		checkValidatorPasses(heightFrom)
	case CheckValidatorSlashesCase:
		checkValidatorSlashes(heightFrom)
	default:
		log.Fatalln("unknown case", usecase)
	}

	return
}
