package main

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/davidsteed/cognito/lib/logs"
)

var Settings Config

type Config struct {
	LogLevel            string `split_words:"true" default:"info"`
	CommonDBTable       string `split_words:"true" default:"NBIT-BM-Common-Task-Management"`
	CashDBTable         string `split_words:"true" default:"NBIT-BM-Cash-Management"`
	TransactionDB       string `split_words:"true" default:"spm-transactions"`
	LocalPort           int    `default:"9091"`
	LocalDynamoEndpoint string
	CounterBucketName   string `split_words:"true"`
	PouchBucketName     string `split_words:"true"`
	TCBucketName        string `split_words:"true"`
}

func Load() error {
	err := envconfig.Process("app", &Settings)
	if err != nil {
		return err
	}
	if Settings.LocalDynamoEndpoint != "" {
		logs.Log.Info("[LOCAL] DynamoDB service enabled")
	}
	return nil
}
