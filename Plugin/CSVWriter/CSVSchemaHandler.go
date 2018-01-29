package CSVWriter

import (
	"github.com/xitongsys/parquet-go/Common"
	"github.com/xitongsys/parquet-go/SchemaHandler"
	"github.com/xitongsys/parquet-go/parquet"
)

//Create a schema handler from CSV metadata
func NewSchemaHandlerFromMetadata(mds []string) *SchemaHandler.SchemaHandler {
	schemaList := make([]*parquet.SchemaElement, 0)
	infos := make([]*Common.Tag, 0)

	rootSchema := parquet.NewSchemaElement()
	rootSchema.Name = "parquet_go_root"
	rootNumChildren := int32(len(mds))
	rootSchema.NumChildren = &rootNumChildren
	rt := parquet.FieldRepetitionType(-1)
	rootSchema.RepetitionType = &rt
	schemaList = append(schemaList, rootSchema)

	for _, md := range mds {
		info := Common.StringToTag(md)
		infos = append(infos, info)

		schema := parquet.NewSchemaElement()
		schema.Name = info.ExName
		numChildren := int32(0)
		schema.NumChildren = &numChildren
		rt := parquet.FieldRepetitionType(1)
		schema.RepetitionType = &rt

		if t, err := parquet.TypeFromString(info.Type); err == nil {
			schema.Type = &t
			if info.Type == "FIXED_LEN_BYTE_ARRAY" {
				schema.TypeLength = &info.Length
			}
		} else {
			name := info.Type
			ct, _ := parquet.ConvertedTypeFromString(name)
			schema.ConvertedType = &ct
			if name == "INT_8" || name == "INT_16" || name == "INT_32" ||
				name == "UINT_8" || name == "UINT_16" || name == "UINT_32" ||
				name == "DATE" || name == "TIME_MILLIS" {
				schema.Type = parquet.TypePtr(parquet.Type_INT32)
			} else if name == "INT_64" || name == "UINT_64" ||
				name == "TIME_MICROS" || name == "TIMESTAMP_MICROS" {
				schema.Type = parquet.TypePtr(parquet.Type_INT64)
			} else if name == "UTF8" {
				schema.Type = parquet.TypePtr(parquet.Type_BYTE_ARRAY)
			} else if name == "INTERVAL" {
				schema.Type = parquet.TypePtr(parquet.Type_FIXED_LEN_BYTE_ARRAY)
				var ln int32 = 12
				schema.TypeLength = &ln
			} else if name == "DECIMAL" {
				schema.Type = parquet.TypePtr(parquet.Type_BYTE_ARRAY)
				schema.Scale = &info.Scale
				schema.Precision = &info.Precision
			}
		}
		schemaList = append(schemaList, schema)
	}
	res := SchemaHandler.NewSchemaHandlerFromSchemaList(schemaList)
	res.Infos = infos
	return res
}
