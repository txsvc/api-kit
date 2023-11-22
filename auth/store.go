package auth

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/txsvc/stdlib/v2"

	"github.com/txsvc/cloudlib/helpers"
	"github.com/txsvc/cloudlib/observer"
	"github.com/txsvc/cloudlib/settings"
)

type (
	authCache struct {
		root string // location on disc
		// different types of lookup tables
		tokenToAuth map[string]*settings.DialSettings
		idToAuth    map[string]*settings.DialSettings
	}
)

var (
	cache *authCache // authorization cache
	mu    sync.Mutex // used to protect the above cache
)

func FlushAuthorizations(root string) {
	mu.Lock()
	defer mu.Unlock()

	observer.LogWithLevel(observer.LevelDebug, fmt.Sprintf("flushing auth cache. root=%s", root))

	cache = &authCache{
		root:        root,
		tokenToAuth: make(map[string]*settings.DialSettings),
		idToAuth:    make(map[string]*settings.DialSettings),
	}

	if root != "" {
		_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if info == nil {
				return nil
			}

			if !info.IsDir() {
				cfg, err := helpers.ReadDialSettings(path)
				if err != nil {
					return err // FIXME: this is never checked on exit !
				}
				_ = cache.Register(cfg)
			}
			return nil
		})
	}
}

func RegisterAuthorization(ds *settings.DialSettings) error {
	return cache.Register(ds)
}

func LookupByToken(token string) (*settings.DialSettings, error) {
	return cache.LookupByToken(token)
}

func UpdateStore(ds *settings.DialSettings) error {
	if _, err := cache.LookupByToken(ds.Credentials.Token); err != nil {
		return err // only allow to write already registered settings
	}
	return cache.writeToStore(ds)
}

func (c *authCache) Register(ds *settings.DialSettings) error {

	observer.LogWithLevel(observer.LevelDebug, fmt.Sprintf("register. t=%s/%s", ds.Credentials.Token, fileName(ds.Credentials)))

	// check if the settings already exists
	if a, ok := c.idToAuth[ds.Credentials.Key()]; ok {
		if a.Credentials.Status == settings.StateAuthorized { // FIXME this is weird, why?
			_ = observer.ReportError(fmt.Errorf("already authorized. t=%s, state=%d", a.Credentials.Token, a.Credentials.Status))
			return ErrAlreadyAuthorized
		}

		// remove from token lookup if the token changed
		if a.Credentials.Token != ds.Credentials.Token {
			delete(c.tokenToAuth, a.Credentials.Token)
		}
	}

	// write to the file store
	path := filepath.Join(c.root, fileName(ds.Credentials))
	if err := helpers.WriteDialSettings(ds, path); err != nil {
		return err
	}

	// update to the cache
	c.tokenToAuth[ds.Credentials.Token] = ds
	c.idToAuth[ds.Credentials.Key()] = ds

	return nil
}

func (c *authCache) LookupByToken(token string) (*settings.DialSettings, error) {
	observer.LogWithLevel(observer.LevelDebug, fmt.Sprintf("lookup. t=%s", token))

	if token == "" {
		return nil, ErrNoToken
	}
	if a, ok := c.tokenToAuth[token]; ok {
		return a, nil
	}
	return nil, nil // FIXME: return an error ?
}

func (c *authCache) writeToStore(ds *settings.DialSettings) error {
	// write to the file store
	path := filepath.Join(c.root, fileName(ds.Credentials))
	if err := helpers.WriteDialSettings(ds, path); err != nil {
		return err
	}
	return nil
}

func fileName(cred *settings.Credentials) string {
	return stdlib.Fingerprint(cred.Key())
}
