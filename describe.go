package main

import (
	"encoding/json"
	"log"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
)

func DescribeTableClient(client *tablestore.TableStoreClient) []string {
	list, err := client.ListTable()
	log.Println("list: ", list, " err: ", err)
	for _, tableName := range list.TableNames {
		describeTableName(client, tableName)
	}
	return list.TableNames
}

func describeTableName(client *tablestore.TableStoreClient, tableName string) {
	request := &tablestore.DescribeTableRequest{}
	request.TableName = tableName
	resp, err := client.DescribeTable(request)
	if err != nil {
		log.Println("describe "+tableName+" error: ", err)
		return
	}
	schemaEntry, _ := json.MarshalIndent(resp.TableMeta.SchemaEntry, "", " ")
	log.Printf("tablestore tableName: %s, PrimaryKey: %s\n", tableName, schemaEntry)
	for _, indexMeta := range resp.IndexMetas {
		log.Printf("indexName: %s, indexType: %d PrimaryKey: %#v\n", indexMeta.IndexName, indexMeta.IndexType, indexMeta.Primarykey)
	}
}
