package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Translations holds all the translation strings for the application
type Translations map[string]interface{}

// I18n represents the internationalization system
type I18n struct {
	CurrentLang string
	LangPath    string
	Translations Translations
	FallbackLang string
}

// NewI18n creates a new I18n instance
func NewI18n(langPath string, defaultLang string) (*I18n, error) {
	i18n := &I18n{
		CurrentLang:  defaultLang,
		LangPath:     langPath,
		Translations: make(Translations),
		FallbackLang: "en",
	}

	// Try to load the specified language
	err := i18n.LoadLanguage(defaultLang)
	if err != nil {
		// If the specified language fails, try to load the fallback language
		fmt.Printf("Warning: Could not load language '%s', falling back to '%s'\n", defaultLang, i18n.FallbackLang)
		err = i18n.LoadLanguage(i18n.FallbackLang)
		if err != nil {
			return nil, fmt.Errorf("failed to load fallback language: %v", err)
		}
	}

	return i18n, nil
}

// LoadLanguage loads the specified language file
func (i *I18n) LoadLanguage(lang string) error {
	filePath := filepath.Join(i.LangPath, lang+".json")
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("could not read language file: %v", err)
	}

	var translations Translations
	err = json.Unmarshal(data, &translations)
	if err != nil {
		return fmt.Errorf("could not parse language file: %v", err)
	}

	i.Translations = translations
	i.CurrentLang = lang
	return nil
}

// Get retrieves a translation string by its key
func (i *I18n) Get(key string) string {
	keys := strings.Split(key, ".")
	var current interface{} = i.Translations

	for _, k := range keys {
		m, ok := current.(map[string]interface{})
		if !ok {
			return key // Return the key if the path is invalid
		}

		current, ok = m[k]
		if !ok {
			return key // Return the key if it doesn't exist
		}
	}

	str, ok := current.(string)
	if !ok {
		return key // Return the key if the value is not a string
	}

	return str
}

// GetWithFormat retrieves a translation string by its key and formats it with the given arguments
func (i *I18n) GetWithFormat(key string, args ...interface{}) string {
	format := i.Get(key)
	return fmt.Sprintf(format, args...)
}

// AvailableLanguages returns a list of available languages
func (i *I18n) AvailableLanguages() ([]string, error) {
	files, err := ioutil.ReadDir(i.LangPath)
	if err != nil {
		return nil, fmt.Errorf("could not read language directory: %v", err)
	}

	languages := []string{}
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") && file.Name() != "template.json" {
			lang := strings.TrimSuffix(file.Name(), ".json")
			languages = append(languages, lang)
		}
	}

	return languages, nil
}

// DetectLanguage tries to detect the system language
func DetectLanguage() string {
	// Try to get the LANG environment variable
	lang := os.Getenv("LANG")
	if lang != "" {
		// Extract the language code (e.g., "en_US.UTF-8" -> "en")
		parts := strings.Split(lang, "_")
		if len(parts) > 0 {
			return parts[0]
		}
	}

	// Default to English if detection fails
	return "en"
}

// InitializeLanguage initializes the i18n system with the specified language
func InitializeLanguage(language string) error {
	// Get the executable directory
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("could not get executable path: %v", err)
	}

	// Determine the language directory path
	langPath := filepath.Join(filepath.Dir(execPath), "lang")

	// If the lang directory doesn't exist in the executable directory,
	// try to find it in the current directory or parent directories
	if _, err := os.Stat(langPath); os.IsNotExist(err) {
		// Try current directory
		cwd, err := os.Getwd()
		if err == nil {
			langPath = filepath.Join(cwd, "lang")
			if _, err := os.Stat(langPath); os.IsNotExist(err) {
				// Try parent directory
				langPath = filepath.Join(cwd, "..", "lang")
				if _, err := os.Stat(langPath); os.IsNotExist(err) {
					// Try parent of parent directory
					langPath = filepath.Join(cwd, "..", "..", "lang")
				}
			}
		}
	}

	// If no language specified, detect the system language
	if language == "" {
		language = DetectLanguage()
	}

	// Initialize the i18n system
	i18nInstance, err := NewI18n(langPath, language)
	if err != nil {
		return fmt.Errorf("could not initialize i18n: %v", err)
	}

	// Set the global i18n instance in help.go
	i18n = i18nInstance
	return nil
}
