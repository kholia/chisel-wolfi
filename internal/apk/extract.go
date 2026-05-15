package apk

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"github.com/canonical/chisel/internal/deb"
)

type ExtractOptions = deb.ExtractOptions
type ExtractInfo = deb.ExtractInfo

func Extract(pkgReader io.ReadSeeker, options *ExtractOptions) error {
	opts := *options
	opts.DataReader = DataReader
	return deb.Extract(pkgReader, &opts)
}

func DataReader(pkgReader io.ReadSeeker) (io.ReadCloser, error) {
	_, err := pkgReader.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(pkgReader)
	for range 3 {
		names, err := readGzipTarStream(reader)
		if err != nil {
			return nil, err
		}
		if containsFile(names, ".PKGINFO") {
			gzipReader, err := gzip.NewReader(reader)
			if err != nil {
				return nil, err
			}
			gzipReader.Multistream(false)
			return gzipReader, nil
		}
	}
	return nil, fmt.Errorf("no data payload")
}

func IndexReader(reader io.Reader) (io.ReadCloser, error) {
	buffered := bufio.NewReader(reader)
	for range 2 {
		gzipReader, err := gzip.NewReader(buffered)
		if err != nil {
			return nil, err
		}
		gzipReader.Multistream(false)
		data, found, err := readFileFromTar(gzipReader, "APKINDEX")
		closeErr := gzipReader.Close()
		if err != nil {
			return nil, err
		}
		if closeErr != nil {
			return nil, closeErr
		}
		if found {
			return io.NopCloser(bytes.NewReader(data)), nil
		}
	}
	return nil, fmt.Errorf("no APKINDEX payload")
}

func ControlChecksum(pkgReader io.ReadSeeker) (string, error) {
	_, err := pkgReader.Seek(0, io.SeekStart)
	if err != nil {
		return "", err
	}
	data, err := io.ReadAll(pkgReader)
	if err != nil {
		return "", err
	}
	_, _ = pkgReader.Seek(0, io.SeekStart)

	reader := bytes.NewReader(data)
	for range 3 {
		start := len(data) - reader.Len()
		names, err := readGzipTarStream(reader)
		if err != nil {
			return "", err
		}
		end := len(data) - reader.Len()
		if containsFile(names, ".PKGINFO") {
			sum := sha1.Sum(data[start:end])
			return "Q1" + base64.StdEncoding.EncodeToString(sum[:]), nil
		}
	}
	return "", fmt.Errorf("no control payload")
}

func readGzipTarStream(reader io.Reader) ([]string, error) {
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}
	gzipReader.Multistream(false)
	names, err := readTarNames(gzipReader)
	closeErr := gzipReader.Close()
	if err != nil {
		return nil, err
	}
	if closeErr != nil {
		return nil, closeErr
	}
	return names, nil
}

func readTarNames(reader io.Reader) ([]string, error) {
	var names []string
	tarReader := tar.NewReader(reader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		names = append(names, cleanTarName(header.Name))
	}
	_, err := io.Copy(io.Discard, reader)
	if err != nil {
		return nil, err
	}
	return names, nil
}

func readFileFromTar(reader io.Reader, name string) ([]byte, bool, error) {
	tarReader := tar.NewReader(reader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, false, err
		}
		if cleanTarName(header.Name) != name {
			continue
		}
		data, err := io.ReadAll(tarReader)
		if err != nil {
			return nil, false, err
		}
		_, err = io.Copy(io.Discard, reader)
		if err != nil {
			return nil, false, err
		}
		return data, true, nil
	}
	_, err := io.Copy(io.Discard, reader)
	if err != nil {
		return nil, false, err
	}
	return nil, false, nil
}

func containsFile(names []string, name string) bool {
	for _, candidate := range names {
		if candidate == name {
			return true
		}
	}
	return false
}

func cleanTarName(name string) string {
	name = strings.TrimPrefix(name, "./")
	return strings.TrimPrefix(name, "/")
}
