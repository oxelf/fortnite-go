package fortnitego

// /EPIC ERRORS
type Error struct {
	ErrorMessage           string
	EpicErrorCode          string   `json:"errorCode"`
	EpicErrorMessage       string   `json:"errorMessage"`
	EpicMessageVars        []string `json:"messageVars"`
	EpicNumericErrorCode   int      `json:"numericErrorCode"`
	EpicOriginatingService string   `json:"originatingService"`
	EpicIntent             string   `json:"intent"`
}
