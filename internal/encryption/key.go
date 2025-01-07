package encryption

import (
	"context"
	"encoding/json"
)

// Key combines the Encrypter & Decrypter interfaces.
type Key interface {
	Encrypter
	Decrypter

	// Version returns info containing to concretely identify
	// the underlying key, eg: key type, name, & version.
	Version(ctx context.Context) (KeyVersion, error)
}

type KeyVersion struct {
	// TODO: generate this as an enum from JSONSchema
	Type    string
	Name    string
	Version string
}

func (v KeyVersion) JSON() string {
	buf, _ := json.Marshal(v)
	return string(buf)
}

// Encrypter is anything that can encrypt a value
type Encrypter interface {
	Encrypt(ctx context.Context, value []byte) ([]byte, error)
}

// Decrypter is anything that can decrypt a value
type Decrypter interface {
	Decrypt(ctx context.Context, cipherText []byte) (*Secret, error)
}

func NewSecret(v string) Secret {
	return Secret{
		value: v,
	}
}

// Secret is a utility type to make it harder to accidentally leak secret
// values in logs. The actual value is unexported inside a struct, making
// harder to leak via reflection, the string value is only ever returned
// on explicit Secret() calls, meaning we can statically analyse secret
// usage and statically detect leaks.
type Secret struct {
	value string
}

// String implements stringer, obfuscating the value
func (s Secret) String() string {
	return "********"
}

// Secret returns the unobfuscated value
func (s Secret) Secret() string {
	return s.value
}

// MarshalJSON overrides the default JSON marshaling implementation, obfuscating
// the value in any marshaled JSON
func (s Secret) MarshalJSON() ([]byte, error) {
	return []byte(s.String()), nil
}
