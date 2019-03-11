package app

import (
	"flag"
	"os"
	"spiderman-agent/utils"
)

type ProcParams struct {
	ConfigPath   string
	Environment  string
	RegisterIp   string
	RegisterName string
}

func LoadParams() (*ProcParams, error) {

	flagSet := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	configPath := flagSet.String("config", "./conf", "config path")
	env := flagSet.String("env", "prod", "project Env")
	registerIp := flagSet.String("ip", utils.GetNodeIp(), "register server ip")
	registerName := flagSet.String("name", "", "register server name")

	flagSet.Parse(os.Args[1:])

	params := &ProcParams{
		ConfigPath:   *configPath,
		Environment:  *env,
		RegisterIp:   *registerIp,
		RegisterName: *registerName,
	}
	if params.RegisterName == "" {
		panic("register server name is null")
	}
	return params, nil
}
