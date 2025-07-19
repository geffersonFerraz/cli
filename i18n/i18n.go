package i18n

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/spf13/cobra"
)

//go:embed translations/*.json
var translationsFS embed.FS

// Locale representa um idioma/localização
type Locale struct {
	Code         string            `json:"code"`
	Name         string            `json:"name"`
	NativeName   string            `json:"native_name"`
	Translations map[string]string `json:"translations"`
}

// Manager gerencia as traduções e idiomas
type Manager struct {
	locales     map[string]*Locale
	current     *Locale
	defaultLang string
	mutex       sync.RWMutex
}

var (
	instance *Manager
	once     sync.Once
)

// GetInstance retorna a instância singleton do gerenciador de i18n
func GetInstance() *Manager {
	once.Do(func() {
		instance = &Manager{
			locales:     make(map[string]*Locale),
			defaultLang: "pt-BR",
		}
		// Carregar traduções automaticamente na inicialização
		if err := instance.LoadLocales(); err != nil {
			// Log do erro, mas não falhar a inicialização
			fmt.Printf("Warning: failed to load translations: %v\n", err)
		}
	})
	return instance
}

// LoadLocales carrega os arquivos de tradução embutidos
func (m *Manager) LoadLocales() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Listar todos os arquivos .json no diretório translations
	files, err := translationsFS.ReadDir("translations")
	if err != nil {
		return fmt.Errorf("erro ao listar arquivos de tradução: %w", err)
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			locale, err := m.loadLocaleFile("translations/" + file.Name())
			if err != nil {
				return fmt.Errorf("erro ao carregar %s: %w", file.Name(), err)
			}
			m.locales[locale.Code] = locale
		}
	}

	// Definir idioma padrão se não houver nenhum carregado
	if len(m.locales) == 0 {
		return fmt.Errorf("nenhum arquivo de tradução encontrado")
	}

	// Definir idioma atual
	m.setCurrentLocale(m.detectLanguage())

	return nil
}

// loadLocaleFile carrega um arquivo de tradução individual do embed
func (m *Manager) loadLocaleFile(filepath string) (*Locale, error) {
	data, err := translationsFS.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var locale Locale
	if err := json.Unmarshal(data, &locale); err != nil {
		return nil, err
	}

	return &locale, nil
}

// detectLanguage detecta o idioma preferido do usuário
func (m *Manager) detectLanguage() string {
	// 1. Verificar variável de ambiente MGC_LANG
	if lang := os.Getenv("MGC_LANG"); lang != "" {
		if m.isValidLocale(lang) {
			return lang
		}
	}

	// 2. Verificar variável de ambiente LANG
	if lang := os.Getenv("LANG"); lang != "" {
		// Extrair código do idioma (ex: pt_BR.UTF-8 -> pt-BR)
		langCode := strings.Split(lang, ".")[0]
		langCode = strings.Replace(langCode, "_", "-", 1)

		if m.isValidLocale(langCode) {
			return langCode
		}

		// Tentar apenas o código principal (pt_BR -> pt)
		mainLang := strings.Split(langCode, "-")[0]
		if m.isValidLocale(mainLang) {
			return mainLang
		}
	}

	// 3. Verificar variável de ambiente LC_ALL
	if lang := os.Getenv("LC_ALL"); lang != "" {
		langCode := strings.Split(lang, ".")[0]
		langCode = strings.Replace(langCode, "_", "-", 1)

		if m.isValidLocale(langCode) {
			return langCode
		}
	}

	// 4. Fallback para idioma padrão
	return m.defaultLang
}

// isValidLocale verifica se um código de idioma é válido
func (m *Manager) isValidLocale(code string) bool {
	_, exists := m.locales[code]
	return exists
}

// setCurrentLocale define o idioma atual
func (m *Manager) setCurrentLocale(code string) {
	if locale, exists := m.locales[code]; exists {
		m.current = locale
	} else {
		// Fallback para o primeiro idioma disponível
		for _, locale := range m.locales {
			m.current = locale
			break
		}
	}
}

// SetLanguage define o idioma atual
func (m *Manager) SetLanguage(code string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if !m.isValidLocale(code) {
		return fmt.Errorf("idioma não suportado: %s", code)
	}

	m.setCurrentLocale(code)
	return nil
}

// GetLanguage retorna o código do idioma atual
func (m *Manager) GetLanguage() string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if m.current == nil {
		return m.defaultLang
	}
	return m.current.Code
}

// T traduz uma chave para o idioma atual
func (m *Manager) T(key string, args ...interface{}) string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if m.current == nil {
		return key
	}

	translation, exists := m.current.Translations[key]
	if !exists {
		// Fallback para a chave original
		return key
	}

	if len(args) > 0 {
		return fmt.Sprintf(translation, args...)
	}

	return translation
}

// GetAvailableLanguages retorna a lista de idiomas disponíveis
func (m *Manager) GetAvailableLanguages() []string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var languages []string
	for code := range m.locales {
		languages = append(languages, code)
	}
	return languages
}

// GetLanguageInfo retorna informações sobre um idioma específico
func (m *Manager) GetLanguageInfo(code string) (*Locale, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	locale, exists := m.locales[code]
	if !exists {
		return nil, fmt.Errorf("idioma não encontrado: %s", code)
	}

	return locale, nil
}

// SetupCobraI18n configura o Cobra para usar internacionalização
func (m *Manager) SetupCobraI18n(cmd *cobra.Command) {
	// Adicionar flag para idioma
	cmd.PersistentFlags().String("lang", "", "Idioma da interface (ex: pt-BR, en-US)")

	// Hook para processar a flag de idioma
	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if lang, _ := cmd.Flags().GetString("lang"); lang != "" {
			return m.SetLanguage(lang)
		}
		return nil
	}
}
