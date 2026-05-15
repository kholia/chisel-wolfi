package apk

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"strings"
)

var ErrNoSignature = errors.New("missing APK RSA signature")

type rsaSignature struct {
	name      string
	hash      crypto.Hash
	signature []byte
}

// DecodeRSAPublicKey decodes an APK repository RSA public key.
func DecodeRSAPublicKey(data []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(bytes.TrimSpace(data))
	if block == nil {
		return nil, fmt.Errorf("cannot decode PEM public key")
	}

	switch block.Type {
	case "PUBLIC KEY":
		key, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		rsaKey, ok := key.(*rsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("PEM data does not contain an RSA public key")
		}
		return rsaKey, nil
	case "RSA PUBLIC KEY":
		return x509.ParsePKCS1PublicKey(block.Bytes)
	case "CERTIFICATE":
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, err
		}
		rsaKey, ok := cert.PublicKey.(*rsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("certificate does not contain an RSA public key")
		}
		return rsaKey, nil
	default:
		return nil, fmt.Errorf("unsupported PEM block type %q", block.Type)
	}
}

// RSAPublicKeyID returns a stable identifier for an APK repository RSA public
// key. The value is intended for release-file sanity checks, not apk-tools
// compatibility.
func RSAPublicKeyID(key *rsa.PublicKey) (string, error) {
	der, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(der)
	return "SHA256:" + hex.EncodeToString(sum[:]), nil
}

// VerifySignature verifies the first APK signature stream against any of the
// provided RSA public keys.
func VerifySignature(data []byte, keys []*rsa.PublicKey) error {
	if len(keys) == 0 {
		return fmt.Errorf("missing APK RSA public keys")
	}

	signatures, signedData, err := splitSignature(data)
	if err != nil {
		return err
	}

	signedCandidates, err := signedDataCandidates(signedData)
	if err != nil {
		return err
	}

	for _, sig := range signatures {
		for _, candidate := range signedCandidates {
			digest, err := hashSignedData(sig.hash, candidate)
			if err != nil {
				return fmt.Errorf("cannot hash signed APK data: %v", err)
			}
			for _, key := range keys {
				err = rsa.VerifyPKCS1v15(key, sig.hash, digest, sig.signature)
				if err == nil {
					return nil
				}
			}
		}
	}

	if len(signatures) == 1 && len(keys) == 1 {
		return fmt.Errorf("cannot verify APK RSA signature %q", signatures[0].name)
	}
	return fmt.Errorf("cannot verify any APK RSA signatures")
}

func signedDataCandidates(data []byte) ([][]byte, error) {
	firstMember, err := firstGzipMember(data)
	if err != nil {
		return nil, err
	}
	if len(firstMember) == len(data) {
		return [][]byte{data}, nil
	}
	return [][]byte{firstMember, data}, nil
}

func firstGzipMember(data []byte) ([]byte, error) {
	reader := bytes.NewReader(data)
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}
	gzipReader.Multistream(false)
	_, copyErr := io.Copy(io.Discard, gzipReader)
	closeErr := gzipReader.Close()
	if copyErr != nil {
		return nil, copyErr
	}
	if closeErr != nil {
		return nil, closeErr
	}
	memberEnd := len(data) - reader.Len()
	return data[:memberEnd], nil
}

func splitSignature(data []byte) ([]rsaSignature, []byte, error) {
	reader := bytes.NewReader(data)
	buffered := bufio.NewReader(reader)

	gzipReader, err := gzip.NewReader(buffered)
	if err != nil {
		return nil, nil, err
	}
	gzipReader.Multistream(false)

	signatures, err := readSignaturesFromTar(gzipReader)
	closeErr := gzipReader.Close()
	if err != nil {
		return nil, nil, err
	}
	if closeErr != nil {
		return nil, nil, closeErr
	}
	if len(signatures) == 0 {
		return nil, nil, ErrNoSignature
	}

	signedStart := len(data) - reader.Len() - buffered.Buffered()
	return signatures, data[signedStart:], nil
}

func readSignaturesFromTar(reader io.Reader) ([]rsaSignature, error) {
	var signatures []rsaSignature
	tarReader := tar.NewReader(reader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		hash, ok := signatureHash(cleanTarName(header.Name))
		if !ok {
			continue
		}

		signature, err := io.ReadAll(tarReader)
		if err != nil {
			return nil, err
		}
		signatures = append(signatures, rsaSignature{
			name:      cleanTarName(header.Name),
			hash:      hash,
			signature: signature,
		})
	}
	_, err := io.Copy(io.Discard, reader)
	if err != nil {
		return nil, err
	}
	return signatures, nil
}

func signatureHash(name string) (crypto.Hash, bool) {
	switch {
	case strings.HasPrefix(name, ".SIGN.RSA."):
		return crypto.SHA1, true
	case strings.HasPrefix(name, ".SIGN.RSA256."):
		return crypto.SHA256, true
	case strings.HasPrefix(name, ".SIGN.RSA512."):
		return crypto.SHA512, true
	default:
		return 0, false
	}
}

func hashSignedData(hash crypto.Hash, data []byte) ([]byte, error) {
	switch hash {
	case crypto.SHA1:
		sum := sha1.Sum(data)
		return sum[:], nil
	case crypto.SHA256:
		sum := sha256.Sum256(data)
		return sum[:], nil
	case crypto.SHA512:
		sum := sha512.Sum512(data)
		return sum[:], nil
	default:
		return nil, fmt.Errorf("unsupported APK RSA signature hash %v", hash)
	}
}
