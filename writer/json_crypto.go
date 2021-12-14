package writer

import (
	"io"

	"github.com/xitongsys/parquet-go-source/writerfile"
	"github.com/xitongsys/parquet-go/layout"
	"github.com/xitongsys/parquet-go/marshal"
	"github.com/xitongsys/parquet-go/parquet"
	"github.com/xitongsys/parquet-go/schema"
	"github.com/xitongsys/parquet-go/source"
)

type JSONEncWriter struct {
	EncryptedParquetWriter
}

func NewJSONEncWriterFromWriter(jsonSchema string, w io.Writer, np int64) (*JSONWriter, error) {
	wf := writerfile.NewWriterFile(w)
	return NewJSONWriter(jsonSchema, wf, np)
}

//Create JSON writer
func NewJSONEncWriter(jsonSchema string, pfile source.ParquetFile, np int64) (*JSONEncWriter, error) {
	var err error
	res := new(JSONEncWriter)
	res.SchemaHandler, err = schema.NewSchemaHandlerFromJSON(jsonSchema)
	if err != nil {
		return res, err
	}

	res.PFile = pfile
	res.PageSize = 8 * 1024              //8K
	res.RowGroupSize = 128 * 1024 * 1024 //128M
	res.CompressionType = parquet.CompressionCodec_SNAPPY
	res.PagesMapBuf = make(map[string][]*layout.Page)
	res.DictRecs = make(map[string]*layout.DictRecType)
	res.NP = np
	res.Footer = parquet.NewFileCryptoMetaData()
	res.Footer.Version = 1
	res.Footer.Schema = append(res.Footer.Schema, res.SchemaHandler.SchemaElements...)
	res.Offset = 4
	_, err = res.PFile.Write([]byte("PAR1"))
	res.MarshalFunc = marshal.MarshalJSON
	return res, err
}
