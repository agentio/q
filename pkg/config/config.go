package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/agentio/q/pkg/gcloud"
)

func readCachedToken() (*Token, error) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	b, err := os.ReadFile(filepath.Join(dirname, ".config", "q", "token.json"))
	if err != nil {
		return nil, err
	}
	var token Token
	err = json.Unmarshal(b, &token)
	if err != nil {
		return nil, err
	}
	created, err := time.Parse(time.RFC3339Nano, token.Created)
	if err != nil {
		return nil, err
	}
	if time.Since(created) > time.Hour {
		return nil, errors.New("token expired")
	}
	return &token, nil
}

func saveCachedToken(token *Token) error {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	b, err := json.Marshal(token)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Join(dirname, ".config", "q"), 0755)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dirname, ".config", "q", "token.json"), b, 0644)
}

func GetADCToken(verbose bool) (string, error) {
	token, err := readCachedToken()
	if err == nil {
		return token.Token, nil
	}
	token = &Token{Created: time.Now().Format(time.RFC3339Nano)}
	token.Token, err = gcloud.GetADCToken(verbose)
	if err != nil {
		return "", err
	}
	return token.Token, saveCachedToken(token)
}

type Token struct {
	Token   string `json:"token"`
	Created string `json:"created"`
}
