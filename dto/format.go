package dto

// TODO: Might change MediaType to enum type(video/audio ONLY)
// TODO: WEBHOOKS?
type Format struct {
	MediaType        string `json:"mediaType"`
	RequiredFileType string `json:"requiredFileType"`
}
