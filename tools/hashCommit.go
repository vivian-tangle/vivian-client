package tools

import (
	"encoding/hex"
	"math/big"

	ristretto "github.com/bwesterb/go-ristretto"
)

// Pederson commitment based on Elliptic Curve

// GenerateParameters generates two points on the curve for computing the commitment
func GenerateParameters() (g, h ristretto.Point) {
	g.Rand()
	h.Rand()
	return g, h
}

// GenerateParametersToString generates two points on the curve for computing the commitment
// The return type is string
func GenerateParametersToString() (gString, hString string) {
	var g, h ristretto.Point
	g.Rand()
	h.Rand()
	gBytes, _ := g.MarshalText()
	hBytes, _ := h.MarshalText()
	gString = hex.EncodeToString(gBytes)
	hString = hex.EncodeToString(hBytes)
	return
}

// GenerateRandom generates a random scalar
func GenerateRandom() (r ristretto.Scalar) {
	r.Rand()
	return r
}

// GenerateRandomToString generates a random scalar
// The return type is string
func GenerateRandomToString() string {
	var r ristretto.Scalar
	r.Rand()
	return r.BigInt().String()
}

// Commit computes the Pederson commitment of the secret, and return the point on the curve
func Commit(g, h ristretto.Point, secret []byte, r ristretto.Scalar) (commit ristretto.Point) {
	var x ristretto.Scalar
	x.Derive(secret)
	// c = xG + rH
	commit.Add(g.ScalarMult(&g, &x), h.ScalarMult(&h, &r))
	return commit
}

// CommitByString computes the Pederson commitment of the secret, and return the point on the curve
// The return type is string
func CommitByString(gString, hString, rString string, secret []byte) (commitString string, err error) {
	var g, h ristretto.Point
	gBytes, _ := hex.DecodeString(gString)
	hBytes, _ := hex.DecodeString(hString)
	err = g.UnmarshalText(gBytes)
	if err != nil {
		return "", err
	}
	err = h.UnmarshalText(hBytes)
	if err != nil {
		return "", err
	}
	var r ristretto.Scalar
	var bigInt big.Int
	bigInt.SetString(rString, 10)
	r.SetBigInt(&bigInt)

	commit := Commit(g, h, secret, r)
	bytes, err := commit.MarshalText()
	if err != nil {
		return "", err
	}
	commitString = hex.EncodeToString(bytes)
	return commitString, err
}

// Open verifies the commitment based on the original parameters
func Open(commit, g, h ristretto.Point, secret []byte, r ristretto.Scalar) bool {
	var x ristretto.Scalar
	x.Derive(secret)
	var calculateCommit ristretto.Point
	calculateCommit.Add(g.ScalarMult(&g, &x), h.ScalarMult(&h, &r))
	return calculateCommit.Equals(&commit)
}

// OpenByString verifies the commitment based on the original parameters
// The input parameters are strings
func OpenByString(commitString, gString, hString, rString string, secret []byte) bool {
	var g, h, verifyCommit ristretto.Point
	verifyCommitBytes, _ := hex.DecodeString(commitString)
	err := verifyCommit.UnmarshalText(verifyCommitBytes)
	if err != nil {
		return false
	}
	gBytes, _ := hex.DecodeString(gString)
	err = g.UnmarshalText(gBytes)
	if err != nil {
		return false
	}
	hBytes, _ := hex.DecodeString(hString)
	err = h.UnmarshalText(hBytes)
	if err != nil {
		return false
	}
	var r ristretto.Scalar
	var bigInt big.Int
	bigInt.SetString(rString, 10)
	r.SetBigInt(&bigInt)

	commit := Commit(g, h, secret, r)
	return verifyCommit.Equals(&commit)
}

// PedersonCommit is the struct for storing the results of Pederson commitment
type PedersonCommit struct {
	Content string
	G       string
	H       string
	R       string
	Commit  string
	TxHash  string
}
