package wst

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"

	"github.com/cloudflare/circl/dh/sidh"
	"github.com/fasthttp/websocket"
	"github.com/snowmerak/wrapper/ws/wserr"
)

func Connect(url string) (func(data []byte) error, func(writer io.Writer) error, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, nil, err
	}

	privateKey := sidh.NewPrivateKey(sidh.Fp751, sidh.KeyVariantSidhB)
	privateKey.Generate(rand.Reader)
	publicKey := sidh.NewPublicKey(sidh.Fp751, sidh.KeyVariantSidhB)
	privateKey.GeneratePublicKey(publicKey)

	var another []byte
	if typ, msg, err := conn.ReadMessage(); err != nil {
		return nil, nil, err
	} else {
		if typ != websocket.BinaryMessage {
			return nil, nil, wserr.NotBinaryMessage()
		}
		another = msg
	}

	out := make([]byte, publicKey.Size())
	publicKey.Export(out)

	if err := conn.WriteMessage(websocket.BinaryMessage, out); err != nil {
		return nil, nil, err
	}

	publicKey = sidh.NewPublicKey(sidh.Fp751, sidh.KeyVariantSidhA)
	if err := publicKey.Import(another); err != nil {
		return nil, nil, err
	}

	secret := make([]byte, privateKey.SharedSecretSize())
	privateKey.DeriveSecret(secret, publicKey)

	hashed := sha256.Sum256(secret)
	block, err := aes.NewCipher(hashed[:])
	if err != nil {
		return nil, nil, err
	}

	aesd, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	sender := func(data []byte) error {
		msg := aesd.Seal(nil, secret[:aesd.NonceSize()], data, nil)
		return conn.WriteMessage(websocket.BinaryMessage, msg)
	}

	receiver := func(writer io.Writer) error {
		typ, msg, err := conn.ReadMessage()
		if err != nil {
			return err
		}
		if typ != websocket.BinaryMessage {
			return wserr.NotBinaryMessage()
		}
		data, err := aesd.Open(nil, secret[:aesd.NonceSize()], msg, nil)
		if err != nil {
			return err
		}
		if _, err := writer.Write(data); err != nil {
			return err
		}
		return nil
	}

	return sender, receiver, nil
}
