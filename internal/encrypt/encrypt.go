package encrypt

import (
	"fmt"
	"strings"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
)

func decodePublicKey(rawKey string) (*packet.PublicKey, error) {
	// open ascii armored public key
	block, err := armor.Decode(strings.NewReader(rawKey))
	if err != nil {
		return nil, err
	}

	if block.Type != openpgp.PublicKeyType {
		return nil, fmt.Errorf("Invalid public key file")
	}

	reader := packet.NewReader(block.Body)
	pkt, err := reader.Next()
	if err != nil {
		return nil, err
	}
	key, ok := pkt.(*packet.PublicKey)
	if !ok {
		return nil, fmt.Errorf("Invalid public key file")
	}
	return key, nil
}

type PublicKeySource interface {
	Fetch() (*packet.PublicKey, error)
}
