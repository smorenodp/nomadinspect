package main

import (
	"flag"
	"log"

	"github.com/smorenodp/nomadinspect/cli"
)

type ListStringVar []string

func (l *ListStringVar) String() string {
	return "TODO"
}

func (l *ListStringVar) Set(value string) error {
	*l = append(*l, value)
	return nil
}

func main() {
	var namespaces, matches, not ListStringVar
	var and bool
	flag.Var(&namespaces, "namespace", "The namespaces to look into")
	flag.Var(&matches, "match", "Matches")
	flag.Var(&not, "not", "Matches that must not be in the job")
	flag.BoolVar(&and, "and", false, "All the matches must be in the job")
	flag.Parse()

	if len(matches)+len(not) == 0 {
		log.Fatal("[ERROR] You have to at least select one match or not")
	}
	m := cli.New(namespaces, matches, and, not)
	m.Run()

}
