package layout

import (
	"testing"

	"github.com/syucream/goparquet/parquet"
)

func TestTableToDataPages(t *testing.T) {
	table := NewEmptyTable()
	table.Values = []interface{}{int32(1), int32(2)}
	table.DefinitionLevels = []int32{0, 0}
	table.RepetitionLevels = []int32{0, 0}
	table.Schema = &parquet.SchemaElement{
		Type:          parquet.TypePtr(parquet.Type_INT32),
		ConvertedType: parquet.ConvertedTypePtr(parquet.ConvertedType_INT_32),
	}

	// Check unless panic for now
	func() {
		defer func() {
			if err := recover(); err != nil {
				t.Error(err)
			}
		}()
		_, _ = TableToDataPages(table, 8*1024, parquet.CompressionCodec_GZIP)
	}()
}

func TestPage_EncodingValues(t *testing.T) {
	cases := []struct {
		tpe *parquet.Type
		enc parquet.Encoding
		v   []interface{}
	}{
		// Plain
		{
			tpe: parquet.TypePtr(parquet.Type_INT64),
			enc: parquet.Encoding_PLAIN,
			v:   []interface{}{int64(42)},
		},

		// RLE
		{
			tpe: parquet.TypePtr(parquet.Type_BYTE_ARRAY),
			enc: parquet.Encoding_RLE,
			v:   []interface{}{string("test")},
		},

		// Delta
		{
			tpe: parquet.TypePtr(parquet.Type_BYTE_ARRAY),
			enc: parquet.Encoding_DELTA_BINARY_PACKED,
			v:   []interface{}{[]byte("test")},
		},

		// Delta Byte Array
		{
			tpe: parquet.TypePtr(parquet.Type_BYTE_ARRAY),
			enc: parquet.Encoding_DELTA_BYTE_ARRAY,
			v:   []interface{}{[]byte("test")},
		},

		// Delta Length Byte Array
		{
			tpe: parquet.TypePtr(parquet.Type_BYTE_ARRAY),
			enc: parquet.Encoding_DELTA_LENGTH_BYTE_ARRAY,
			v:   []interface{}{[]byte("test")},
		},
	}

	for _, c := range cases {
		page := NewDataPage()
		page.Schema = &parquet.SchemaElement{
			Type: c.tpe,
		}
		page.encoding = c.enc
		page.bitWidths = int32(len(c.v))

		// Check unless panic for now
		func() {
			defer func() {
				if err := recover(); err != nil {
					t.Error(err)
				}
			}()
			_ = page.EncodingValues(c.v)
		}()
	}
}

func TestPage_UseDictionaryEncoding(t *testing.T) {
	cases := []struct {
		enc      parquet.Encoding
		expected bool
	}{
		// With dict
		{
			enc:      parquet.Encoding_PLAIN_DICTIONARY,
			expected: true,
		},
		{
			enc:      parquet.Encoding_RLE_DICTIONARY,
			expected: true,
		},

		// Without dict
		{
			enc:      parquet.Encoding_PLAIN,
			expected: false,
		},
		{
			enc:      parquet.Encoding_RLE,
			expected: false,
		},
		{
			enc:      parquet.Encoding_DELTA_BINARY_PACKED,
			expected: false,
		},
		{
			enc:      parquet.Encoding_DELTA_BYTE_ARRAY,
			expected: false,
		},
		{
			enc:      parquet.Encoding_DELTA_LENGTH_BYTE_ARRAY,
			expected: false,
		},
	}

	for _, c := range cases {
		page := NewDataPage()
		page.encoding = c.enc

		if page.UseDictionaryEncoding() != c.expected {
			t.Errorf("Expected %v, but actual %v", c.expected, page.UseDictionaryEncoding())
		}
	}
}
