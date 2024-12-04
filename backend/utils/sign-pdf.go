package utils

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"mime/multipart"
	"io"
	"fmt"
	"strings"
)

func GenerateEncryptedHash(msg []byte, publicKey *rsa.PublicKey) (string, error) {
	msgHash := sha256.Sum256(msg)

	encryptedMsgHash, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, msgHash[:], nil)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(encryptedMsgHash), nil
}

func VerifyDigitalSignature(signature string, fileSignature string) (bool, error) {

	return signature == fileSignature, nil
}

func AppendDigitalSignature(signature string, fileHeader *multipart.FileHeader) ([]byte, error) {
	// fileHeader.Header.Add("Signature", signature)

	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	buf.WriteString("\n-----SIGNATURE-----\n")
	buf.WriteString(signature)

	return buf.Bytes(), nil
}

func VerifyEmbeddedSignature(fileHeader *multipart.FileHeader, username string, publicKey string) (bool, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return false, err
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		return false, err
	}

	fileStr := string(fileData)

	parts := strings.Split(fileStr, "\n-----SIGNATURE-----\n")
	if len(parts) != 2 {
		return false, fmt.Errorf("signature not found")
	}

	signature := strings.TrimSpace(parts[1])
	fmt.Println("Signature: ", signature)

	signatureParts := strings.SplitN(signature, "\n", 2)
	fmt.Println("Signature Parts: ", signatureParts)
	if len(signatureParts) != 2 {
		return false, fmt.Errorf("invalid signature format")
	}

	usernameSignature := strings.TrimSpace(signatureParts[0])
	publicKeySignature := strings.TrimSpace(signatureParts[1])

	if username != usernameSignature {
		return false, fmt.Errorf("username does not match")
	}

	if strings.TrimSpace(publicKey) != publicKeySignature {
		return false, fmt.Errorf("public key does not match")
	}

	return true, nil
}