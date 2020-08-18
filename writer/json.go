package writer

import (
	"io"

	"github.com/syucream/goparquet/layout"
	"github.com/syucream/goparquet/marshal"
	"github.com/syucream/goparquet/parquet"
	"github.com/syucream/goparquet/schema"
)

type JSONWriter struct {
	ParquetWriter
}

//Create JSON writer
func NewJSONWriter(jsonSchema string, w io.WriteCloser, np int64) (*JSONWriter, error) {
	var err error
	res := new(JSONWriter)
	res.SchemaHandler, err = schema.NewSchemaHandlerFromJSON(jsonSchema)
	if err != nil {
		return res, err
	}

	res.w = w
	res.PageSize = 8 * 1024              //8K
	res.RowGroupSize = 128 * 1024 * 1024 //128M
	res.CompressionType = parquet.CompressionCodec_SNAPPY
	res.PagesMapBuf = make(map[string][]*layout.Page)
	res.DictRecs = make(map[string]*layout.DictRecType)
	res.NP = np
	res.Footer = parquet.NewFileMetaData()
	res.Footer.Version = 1
	res.Footer.Schema = append(res.Footer.Schema, res.SchemaHandler.SchemaElements...)
	res.Offset = 4
	_, err = res.w.Write(magic)
	res.MarshalFunc = marshal.MarshalJSON
	return res, err
}
