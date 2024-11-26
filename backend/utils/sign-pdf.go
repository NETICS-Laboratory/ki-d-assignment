package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func SignPDFWithOpenSSL(pdfPath string, privateKeyPath string, outputSignedPath string) error {
	fmt.Printf("Signing PDF: input=%s, privateKey=%s, output=%s\n", pdfPath, privateKeyPath, outputSignedPath)

	if _, err := os.Stat(pdfPath); os.IsNotExist(err) {
		return fmt.Errorf("input PDF file does not exist: %s", pdfPath)
	}
	if _, err := os.Stat(privateKeyPath); os.IsNotExist(err) {
		return fmt.Errorf("private key file does not exist: %s", privateKeyPath)
	}

	if err := os.MkdirAll(filepath.Dir(outputSignedPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory for signed file: %v", err)
	}

	cmd := exec.Command("openssl", "dgst", "-sha256", "-sign", privateKeyPath, "-out", outputSignedPath, pdfPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to sign PDF: %v, output: %s", err, string(output))
	}

	fmt.Println("PDF signed successfully!")
	return nil
}
