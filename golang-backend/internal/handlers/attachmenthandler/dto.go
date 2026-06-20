package attachmenthandler

import "github.com/aportela/doneo/internal/handlers/userhandler"

type AttachmentResponse struct {
	ID          string                       `json:"id"`
	CreatedBy   userhandler.UserBaseResponse `json:"createdBy"`
	CreatedAt   int64                        `json:"createdAt"`
	Filename    string                       `json:"name"`
	ContentType string                       `json:"contentType"`
	Size        uint32                       `json:"size"`
}

type searchResponse struct {
	Attachments []AttachmentResponse `json:"attachments"`
}
