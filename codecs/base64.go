package codecs

import "encoding/base64"

type Base64 struct {
	encoder *base64.Encoding
}

func (c *Codec) Base64() *Base64 {
	return &Base64{
		encoder: base64.StdEncoding,
	}
}

func (b *Base64) Encode(s string) string {
	return b.encoder.EncodeToString([]byte(s))
}

func (b *Base64) Decode(s string) ([]byte, error) {
	return b.encoder.DecodeString(s)
}

// DecodeBytes decodes the given base64 encoded byte slice and returns the decoded byte slice.
func (b *Base64) DecodeBytes(b64 []byte) ([]byte, error) {
	db := make([]byte, base64.StdEncoding.DecodedLen(len(b64)))
	n, err := base64.StdEncoding.Decode(db, b64)
	if err != nil {
		return nil, err
	}
	return db[:n], nil
}

// EncodeBytes takes a byte slice and returns its base64 encoded byte slice.
func (b *Base64) EncodeBytes(s []byte) []byte {
	eb := make([]byte, base64.StdEncoding.EncodedLen(len(s)))
	base64.StdEncoding.Encode(eb, s)
	return eb
}
