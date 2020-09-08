package main

import (
	"regexp"
	"strconv"
)

var (
	clauses        []string
	volatileMemory map[string]string
	rgx            = regexp.MustCompile(`^Run\((.*?)\)`)
)

func InitializeKB() {
	volatileMemory = make(map[string]string)
}

func AddClause(clause string) bool {
	clauses = append(clauses, clause)
	return true
}

func RemoveClause(targetClause string) bool {
	for i, clause := range clauses {
		if clause == targetClause {
			clauses = append(clauses[:i], clauses[i+1:]...)
			return true
		}
	}
	return false
}

func PrintClauses() []string {
	return clauses
}

func PrintMemory() map[string]string {
	return volatileMemory
}

func Memorize(subject string, value string) {
	volatileMemory[subject] = value

	// TODO: this will be removed as we bring KB here
	if subject == "rain_gauge" {
		intValue, _ := strconv.Atoi(value)
		if intValue > 3 {
			AddClause("Rain(Now)")
			AddClause("Run(Cloud)")
		} else {
			RemoveClause("Rain(Now)")
			RemoveClause("Run(Cloud)")
		}
	}
}

func ClearMemory(subject string) {
	delete(volatileMemory, subject)
}

func Ask() (plugins []string) {
	plugins = make([]string, 5)
	for i := 0; i < len(clauses); i++ {
		rs := rgx.FindStringSubmatch(clauses[i])
		if len(rs) >= 1 {
			plugins = append(plugins, rs[1])
		}
	}
	return
}
