// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Phoobobo

package v1

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/hay-kot/httpkit/errchain"
	"github.com/hay-kot/httpkit/server"
	"github.com/sysadminsmedia/homebox/backend/internal/core/services"
	"github.com/sysadminsmedia/homebox/backend/internal/sys/validate"
)

var aiMu sync.Mutex
var aiEndpoints = map[string]string{
	"openrouter": "https://openrouter.ai/api/v1",
	"bailian":    "https://dashscope.aliyuncs.com/compatible-mode/v1",
	"deepseek":   "https://api.deepseek.com/v1",
}

type aiStored struct {
	Provider       string   `json:"provider"`
	Model          string   `json:"model"`
	FallbackModels []string `json:"fallbackModels"`
	Secret         string   `json:"encryptedKey"`
}
type aiInput struct {
	Provider       string   `json:"provider"`
	Model          string   `json:"model"`
	FallbackModels []string `json:"fallbackModels"`
	APIKey         string   `json:"apiKey"`
	ClearKey       bool     `json:"clearKey"`
}
type aiView struct {
	Provider       string   `json:"provider"`
	Model          string   `json:"model"`
	FallbackModels []string `json:"fallbackModels"`
	BaseURL        string   `json:"baseUrl"`
	HasKey         bool     `json:"hasKey"`
}

func aiPublic(s aiStored) aiView {
	return aiView{Provider: s.Provider, Model: s.Model, FallbackModels: s.FallbackModels, BaseURL: aiEndpoints[s.Provider], HasKey: s.Secret != ""}
}
func validateAI(in aiInput) error {
	if _, ok := aiEndpoints[in.Provider]; !ok {
		return errors.New("unsupported provider")
	}
	if strings.TrimSpace(in.Model) == "" || len(in.Model) > 200 || strings.ContainsAny(in.Model, "\r\n") {
		return errors.New("model ID is required (maximum 200 characters)")
	}
	if len(in.FallbackModels) > 2 {
		return errors.New("at most two fallback models are allowed")
	}
	seen := map[string]bool{in.Model: true}
	for _, m := range in.FallbackModels {
		if strings.TrimSpace(m) == "" || len(m) > 200 || strings.ContainsAny(m, "\r\n") || seen[m] {
			return errors.New("fallback models must be unique non-empty model IDs (maximum 200 characters)")
		}
		seen[m] = true
	}
	if len(in.APIKey) > 4096 || strings.ContainsAny(in.APIKey, "\r\n") {
		return errors.New("invalid API key")
	}
	if in.ClearKey && in.APIKey != "" {
		return errors.New("cannot replace and clear key together")
	}
	return nil
}
func aiSeal(key []byte, text string) (string, error) {
	b, e := aes.NewCipher(key)
	if e != nil {
		return "", e
	}
	g, e := cipher.NewGCM(b)
	if e != nil {
		return "", e
	}
	n := make([]byte, g.NonceSize())
	if _, e = rand.Read(n); e != nil {
		return "", e
	}
	return base64.StdEncoding.EncodeToString(g.Seal(n, n, []byte(text), []byte("zainaer-ai-v1"))), nil
}
func aiOpen(key []byte, text string) (string, error) {
	b, e := aes.NewCipher(key)
	if e != nil {
		return "", e
	}
	g, e := cipher.NewGCM(b)
	if e != nil {
		return "", e
	}
	v, e := base64.StdEncoding.DecodeString(text)
	if e != nil || len(v) < g.NonceSize() {
		return "", errors.New("invalid encrypted key")
	}
	p, e := g.Open(nil, v[:g.NonceSize()], v[g.NonceSize():], []byte("zainaer-ai-v1"))
	return string(p), e
}
func aiWrite(path string, data []byte) error {
	f, e := os.CreateTemp(filepath.Dir(path), ".ai-*")
	if e != nil {
		return e
	}
	defer os.Remove(f.Name())
	if _, e = f.Write(data); e != nil {
		f.Close()
		return e
	}
	if e = f.Sync(); e != nil {
		f.Close()
		return e
	}
	if e = f.Close(); e != nil {
		return e
	}
	return os.Rename(f.Name(), path)
}
func aiRead(path string) (aiStored, error) {
	s := aiStored{Provider: "openrouter"}
	b, e := os.ReadFile(path)
	if os.IsNotExist(e) {
		return s, nil
	}
	if e != nil {
		return s, e
	}
	e = json.Unmarshal(b, &s)
	return s, e
}

// Called under aiMu: master key never leaves the service's private storage.
func aiMaster(dir string) ([]byte, error) {
	path := filepath.Join(dir, "master.key")
	b, e := os.ReadFile(path)
	if e == nil {
		if len(b) != 32 {
			return nil, errors.New("invalid master key")
		}
		return b, nil
	}
	if !os.IsNotExist(e) {
		return nil, e
	}
	b = make([]byte, 32)
	if _, e = rand.Read(b); e != nil {
		return nil, e
	}
	return b, aiWrite(path, b)
}
func (ctrl *V1Controller) aiPath(r *http.Request) (string, error) {
	if r.Method != http.MethodGet && (r.Header.Get("Sec-Fetch-Site") == "cross-site" || !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json")) {
		return "", validate.NewRequestError(errors.New("same-origin JSON request required"), http.StatusForbidden)
	}

	ctx := services.NewContext(r.Context())
	ok, e := ctrl.repo.Groups.IsOwnerOf(ctx, ctx.UID, ctx.GID)
	if e != nil {
		return "", validate.NewRequestError(errors.New("cannot verify collection owner"), 500)
	}
	if !ok || ctrl.isDemo {
		return "", validate.NewRequestError(errors.New("only collection owners can configure models"), 403)
	}
	return filepath.Join(ctrl.config.Storage.PrefixPath, "ai", ctx.GID.String()+".json"), nil
}
func (ctrl *V1Controller) HandleAISettings() errchain.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		path, e := ctrl.aiPath(r)
		if e != nil {
			return e
		}
		w.Header().Set("Cache-Control", "no-store")
		aiMu.Lock()
		defer aiMu.Unlock()
		s, e := aiRead(path)
		if e != nil {
			return validate.NewRequestError(errors.New("cannot read model configuration"), 500)
		}
		if r.Method == http.MethodPut {
			r.Body = http.MaxBytesReader(w, r.Body, 8192)
			var in aiInput
			if e = json.NewDecoder(r.Body).Decode(&in); e != nil {
				return validate.NewRequestError(errors.New("invalid request"), 400)
			}
			in.APIKey = strings.TrimSpace(in.APIKey)
			in.Model = strings.TrimSpace(in.Model)
			for i := range in.FallbackModels {
				in.FallbackModels[i] = strings.TrimSpace(in.FallbackModels[i])
			}
			if e = validateAI(in); e != nil {
				return validate.NewRequestError(e, 400)
			}
			if in.Provider != s.Provider && in.APIKey == "" && !in.ClearKey {
				return validate.NewRequestError(errors.New("enter a new key when switching providers"), 400)
			}
			if e = os.MkdirAll(filepath.Dir(path), 0700); e != nil {
				return validate.NewRequestError(errors.New("cannot create private storage"), 500)
			}
			s.Provider = in.Provider
			s.Model = in.Model
			s.FallbackModels = in.FallbackModels
			if in.ClearKey {
				s.Secret = ""
			}
			if in.APIKey != "" {
				k, err := aiMaster(filepath.Dir(path))
				if err != nil {
					return validate.NewRequestError(errors.New("cannot load encryption key"), 500)
				}
				s.Secret, e = aiSeal(k, in.APIKey)
				if e != nil {
					return validate.NewRequestError(errors.New("cannot encrypt API key"), 500)
				}
			}
			b, _ := json.Marshal(s)
			if e = aiWrite(path, b); e != nil {
				return validate.NewRequestError(errors.New("cannot save model configuration"), 500)
			}
		}
		return server.JSON(w, 200, aiPublic(s))
	}
}
func (ctrl *V1Controller) HandleAITest() errchain.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		path, e := ctrl.aiPath(r)
		if e != nil {
			return e
		}
		w.Header().Set("Cache-Control", "no-store")
		aiMu.Lock()
		s, e := aiRead(path)
		var key string
		if e == nil && s.Secret != "" {
			var master []byte
			master, e = os.ReadFile(filepath.Join(filepath.Dir(path), "master.key"))
			if e == nil {
				key, e = aiOpen(master, s.Secret)
			}
		}
		aiMu.Unlock()
		if e != nil {
			return validate.NewRequestError(errors.New("cannot read saved credentials"), 500)
		}
		_, ok := aiEndpoints[s.Provider]
		if !ok || key == "" || s.Model == "" {
			return validate.NewRequestError(errors.New("save provider, model and API key first"), 400)
		}
		msg, e := aiCompleteWithFallback(r.Context(), s, key, []aiMessage{{Role: "user", Content: "Reply only with OK."}}, false, 32)
		if e != nil {
			return validate.NewRequestError(e, 502)
		}
		return server.JSON(w, 200, map[string]any{"ok": true, "model": msg.ModelUsed})
	}
}
