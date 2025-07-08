package wrappers

type GRPCError struct {
	Code int    `json:"code"`
	Err  string `json:"error"`
}
