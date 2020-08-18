package layout

import (
	"testing"

	"github.com/syucream/parquet-go/parquet"
)

func TestTableToDictPages(t *testing.T) {
	table := NewEmptyTable()
	table.Values = []interface{}{"test", "test"}
	table.DefinitionLevels = []int32{0, 0}
	table.RepetitionLevels = []int32{0, 0}
	table.Schema = &parquet.SchemaElement{
		Type:          parquet.TypePtr(parquet.Type_BYTE_ARRAY),
		ConvertedType: parquet.ConvertedTypePtr(parquet.ConvertedType_UTF8),
	}

	// Check unless panic for now
	func() {
		defer func() {
			if err := recover(); err != nil {
				t.Error(err)
			}
		}()
		dictRec := NewDictRec(parquet.Type_BYTE_ARRAY)
		_, _ = TableToDictDataPages(dictRec, table, 8*1024, 32, parquet.CompressionCodec_GZIP)
	}()
}
