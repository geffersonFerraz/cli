package cmd

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"mgccli/beautiful"
	"mgccli/cmd/gen"
	"mgccli/cmd/static"
	"mgccli/i18n"
	"runtime"

	sdk "github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func RootCmd(ctx context.Context, version string, manager *i18n.Manager) *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:     "mgc",
		Short:   manager.T("cli.short_description"),
		Long:    manager.T("cli.long_description"),
		Version: version,
	}

	rootCmd.SetContext(ctx)

	rootCmd.AddGroup(&cobra.Group{
		ID:    "products",
		Title: manager.T("cli.products_group"),
	})
	rootCmd.AddGroup(&cobra.Group{
		ID:    "settings",
		Title: manager.T("cli.settings_group"),
	})
	rootCmd.AddGroup(&cobra.Group{
		ID:    "other",
		Title: manager.T("cli.other_group"),
	})
	rootCmd.SetHelpCommandGroupID("other")
	rootCmd.SetCompletionCommandGroupID("other")

	addApiKeyFlag(rootCmd)
	addLogDebugFlag(rootCmd)
	addNoConfirmationFlag(rootCmd)
	addRawOutputFlag(rootCmd)
	addLangFlag(rootCmd)

	// Init SDK
	apiKey := os.Getenv("CLI_API_KEY")
	if apiKey == "" {
		apiKey = getApiKeyFlag(rootCmd)
		if apiKey == "" {
			log.Fatal(manager.T("cli.api_key_required"))
		}
	}

	sdkOptions := []sdk.Option{}
	debugLevel := getLogDebugFlag(rootCmd)
	sdkOptions = append(sdkOptions, sdk.WithLogger(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.Level(debugLevel)}))))
	sdkOptions = append(sdkOptions, sdk.WithUserAgent(fmt.Sprintf("CommunityCLI/%s (%s; %s)", version, runtime.GOOS, runtime.GOARCH)))

	sdkCoreConfig := sdk.NewMgcClient(apiKey,
		sdkOptions...,
	)

	static.RootStatic(rootCmd, *sdkCoreConfig)
	gen.RootGen(ctx, rootCmd, *sdkCoreConfig)

	// Adicionar comando i18n
	rootCmd.AddCommand(i18nCmd)

	// Aplicar embelezamento
	beautifulPrint(rootCmd)

	return rootCmd
}

func beautifulPrint(cmd *cobra.Command) {
	manager := i18n.GetInstance()

	// Configurar templates personalizados
	cmd.SetHelpTemplate(helpTemplate())
	cmd.SetUsageTemplate(usageTemplate(manager))

	// Configurar função de formatação personalizada
	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		// Cabeçalho colorido
		headerColor := color.New(color.FgCyan, color.Bold)
		headerColor.Printf("%s\n\n", cmd.Short)

		// Descrição longa
		if cmd.Long != "" {
			descColor := color.New(color.FgWhite)
			descColor.Println(cmd.Long)
			fmt.Println()
		}

		// Uso do comando
		if cmd.Runnable() {
			usageColor := color.New(color.FgYellow, color.Bold)
			usageColor.Print(manager.T("cli.usage") + ": ")
			usageText := color.New(color.FgWhite)
			usageText.Printf("%s", cmd.UseLine())
			fmt.Println()
		}

		// Comandos disponíveis organizados por grupos
		if cmd.HasAvailableSubCommands() {
			fmt.Println()
			subCmdColor := color.New(color.FgGreen, color.Bold)
			subCmdColor.Println(manager.T("cli.available_commands") + ":")

			// Organizar comandos por grupos
			commandsByGroup := make(map[string][]*cobra.Command)
			var ungroupedCommands []*cobra.Command

			for _, subCmd := range cmd.Commands() {
				if subCmd.IsAvailableCommand() || subCmd.Name() == "help" {
					if subCmd.GroupID != "" {
						commandsByGroup[subCmd.GroupID] = append(commandsByGroup[subCmd.GroupID], subCmd)
					} else {
						ungroupedCommands = append(ungroupedCommands, subCmd)
					}
				}
			}

			// Exibir comandos agrupados
			for _, group := range cmd.Groups() {
				if commands, exists := commandsByGroup[group.ID]; exists && len(commands) > 0 {
					fmt.Println()
					groupColor := color.New(color.FgMagenta, color.Bold)
					groupColor.Printf("%s\n", group.Title)

					for _, subCmd := range commands {
						cmdName := color.New(color.FgCyan, color.Bold)
						cmdName.Printf("  %-20s", subCmd.Name())
						cmdDesc := color.New(color.FgWhite)
						cmdDesc.Printf("%s\n", subCmd.Short)
					}
				}
			}

			// Exibir comandos sem grupo
			if len(ungroupedCommands) > 0 {
				fmt.Println()
				ungroupedColor := color.New(color.FgMagenta, color.Bold)
				ungroupedColor.Println("🔧 " + manager.T("cli.other_commands") + ":")

				for _, subCmd := range ungroupedCommands {
					cmdName := color.New(color.FgCyan, color.Bold)
					cmdName.Printf("  %-20s", subCmd.Name())
					cmdDesc := color.New(color.FgWhite)
					cmdDesc.Printf("%s\n", subCmd.Short)
				}
			}
		}

		// Flags locais
		if cmd.HasAvailableLocalFlags() {
			fmt.Println()
			flagColor := color.New(color.FgMagenta, color.Bold)
			flagColor.Println(manager.T("cli.local_flags") + ":")
			cmd.LocalFlags().VisitAll(func(flag *pflag.Flag) {
				if !flag.Hidden {
					flagName := color.New(color.FgYellow)
					flagName.Printf("  --%-15s", flag.Name)
					if flag.Shorthand != "" {
						shorthand := color.New(color.FgYellow)
						shorthand.Printf(" -%s", flag.Shorthand)
					}
					flagDesc := color.New(color.FgWhite)
					flagDesc.Printf(" %s\n", flag.Usage)
				}
			})
		}

		// Flags herdadas
		if cmd.HasAvailableInheritedFlags() {
			fmt.Println()
			flagColor := color.New(color.FgMagenta, color.Bold)
			flagColor.Println(manager.T("cli.global_flags") + ":")
			cmd.InheritedFlags().VisitAll(func(flag *pflag.Flag) {
				if !flag.Hidden {
					flagName := color.New(color.FgYellow)
					flagName.Printf("  --%-15s", flag.Name)
					if flag.Shorthand != "" {
						shorthand := color.New(color.FgYellow)
						shorthand.Printf(" -%s", flag.Shorthand)
					}
					flagDesc := color.New(color.FgWhite)
					flagDesc.Printf(" %s\n", flag.Usage)
				}
			})
		}

		// Exemplos
		if cmd.HasExample() {
			fmt.Println()
			exampleColor := color.New(color.FgGreen, color.Bold)
			exampleColor.Println(manager.T("cli.examples") + ":")
			exampleText := color.New(color.FgWhite)
			exampleText.Printf("%s\n", cmd.Example)
		}

		// Footer
		if cmd.HasAvailableSubCommands() {
			fmt.Println()
			footerColor := color.New(color.FgBlue, color.Italic)
			footerColor.Printf(manager.T("cli.help_more_info")+"\n", cmd.CommandPath())
		}
	})

	// Configurar função de erro personalizada
	cmd.SetUsageFunc(func(cmd *cobra.Command) error {
		errorColor := color.New(color.FgRed, color.Bold)
		errorColor.Print("❌ " + manager.T("cli.error") + ": ")
		errorText := color.New(color.FgWhite)
		errorText.Printf("Uso incorreto do comando\n\n")

		usageColor := color.New(color.FgYellow, color.Bold)
		usageColor.Print(manager.T("cli.usage") + ": ")
		usageText := color.New(color.FgWhite)
		usageText.Printf("%s\n", cmd.UseLine())

		if cmd.HasAvailableSubCommands() {
			fmt.Println()
			subCmdColor := color.New(color.FgGreen, color.Bold)
			subCmdColor.Println(manager.T("cli.available_commands") + ":")

			// Organizar comandos por grupos
			commandsByGroup := make(map[string][]*cobra.Command)
			var ungroupedCommands []*cobra.Command

			for _, subCmd := range cmd.Commands() {
				if subCmd.IsAvailableCommand() {
					if subCmd.GroupID != "" {
						commandsByGroup[subCmd.GroupID] = append(commandsByGroup[subCmd.GroupID], subCmd)
					} else {
						ungroupedCommands = append(ungroupedCommands, subCmd)
					}
				}
			}

			// Exibir comandos agrupados
			for _, group := range cmd.Groups() {
				if commands, exists := commandsByGroup[group.ID]; exists && len(commands) > 0 {
					fmt.Println()
					groupColor := color.New(color.FgMagenta, color.Bold)
					groupColor.Printf("%s\n", group.Title)

					for _, subCmd := range commands {
						cmdName := color.New(color.FgCyan)
						cmdName.Printf("  %-20s", subCmd.Name())
						cmdDesc := color.New(color.FgWhite)
						cmdDesc.Printf("%s\n", subCmd.Short)
					}
				}
			}

			// Exibir comandos sem grupo
			if len(ungroupedCommands) > 0 {
				fmt.Println()
				ungroupedColor := color.New(color.FgMagenta, color.Bold)
				ungroupedColor.Println("🔧 " + manager.T("cli.other_commands") + ":")

				for _, subCmd := range ungroupedCommands {
					cmdName := color.New(color.FgCyan)
					cmdName.Printf("  %-20s", subCmd.Name())
					cmdDesc := color.New(color.FgWhite)
					cmdDesc.Printf("%s\n", subCmd.Short)
				}
			}
		}

		return nil
	})

	// Configurar função de execução personalizada para interceptar outputs
	originalRunE := cmd.RunE
	if originalRunE != nil {
		cmd.RunE = func(cmd *cobra.Command, args []string) error {
			// Interceptar e embelezar outputs antes de executar
			rawMode := getRawOutputFlag(cmd)
			beautifulOutput := beautiful.NewOutput(rawMode)

			// Executar o comando original
			err := originalRunE(cmd, args)

			// Se houver erro, embelezar a mensagem
			if err != nil {
				beautifulOutput.PrintError(err.Error())
			}

			return err
		}
	}

	// Aplicar recursivamente para todos os subcomandos
	for _, subCmd := range cmd.Commands() {
		beautifulPrint(subCmd)
	}
}

func usageTemplate(manager *i18n.Manager) string {
	usageTemplate := `{{if .Runnable}}` + manager.T("cli.usage") + `:{{if .HasAvailableFlags}} [FLAGS]{{end}}{{if .HasAvailableSubCommands}} [COMANDO]{{end}}{{if gt .Aliases 0}}

	` + manager.T("cli.aliases") + `:
	  {{.NameAndAliases}}{{end}}{{if .HasExample}}
	
	` + manager.T("cli.examples") + `:
	{{.Example}}{{end}}{{if .HasAvailableSubCommands}}
	
	` + manager.T("cli.available_commands") + `:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
	  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}
	
	` + manager.T("cli.local_flags") + `:
	{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}
	
	` + manager.T("cli.global_flags") + `:
	{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}
	
	` + manager.T("cli.additional_help_commands") + `:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
	  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}
	
	` + manager.T("cli.help_more_info") + `.{{end}}
	`
	return usageTemplate
}

func helpTemplate() string {
	return `{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`
}
