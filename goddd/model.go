package goddd

// Summary a container to encode the summary type's metric from goddd / go-kit
type Summary struct {
	SampleCount uint64         `json:"sampleCount"`
	SampleSum   float64        `json:"sampleSum"`
	Quantile050 float64        `json:"quantile050"`
	Quantile090 float64        `json:"quantile090"`
	Quantile099 float64        `json:"quantile099"`
	Label       []*LabelStruct `json:"label"`
}

// LabelStruct smallest unit of label of metric
type LabelStruct struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
