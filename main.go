package main

import (
	"flag"

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
	var namespaces ListStringVar
	var matches ListStringVar
	var and bool
	flag.Var(&namespaces, "namespace", "The namespaces to look into")
	flag.Var(&matches, "match", "Matches")
	flag.BoolVar(&and, "and", false, "All the matches must be in the job")
	flag.Parse()

	m := cli.New(namespaces, matches, and)
	m.Run()

}
