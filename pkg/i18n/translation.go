package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2"
)

const ctxLocaleKey = "__locale"

type i18nStore struct {
	mu     sync.RWMutex
	table  map[string]map[string]string // locale -> flat "a.b.c" -> value
	defLoc string
}

var (
	store    = &i18nStore{table: map[string]map[string]string{}, defLoc: "id"}
	initOnce sync.Once
	initErr  error
)

func Init(dir, defaultLocale string) error {
	initOnce.Do(func() {
		store.defLoc = strings.ToLower(defaultLocale)
		ents, err := os.ReadDir(dir)
		if err != nil {
			initErr = fmt.Errorf("i18n: readdir: %w", err)
			return
		}
		loaded := 0
		for _, e := range ents {
			if e.IsDir() || filepath.Ext(e.Name()) != ".json" {
				continue
			}
			loc := strings.TrimSuffix(e.Name(), ".json")
			b, err := os.ReadFile(filepath.Join(dir, e.Name()))
			if err != nil {
				initErr = fmt.Errorf("i18n: read %s: %w", e.Name(), err)
				return
			}
			var m map[string]any
			if err := json.Unmarshal(b, &m); err != nil {
				initErr = fmt.Errorf("i18n: json %s: %w", e.Name(), err)
				return
			}
			flat := map[string]string{}
			flatten("", m, flat)

			store.mu.Lock()
			store.table[strings.ToLower(loc)] = flat
			store.mu.Unlock()
			loaded++
		}
		if loaded == 0 {
			initErr = fmt.Errorf("i18n: no json files found in %s", dir)
		}
	})
	return initErr
}

func flatten(prefix string, v any, out map[string]string) {
	switch t := v.(type) {
	case map[string]any:
		for k, vv := range t {
			if prefix == "" {
				flatten(k, vv, out)
			} else {
				flatten(prefix+"."+k, vv, out)
			}
		}
	case string:
		out[prefix] = t
	case float64, bool, nil:
		out[prefix] = fmt.Sprint(t)
	case []any:
		b, _ := json.Marshal(t)
		out[prefix] = string(b)
	}
}

func InitLocale(defaultLocale string) fiber.Handler {
	store.defLoc = strings.ToLower(defaultLocale)
	return func(c *fiber.Ctx) error {
		if l := c.Query("lang"); l != "" {
			c.Locals(ctxLocaleKey, strings.ToLower(l))
			return c.Next()
		}
		if l := c.Cookies("lang"); l != "" {
			c.Locals(ctxLocaleKey, strings.ToLower(l))
			return c.Next()
		}
		if al := c.Get("App-Language"); al != "" {
			p := strings.Split(al, ",")[0]
			p = strings.ToLower(strings.TrimSpace(strings.SplitN(p, ";", 2)[0]))
			if i := strings.Index(p, "-"); i > 0 {
				p = p[:i]
			}
			if p != "" {
				c.Locals(ctxLocaleKey, p)
				return c.Next()
			}
		}
		c.Locals(ctxLocaleKey, store.defLoc)
		return c.Next()
	}
}

func T(c *fiber.Ctx, key string, data ...map[string]any) string {
	if initErr != nil {
		return key
	}
	locale, _ := c.Locals(ctxLocaleKey).(string)
	if locale == "" {
		locale = store.defLoc
	}
	txt := lookup(locale, key)
	if len(data) > 0 && len(data[0]) > 0 {
		for k, v := range data[0] {
			txt = strings.ReplaceAll(txt, "{{"+k+"}}", fmt.Sprint(v))
		}
	}
	return txt
}

func lookup(locale, key string) string {
	locale = strings.ToLower(locale)
	store.mu.RLock()
	defer store.mu.RUnlock()

	if m, ok := store.table[locale]; ok {
		if v, ok := m[key]; ok {
			return v
		}
	}
	if locale != store.defLoc {
		if m, ok := store.table[store.defLoc]; ok {
			if v, ok := m[key]; ok {
				return v
			}
		}
	}
	return key
}
