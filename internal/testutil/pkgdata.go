package testutil

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"io"
	"strings"
	"time"

	"github.com/blakesmith/ar"
	"github.com/klauspost/compress/zstd"
)

var PackageData = map[string][]byte{}

var TestPackageEntries = []TarEntry{
	Dir(0755, "./"),
	Dir(0755, "./dir/"),
	Reg(0644, "./dir/file", "12u3q0wej	ajsd"),
	Reg(0644, "./dir/other-file", "kasjdf0"),
	Dir(0755, "./dir/nested/"),
	Reg(0644, "./dir/nested/file", "0jqei"),
	Reg(0644, "./dir/nested/other-file", "1"),
	Dir(0755, "./dir/several/"),
	Dir(0755, "./dir/several/levels/"),
	Dir(0755, "./dir/several/levels/deep/"),
	Reg(0644, "./dir/several/levels/deep/file", "129i381		"),
	Dir(0755, "./other-dir/"),
	Dir(01777, "./parent/"),
	Dir(0764, "./parent/permissions/"),
	Reg(0755, "./parent/permissions/file", "ajse0"),
}

var OtherPackageEntries = []TarEntry{
	Dir(0755, "./"),
	Reg(0644, "./file", "masfdko"),
}

func init() {
	PackageData["test-package"] = MustMakeDeb(TestPackageEntries)
	PackageData["other-package"] = MustMakeDeb(OtherPackageEntries)
}

type TarEntry struct {
	Header  tar.Header
	NoFixup bool
	Content []byte
}

var zeroTime time.Time
var epochStartTime time.Time = time.Unix(0, 0)

func fixupTarEntry(entry *TarEntry) {
	if entry.NoFixup {
		return
	}
	hdr := &entry.Header
	if hdr.Typeflag == 0 {
		if hdr.Linkname != "" {
			hdr.Typeflag = tar.TypeSymlink
		} else if strings.HasSuffix(hdr.Name, "/") {
			hdr.Typeflag = tar.TypeDir
		} else {
			hdr.Typeflag = tar.TypeReg
		}
	}
	if hdr.Mode == 0 {
		switch hdr.Typeflag {
		case tar.TypeDir:
			hdr.Mode = 0755
		case tar.TypeSymlink:
			hdr.Mode = 0777
		default:
			hdr.Mode = 0644
		}
	}
	if hdr.Size == 0 && entry.Content != nil {
		hdr.Size = int64(len(entry.Content))
	}
	if hdr.Uid == 0 && hdr.Uname == "" {
		hdr.Uname = "root"
	}
	if hdr.Gid == 0 && hdr.Gname == "" {
		hdr.Gname = "root"
	}
	if hdr.ModTime.Equal(zeroTime) {
		hdr.ModTime = epochStartTime
	}
	if hdr.Format == 0 {
		hdr.Format = tar.FormatGNU
	}
}

func makeTar(entries []TarEntry) ([]byte, error) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	for _, entry := range entries {
		fixupTarEntry(&entry)
		if err := tw.WriteHeader(&entry.Header); err != nil {
			return nil, err
		}
		if entry.Content != nil {
			if _, err := tw.Write(entry.Content); err != nil {
				return nil, err
			}
		}
	}
	return buf.Bytes(), nil
}

func compressBytesZstd(input []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer, err := zstd.NewWriter(&buf)
	if err != nil {
		return nil, err
	}
	if _, err = writer.Write(input); err != nil {
		return nil, err
	}
	if err = writer.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func compressBytesGzip(input []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	if _, err := writer.Write(input); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func MakeDeb(entries []TarEntry) ([]byte, error) {
	var buf bytes.Buffer

	tarData, err := makeTar(entries)
	if err != nil {
		return nil, err
	}
	compTarData, err := compressBytesZstd(tarData)
	if err != nil {
		return nil, err
	}

	writer := ar.NewWriter(&buf)
	if err := writer.WriteGlobalHeader(); err != nil {
		return nil, err
	}
	dataHeader := ar.Header{
		Name: "data.tar.zst",
		Mode: 0644,
		Size: int64(len(compTarData)),
	}
	if err := writer.WriteHeader(&dataHeader); err != nil {
		return nil, err
	}
	if _, err = writer.Write(compTarData); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func MustMakeDeb(entries []TarEntry) []byte {
	data, err := MakeDeb(entries)
	if err != nil {
		panic(err)
	}
	return data
}

func MakeAPK(entries []TarEntry, signed bool) ([]byte, error) {
	var buf bytes.Buffer

	if signed {
		signatureTar, err := makeTar([]TarEntry{
			Reg(0644, ".SIGN.RSA.test.rsa.pub", "signature"),
		})
		if err != nil {
			return nil, err
		}
		signatureData, err := compressBytesGzip(signatureTar)
		if err != nil {
			return nil, err
		}
		buf.Write(signatureData)
	}

	payload, err := makeAPKPayload(entries)
	if err != nil {
		return nil, err
	}
	buf.Write(payload)

	return buf.Bytes(), nil
}

func makeAPKPayload(entries []TarEntry) ([]byte, error) {
	var buf bytes.Buffer

	controlTar, err := makeTar([]TarEntry{
		Reg(0644, ".PKGINFO", "pkgname = test-package\npkgver = 1.0-r0\narch = x86_64\n"),
	})
	if err != nil {
		return nil, err
	}
	controlData, err := compressBytesGzip(controlTar)
	if err != nil {
		return nil, err
	}
	buf.Write(controlData)

	dataTar, err := makeTar(entries)
	if err != nil {
		return nil, err
	}
	data, err := compressBytesGzip(dataTar)
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	return buf.Bytes(), nil
}

func MakeSignedAPK(entries []TarEntry, key *rsa.PrivateKey, keyName string) ([]byte, error) {
	payload, err := makeAPKPayload(entries)
	if err != nil {
		return nil, err
	}
	controlData, err := firstGzipMember(payload)
	if err != nil {
		return nil, err
	}
	signature, err := makeAPKSignature(controlData, key, keyName)
	if err != nil {
		return nil, err
	}
	return append(signature, payload...), nil
}

func MustMakeSignedAPK(entries []TarEntry, key *rsa.PrivateKey, keyName string) []byte {
	data, err := MakeSignedAPK(entries, key, keyName)
	if err != nil {
		panic(err)
	}
	return data
}

func MustMakeAPK(entries []TarEntry, signed bool) []byte {
	data, err := MakeAPK(entries, signed)
	if err != nil {
		panic(err)
	}
	return data
}

func MakeAPKIndex(index string, signed bool) ([]byte, error) {
	var buf bytes.Buffer

	if signed {
		signatureTar, err := makeTar([]TarEntry{
			Reg(0644, ".SIGN.RSA.test.rsa.pub", "signature"),
		})
		if err != nil {
			return nil, err
		}
		signatureData, err := compressBytesGzip(signatureTar)
		if err != nil {
			return nil, err
		}
		buf.Write(signatureData)
	}

	indexTar, err := makeTar([]TarEntry{
		Reg(0644, "DESCRIPTION", "test apk index\n"),
		Reg(0644, "APKINDEX", index),
	})
	if err != nil {
		return nil, err
	}
	indexData, err := compressBytesGzip(indexTar)
	if err != nil {
		return nil, err
	}
	buf.Write(indexData)

	return buf.Bytes(), nil
}

func MakeSignedAPKIndex(index string, key *rsa.PrivateKey, keyName string) ([]byte, error) {
	indexData, err := makeAPKIndexPayload(index)
	if err != nil {
		return nil, err
	}
	signature, err := makeAPKSignature(indexData, key, keyName)
	if err != nil {
		return nil, err
	}
	return append(signature, indexData...), nil
}

func MustMakeSignedAPKIndex(index string, key *rsa.PrivateKey, keyName string) []byte {
	data, err := MakeSignedAPKIndex(index, key, keyName)
	if err != nil {
		panic(err)
	}
	return data
}

func makeAPKIndexPayload(index string) ([]byte, error) {
	indexTar, err := makeTar([]TarEntry{
		Reg(0644, "DESCRIPTION", "test apk index\n"),
		Reg(0644, "APKINDEX", index),
	})
	if err != nil {
		return nil, err
	}
	return compressBytesGzip(indexTar)
}

func makeAPKSignature(payload []byte, key *rsa.PrivateKey, keyName string) ([]byte, error) {
	if keyName == "" {
		keyName = "test.rsa.pub"
	}
	sum := sha256.Sum256(payload)
	signature, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, sum[:])
	if err != nil {
		return nil, err
	}
	signatureTar, err := makeTar([]TarEntry{{
		Header: tar.Header{
			Typeflag: tar.TypeReg,
			Name:     ".SIGN.RSA256." + keyName,
			Mode:     0644,
		},
		Content: signature,
	}})
	if err != nil {
		return nil, err
	}
	return compressBytesGzip(signatureTar)
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

func MustMakeAPKIndex(index string, signed bool) []byte {
	data, err := MakeAPKIndex(index, signed)
	if err != nil {
		panic(err)
	}
	return data
}

// Reg is a shortcut for creating a regular file TarEntry structure (with
// tar.Typeflag set tar.TypeReg). Reg stands for "REGular file".
func Reg(mode int64, path, content string) TarEntry {
	return TarEntry{
		Header: tar.Header{
			Typeflag: tar.TypeReg,
			Name:     path,
			Mode:     mode,
		},
		Content: []byte(content),
	}
}

// Dir is a shortcut for creating a directory TarEntry structure (with
// tar.Typeflag set to tar.TypeDir). Dir stands for "DIRectory".
func Dir(mode int64, path string) TarEntry {
	return TarEntry{
		Header: tar.Header{
			Typeflag: tar.TypeDir,
			Name:     path,
			Mode:     mode,
		},
	}
}

// Lnk is a shortcut for creating a symbolic link TarEntry structure (with
// tar.Typeflag set to tar.TypeSymlink). Lnk stands for "symbolic LiNK".
func Lnk(mode int64, path, target string) TarEntry {
	return TarEntry{
		Header: tar.Header{
			Typeflag: tar.TypeSymlink,
			Name:     path,
			Mode:     mode,
			Linkname: target,
		},
	}
}

// Hrd is a shortcut for creating a hard link TarEntry structure (with
// tar.Typeflag set to tar.TypeLink). Hrd stands for "HaRD link".
func Hrd(mode int64, path, target string) TarEntry {
	return TarEntry{
		Header: tar.Header{
			Typeflag: tar.TypeLink,
			Name:     path,
			Mode:     mode,
			Linkname: target,
		},
	}
}
