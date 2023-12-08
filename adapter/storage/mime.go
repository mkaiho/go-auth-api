package storage

type MimeType string

const (
	MimeTypeUnsupported MimeType = ""
	MimeTypeJSON        MimeType = "application/json"
	MimeTypeJWK         MimeType = "application/jwk+json"
	MimeTypePem         MimeType = "application/x-pem-file"
)

func (m MimeType) String() string {
	return string(m)
}
