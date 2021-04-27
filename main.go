package main

import (
	"log"
	"os"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
)

// TablestoreConfig 实例结构
type TablestoreConfig struct {
	Endpoint  string
	Instance  string
	AkID      string
	AkSecret  string
}

var readConfig = TablestoreConfig{
	Endpoint: os.Getenv("read_endpoint"),
	Instance: os.Getenv("read_instance"),
	AkID: os.Getenv("read_ak"),
	AkSecret: os.Getenv("read_sk"),
}

var writeConfig = TablestoreConfig{
	Endpoint: os.Getenv("write_endpoint"),
	Instance: os.Getenv("write_instance"),
	AkID: os.Getenv("write_ak"),
	AkSecret: os.Getenv("write_sk"),
}

// TODO: 删除表结构时，如果有多元索引，需要先删除多元索引，才能删除表结构
func delete(client *tablestore.TableStoreClient) {
	list, err := client.ListTable()
	log.Println("list: ", list, " err: ", err)
	for _, tableName := range list.TableNames {
		log.Println("table: ", tableName)
		res, err := client.DeleteTable(&tablestore.DeleteTableRequest{TableName: tableName})
		log.Println("delete table: ", res, " err: ", err)
	}
}

func copy(readClient *tablestore.TableStoreClient, writeClient *tablestore.TableStoreClient) {
	list, err := readClient.ListTable()
	log.Println("read table list: ", list, "err: ", err)
	for _, tableName := range list.TableNames {
		copyTable(readClient, writeClient, tableName)
	}
}

func copyTable(readClient *tablestore.TableStoreClient, writeClient *tablestore.TableStoreClient, tableName string) {
	sourceTable, _ := readClient.DescribeTable(&tablestore.DescribeTableRequest{TableName: tableName})
	_, err := writeClient.DescribeTable(&tablestore.DescribeTableRequest{TableName: tableName})
	if err != nil {
		createRequest := tablestore.CreateTableRequest{
			TableMeta:          sourceTable.TableMeta,
			TableOption:        sourceTable.TableOption,
			ReservedThroughput: sourceTable.ReservedThroughput,
			StreamSpec:         &tablestore.StreamSpecification{
				EnableStream:   sourceTable.StreamDetails.EnableStream,
				ExpirationTime: sourceTable.StreamDetails.ExpirationTime,
			},
			IndexMetas:         sourceTable.IndexMetas,
		}
		res, err := writeClient.CreateTable(&createRequest)
		if err != nil {
			log.Println("write client create table fail: ", err)
		} else {
			log.Println("write client create table success: ", res)
			copySearchIndex(readClient, writeClient, tableName)
		}
	} else {
		log.Println("write client table: " + tableName + " exist")
	}
}

func copySearchIndex(readClient *tablestore.TableStoreClient, writeClient *tablestore.TableStoreClient, tableName string) {
	request := &tablestore.ListSearchIndexRequest{}
    request.TableName = tableName
    resp, err := readClient.ListSearchIndex(request)
    if err != nil {
        log.Println("list " + tableName + " searchIndex error: ", err)
        return
    }
    for _, info := range resp.IndexInfo {
		createSearchIndex(readClient, writeClient, *info)
    }
    log.Println("list " + tableName + " search index finished, requestId:", resp.ResponseInfo.RequestId)
}

func createSearchIndex(readClient *tablestore.TableStoreClient, writeClient *tablestore.TableStoreClient, indexInfo tablestore.IndexInfo) {
    describeRequest := &tablestore.DescribeSearchIndexRequest{}
    describeRequest.TableName = indexInfo.TableName
    describeRequest.IndexName = indexInfo.IndexName
    resp, err := readClient.DescribeSearchIndex(describeRequest)
    if err != nil {
        log.Println("create table name: " + indexInfo.TableName + " index name: " + indexInfo.IndexName + " search index error: ", err)
        return
    }
	schemas := []*tablestore.FieldSchema{}
	createRequest := &tablestore.CreateSearchIndexRequest{}
	createRequest.TableName = indexInfo.TableName
	createRequest.IndexName = indexInfo.IndexName

    for _, schema := range resp.Schema.FieldSchemas {	
		field1 := &tablestore.FieldSchema{
			FieldName: schema.FieldName,
			FieldType: schema.FieldType,
			Index:     schema.Index,
			EnableSortAndAgg: schema.EnableSortAndAgg,
		}
		schemas = append(schemas, field1)
    }
	createRequest.IndexSchema = &tablestore.IndexSchema{
		FieldSchemas: schemas,
	}
	_, cerr := writeClient.CreateSearchIndex(createRequest)
	if cerr != nil {
		log.Println("error :", err)
		return
	}
    log.Println("CreateSearchIndex finished, requestId:", resp.ResponseInfo.RequestId)
}

func deleteSearchIndex(client *tablestore.TableStoreClient, tableName string) {
	request := &tablestore.ListSearchIndexRequest{}
    request.TableName = tableName
    resp, err := client.ListSearchIndex(request)
    if err != nil {
        log.Println("list " + tableName + " searchIndex error: ", err)
        return
    }
    for _, info := range resp.IndexInfo {
		request := &tablestore.DeleteSearchIndexRequest{}
		request.IndexName = info.IndexName
		request.TableName = tableName;
		client.DeleteSearchIndex(request)
    }
    log.Println("list " + tableName + " search index finished, requestId:", resp.ResponseInfo.RequestId)
}

func main()  {
	log.Printf("readClient: %+v\n", readConfig)
	log.Printf("writeClient: %+v\n", writeConfig)
	readClient := tablestore.NewClient(readConfig.Endpoint, readConfig.Instance, readConfig.AkID, readConfig.AkSecret)
	writeClient := tablestore.NewClient(writeConfig.Endpoint, writeConfig.Instance, writeConfig.AkID, writeConfig.AkSecret)
	// deleteTables(writeTable)
	copy(readClient, writeClient)
}