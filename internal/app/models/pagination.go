package models

// ListContactsRequest represents the paginated list request parameters
type ListContactsRequest struct {
	Query  string `form:"q"`
	Page   int    `form:"page,default=1"`
	Limit  int    `form:"limit,default=10"`
	Offset int    `form:"-"`
}
