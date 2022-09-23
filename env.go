package main

import "os"

// TablestoreConfig 实例结构
type TablestoreConfig struct {
	Endpoint string
	Instance string
}

type AkSk struct {
	AkID     string
	AkSecret string
}

var ReadAkSk = AkSk{
	AkID:     os.Getenv("read_ak"),
	AkSecret: os.Getenv("read_ak"),
}

var ReadStores = []TablestoreConfig{
	{
		Endpoint: os.Getenv("read_endpoint1"),
		Instance: os.Getenv("read_instance1"),
	},
	{
		Endpoint: os.Getenv("read_endpoint2"),
		Instance: os.Getenv("read_instance2"),
	},
	{
		Endpoint: os.Getenv("read_endpoint3"),
		Instance: os.Getenv("read_instance3"),
	},
}

var WriteAkSk = AkSk{
	AkID:     os.Getenv("write_ak"),
	AkSecret: os.Getenv("write_sk"),
}

var WriteStores = []TablestoreConfig{
	{
		Endpoint: os.Getenv("write_endpoint1"),
		Instance: os.Getenv("write_instance1"),
	},
	{
		Endpoint: os.Getenv("write_endpoint2"),
		Instance: os.Getenv("write_instance2"),
	},
	{
		Endpoint: os.Getenv("write_endpoint3"),
		Instance: os.Getenv("write_instance3"),
	},
}
