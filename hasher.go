package main

import (
	// For future research, you can use
	"bytes"

	"github.com/minio/sha256-simd"
	// However, this relies on assembly code and may not build properly
	// "crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

var (
	ErrParsingDifficulty   = errors.New("Error Parsing Difficulty")
	ErrCreatingRadixNumber = errors.New("Error Creating Radix Number")

	MaxRetries        = errors.New("Unsolvable Challenge: Retry Limit Exceeded")
	HashLimitExceeded = errors.New("Unsolvable Challenge: Hash Limit Exceeded")
)

func main() {
	fmt.Println(GrabHash(`{"do":["sid|855a6699-bb9a-11ed-bbff-626c70424d52","pnf|cu","cls|32060492575082180333|23765544575816579762","sts|1678050785554","wcs|cg2gboauqq6urpbp0480","drc|6412","cts|855a6a77-bb9a-11ed-bbff-626c70424d52|false","cs|99006a3a7a4e999df3c2e5de2b63611201318ad0dbbb050765d0c8a44bf201c3","vid|67932318-bb9a-11ed-89dd-436342496472|31536000|false","cp|1|35bc35ef77d4d787c8f08d309e88a5c3c2d0ef6cbc05abc4b5f0b097b8|d6af800afdeb6718b9207fd11ff6ea7718142f6d81ca79fd0d9f0da42168c031|25|false","ci|1|85628f20-bb9a-11ed-b230-9124f8607fdd|1056|ee0addfef8ccf455eddcf27b24c993887a5731ec83bb00b4d81e01f1f4e2314270372b7f80313bbe706fde46b9ab548dcfd596d6d8aee51352984821ed4b3a1a󠄻󠄺󠄸󠄸|0|NA","sff|cc|60|U2FtZVNpdGU9TGF4Ow==","sff|idp_c|60|1,s","sff|rf|60|1","sff|fp|60|1"]}`, 30))
}

func GrabHash(doResponse string, maxDifficulty int) (string, error) {
	chalValues := fetchCapChallenge(doResponse)
	if chalValues == nil {
		return "", ErrParsingDifficulty
	}

	solve, err := solveChallenge(chalValues, maxDifficulty)
	return solve, err
}

func timestampNow() int64 {
	return time.Now().UnixMilli()
}

func fetchCapChallenge(responseOne string) []string {
	if !strings.Contains(responseOne, "do") {
		responseOne = fmt.Sprintf("{\"do\": %s}", responseOne)
	}

	// rewrite this at some point
	for _, item := range gjson.Get(responseOne, "do").Array() {
		values := strings.Split(item.String(), "|")
		if len(values) >= 5 && values[0] == "cp" {
			return values[2:5]
		}
	}

	return nil
}

// Original version
func solveChallenge(chalValues []string, maxDifficulty int) (string, error) {
	startHash, targetHashS, difficulty := chalValues[0], chalValues[1], chalValues[2]

	targetHash := make([]byte, 32)
	hex.Decode(targetHash, []byte(targetHashS))

	intDifficulty, err := strconv.ParseInt(difficulty, 10, 64)
	if err != nil {
		return "", ErrParsingDifficulty
	}

	floatDifficulty, err := strconv.ParseFloat(difficulty, 64)
	if err != nil {
		return "", ErrParsingDifficulty
	}

	if maxDifficulty != 0 && intDifficulty > int64(maxDifficulty) {
		return "", HashLimitExceeded
	}

	aa := int64(math.Floor(floatDifficulty / 4))
	shiftedNumber := int64((1 << (4 * aa)) - 1)

	hexNumber := string([]byte{startHash[len(startHash)-1]})
	startHashTrimmedLast := startHash[:len(startHash)-1]

	radix16HexNumber, err := strconv.ParseInt(hexNumber, 16, 64)
	if err != nil {
		return "", ErrCreatingRadixNumber
	}

	zeroString := ""
	for i := 0; int64(i) < aa; i++ {
		zeroString += "0"
	}

	startHashTwo := int64(1 << intDifficulty)

	var ff int64 = 0
	var gg int64 = 1
	var hh int64 = 250

	var solvedHash string
	var notFound bool = true

	tries := 0

	hasher := sha256.New()

	//this part will be replaced//
	for {
		if !notFound {
			break
		}

		tries++

		var ts int64

		for c := hh * gg; ff < startHashTwo; ff++ {
			c--
			ts = timestampNow()
			tries++

			if tries > 10000000000 {
				return "", MaxRetries
			}

			p1 := strconv.FormatInt(radix16HexNumber+(ff>>(aa<<2)), 16)
			basep2 := zeroString + strconv.FormatInt(ff&shiftedNumber, 16)
			p2 := basep2[int64(len(basep2))-aa:]

			g := startHashTrimmedLast + p1 + p2

			hasher.Reset()
			hasher.Write([]byte(g))

			sum := hasher.Sum(nil)

			isEqual := true
			for i, v := range targetHash {
				if sum[i] != v {
					isEqual = false
					break
				}
			}

			if isEqual {
				solvedHash = g
				notFound = false
				break
			}
		}

		if timestampNow()-ts <= 50 {
			gg++
		} else {
			gg--
			gg = int64(math.Max(float64(gg), 1))
		}

	}

	// until here //

	return solvedHash, nil
}

// Coroutine version
func solveChallengeWithCoroutines(chalValues []string, maxDifficulty int) (string, error) {
	startHash, targetHashS, difficulty := chalValues[0], chalValues[1], chalValues[2]

	targetHash := make([]byte, 32)
	hex.Decode(targetHash, []byte(targetHashS))

	intDifficulty, err := strconv.ParseInt(difficulty, 10, 64)
	if err != nil {
		return "", ErrParsingDifficulty
	}

	floatDifficulty, err := strconv.ParseFloat(difficulty, 64)
	if err != nil {
		return "", ErrParsingDifficulty
	}

	if maxDifficulty != 0 && intDifficulty > int64(maxDifficulty) {
		return "", HashLimitExceeded
	}

	aa := int64(math.Floor(floatDifficulty / 4))
	shiftedNumber := (1 << (4 * aa)) - 1

	hexNumber := string([]byte{startHash[len(startHash)-1]})
	startHashTrimmedLast := startHash[:len(startHash)-1]

	radix16HexNumber, err := strconv.ParseInt(hexNumber, 16, 64)
	if err != nil {
		return "", ErrCreatingRadixNumber
	}

	zeroString := ""
	for i := int64(0); i < aa; i++ {
		zeroString += "0"
	}

	startHashTwo := int64(1 << intDifficulty)

	var solvedHash string
	var notFound bool = true

	var ff int64 = 0
	var gg int64 = 1
	var hh int64 = 250

	tries := 0

	hasher := sha256.New()

	hashCh := make(chan int64)

	// Generate hash values and send them to the worker goroutines
	go func() {
		for i := int64(0); i < startHashTwo; i++ {
			hashCh <- i
		}
		close(hashCh)
	}()

	// Receive the generated hashes from the channel and check if any of them match the target hash
	for i := int64(0); i < 10000000000 && notFound; i++ {
		select {
		case _, ok := <-hashCh:
			if !ok {
				break
			}

			tries++

			if tries > 10000000000 {
				return "", MaxRetries
			}

			var ts int64

			for c := hh * gg; ff < startHashTwo; ff++ {
				c--
				ts = timestampNow()

				p1 := strconv.FormatInt(radix16HexNumber+(ff>>(aa<<2)), 16)
				basep2 := zeroString + strconv.FormatInt(ff&int64(shiftedNumber), 16)
				p2 := basep2[int64(len(basep2))-aa:]

				g := startHashTrimmedLast + p1 + p2

				hasher.Reset()
				hasher.Write([]byte(g))

				sum := hasher.Sum(nil)

				if bytes.Equal(sum, targetHash) {
					solvedHash = g
					notFound = false
					break
				}
			}

			if timestampNow()-ts <= 50 {
				gg++
			} else {
				gg--
				gg = int64(math.Max(float64(gg), 1))
			}

		default:
			// Do nothing and continue waiting for more hashes to process
		}
	}

	if notFound {
		return "not found", nil
	}

	return solvedHash, nil
}
