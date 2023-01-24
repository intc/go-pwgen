package pwgen

/* pwgen.go - Minimalistic Go adaptation from Theodore Ts'o's pwgen/pw_phonemes.c
 * Copyright (C) 2023 by Antti Antinoja
 * Lisence: MIT
 */

import (
	"bytes"
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	mrand "math/rand"
	"os"
)

const (
	Consonant byte = 0x1
	Vowel     byte = 0x2
	Dipthong  byte = 0x4
	NotFirst  byte = 0x8

	ShowDebug = false
)

type pwE struct {
	bSyl  []byte
	flags byte
}

var numElements int = 0

var nrArr = [10]byte{0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39}

var elements = []pwE{}

var elementsStd = []pwE{
	pwE{bSyl: []byte("a"), flags: Vowel},
	pwE{bSyl: []byte("ae"), flags: Vowel | Dipthong},
	pwE{bSyl: []byte("ah"), flags: Vowel | Dipthong},
	pwE{bSyl: []byte("ai"), flags: Vowel | Dipthong},
	pwE{bSyl: []byte("b"), flags: Consonant},
	pwE{bSyl: []byte("c"), flags: Consonant},
	pwE{bSyl: []byte("ch"), flags: Consonant | Dipthong},
	pwE{bSyl: []byte("d"), flags: Consonant},
	pwE{bSyl: []byte("e"), flags: Vowel},
	pwE{bSyl: []byte("ee"), flags: Vowel | Dipthong},
	pwE{bSyl: []byte("ei"), flags: Vowel | Dipthong},
	pwE{bSyl: []byte("f"), flags: Consonant},
	pwE{bSyl: []byte("g"), flags: Consonant},
	pwE{bSyl: []byte("gh"), flags: Consonant | Dipthong | NotFirst},
	pwE{bSyl: []byte("h"), flags: Consonant},
	pwE{bSyl: []byte("i"), flags: Vowel},
	pwE{bSyl: []byte("ie"), flags: Vowel | Dipthong},
	pwE{bSyl: []byte("j"), flags: Consonant},
	pwE{bSyl: []byte("k"), flags: Consonant},
	pwE{bSyl: []byte("l"), flags: Consonant},
	pwE{bSyl: []byte("m"), flags: Consonant},
	pwE{bSyl: []byte("n"), flags: Consonant},
	pwE{bSyl: []byte("ng"), flags: Consonant | Dipthong | NotFirst},
	pwE{bSyl: []byte("o"), flags: Vowel},
	pwE{bSyl: []byte("oh"), flags: Vowel | Dipthong},
	pwE{bSyl: []byte("oo"), flags: Vowel | Dipthong},
	pwE{bSyl: []byte("p"), flags: Consonant},
	pwE{bSyl: []byte("ph"), flags: Consonant | Dipthong},
	pwE{bSyl: []byte("qu"), flags: Consonant | Dipthong},
	pwE{bSyl: []byte("r"), flags: Consonant},
	pwE{bSyl: []byte("s"), flags: Consonant},
	pwE{bSyl: []byte("sh"), flags: Consonant | Dipthong},
	pwE{bSyl: []byte("t"), flags: Consonant},
	pwE{bSyl: []byte("th"), flags: Consonant | Dipthong},
	pwE{bSyl: []byte("u"), flags: Vowel},
	pwE{bSyl: []byte("v"), flags: Consonant},
	pwE{bSyl: []byte("w"), flags: Consonant},
	pwE{bSyl: []byte("x"), flags: Consonant},
	pwE{bSyl: []byte("y"), flags: Consonant},
	pwE{bSyl: []byte("z"), flags: Consonant},
}

var elementsNoIl = []pwE{
	pwE{bSyl: []byte("a"), flags: Vowel},
	pwE{bSyl: []byte("ae"), flags: Vowel | Dipthong},
	pwE{bSyl: []byte("ah"), flags: Vowel | Dipthong},
	pwE{bSyl: []byte("ay"), flags: Vowel | Dipthong},
	pwE{bSyl: []byte("b"), flags: Consonant},
	pwE{bSyl: []byte("c"), flags: Consonant},
	pwE{bSyl: []byte("ch"), flags: Consonant | Dipthong},
	pwE{bSyl: []byte("d"), flags: Consonant},
	pwE{bSyl: []byte("e"), flags: Vowel},
	pwE{bSyl: []byte("ee"), flags: Vowel | Dipthong},
	pwE{bSyl: []byte("eh"), flags: Vowel | Dipthong},
	pwE{bSyl: []byte("ey"), flags: Vowel | Dipthong},
	pwE{bSyl: []byte("f"), flags: Consonant},
	pwE{bSyl: []byte("g"), flags: Consonant},
	pwE{bSyl: []byte("gh"), flags: Consonant | Dipthong | NotFirst},
	pwE{bSyl: []byte("h"), flags: Consonant},
	pwE{bSyl: []byte("ye"), flags: Vowel | Dipthong},
	pwE{bSyl: []byte("j"), flags: Consonant},
	pwE{bSyl: []byte("k"), flags: Consonant},
	pwE{bSyl: []byte("m"), flags: Consonant},
	pwE{bSyl: []byte("n"), flags: Consonant},
	pwE{bSyl: []byte("ng"), flags: Consonant | Dipthong | NotFirst},
	pwE{bSyl: []byte("o"), flags: Vowel},
	pwE{bSyl: []byte("oh"), flags: Vowel | Dipthong},
	pwE{bSyl: []byte("oo"), flags: Vowel | Dipthong},
	pwE{bSyl: []byte("p"), flags: Consonant},
	pwE{bSyl: []byte("ph"), flags: Consonant | Dipthong},
	pwE{bSyl: []byte("qu"), flags: Consonant | Dipthong},
	pwE{bSyl: []byte("r"), flags: Consonant},
	pwE{bSyl: []byte("s"), flags: Consonant},
	pwE{bSyl: []byte("sh"), flags: Consonant | Dipthong},
	pwE{bSyl: []byte("t"), flags: Consonant},
	pwE{bSyl: []byte("th"), flags: Consonant | Dipthong},
	pwE{bSyl: []byte("u"), flags: Vowel},
	pwE{bSyl: []byte("uo"), flags: Vowel | Dipthong},
	pwE{bSyl: []byte("v"), flags: Consonant},
	pwE{bSyl: []byte("w"), flags: Consonant},
	pwE{bSyl: []byte("x"), flags: Consonant},
	pwE{bSyl: []byte("y"), flags: Consonant},
	pwE{bSyl: []byte("z"), flags: Consonant},
}

func init() {
	var b [8]byte
	count, err := crand.Read(b[:])
	if err != nil || count < 8 {
		fmt.Printf("FATAL: pwgen.init(): Can't read from random source")
		os.Exit(1)
	}
	mrand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
	if ShowDebug {
		fmt.Printf("Element flags:\n Vowel:     %08b\n Consonant: %08b\n Dipthong:  %08b\n NotFirst:  %08b\n\n",
			Vowel,
			Consonant,
			Dipthong,
			NotFirst,
		)
		for _, el := range elements {
			fmt.Printf("%3s: flags: %08b\n", string(el.bSyl), el.flags)
		}
	}
	elements = elementsStd
	numElements = len(elements)
}

func pwNumber(max int) int {
	return mrand.Intn(max)
}

func phonemeRun(size int) (*[]byte, error) {
	var res []byte
	strBuf := new(bytes.Buffer)
	c := 0
	first := true
	prev := byte(0)
	shouldBe := Consonant

	if pwNumber(2) == 1 {
		shouldBe = Vowel
	}

	for c < size {
		i := pwNumber(numElements)
		bSyl := elements[i].bSyl
		bSylLen := cap(bSyl)
		flags := elements[i].flags
		// filter on the basic type of the next element
		if flags&shouldBe == 0 {
			continue
		}
		if first && flags&NotFirst != 0 {
			continue
		}
		// don't allow VOWEL followed a Vowel/Dipthong pair
		if prev&Vowel != 0 && flags&Vowel != 0 && flags&Dipthong != 0 {
			continue
		}
		// syllable is too long -> skip
		if bSylLen > (size - c) {
			continue
		}
		// syllable ok, add it to our buffer
		_, err := strBuf.Write(bSyl)
		if err != nil {
			return nil, err
		}
		// update counter
		c = c + bSylLen
		// are we done?
		if c >= size {
			break
		}
		// numbers
		if !first && pwNumber(10) < 3 {
			ch := nrArr[pwNumber(10)]
			err = strBuf.WriteByte(ch)
			if err != nil {
				return nil, err
			}
			first = true
			prev = 0
			shouldBe = Consonant
			if pwNumber(2) == 1 {
				shouldBe = Vowel
			}
			c = c + 1
			continue
		}
		// what next?
		if shouldBe == Consonant {
			shouldBe = Vowel
		} else {
			if prev&Vowel != 0 || flags&Dipthong != 0 || pwNumber(10) > 3 {
				shouldBe = Consonant
			} else {
				shouldBe = Vowel
			}
		}
		prev = flags
		first = false
	}
	res = strBuf.Bytes()
	return &res, nil
}

func ActivateNoIlSet() {
	elements = elementsNoIl
	numElements = len(elements)
}

func PhonemeGen(size int) (*string, error) {
	res := ""
	for {
		bArr, err := phonemeRun(size)
		if err != nil {
			return nil, err
		}
		if bArr != nil {
			res = string(*bArr)
			break
		}
	}
	return &res, nil
}
