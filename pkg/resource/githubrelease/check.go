package githubrelease

// CheckRequest contains the information needed to check if a new github release is available
type CheckRequest struct {
	Version Version `json:"version"`
	Source  Source  `json:"source"`
}

// CheckResponse defines the response from the github release check operation
type CheckResponse struct {
	Versions []*Version `json:"versions"`
}
