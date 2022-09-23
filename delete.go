package main

import (
	"log"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
)

// TODO: 删除表结构时，如果有多元索引，需要先删除多元索引，才能删除表结构
func Delete(client *tablestore.TableStoreClient) {
	list, err := client.ListTable()
	log.Println("list: ", list, " err: ", err)
	for _, tableName := range list.TableNames {
		log.Println("table: ", tableName)
		res, err := client.DeleteTable(&tablestore.DeleteTableRequest{TableName: tableName})
		log.Println("delete table: ", res, " err: ", err)
	}
}

func deleteSearchIndex(client *tablestore.TableStoreClient, tableName string) {
	request := &tablestore.ListSearchIndexRequest{}
	request.TableName = tableName
	resp, err := client.ListSearchIndex(request)
	if err != nil {
		log.Println("list "+tableName+" searchIndex error: ", err)
		return
	}
	for _, info := range resp.IndexInfo {
		request := &tablestore.DeleteSearchIndexRequest{}
		request.IndexName = info.IndexName
		request.TableName = tableName
		client.DeleteSearchIndex(request)
	}
	log.Println("list "+tableName+" search index finished, requestId:", resp.ResponseInfo.RequestId)
}
