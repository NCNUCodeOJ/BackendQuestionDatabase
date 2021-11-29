package views

type sampleTemplate struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

type testCaseTemplate struct {
	StrippedOutputMD5 string `json:"stripped_output_md5"`
	OutputSize        int    `json:"output_size"`
	OutputMD5         string `json:"output_md5"`
	InputName         string `json:"input_name"`
	InputSize         int    `json:"input_size"`
	OutputName        string `json:"output_name"`
}

type problemAPIRequest struct {
	ProblemName       *string           `json:"problem_name"`
	Description       *string           `json:"description"`
	InputDescription  *string           `json:"input_description"`
	OutputDescription *string           `json:"output_description"`
	MemoryLimit       *uint             `json:"memory_limit"`
	CPUTime           *uint             `json:"cpu_time"`
	Layer             *uint8            `json:"layer"`
	Sample            []*sampleTemplate `json:"sample"`
	TagsList          []string          `json:"tags_list"`
	ProgramName       *string           `json:"program_name"`
}

type submissionAPIRequest struct {
	SourceCode *string `json:"source_code"`
	Language   *string `json:"language"`
}
