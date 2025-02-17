package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"golang.org/x/exp/maps"
)

var pythonExecutable string

func main() {
	flag.StringVar(&pythonExecutable, "python-executable", "", "path to the Python executable")
	flag.Parse()

	if pythonExecutable == "" {
		flag.Usage()
		os.Exit(2)
	}

	if _, err := os.Stat(pythonExecutable); err != nil {
		log.Fatalf("invalid value for -python-executable: %v", err)
	}

	var languages []Language
	if precommitConfig, err := FindPrecommitConfig(); err == nil {
		languages, err = EFMConfigFromPrecommit(precommitConfig)
		if err != nil {
			log.Fatalf("failed to load languages from %q: %v", precommitConfig, err)
		}
	}

	languages = removeDuplicates(languages)
	err := json.NewEncoder(os.Stdout).Encode(languages)
	if err != nil {
		log.Fatal(err)
	}
}

func removeDuplicates(langs []Language) []Language {
	m := map[string]Language{}
	for _, lang := range langs {
		if lang.FormatCommand != "" {
			m["format#"+lang.FormatCommand] = lang
		} else if lang.LintCommand != "" {
			m["lint#"+lang.LintCommand] = lang
		}
	}

	return maps.Values(m)
}
