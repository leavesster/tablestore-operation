package main

import (
	"log"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
)

func main() {
	log.Printf("readStores: %+v writeStores: %+v\n", ReadStores, WriteStores)

	for _, store := range ReadStores {
		client := tablestore.NewClient(store.Endpoint, store.Instance, ReadAkSk.AkID, ReadAkSk.AkSecret)
		DescribeTableClient(client)
	}
}
