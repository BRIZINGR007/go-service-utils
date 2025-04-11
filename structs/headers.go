package structs

type Headers struct {
	CorrelationId string `json:"correlationid" binding:"required"`
	UserId        string `json:"userId" binding:"required"`
	Authorization string `json:"authorization" binding:"required"`
}
