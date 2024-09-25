/*
Package keys provides Ed25519 cryptographic key management and signature operations.

To future me reading this code: You wrote it after a few shots of whiskey.
Only God knows what you did. Good luck!

Key features:
- Ed25519 key pair generation and management
- Signing and signature verification
- Conversion between keys/signatures and hex/JSON formats
- Loading private keys from hex strings

Types:
- Ed25519PrivateKey, Ed25519PublicKey, Ed25519Signature
*/

package keys

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"unsafe"
)

// Ed25519PrivateKey, Ed25519PublicKey, and Ed25519Signature define the key and signature types for the ed25519 algorithm.
type (
	Ed25519PrivateKey [ed25519.PrivateKeySize]byte
	Ed25519PublicKey  [ed25519.PublicKeySize]byte
	Ed25519Signature  [ed25519.SignatureSize]byte
)

// Json returns the JSON-encoded form of the signature.
func (e Ed25519Signature) Json() ([]byte, error) {
	sign, err := e.String()
	if err != nil {
		return nil, fmt.Errorf("failed to convert signature to string: %w", err)
	}
	return json.Marshal(sign)
}

// String returns the hexadecimal string representation of the signature.
func (e Ed25519Signature) String() (string, error) {
	return hex.EncodeToString(e[:]), nil
}

// Json returns the JSON-encoded form of the public key.
func (e Ed25519PublicKey) Json() ([]byte, error) {
	key, err := e.String()
	if err != nil {
		return nil, fmt.Errorf("failed to convert public key to string: %w", err)
	}
	return json.Marshal(key)
}

// String returns the hexadecimal string representation of the public key.
func (e Ed25519PublicKey) String() (string, error) {
	return hex.EncodeToString(e[:]), nil
}

// Verify checks if the provided signature is valid for the data using the public key.
func (e Ed25519PublicKey) Verify(data []byte, signature Signature) (bool, error) {
	ed25519Sig, ok := signature.(Ed25519Signature)
	if !ok {
		return false, errors.New("invalid signature type")
	}
	// Verify the signature using ed25519 algorithm.
	return ed25519.Verify(e[:], data, ed25519Sig[:]), nil
}

// Public returns the public key associated with the private key.
func (e Ed25519PrivateKey) Public() (PublicKey, error) {
	var publicKey Ed25519PublicKey
	pub := (ed25519.PrivateKey)(e[:]).Public().(ed25519.PublicKey)
	copy(publicKey[:], pub)
	return publicKey, nil
}

// Sign signs the provided data and returns the signature.
func (e Ed25519PrivateKey) Sign(data []byte) (Signature, error) {
	signatureBytes := ed25519.Sign(e[:], data)
	return ed25519UnmarshalSignature(signatureBytes), nil
}

// ed25519UnmarshalSignature converts a byte slice into an Ed25519Signature.
func ed25519UnmarshalSignature(data []byte) Ed25519Signature {
	_ = data[ed25519.SignatureSize-1] // ensure the data length is at least SignatureSize
	return *(*Ed25519Signature)(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&data)).Data))
}

// GenerateKeys generates a new ed25519 public/private key pair.
func GenerateKeys(rand io.Reader) (PublicKey, PrivateKey, error) {
	pub, priv, err := ed25519.GenerateKey(rand)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate keys: %w", err)
	}
	var pubKey Ed25519PublicKey
	var privKey Ed25519PrivateKey
	copy(pubKey[:], pub)
	copy(privKey[:], priv)
	return pubKey, privKey, nil
}

// LoadKeysFromHex loads a private key from its hexadecimal string representation.
func LoadKeysFromHex(secretHex string) (PrivateKey, error) {
	secret, err := hex.DecodeString(secretHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key from hex: %w", err)
	}

	if len(secret) != ed25519.PrivateKeySize {
		return nil, fmt.Errorf("invalid private key size: got %d bytes, expected %d: %w",
			len(secret), ed25519.PrivateKeySize, io.ErrUnexpectedEOF,
		)
	}

	var privateKey Ed25519PrivateKey
	copy(privateKey[:], secret)

	return &privateKey, nil
}
