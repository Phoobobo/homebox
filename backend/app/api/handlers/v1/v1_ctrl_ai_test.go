// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Phoobobo

package v1

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestAISecretRoundTrip(t *testing.T) {
	key := bytes.Repeat([]byte{7}, 32)
	a, e := aiSeal(key, "sk-test-secret")
	if e != nil {
		t.Fatal(e)
	}
	b, e := aiSeal(key, "sk-test-secret")
	if e != nil {
		t.Fatal(e)
	}
	if a == b {
		t.Fatal("nonce was reused")
	}
	got, e := aiOpen(key, a)
	if e != nil || got != "sk-test-secret" {
		t.Fatal("round trip failed")
	}
	if _, e = aiOpen(bytes.Repeat([]byte{8}, 32), a); e == nil {
		t.Fatal("wrong key accepted")
	}
	if _, e = aiOpen(key, "invalid"); e == nil {
		t.Fatal("malformed ciphertext accepted")
	}
	data, _ := json.Marshal(aiPublic(aiStored{Provider: "openrouter", Model: "vendor/model", Secret: a}))
	if bytes.Contains(data, []byte(a)) || bytes.Contains(data, []byte("sk-test-secret")) {
		t.Fatal("secret exposed")
	}
}
func TestAIValidation(t *testing.T) {
	for _, in := range []aiInput{{Provider: "http://localhost", Model: "m"}, {Provider: "openrouter"}, {Provider: "openrouter", Model: "m", APIKey: "a\nb"}, {Provider: "openrouter", Model: "m", APIKey: "x", ClearKey: true}} {
		if validateAI(in) == nil {
			t.Fatal("invalid settings accepted")
		}
	}
	if e := validateAI(aiInput{Provider: "openrouter", Model: "vendor/model", APIKey: "sk-test"}); e != nil {
		t.Fatal(e)
	}
}
func TestAIPrivatePersistence(t *testing.T) {
	dir := t.TempDir()
	k, e := aiMaster(dir)
	if e != nil {
		t.Fatal(e)
	}
	k2, e := aiMaster(dir)
	if e != nil || !bytes.Equal(k, k2) {
		t.Fatal("master key not stable")
	}
	secret, _ := aiSeal(k, "private-key")
	s := aiStored{Provider: "openrouter", Model: "vendor/model", Secret: secret}
	b, _ := json.Marshal(s)
	path := filepath.Join(dir, "settings.json")
	if e = aiWrite(path, b); e != nil {
		t.Fatal(e)
	}
	got, e := aiRead(path)
	if e != nil || !reflect.DeepEqual(got, s) {
		t.Fatal("settings not persisted")
	}
	for _, p := range []string{path, filepath.Join(dir, "master.key")} {
		st, e := os.Stat(p)
		if e != nil {
			t.Fatal(e)
		}
		if st.Mode().Perm() != 0600 {
			t.Fatal("secret file permissions must be 0600")
		}
	}
}
