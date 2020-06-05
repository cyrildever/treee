package response

// PostLeaf ...
type PostLeaf struct {
	Code 	int		`json:"code"`
	Result	string	`json:"result,omitempty"`
	Error	string	`json:"error,omitempty"`
}
