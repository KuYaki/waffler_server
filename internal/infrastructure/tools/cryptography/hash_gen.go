package cryptography

import (
	"fmt"
	"github.com/google/uuid"
)

type Hasher interface {
	GenHashString(in []byte, kind int) string
	GenHash(in []byte, kind int) []byte
}

const (
	UUID = iota + 1
)

type Hash struct {
	uuid Generator
}

func NewHash(uuid Generator) *Hash {
	return &Hash{uuid: uuid}
}

func (h *Hash) GenHashString(in []byte, kind int) string {
	switch kind {
	case UUID:
		return h.uuid.String()
	}

	return ""
}

func (h *Hash) GenHash(in []byte, kind int) []byte {
	switch kind {
	case UUID:
		return h.uuid.Bytes()
	}

	return nil
}

type Generator interface {
	fmt.Stringer
	Bytes() []byte
}

type UUIDGenerator struct {
	uuid uuid.UUID
}

func NewUUIDGenerator() *UUIDGenerator {
	uuID, _ := uuid.NewUUID()
	return &UUIDGenerator{
		uuid: uuID,
	}
}

func (u *UUIDGenerator) String() string {
	return u.uuid.String()
}

func (u *UUIDGenerator) Bytes() []byte {
	return func() []byte {
		uuidBinary, _ := u.uuid.MarshalBinary()

		return uuidBinary
	}()
}
