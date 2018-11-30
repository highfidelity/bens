// Copyright 2018 High Fidelity, Inc.
//
// Distributed under the Apache License, Version 2.0.
// See the accompanying file LICENSE or http://www.apache.org/licenses/LICENSE-2.0.html

package key

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var passFileContent = []byte("3gPyttqJ3luMmeok/npIiF+x/k61+B2r8gPZhUmvpFfk\n")

var publicKeyFileContent = []byte(`
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA6DY0TwTvEZt5j6vh/YoA
YL5OajVZ/J4ocnPLoI7Hol0ThjZoClZDjfe7Os3kXnrunm6AEcMOEizP/XSjRBBC
c12yilWHnb/pWA3Ko7Mu4xq8BKXBdFoUr0CGpI4jO4vFeFEk3dfXhSpqaK78AltW
VTM2AfgxRVpEZEw4+R5sqT9rmTsRYwSXbK6ImlJD58x/owvwFLBnPkCTuguEi5p2
L9yeoSC3r7bvsePfcxqGrgxDYi7b8+Ugx5F7im+pPkbRAkvVgrzjgHP/aS8w76MS
6rApqBp8C+CT4lLJs7a7CqMHvxml5+XmjMc8OLC2hdNZGFuKjCcvCffNjrFGzEHq
RQIDAQAB
-----END PUBLIC KEY-----
`)

var privateKeyFileContent = []byte(`
-----BEGIN RSA PRIVATE KEY-----
Proc-Type: 4,ENCRYPTED
DEK-Info: AES-128-CBC,154C7B300BCFDC2B278F682E58DCFED5

dTHbTam/RZMH/X7KM00Jq3+a/aP1SoaZccjlUM2x0vgjnQ0To7gvVhtpc4yOtqZx
PfsU7Zk+cptA4y9QOmBfyZBvPWwXAWLcMXqdN6+VuC4SyOl7wgXlddw4wpLy5JWt
mdXHP8d1gEMk2D6NZcIWUKvVPB2kZ5AFnpRuuz+dC6R1FkMUzPKbbHPim5UJWDpW
uTTARiPCO8+AqdtXTkt8qhjTVl6Vdg7oDGZyDRpMHuY5Pc0brgSMX18dgUQMwDtz
amhSL+mNc93wIvM3uRwSc4vHcGxtjCved/H95q/uCaT3A6wtvVSX29+v47Qt9y7a
20eNksdxZhrtCr0gIlfiVXJnNMPaw5vltoNKbkXzGzYF+1M5QogXynnjrw76r1Xs
N4XtuB3vQQPtBh4uVjGvJtNIxTiQVj6vd1m7jzTogOLlqfkjuu4b81RWFma8jwee
2HG4W9dqHl8b+DW1PVutemWEPGEeaSoKJcPP9+2orLs4dZiNh3D98Jxb+JCJnKcm
AkR2CWtRpwT/yUwVZmdgS3ZmGW6/uWrFeQ0JZIEq4t96JQNzD3XITjjB0lyCuocK
rxhWs2dm95S8x+LkmWYxlqbVrEW3RdDr/GJwzH/niaMNgoN/R7mpmzJ/zFcpd1ih
FbHGC3MD4BL0CNN6WMIIaL3NRKFMcjDCd2GOmIFI8ZY6r05njBA+KaknOTghsAW/
OUyYSmuwkew4696Ot1xvo83JATEbg3zJ+lqRTy4IwAIy1DO2QqK8xAbdBByCyoDM
8mQ7k8xhjTEmEN/YuNWzNBZXJkTx3atJC97X47OldmrEvMtpivXLHGKD8RwUXtOg
6dXkjB2xaLvcMIecIRfC8AxYefsAk4eFE59dLt9eQhFgqrhPsCOHJFkXPmZOd7kW
+C0fUV6ACzOCqVQ3NeDLWHzmn330RGaqyk67kBnSrUaD0QflX0fPCZ31cUf/5GDm
Th6Tid5evr981yecqf9e/MWe4xGTDcLQ8zuoX83Ve1yOvfqkhrp8nom/t+qmh57X
k6xw8X/eFXD1ND3dqjZIDhx+u7gNBVoItkNePR5+LOP3LzeEq5hIOB35trJEnXl3
47NlsK/3fwskFnGLl98xVk1KYOSDZ7IOOOLmZ+tovwOb8PZgWEsVR5kcR+xnL07h
7rtiRqq3AG16SWVOBmzOO9ruI5GiDriQQ4628qCy706dKMF0lcXZ8sPhtYJCtGv4
ZjD1n8rbX6TessIyD/axTI8s2BLrvuY2XnqgOmpw27yKUPZw/aShXM5RmP+j6k9G
e17DoAABrWqbpHzmNG0zg/5cDB7lkbty01ZLdOd2Qm/JB884eoCid8gh4kkdY1nb
+XmLT8QizwZKf0WN3GEoEnzu1GX4XFSduLf/w1/5R0pOMGzo1Fs1QB6BMHszk9Vb
V0CCkh1pmG/ApdO2RebN289AA5Z6kLZJ1cTpUb0Y1lU/3EZ9SZifzmxgUkQrY1hX
558HO70/ImaHUyfkzTmJMRHR1MhEzlUeouM7qmA098wchfNyFDp84ydVG6rQX6nd
eKzJhr8b1hgAAj+GXDDe3ZwH2hCCfliFu3GhNjJc6xG87HgyT3MsDSwU08UBUSlM
-----END RSA PRIVATE KEY-----
`)

var fooBarEncryptedContent = string("Bn0YKTm/pF5OdkoNOKy7fGuLXXmwF2Nc3pIYp24yTMacfvP7vcsmAli2sxi+VJ25HrrrlTmQbGgUiEUV4BneIpByhsrwOv/DhXti9fqF0zwSFawdvS4qvk3UmJBiq/k7k7rQb9UtwIXY9zb/t9hGIeOJkfMzosfgubkgs/ZhvBPnUOiqnDqWf9NLNjH6GwTirfcVJH2ih3gfHGLpW3ehi19VNPRIOeFVDgWpHhIjYzGol8e31bu9M+8/5pHm1bXQ8JvMUhYGFkY4I++/+UuVB67ydVN5YrPSnQRggwOHF1klwzJvf89tnyGlKXeNCWuzJWBwEB69RDQJYxmQZ8M1Qg==")

func writeTestFiles() (string, string, string, string, error) {
	dir, err := ioutil.TempDir("", "key_test")
	if err != nil {
		return "", "", "", "", err
	}

	var failed bool = false
	passPath := filepath.Join(dir, "pass.txt")
	if ioutil.WriteFile(passPath, passFileContent, 0600) != nil {
		failed = true
	}
	priKeyPath := filepath.Join(dir, "pri.key")
	if ioutil.WriteFile(priKeyPath, privateKeyFileContent, 0600) != nil {
		failed = true
	}
	pubKeyPath := filepath.Join(dir, "pub.key")
	if ioutil.WriteFile(pubKeyPath, publicKeyFileContent, 0644) != nil {
		failed = true
	}
	if failed {
		os.RemoveAll(dir)
		return "", "", "", "", err
	}

	return dir, passPath, priKeyPath, pubKeyPath, nil
}

func TestDecrypt(t *testing.T) {
	dir, passPath, priKeyPath, pubKeyPath, err := writeTestFiles()
	if err != nil {
		t.Fatalf("couldn't write test files: %v", err)
	}
	defer os.RemoveAll(dir)

	k, err := New(passPath, priKeyPath, pubKeyPath)
	if err != nil {
		t.Fatalf("couldn't create key struct with New: %v", err)
	}
	fooBar, err := k.Decrypt(fooBarEncryptedContent)
	if err != nil {
		t.Fatalf("couldn't decrypt value with Decrypt: %v", err)
	}
	if fooBar != "foobar" {
		t.Fatalf("%s != foobar", fooBar)
	}

	return
}

func TestEncrypt(t *testing.T) {
	dir, passPath, priKeyPath, pubKeyPath, err := writeTestFiles()
	if err != nil {
		t.Fatalf("couldn't write test files: %v", err)
	}
	defer os.RemoveAll(dir)

	k, err := New(passPath, priKeyPath, pubKeyPath)
	if err != nil {
		t.Fatalf("couldn't create key struct with New: %v", err)
	}

	plainText := "foo"
	cipherText, err := k.Encrypt(plainText)
	if err != nil {
		t.Fatalf("couln't Encrypt plainText: %v", err)
	}
	if cipherText == "" {
		t.Fatal("cipher text is empty")
	}
	if cipherText == plainText {
		t.Fatal("cipher text didn't change")
	}
	_, err = base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		t.Fatalf("cipher text (%s) isn't base64 encoded: %v", cipherText, err)
	}

	decryptedPlainText, err := k.Decrypt(cipherText)
	if err != nil {
		t.Fatalf("couldn't decrypted cipher text: %v", err)
	}
	if decryptedPlainText != plainText {
		t.Fatalf("Decrypted text (%s) didn't match origin plain text (%s)", decryptedPlainText, plainText)
	}
}

func TestNewWithEmptyPassFile(t *testing.T) {
	dir, _, priKeyPath, pubKeyPath, err := writeTestFiles()
	if err != nil {
		t.Fatalf("couldn't write test files: %v", err)
	}
	defer os.RemoveAll(dir)

	_, err = New("", priKeyPath, pubKeyPath)
	if err != nil {
		t.Fatalf("new should handle empty pass by ignoring pass: %v", err)
	}
}
