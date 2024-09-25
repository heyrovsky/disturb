/*
Package keys defines interfaces for cryptographic signatures and key operations.

This package provides three main interfaces:

1. Signature: Represents a cryptographic signature with methods for string and JSON representation.

 2. PublicKey: Represents a public key with methods for signature verification,
    string representation, and JSON formatting.

 3. PrivateKey: Represents a private key with methods for signing data and retrieving
    the associated public key.

These interfaces allow for flexible implementation of various cryptographic algorithms
while providing a consistent API for key and signature operations.

Usage of these interfaces enables easier testing, mocking, and switching between
different cryptographic implementations as needed.
*/
package keys

// Signature represents a cryptographic signature.
type Signature interface {
	// String returns the hexadecimal representation of the signature.
	String() (string, error)

	// Json returns the signature in JSON format.
	Json() ([]byte, error)
}

// PublicKey represents a public key for cryptographic operations.
type PublicKey interface {
	// Verify checks if the signature is valid for the given data using this public key.
	Verify(data []byte, signature Signature) (bool, error)

	// String returns the hexadecimal representation of the public key.
	String() (string, error)

	// Json returns the public key in JSON format.
	Json() ([]byte, error)
}

// PrivateKey represents a private key for cryptographic operations.
type PrivateKey interface {
	// Sign creates a cryptographic signature for the given data.
	Sign(data []byte) (Signature, error)

	// Public returns the associated public key.
	Public() (PublicKey, error)
}
