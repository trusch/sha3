package main

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/sha3"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

var (
	hashType     = pflag.StringP("type", "t", "shake256", "one of 'sum224', 'sum256', 'sum384', 'sum512', 'shake128', 'shake256'")
	outputLength = pflag.IntP("length", "l", 32, "output length in bytes if using a shake hash")
	help         = pflag.BoolP("help", "h", false, "show usage info")
	check        = pflag.BoolP("check", "c", false, "check sum files")
	workers      = pflag.IntP("workers", "w", 8, "number of workers for recursive hashing")
)

func main() {
	pflag.Parse()
	if *help {
		pflag.Usage()
		os.Exit(1)
	}
	args := pflag.Args()
	if len(args) == 0 {
		args = []string{"/dev/stdin"}
	}

	for _, arg := range args {
		if *check {
			err := checkSumFile(arg, *hashType, *outputLength)
			if err != nil {
				logrus.Fatal(err)
			}
		} else {
			if info, err := os.Stat(arg); err == nil && info.IsDir() {
				hashDirectory(arg)
			} else {
				hash, err := hashFile(arg, *hashType, *outputLength)
				if err != nil {
					logrus.Fatal(err)
				}
				fmt.Printf("%x  %s\n", hash, arg)
			}
		}
	}
}

func hashFile(file string, hashType string, outputLength int) ([]byte, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	switch hashType {
	case "sum224", "sum256", "sum384", "sum512":
		return computeSumHash(f, hashType)
	case "shake128", "shake256":
		return computeShakeHash(f, hashType, outputLength)
	}
	return nil, errors.New("unknown hash type")
}

func computeSumHash(f io.Reader, hashType string) ([]byte, error) {
	var hash hash.Hash
	switch hashType {
	case "sum224":
		hash = sha3.New224()
	case "sum256":
		hash = sha3.New256()
	case "sum384":
		hash = sha3.New384()
	case "sum512":
		hash = sha3.New512()
	}
	_, err := io.Copy(hash, f)
	if err != nil {
		return nil, err
	}
	return hash.Sum(nil), nil
}

func computeShakeHash(f io.Reader, hashType string, outputLength int) ([]byte, error) {
	var hash sha3.ShakeHash
	switch hashType {
	case "shake128":
		hash = sha3.NewShake128()
	case "shake256":
		hash = sha3.NewShake256()
	}
	_, err := io.Copy(hash, f)
	if err != nil {
		return nil, err
	}
	output := make([]byte, outputLength)
	_, err = hash.Read(output[0:])
	if err != nil {
		return nil, err
	}
	return output[:outputLength], nil
}

func checkSumFile(file string, hashType string, outputLength int) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "  ")
		if len(parts) != 2 {
			return errors.New("malformed sum file")
		}
		hash, err := hex.DecodeString(parts[0])
		if err != nil {
			return err
		}
		err = checkHash(parts[1], hashType, outputLength, hash)
		if err != nil {
			fmt.Printf("%v: FAIL\n", parts[1])
			os.Exit(1)
		} else {
			fmt.Printf("%v: OK\n", parts[1])
		}
	}
	return nil
}

func checkHash(file string, hashType string, outputLength int, hash []byte) error {
	computedHash, err := hashFile(file, hashType, outputLength)
	if err != nil {
		return err
	}
	if len(computedHash) != len(hash) {
		return errors.New("hash length doesn't match")
	}
	for k, v := range computedHash {
		if hash[k] != v {
			return errors.New("hash mismatch")
		}
	}
	return nil
}

func forAllFiles(dir string, fn func(path string) error) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return fn(path)
		}
		return nil
	})
}

func hashDirectory(dir string) {
	pathes := make(chan string, 512)
	printerInput := make(chan string, 512)

	// produce pathes
	go func() {
		defer close(pathes)
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				logrus.Warn(err)
				return nil
			}
			if !info.IsDir() {
				pathes <- path
			}
			return nil
		})
		if err != nil {
			logrus.Error(err)
		}
	}()

	// consume pathes in hash workers
	hashWorkerCount := *workers
	done := make(chan struct{}, hashWorkerCount)
	for i := 0; i < hashWorkerCount; i++ {
		go func() {
			hashWorker(pathes, printerInput)
			done <- struct{}{}
		}()
	}
	// wait for them and close the printerInput afterwards
	go func() {
		for i := 0; i < hashWorkerCount; i++ {
			<-done
		}
		close(printerInput)
	}()

	// consume the printer input
	for line := range printerInput {
		fmt.Println(line)
	}
}

func hashWorker(in <-chan string, out chan<- string) {
	for path := range in {
		hash, err := hashFile(path, *hashType, *outputLength)
		if err != nil {
			logrus.Warn(err)
			continue
		}
		out <- fmt.Sprintf("%x  %s", hash, path)
	}
}
