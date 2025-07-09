package static

import (
	"mgccli/cmd/static/config"

	sdk "github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/spf13/cobra"
)

func RootStatic(parent *cobra.Command, sdkCoreConfig sdk.CoreClient) {

	config.ConfigCmd(parent, sdkCoreConfig)

}
