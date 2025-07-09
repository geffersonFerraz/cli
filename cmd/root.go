package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"mgccli/cmd/gen"
	"mgccli/cmd/static"
	"runtime"

	sdk "github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	apiKey := os.Getenv("MGC_API_KEY")
	if apiKey == "" {
		log.Fatal("Env MGC_API_KEY is required!")
	}
	sdkCoreConfig := sdk.NewMgcClient(apiKey,
		sdk.WithUserAgent(fmt.Sprintf("MgcCLI2/%s (%s; %s)", "0.5.0", runtime.GOOS, runtime.GOARCH)),
	)

	var rootCmd = &cobra.Command{
		Use:   "mgc",
		Short: "MGC CLI",
		Long:  `MGC CLI`,
	}

	rootCmd.AddGroup(&cobra.Group{
		ID:    "products",
		Title: "Products:",
	})
	rootCmd.AddGroup(&cobra.Group{
		ID:    "settings",
		Title: "Settings:",
	})
	rootCmd.AddGroup(&cobra.Group{
		ID:    "other",
		Title: "Other commands:",
	})
	rootCmd.SetHelpCommandGroupID("other")
	rootCmd.SetCompletionCommandGroupID("other")

	addApiKeyFlag(rootCmd)
	addLogDebugFlag(rootCmd)
	addOutputFlag(rootCmd)
	addNoConfirmationFlag(rootCmd)
	addRawOutputFlag(rootCmd)

	ctx := context.Background()

	static.RootStatic(rootCmd, *sdkCoreConfig)
	gen.RootGen(ctx, rootCmd, *sdkCoreConfig)

	return rootCmd
}
