package main

import (
	"flag"
	"time"

	"github.com/iccolo/rma/analyzer"
)

var (
	host       string
	port       uint
	password   string
	count      uint
	limit      uint64
	match      string
	types      string
	separators string
	cluster    bool
	pause      time.Duration
)

func init() {
	flag.StringVar(&host, "h", "127.0.0.1", "host")
	flag.UintVar(&port, "p", 6379, "port")
	flag.StringVar(&password, "a", "", "password")
	flag.UintVar(&count, "count", 10000, "count")
	flag.Uint64Var(&limit, "l", 100000, "limit")
	flag.StringVar(&match, "m", "*", "match")
	flag.StringVar(&types, "t", "", "types")
	flag.StringVar(&separators, "s", ":", "separators")
	flag.BoolVar(&cluster, "c", true, "cluster")
	flag.DurationVar(&pause, "pause", 1000, "pause")
}

func main() {
	flag.Parse()
	a := &analyzer.Analyzer{
		Host:       host,
		Port:       port,
		Password:   password,
		Count:      count,
		Limit:      limit,
		Match:      match,
		Types:      types,
		Separators: separators,
		Cluster:    cluster,
		Pause:      pause,
	}
	tree := a.Run()
	tree.Print()
}
