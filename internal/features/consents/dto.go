package consents

type LatestDocumentRequest struct {
	Type string `uri:"type" binding:"required,oneof=terms privacy"`
}

type SaveConsentRequest struct {
	DocID int `json:"doc_id" binding:"required"`
}
