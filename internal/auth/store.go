package auth

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/txsvc/apikit/internal"
	"github.com/txsvc/apikit/internal/settings"
	"github.com/txsvc/stdlib/v2"
)

type (
	authCache struct {
		root string // location on disc
		// different types of lookup tables
		tokenToAuth map[string]*settings.Settings
		idToAuth    map[string]*settings.Settings
	}
)

var (
	cache *authCache // authorization cache
	mu    sync.Mutex // used to protect the above cache
)

func FlushAuthorizations(root string) {
	mu.Lock()
	defer mu.Unlock()

	_log.Debugf("flushing auth cache. root=%s", root)

	cache = &authCache{
		root:        root,
		tokenToAuth: make(map[string]*settings.Settings),
		idToAuth:    make(map[string]*settings.Settings),
	}

	if root != "" {
		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				cfg, err := internal.ReadSettingsFromFile(path)
				if err != nil {
					return err // FIXME: this is never checked on exit !
				}
				cache.Register(cfg)
			}
			return nil
		})
	}
}

func RegisterAuthorization(cfg *settings.Settings) error {
	return cache.Register(cfg)
}

func LookupByToken(token string) (*settings.Settings, error) {
	return cache.LookupByToken(token)
}

func UpdateStore(cfg *settings.Settings) error {
	if _, err := cache.LookupByToken(cfg.Credentials.Token); err != nil {
		return err // only allow to write already registered settings
	}
	return cache.writeToStore(cfg)
}

func (c *authCache) Register(cfg *settings.Settings) error {

	_log.Debugf("register. t=%s/%s", cfg.Credentials.Token, fileName(cfg.Credentials))

	// check if the settings already exists
	if a, ok := c.idToAuth[cfg.Credentials.Key()]; ok {
		if a.Status == settings.StateAuthorized {
			_log.Errorf("already authorized. t=%s, state=%d", a.Credentials.Token, a.Status)
			return ErrAlreadyAuthorized
		}

		// remove from token lookup if the token changed
		if a.Credentials.Token != cfg.Credentials.Token {
			delete(c.tokenToAuth, a.Credentials.Token)
		}
	}

	// write to the file store
	path := filepath.Join(c.root, fileName(cfg.Credentials))
	if err := cfg.WriteToFile(path); err != nil {
		return err
	}

	// update to the cache
	c.tokenToAuth[cfg.Credentials.Token] = cfg
	c.idToAuth[cfg.Credentials.Key()] = cfg

	return nil
}

func (c *authCache) LookupByToken(token string) (*settings.Settings, error) {
	_log.Debugf("lookup. t=%s", token)

	if token == "" {
		return nil, ErrNoToken
	}
	if a, ok := c.tokenToAuth[token]; ok {
		return a, nil
	}
	return nil, nil // FIXME: return an error ?
}

func (c *authCache) writeToStore(cfg *settings.Settings) error {
	// write to the file store
	path := filepath.Join(c.root, fileName(cfg.Credentials))
	if err := cfg.WriteToFile(path); err != nil {
		return err
	}
	return nil
}

func fileName(cred *settings.Credentials) string {
	return stdlib.Fingerprint(cred.Key())
}
