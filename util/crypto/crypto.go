package crypto

import (
	"crypto/ecdsa"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
)

// ValidatePersonalSignature checks whether (eth.personal.sign) signature,
// payload and pubkey are matched.
// Pubkey and signature should be without "0x".
func ValidatePersonalSignature(payload string, signature []byte, pubkey *ecdsa.PublicKey) bool {
	// Recover pubkey from signature
	if len(signature) != 65 {
		logrus.Warnf("Error: Signature length invalid: %d instead of 65", len(signature))
		return false
	}
	if signature[64] == 27 || signature[64] == 28 {
		signature[64] -= 27
	}

	if signature[64] != 0 && signature[64] != 1 {
		logrus.Warnf("Error: Signature Recovery ID not supported: %d", signature[64])
		return false
	}

	pubkeyRecovered, err := crypto.SigToPub(signPersonalHash([]byte(payload)), signature)
	if err != nil {
		logrus.Warnf("Error when recovering pubkey from signature: %s", err.Error())
		return false
	}

	return crypto.PubkeyToAddress(*pubkey) == crypto.PubkeyToAddress(*pubkeyRecovered)
}

// GenerateKeypair generates a keypair.
// For test purpose only.
func GenerateKeypair() (publicKey *ecdsa.PublicKey, privateKey *ecdsa.PrivateKey) {
	privateKey, _ = crypto.GenerateKey()
	publicKey = &privateKey.PublicKey
	return publicKey, privateKey
}

// SignPersonal signs a payload using given secret key.
// For test purpose only.
func SignPersonal(payload []byte, sk *ecdsa.PrivateKey) (signature []byte, err error) {
	hash := signPersonalHash(payload)
	signature, err = crypto.Sign(hash, sk)
	if err != nil {
		return nil, xerrors.Errorf("%w", err)
	}

	return signature, nil
}

func signPersonalHash(data []byte) []byte {
	messsage := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256([]byte(messsage))
}

// StringToPubkey is compatible with comressed / uncompressed pubkey
// hex, and with / without '0x' head.
func StringToPubkey(pk_str string) (*ecdsa.PublicKey, error) {
	pk_str_parsed := strings.TrimPrefix(pk_str, "0x")
	pk_str_parsed = strings.ToLower(pk_str_parsed)
	pk_bytes := common.Hex2Bytes(pk_str_parsed)
	return BytesToPubKey(pk_bytes)
}

// BytesToPubKey is compatible with comressed / uncompressed pubkey
// bytes.
func BytesToPubKey(pk_bytes []byte) (*ecdsa.PublicKey, error) {
	var result *ecdsa.PublicKey
	var err error
	if len(pk_bytes) == 33 { // compressed
		result, err = crypto.DecompressPubkey(pk_bytes)
	} else {
		result, err = crypto.UnmarshalPubkey(pk_bytes)
	}
	if err != nil {
		return nil, xerrors.Errorf("%w", err)
	}
	return result, nil
}
