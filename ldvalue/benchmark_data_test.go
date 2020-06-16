package ldvalue

var (
	benchmarkStringValue        = "value"
	benchmarkStringPointer      = &benchmarkStringValue
	benchmarkOptStringWithValue = NewOptionalString(benchmarkStringValue)

	benchmarkSerializeNullValue   = Null()
	benchmarkSerializeBoolValue   = Bool(true)
	benchmarkSerializeIntValue    = Int(1000)
	benchmarkSerializeFloatValue  = Float64(1000.5)
	benchmarkSerializeStringValue = String("value")
	benchmarkSerializeArrayValue  = ArrayOf(String("a"), String("b"), String("c"))
	benchmarkSerializeObjectValue = ObjectBuild().Set("a", Int(1)).Set("b", Int(2)).Set("c", Int(3)).Build()

	benchmarkOptStringResult OptionalString
	benchmarkStringResult    string
	benchmarkValueResult     Value
	benchmarkBoolResult      bool
	benchmarkIntResult       int
	benchmarkFloat64Result   float64
	benchmarkJSONResult      []byte
)
