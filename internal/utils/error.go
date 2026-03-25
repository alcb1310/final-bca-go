package utils

type BcaError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (b *BcaError) Error() string {
	return b.Message
}
