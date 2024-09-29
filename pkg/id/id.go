/*
Package id provides the core functionality for representing and managing peer identities
within a distributed system. This package enables the creation, serialization, and deserialization
of peer identities, encapsulating important details such as cryptographic public keys, IP addresses,
and port numbers.


- `ID`: Encapsulates the public key, IP address, and port of a peer in the network.
- `Marshal`: Serializes the ID to a byte slice.
- `UnmarshalID`: Deserializes a byte slice into an ID struct.
- `String`: Returns the JSON string representation of the ID for easy transmission and logging.

*/

package id

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"net"
	"strconv"

	"github.com/heyrovsky/disturbdb/pkg/keys"
	"github.com/heyrovsky/disturbdb/utils"
)

// StringIDRep represents the JSON-friendly format of an ID.

type StringIDRep struct {
	PublicKey string `json:"public_key"`
	Address   string `json:"address"`
}

//ID represents a peer's identity in the network.

type ID struct {
	ID   keys.PublicKey // Public key of the peer, used for identity and cryptographic verification.
	Host net.IP         // IP address of the peer (can be IPv4 or IPv6).
	Port uint16         // Port number that the peer is listening on.
}

// NewID constructs and returns a new ID struct.

func NewID(id keys.PublicKey, host net.IP, port uint16) ID {
	// Construct and return a new ID instance.
	return ID{
		ID:   id,
		Host: host,
		Port: port,
	}
}

// Address returns a string representation of the peer’s address, combining the IP address and port.
func (i *ID) Address() string {
	// Normalize the IP address and combine it with the port number.
	return net.JoinHostPort(utils.NormalizeIP(i.Host), strconv.FormatUint(uint64(i.Port), 10))
}

// Size returns the total size of the ID struct in bytes.
func (i *ID) Size() int {
	// Calculate and return the size of the ID.
	return i.ID.Size() + net.IPv4len + 2
}

// String returns the JSON string representation of the ID.
func (e ID) String() (string, error) {
	// Convert the public key to a hexadecimal string.
	hexkey, _ := e.ID.String()

	// Create a StringIDRep struct to hold the public key and address.
	rep := StringIDRep{
		PublicKey: hexkey,
		Address:   e.Address(),
	}

	// Serialize the struct into JSON format.
	jsonBytes, err := json.Marshal(rep)
	if err != nil {
		return "", err
	}

	// Return the JSON string.
	return string(jsonBytes), nil
}

// Marshal serializes the ID struct into a byte slice for network transmission.
func (e ID) Marshal() []byte {
	// Serialize the public key, IP address, and port number into byte slices.
	pubKeyBytes := e.ID.Bytes()  // Public key bytes
	hostBytes := e.Host.To16()   // Convert IP address to 16-byte format
	portBytes := make([]byte, 2) // Create 2 bytes for the port number
	binary.BigEndian.PutUint16(portBytes, e.Port)

	// Concatenate the byte slices and return the result.
	return append(append(pubKeyBytes, hostBytes...), portBytes...)
}

// UnmarshalID deserializes a byte slice into an ID struct.
func UnmarshalID(buf []byte, pubkey keys.PublicKey) (ID, error) {
	// Calculate the expected size of the buffer based on the size of the public key, IPv6 length, and port size.
	expectedSize := pubkey.Size() + net.IPv6len + 2
	if len(buf) < expectedSize {
		// Return an error if the buffer length is less than expected.
		return ID{}, io.ErrUnexpectedEOF
	}

	// Extract the public key from the first part of the buffer.
	pubKey, err := pubkey.UnmarshalPublicKeyFromByte(buf[:pubkey.Size()])
	if err != nil {
		return ID{}, err
	}

	// Extract the IP address from the buffer.
	hostBytes := buf[pubkey.Size() : pubkey.Size()+net.IPv6len]
	host := net.IP(hostBytes)

	// Extract the port number from the last 2 bytes of the buffer.
	port := binary.BigEndian.Uint16(buf[pubkey.Size()+net.IPv6len:])

	// Return the reconstructed ID struct.
	return ID{
		ID:   pubKey,
		Host: host,
		Port: port,
	}, nil
}
