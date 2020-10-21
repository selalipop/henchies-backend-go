package main

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type Arguments struct {
	RedisConnectURL string `required:"true"`
}

func getArgumentsFromEnv() Arguments {
	var args Arguments
	const argumentsPrefix = "henchies"
	err := envconfig.Process(argumentsPrefix, &args)
	if err != nil {
		logrus.Panicf("failed to parse enviornment arguments %v", err)
	}
	return args
}
