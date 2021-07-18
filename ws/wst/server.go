package wst

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"

	"github.com/cloudflare/circl/dh/sidh"
	"github.com/fasthttp/websocket"
	"github.com/snowmerak/wrapper/ws"
	"github.com/snowmerak/wrapper/ws/wserr"
	"github.com/valyala/fasthttp"
)

func Handshake(conn *websocket.Conn, writer io.Writer) (cipher.AEAD, error) {
	privateKey := sidh.NewPrivateKey(sidh.Fp751, sidh.KeyVariantSidhA)
	privateKey.Generate(rand.Reader)
	publicKey := sidh.NewPublicKey(sidh.Fp751, sidh.KeyVariantSidhA)
	privateKey.GeneratePublicKey(publicKey)

	out := make([]byte, publicKey.Size())
	publicKey.Export(out)

	if err := conn.WriteMessage(websocket.BinaryMessage, out); err != nil {
		return nil, err
	}

	var another []byte
	if typ, msg, err := conn.ReadMessage(); err != nil {
		return nil, err
	} else {
		if typ != websocket.BinaryMessage {
			return nil, wserr.NotBinaryMessage()
		}
		another = msg
	}

	publicKey = sidh.NewPublicKey(sidh.Fp751, sidh.KeyVariantSidhB)
	if err := publicKey.Import(another); err != nil {
		return nil, err
	}

	secret := make([]byte, privateKey.SharedSecretSize())
	privateKey.DeriveSecret(secret, publicKey)

	hashed := sha256.Sum256(secret)
	block, err := aes.NewCipher(hashed[:])
	if err != nil {
		return nil, err
	}

	aesd, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if _, err := writer.Write(secret); err != nil {
		return nil, err
	}

	return aesd, nil
}

func SendBinaryMessage(conn *websocket.Conn, aead cipher.AEAD, secret, data []byte) error {
	encrypted := aead.Seal(nil, secret[:aead.NonceSize()], data, nil)
	return conn.WriteMessage(websocket.BinaryMessage, encrypted)
}

func ReceiveBinaryMessage(conn *websocket.Conn, aead cipher.AEAD, secret []byte, writer io.Writer) error {
	typ, data, err := conn.ReadMessage()
	if err != nil {
		return err
	}
	if typ != websocket.BinaryMessage {
		return wserr.NotBinaryMessage()
	}
	plain, err := aead.Open(nil, secret[:aead.NonceSize()], data, nil)
	if err != nil {
		return err
	}
	_, err = writer.Write(plain)
	return err
}

func Listen(ctx *fasthttp.RequestCtx, handler func(conn *websocket.Conn)) error {
	return ws.Upgrader.Upgrade(ctx, handler)
}
