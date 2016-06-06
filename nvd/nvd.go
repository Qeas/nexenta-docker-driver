package main

import (
	"github.com/qeas/nexenta-docker-volume/nvdv/nvdapi"
	"github.com/qeas/nexenta-docker-volume/nvdv/nvdcli"
	"os"
)

const (
	VERSION = "0.0.1"
)

var (
	client *nvdapi.Client
)

func main() {
	ncli := nvdli.NewCli(VERSION)
	ncli.Run(os.Args)
}

