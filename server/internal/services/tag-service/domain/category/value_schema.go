package category

type ValuePropertySchema struct {
	Type string
}

func (prop ValuePropertySchema) Validate() bool {

}

type RangePropertySchema struct {
	ValuePropertySchema
	Max int64
	Min int64
}

func aa() {
	schema := RangePropertySchema{
		ValuePropertySchema: ValuePropertySchema{},
		Max:                 0,
		Min:                 0,
	}

	schema.Validate()
}
