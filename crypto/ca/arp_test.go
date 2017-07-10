package ca

import (
	"encoding/hex"
	"testing"

	"chain/crypto/ed25519/ecmath"
)

func TestAssetRangeProof(t *testing.T) {
	msg := []byte("message")

	type acPair struct {
		assetIDHex string
		aekHex     string
	}

	cases := []struct {
		ac          []acPair
		aekPrimeHex string
	}{
		{
			ac: []acPair{
				{
					"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
					"012345",
				},
				{
					"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
					"6789",
				},
				{
					"cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc",
					"abcd",
				},
			},
			aekPrimeHex: "efdcab",
		},
	}

	for _, ca := range cases {
		var assetID AssetID
		hex.Decode(assetID[:], []byte(ca.ac[0].assetIDHex))
		var aek []byte
		if ca.aekPrimeHex != "" {
			aek = make([]byte, hex.DecodedLen(len(ca.aekPrimeHex)))
			hex.Decode(aek, []byte(ca.aekPrimeHex))
		}
		acPrime, cPrime := CreateAssetCommitment(assetID, aek)
		if cPrime == nil {
			cPrime = &ecmath.Zero
		}

		var (
			ac []*AssetCommitment
			c  *ecmath.Scalar
		)
		for i, pair := range ca.ac {
			hex.Decode(assetID[:], []byte(pair.assetIDHex))
			aek = nil
			if pair.aekHex != "" {
				aek = make([]byte, hex.DecodedLen(len(pair.aekHex)))
				hex.Decode(aek[:], []byte(pair.aekHex))
			}
			thisAC, thisC := CreateAssetCommitment(assetID, aek)
			if thisC == nil {
				thisC = &ecmath.Zero
			}
			ac = append(ac, thisAC)
			if i == 0 {
				c = thisC
			}
		}

		arp := CreateAssetRangeProof(msg, ac, acPrime, 0, *c, *cPrime)
		if !arp.Validate(msg, acPrime) {
			t.Error("what we have here is failure to validate")
		}
		if arp.Validate(msg[1:], acPrime) {
			t.Error("validated invalid proof")
		}
		if arp.Validate(msg, ac[0]) {
			t.Error("validated invalid proof")
		}
	}
}