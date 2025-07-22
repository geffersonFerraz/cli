package cmd

import (
	"fmt"
	"os"

	"gfcli/i18n"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// i18nCmd representa o comando de internacionalização
var i18nCmd = &cobra.Command{
	Use:     "i18n",
	Short:   "Gerenciar idiomas da interface",
	GroupID: "other",
	Long: `Comando para gerenciar idiomas da interface da CLI.
Permite listar idiomas disponíveis, definir o idioma atual e obter informações sobre traduções.`,
}

// i18nListCmd lista os idiomas disponíveis
var i18nListCmd = &cobra.Command{
	Use:   "list",
	Short: "Listar idiomas disponíveis",
	Long:  "Lista todos os idiomas disponíveis na CLI com suas informações.",
	RunE: func(cmd *cobra.Command, args []string) error {
		manager := i18n.GetInstance()
		languages := manager.GetAvailableLanguages()
		currentLang := manager.GetLanguage()

		if len(languages) == 0 {
			fmt.Println("Nenhum idioma disponível.")
			return nil
		}

		headerColor := color.New(color.FgCyan, color.Bold)
		headerColor.Println("🌍 Idiomas Disponíveis:")
		fmt.Println()

		for _, code := range languages {
			info, err := manager.GetLanguageInfo(code)
			if err != nil {
				continue
			}

			// Destacar idioma atual
			if code == currentLang {
				currentColor := color.New(color.FgGreen, color.Bold)
				currentColor.Printf("  ✓ %s (%s)\n", info.NativeName, code)
			} else {
				fmt.Printf("    %s (%s)\n", info.NativeName, code)
			}
		}

		fmt.Println()
		noteColor := color.New(color.FgYellow)
		noteColor.Printf("Idioma atual: %s\n", currentLang)
		noteColor.Println("Use 'mgc i18n set <código>' para alterar o idioma.")

		return nil
	},
}

// i18nSetCmd define o idioma atual
var i18nSetCmd = &cobra.Command{
	Use:   "set [código]",
	Short: "Definir idioma atual",
	Long:  "Define o idioma atual da interface da CLI.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		code := args[0]
		manager := i18n.GetInstance()

		// Verificar se o idioma existe
		info, err := manager.GetLanguageInfo(code)
		if err != nil {
			return fmt.Errorf("idioma não encontrado: %s", code)
		}

		// Definir o idioma
		if err := manager.SetLanguage(code); err != nil {
			return err
		}

		successColor := color.New(color.FgGreen, color.Bold)
		successColor.Printf("✅ Idioma alterado para: %s (%s)\n", info.NativeName, code)

		// Mostrar como persistir a configuração
		fmt.Println()
		noteColor := color.New(color.FgYellow)
		noteColor.Println("Para persistir esta configuração, você pode:")
		fmt.Println("  1. Definir a variável de ambiente CLI_LANG:")
		fmt.Printf("     export CLI_LANG=%s\n", code)
		fmt.Println("  2. Usar a flag --lang em cada comando:")
		fmt.Printf("     cli --lang=%s [comando]\n", code)

		return nil
	},
}

// i18nInfoCmd mostra informações sobre um idioma específico
var i18nInfoCmd = &cobra.Command{
	Use:   "info [código]",
	Short: "Mostrar informações sobre um idioma",
	Long:  "Mostra informações detalhadas sobre um idioma específico.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		code := args[0]
		manager := i18n.GetInstance()

		info, err := manager.GetLanguageInfo(code)
		if err != nil {
			return fmt.Errorf("idioma não encontrado: %s", code)
		}

		headerColor := color.New(color.FgCyan, color.Bold)
		headerColor.Printf("📋 Informações do Idioma: %s\n\n", info.NativeName)

		fmt.Printf("Código: %s\n", info.Code)
		fmt.Printf("Nome: %s\n", info.Name)
		fmt.Printf("Total de Traduções: %d\n", len(info.Translations))

		// Mostrar algumas traduções de exemplo
		if len(info.Translations) > 0 {
			fmt.Println("\nExemplos de traduções:")
			count := 0
			for key, translation := range info.Translations {
				if count >= 5 { // Limitar a 5 exemplos
					break
				}
				fmt.Printf("  %s: %s\n", key, translation)
				count++
			}
		}

		return nil
	},
}

// i18nCurrentCmd mostra o idioma atual
var i18nCurrentCmd = &cobra.Command{
	Use:   "current",
	Short: "Mostrar idioma atual",
	Long:  "Mostra informações sobre o idioma atualmente em uso.",
	RunE: func(cmd *cobra.Command, args []string) error {
		manager := i18n.GetInstance()
		currentLang := manager.GetLanguage()

		info, err := manager.GetLanguageInfo(currentLang)
		if err != nil {
			return fmt.Errorf("erro ao obter informações do idioma atual: %s", err)
		}

		headerColor := color.New(color.FgCyan, color.Bold)
		headerColor.Println("🎯 Idioma Atual:")
		fmt.Println()

		fmt.Printf("Código: %s\n", info.Code)
		fmt.Printf("Nome: %s\n", info.Name)

		// Mostrar como o idioma foi detectado
		fmt.Println()
		noteColor := color.New(color.FgYellow)
		noteColor.Println("Detecção de idioma:")

		if os.Getenv("CLI_LANG") != "" {
			fmt.Printf("  Definido por CLI_LANG: %s\n", os.Getenv("CLI_LANG"))
		} else if os.Getenv("LANG") != "" {
			fmt.Printf("  Detectado de LANG: %s\n", os.Getenv("LANG"))
		} else if os.Getenv("LC_ALL") != "" {
			fmt.Printf("  Detectado de LC_ALL: %s\n", os.Getenv("LC_ALL"))
		} else {
			fmt.Println("  Usando idioma padrão")
		}

		return nil
	},
}

func init() {
	// Adicionar subcomandos ao comando i18n
	i18nCmd.AddCommand(i18nListCmd)
	i18nCmd.AddCommand(i18nSetCmd)
	i18nCmd.AddCommand(i18nInfoCmd)
	i18nCmd.AddCommand(i18nCurrentCmd)
}
