package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type stringList []string

func (s *stringList) String() string {
	return fmt.Sprintf("%v", *s)
}

func (s *stringList) Set(value string) error {
	*s = strings.Split(value, ",")
	return nil
}

func main() {
	borkerCommand := flag.NewFlagSet("broker", flag.ExitOnError)
	dbCommand := flag.NewFlagSet("db", flag.ExitOnError)

	borkerAdd := borkerCommand.String("add", "", "add uri to broker")
	borkerExport := borkerCommand.String("export", "", " export data in broker to database . (Required)")

	var brokerStringList dbStringList
	borkerCommand.Var(&brokerStringList, "brokerStringList", "A comma seperated list of substrings to be counted.")

	dbAdd := dbCommand.String("add", "", "add connection to database. ")
	dbRemove := dbCommand.String("remove", "", "remove connection to database")


	if len(os.Args) < 2 {
		fmt.Println("list or count subcommand is required")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "list":
		listCommand.Parse(os.Args[2:])
	case "count":
		countCommand.Parse(os.Args[2:])
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	if listCommand.Parsed() {
		// Required Flags
		if *listTextPtr == "" {
			listCommand.PrintDefaults()
			os.Exit(1)
		}
		//Choice flag
		metricChoices := map[string]bool{"chars": true, "words": true, "lines": true}
		if _, validChoice := metricChoices[*listMetricPtr]; !validChoice {
			listCommand.PrintDefaults()
			os.Exit(1)
		}
		// Print
		fmt.Printf("textPtr: %s, metricPtr: %s, uniquePtr: %t\n",
			*listTextPtr,
			*listMetricPtr,
			*listUniquePtr,
		)
	}

	if countCommand.Parsed() {
		// Required Flags
		if *countTextPtr == "" {
			countCommand.PrintDefaults()
			os.Exit(1)
		}
		// If the metric flag is substring, the substring or substringList flag is required
		if *countMetricPtr == "substring" && *countSubstringPtr == "" && (&countStringList).String() == "[]" {
			countCommand.PrintDefaults()
			os.Exit(1)
		}
		//If the metric flag is not substring, the substring flag must not be used
		if *countMetricPtr != "substring" && (*countSubstringPtr != "" || (&countStringList).String() != "[]") {
			fmt.Println("--substring and --substringList may only be used with --metric=substring.")
			countCommand.PrintDefaults()
			os.Exit(1)
		}
		//Choice flag
		metricChoices := map[string]bool{"chars": true, "words": true, "lines": true, "substring": true}
		if _, validChoice := metricChoices[*listMetricPtr]; !validChoice {
			countCommand.PrintDefaults()
			os.Exit(1)
		}
		//Print
		fmt.Printf("textPtr: %s, metricPtr: %s, substringPtr: %v, substringListPtr: %v, uniquePtr: %t\n",
			*countTextPtr,
			*countMetricPtr,
			*countSubstringPtr,
			(&countStringList).String(),
			*countUniquePtr,
		)
	}
}