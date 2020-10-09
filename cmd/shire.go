package cmd

import (
	"encoding/base64"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var shireCmd = &cobra.Command{
	Use:    "shire",
	Hidden: true,
	Short:  "A delightful refreshment",
	Long:   `Use 'shire' to produce output compliant with RFC 2324.`,
	Run: func(cmd *cobra.Command, args []string) {
		bytes, _ := base64.StdEncoding.DecodeString("ICAgICAgICAgICAgIDssJwogICAgIF9vXyAgICA7OjsnCiAsLS4nLS0tYC5fXyA7CigoamA9PT09PScsLScKIGAtXCAgICAgLwogICAgYC09LScgICAgIGhqdw==")
		shireStr := string(bytes)
		fmt.Println(shireStr)

		if viper.GetBool("verbose") {
			bytes, _ = base64.StdEncoding.DecodeString("QXJ0IGJ5IEhheWxleSBKYW5lIFdha2Vuc2hhdwpTb3VyY2U6IGh0dHBzOi8vd3d3LmFzY2lpYXJ0LmV1L2Zvb2QtYW5kLWRyaW5rcy9jb2ZmZWUtYW5kLXRlYQ==")
			credits := string(bytes)
			fmt.Println(credits)
		}
	},
}

func init() {
	rootCmd.AddCommand(shireCmd)
}
