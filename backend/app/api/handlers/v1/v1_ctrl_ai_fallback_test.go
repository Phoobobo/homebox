// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Phoobobo

package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAIModelFallbackOrder(t *testing.T) {
	calls := []string{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			Model string `json:"model"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		calls = append(calls, payload.Model)
		if payload.Model == "primary" {
			w.WriteHeader(503)
			return
		}
		if payload.Model == "backup1" {
			w.WriteHeader(429)
			return
		}
		_, _ = w.Write([]byte(`{"choices":[{"message":{"role":"assistant","content":"OK"}}]}`))
	}))
	defer srv.Close()
	aiEndpoints["fallback-test"] = srv.URL
	defer delete(aiEndpoints, "fallback-test")
	s := aiStored{Provider: "fallback-test", Model: "primary", FallbackModels: []string{"backup1", "backup2"}}
	msg, e := aiCompletion(context.Background(), s, "test", []aiMessage{{Role: "user", Content: "hi"}})
	require.NoError(t, e)
	require.Equal(t, []string{"primary", "backup1", "backup2"}, calls)
	require.Equal(t, "backup2", msg.ModelUsed)
	calls = nil
	s.Model = "working"
	_, e = aiCompletion(context.Background(), s, "test", nil)
	require.NoError(t, e)
	require.Equal(t, []string{"working"}, calls, "successful primary must not invoke backups")
	calls = nil
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, e = aiCompletion(ctx, s, "test", nil)
	require.Error(t, e)
	require.Empty(t, calls)
}
func TestAIModelFallbackValidation(t *testing.T) {
	require.NoError(t, validateAI(aiInput{Provider: "openrouter", Model: "main", FallbackModels: []string{"one", "two"}}))
	for _, models := range [][]string{{"main"}, {"one", "one"}, {"one", "two", "three"}, {""}, {"bad\nmodel"}} {
		require.Error(t, validateAI(aiInput{Provider: "openrouter", Model: "main", FallbackModels: models}))
	}
	require.Equal(t, []string{"main", "one", "two"}, aiModelOrder(aiStored{Model: "main", FallbackModels: []string{"one", "two"}}))
}
func TestAIEmptyCompletionTriggersFallback(t *testing.T) {
	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if calls == 1 {
			_, _ = w.Write([]byte(`{"choices":[{"message":{"content":""}}]}`))
			return
		}
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"OK"}}]}`))
	}))
	defer srv.Close()
	aiEndpoints["empty-test"] = srv.URL
	defer delete(aiEndpoints, "empty-test")
	m, e := aiCompletion(context.Background(), aiStored{Provider: "empty-test", Model: "main", FallbackModels: []string{"backup"}}, "key", nil)
	require.NoError(t, e)
	require.Equal(t, "backup", m.ModelUsed)
	require.Equal(t, 2, calls)
}
