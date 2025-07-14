package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

var (
	messages = make(map[string]map[string]string)
	once     sync.Once
)

// LoadI18nMessages loads all language JSON files from a folder
func LoadI18nMessages(folder string) error {
	var err error
	once.Do(func() {
		languages := []string{"en", "vi"}

		for _, lang := range languages {
			path := filepath.Join(folder, lang+".json")
			file, e := os.Open(path)
			if e != nil {
				err = fmt.Errorf("failed to open i18n file %s: %w", lang, e)
				return
			}
			defer file.Close()

			var msgs map[string]string
			decoder := json.NewDecoder(file)
			if e := decoder.Decode(&msgs); e != nil {
				err = fmt.Errorf("failed to decode i18n file %s: %w", lang, e)
				return
			}
			messages[lang] = msgs
		}
	})
	return err
}

// GetI18nMessage returns the translated message for a given key and lang
func GetI18nMessage(key, lang string) string {
	if langMsgs, ok := messages[lang]; ok {
		if msg, ok := langMsgs[key]; ok {
			return msg
		}
	}
	// fallback to English
	if msg, ok := messages["en"][key]; ok {
		return msg
	}
	return "Something went wrong"
}
