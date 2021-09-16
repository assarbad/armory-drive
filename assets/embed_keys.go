// Copyright (c) F-Secure Corporation
// https://foundry.f-secure.com
//
// Use of this source code is governed by the license
// that can be found in the LICENSE file.

//go:build linux && ignore
// +build linux,ignore

package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func init() {
	log.SetFlags(0)
	log.SetOutput(os.Stdout)
}

func read(env string) (buf []byte, err error) {
	p := os.Getenv(env)

	if len(p) == 0 {
		return nil, fmt.Errorf("%s must be defined", env)
	}

	if buf, err = os.ReadFile(p); err != nil {
		return
	}

	if len(buf) == 0 {
		return nil, fmt.Errorf("%s is empty", p)
	}

	return
}

func main() {
	out, err := os.Create("tmp_keys.go")

	if err != nil {
		log.Fatal(err)
	}

	out.WriteString(`
package assets

func init() {
`)

	if p := os.Getenv("DISABLE_FR_AUTH"); len(p) != 0 {
		out.WriteString(`}`)
		return
	}

	FRPublicKey, err := read("FR_PUBKEY")

	if err != nil {
		log.Fatal(err)
	}

	LogPublicKey, err := read("LOG_PUBKEY")

	if err != nil {
		log.Fatal(err)
	}

	out.WriteString(fmt.Sprintf("\tFRPublicKey = []byte(%s)\n", strconv.Quote(string(FRPublicKey))))
	out.WriteString(fmt.Sprintf("\tLogPublicKey = []byte(%s)\n", strconv.Quote(string(LogPublicKey))))
	out.WriteString(`}`)
}
