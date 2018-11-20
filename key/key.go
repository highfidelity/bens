// Copyright 2018 High Fidelity, Inc.
//
// Distributed under the Apache License, Version 2.0.
// See the accompanying file LICENSE or http://www.apache.org/licenses/LICENSE-2.0.html

package key

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
)

func readPassFile(passPath string) ([]byte, error) {
	// Read the pass file
	pass, err := ioutil.ReadFile(passPath)
	if err != nil {
		log.Fatal(err)
	}
	pass = bytes.TrimSpace(pass)
	return pass, nil
}

func decodePemFile(path string) (pem.Block, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return pem.Block{}, err
	}
	block, rest := pem.Decode(data)
	if block == nil {
		return pem.Block{}, fmt.Errorf("PEM file contain any blocks")
	}
	if len(rest) > 0 {
		return pem.Block{}, fmt.Errorf("PEM file contained more than one block")
	}
	return *block, nil
}

func zero(data []byte) {
	go func() {
		for i := 0; i < len(data); i++ {
			data[i] = 0
		}
	}()
}

func decryptPrivateKey(path string, pass []byte) (*rsa.PrivateKey, error) {
	// Read the PEM Privdate Key
	block, err := decodePemFile(path)
	if err != nil {
		return nil, err
	}
	if block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("block must be private key")
	}

	// Decrypt the Private Key
	if !x509.IsEncryptedPEMBlock(&block) {
		zero(block.Bytes)
		return nil, fmt.Errorf("unencrypted private key used!")
	}
	der, err := x509.DecryptPEMBlock(&block, pass)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %v", err)
	}
	if err = mlock(der); err != nil {
		return nil, fmt.Errorf("couldn't lock key in memory: %v", err)
	}
	bytes := pem.EncodeToMemory(&pem.Block{Type: block.Type, Bytes: der})
	if bytes == nil {
		return nil, fmt.Errorf("couldn't encode decrypted block")
	}
	if err = mlock(bytes); err != nil {
		return nil, fmt.Errorf("couldn't lock key in memory: %v", err)
	}
	zero(der)
	p, rest := pem.Decode(bytes)
	if p == nil || len(rest) > 0 {
		return nil, fmt.Errorf("couldn't decode decrypted block")
	}
	if err = mlock(p.Bytes); err != nil {
		return nil, fmt.Errorf("couldn't lock key in memory: %v", err)
	}
	zero(bytes)
	pri, err := x509.ParsePKCS1PrivateKey(p.Bytes)
	if err != nil {
		return nil, fmt.Errorf("couldn't decode private key")
	}
	zero(p.Bytes)

	return pri, nil
}

type Key struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewWithPass(pass []byte, priKeyPath, pubKeyPath string) (Key, error) {
	if err := mlock(pass); err != nil {
		return Key{}, err
	}
	privateKey, err := decryptPrivateKey(priKeyPath, pass)
	if err != nil {
		return Key{}, err
	}
	zero(pass)

	block, err := decodePemFile(pubKeyPath)
	if err != nil {
		return Key{}, fmt.Errorf("couldn't decode public key: %v", err)
	}
	if block.Type != "PUBLIC KEY" {
		return Key{}, fmt.Errorf("not a public key: %s", pubKeyPath)
	}
	var pub *rsa.PublicKey
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return Key{}, fmt.Errorf("couldn't parse public key: %v", err)
	}
	switch publicKey := publicKey.(type) {
	case *rsa.PublicKey:
		pub = publicKey
	default:
		return Key{}, fmt.Errorf("key must be rsa public")
	}

	return Key{privateKey: privateKey, publicKey: pub}, nil
}

func New(passPath, priKeyPath, pubKeyPath string) (Key, error) {
	pass, err := readPassFile(passPath)
	if err != nil {
		return Key{}, err
	}
	return NewWithPass(pass, priKeyPath, pubKeyPath)
}

func (k Key) Encrypt(plainText string) (string, error) {
	if k.publicKey == nil {
		return "", fmt.Errorf("no public key associated")
	}
	cipherText, err := rsa.EncryptOAEP(sha1.New(), rand.Reader, k.publicKey, []byte(plainText), []byte(""))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func (k Key) Decrypt(base64CipherText string) (string, error) {
	if k.publicKey == nil {
		return "", fmt.Errorf("no private key associated")
	}
	ciphertext, err := base64.StdEncoding.DecodeString(base64CipherText)
	if err != nil {
		return "", err
	}
	ciphertext = bytes.TrimSpace(ciphertext)
	b, err := rsa.DecryptOAEP(sha1.New(), rand.Reader, k.privateKey, ciphertext, []byte(""))
	if err != nil {
		return "", err
	}
	b = bytes.TrimSpace(b)
	return string(b), err
}
