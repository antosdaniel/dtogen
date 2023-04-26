package input

type StructWithBaseTypes struct {
	Bool    bool
	PtrBool *bool

	String    string
	PtrString *string

	Int      int
	PtrInt   *int
	Int8     int8
	PtrInt8  *int8
	Int16    int16
	PtrInt16 *int16
	Int32    int32
	PtrInt32 *int32
	Int64    int64
	PtrInt64 *int64

	Uint      uint
	PtrUint   *uint
	Uint8     uint8
	PtrUint8  *uint8
	Uint16    uint16
	PtrUint16 *uint16
	Uint32    uint32
	PtrUint32 *uint32
	Uint64    uint64
	PtrUint64 *uint64

	Float32    float32
	PtrFloat32 *float32
	Float64    float64
	PtrFloat64 *float64
}
