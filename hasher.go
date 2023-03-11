package main

/*
#include "solveHash.h"
#include <stdlib.h>
*/
import "C"
import (
	// For future research, you can use
	// "github.com/minio/sha256-simd"
	// However, this relies on assembly code and may not build properly

	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/tidwall/gjson"
)

var (
	ErrParsingDifficulty   = errors.New("Error Parsing Difficulty")
	ErrCreatingRadixNumber = errors.New("Error Creating Radix Number")

	MaxRetries        = errors.New("Unsolvable Challenge: Retry Limit Exceeded")
	HashLimitExceeded = errors.New("Unsolvable Challenge: Hash Limit Exceeded")
)

// Wrapper function to call the C function
func solveHash(targetHash []byte, startHashTwo, radix16HexNumber, shiftedNumber, hh, aa, ff int64, startHashTrimmedLast string) (string, error) {
	startHash, _ := strconv.ParseUint(startHashTrimmedLast, 10, 64)
	cs := C.solveHash(
		(*C.uint8_t)(C.CBytes(targetHash)),
		C.uint64_t(startHashTwo),
		C.uint64_t(startHash),
		C.uint64_t(radix16HexNumber),
		C.uint64_t(shiftedNumber),
		C.uint64_t(hh),
		C.uint64_t(aa),
		C.uint64_t(ff),
	)
	if cs == nil {
		return "", MaxRetries
	}
	result := C.GoString(cs)
	C.free(unsafe.Pointer(cs))
	return result, nil
}

func main() {
	fmt.Println(GrabHash(`{"do":["sid|855a6699-bb9a-11ed-bbff-626c70424d52","pnf|cu","cls|32060492575082180333|23765544575816579762","sts|1678050785554","wcs|cg2gboauqq6urpbp0480","drc|6412","cts|855a6a77-bb9a-11ed-bbff-626c70424d52|false","cs|99006a3a7a4e999df3c2e5de2b63611201318ad0dbbb050765d0c8a44bf201c3","vid|67932318-bb9a-11ed-89dd-436342496472|31536000|false","cp|1|35bc35ef77d4d787c8f08d309e88a5c3c2d0ef6cbc05abc4b5f0b097b8|d6af800afdeb6718b9207fd11ff6ea7718142f6d81ca79fd0d9f0da42168c031|25|false","ci|1|85628f20-bb9a-11ed-b230-9124f8607fdd|1056|ee0addfef8ccf455eddcf27b24c993887a5731ec83bb00b4d81e01f1f4e2314270372b7f80313bbe706fde46b9ab548dcfd596d6d8aee51352984821ed4b3a1a󠄻󠄺󠄸󠄸|0|NA","sff|cc|60|U2FtZVNpdGU9TGF4Ow==","sff|idp_c|60|1,s","sff|rf|60|1","sff|fp|60|1"]}`, 30))
}

func GrabHash(doResponse string, maxDifficulty int) (string, error) {
	chalValues := fetchCapChallenge(doResponse)
	if chalValues == nil {
		return "", ErrParsingDifficulty
	}

    solve, err := solveChallenge(chalValues, int64(maxDifficulty))
	return solve, err
}

func solveChallenge(chalValues []string, maxDifficulty int64) (string, error) {
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
	var hh int64 = 250


	solvedHash, err := solveHash(targetHash, startHashTwo, radix16HexNumber, shiftedNumber, hh, aa, ff, startHashTrimmedLast)

	return solvedHash, nil
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
