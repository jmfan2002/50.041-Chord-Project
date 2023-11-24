package structs

type SetValueReqBody struct {
	Key   string
	Value string
	Nonce string
	// Optional
	PreviousNodeHash string
}

func NewSetValueReqBody(key string, value string, nonce string) SetValueReqBody {
	return SetValueReqBody{
		Key:              key,
		Value:            value,
		Nonce:            nonce,
		PreviousNodeHash: "nil",
	}
}
