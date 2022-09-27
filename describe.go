package main

import (
	"encoding/json"
	"fmt"
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
	// schemaEntry, _ := json.MarshalIndent(resp.TableMeta.SchemaEntry, "", " ")
	fmt.Printf("tableName: %s, PrimaryKey: %s\n", tableName, (TableMetaAlias(*resp.TableMeta).toJSON()))

	for _, indexMeta := range resp.IndexMetas {
		fmt.Printf("tableName:%s, indexName: %s, indexType: %d PrimaryKey: %#v\n", tableName, indexMeta.IndexName, indexMeta.IndexType, indexMeta.Primarykey)
	}
}

type TableMetaAlias tablestore.TableMeta
type TableMetaPrimaryType tablestore.PrimaryKeyType
type TableMetaPrimaryOption tablestore.PrimaryKeyOption
type TableMetaColumnType tablestore.DefinedColumnType

type PkKeySchema struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Option string `json:"option"`
}

type ColumnSchema struct {
	Name       string
	ColumnType string
}

func (keyType TableMetaPrimaryType) typeString() string {
	var s = (tablestore.PrimaryKeyType)(keyType)
	switch s {
	case tablestore.PrimaryKeyType_INTEGER:
		return "INTEGER"
	case tablestore.PrimaryKeyType_STRING:
		return "STRING"
	case tablestore.PrimaryKeyType_BINARY:
		return "BINARY"
	}
	return ""
}

func (keyOption *TableMetaPrimaryOption) optionString() string {
	if keyOption == nil {
		return ""
	}
	var s = (tablestore.PrimaryKeyOption)(*keyOption)
	switch s {
	case tablestore.NONE:
		return "AUTO_INCREMENT"
	case tablestore.AUTO_INCREMENT:
		return "DEFAULT"
	case tablestore.MIN:
		return "MIN"
	case tablestore.MAX:
		return "MAX"
	}
	return ""
}

func (columnType TableMetaColumnType) columnType() string {
	var s = (tablestore.DefinedColumnType)(columnType)
	switch s {
	case tablestore.DefinedColumn_INTEGER:
		return "int64"
	case tablestore.DefinedColumn_DOUBLE:
		return "double"
	case tablestore.DefinedColumn_BOOLEAN:
		return "boolean"
	case tablestore.DefinedColumn_STRING:
		return "string"
	case tablestore.DefinedColumn_BINARY:
		return "binary"
	}
	return ""
}

func (t TableMetaAlias) toJSON() []byte {
	var pks []PkKeySchema
	for _, pk := range t.SchemaEntry {
		pks = append(pks, PkKeySchema{
			Name:   *pk.Name,
			Type:   (TableMetaPrimaryType(*pk.Type).typeString()),
			Option: ((*TableMetaPrimaryOption)(pk.Option).optionString()),
		})
	}
	var cks []ColumnSchema
	for _, ck := range t.DefinedColumns {
		cks = append(cks, ColumnSchema{
			Name:       ck.Name,
			ColumnType: (TableMetaColumnType(ck.ColumnType).columnType()),
		})
	}
	meta := struct {
		TableName   string
		PrimaryKeys []PkKeySchema
		ColumnKeys  []ColumnSchema
	}{
		TableName:   t.TableName,
		PrimaryKeys: pks,
		ColumnKeys:  cks,
	}
	str, _ := json.MarshalIndent(meta, "", " ")
	return str
}
