package crypto

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SignVerify(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		pk, sk := GenerateKeypair()
		fmt.Println("pk：", pk)
		fmt.Println("sk：", sk)
		pkstr := "0x" + CompressedPubkeyHex(pk)
		fmt.Println("pkstr：", pkstr)
		parsed_pubkey, err := StringToPubkey(pkstr)
		assert.Nil(t, err)
		fmt.Println("parsed_pubkey", parsed_pubkey)
	})
	t.Run("success", func(t *testing.T) {
		payload := "test123"
		pk, sk := GenerateKeypair()
		fmt.Println("pk：", pk)
		fmt.Println("sk：", sk)
		signature, err := SignPersonal([]byte(payload), sk)
		assert.Nil(t, err)

		err = ValidatePersonalSignature(payload, signature, pk)
		assert.Nil(t, err)
	})

	t.Run("fail if pubkey mismatch", func(t *testing.T) {
		payload := "test123"
		_, sk := GenerateKeypair()
		signature, _ := SignPersonal([]byte(payload), sk)

		new_pk, _ := GenerateKeypair()
		err := ValidatePersonalSignature(payload, signature, new_pk)
		assert.NotNil(t, err)
	})

	t.Run("fail if payload mismatch", func(t *testing.T) {
		payload := "test123"
		pk, sk := GenerateKeypair()
		signature, _ := SignPersonal([]byte(payload), sk)

		err := ValidatePersonalSignature("foobar", signature, pk)
		assert.NotNil(t, err)
	})

	t.Run("fail if signature mismatch", func(t *testing.T) {
		pk, sk := GenerateKeypair()
		signature, _ := SignPersonal([]byte("foobar"), sk)

		err := ValidatePersonalSignature("test123", signature, pk)
		assert.NotNil(t, err)
	})
}
