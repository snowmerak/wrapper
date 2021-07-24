package auth

import (
	"crypto/ecdsa"
	"crypto/x509"
	"errors"
	"math/big"

	"google.golang.org/protobuf/proto"
)

func SerializePrivateKey(pr *ecdsa.PrivateKey) ([]byte, error) {
	return x509.MarshalECPrivateKey(pr)
}

func DeserializePrivateKey(bs []byte) (*ecdsa.PrivateKey, error) {
	return x509.ParseECPrivateKey(bs)
}

func SerializePublicKey(pb *ecdsa.PublicKey) ([]byte, error) {
	return x509.MarshalPKIXPublicKey(pb)
}

func DeserializePublicKey(bs []byte) (*ecdsa.PublicKey, error) {
	pb, err := x509.ParsePKIXPublicKey(bs)
	if err != nil {
		return nil, err
	}
	rs, ok := pb.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("converted value is not ecdsa public key pointer")
	}
	return rs, nil
}

func SerializeSignature(r *big.Int, s *big.Int) ([]byte, error) {
	sig := &Signature{
		R: r.Text(16),
		S: s.Text(16),
	}
	return proto.Marshal(sig)
}

func DeserializeSignature(data []byte) (*big.Int, *big.Int, error) {
	sig := new(Signature)
	if err := proto.Unmarshal(data, sig); err != nil {
		return nil, nil, err
	}
	r, ok := new(big.Int).SetString(sig.R, 16)
	if !ok {
		return nil, nil, errors.New("bad value of R")
	}
	s, ok := new(big.Int).SetString(sig.S, 16)
	if !ok {
		return nil, nil, errors.New("bad value of S")
	}
	return r, s, nil
}
