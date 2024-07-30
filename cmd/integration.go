package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/metafates/mangal/icon"
	"github.com/metafates/mangal/integration/anilist"
	"github.com/metafates/mangal/key"
	"github.com/metafates/mangal/log"
	"github.com/metafates/mangal/open"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(integrationCmd)
	integrationCmd.AddCommand(integrationAnilistCmd)
	integrationAnilistCmd.Flags().BoolP("disable", "d", false, "Disable Anilist integration")
}

var integrationCmd = &cobra.Command{
	Use:   "integration",
	Short: "Integration with other services",
	Long:  `Integration with other services`,
}

var integrationAnilistCmd = &cobra.Command{
	Use:   "anilist",
	Short: "Integration with Anilist",
	Long: `Integration with Anilist.
See https://github.com/metafates/mangal/wiki/Anilist-Integration for more information`,
	Run: func(cmd *cobra.Command, args []string) {
		if lo.Must(cmd.Flags().GetBool("disable")) {
			viper.Set(key.AnilistEnable, false)
			viper.Set(key.AnilistCode, "")
			viper.Set(key.AnilistSecret, "")
			viper.Set(key.AnilistID, "")
			log.Info("Anilist integration disabled")
			handleErr(viper.WriteConfig())
		}

		if !viper.GetBool(key.AnilistEnable) {
			confirm := survey.Confirm{
				Message: "Anilist is disabled. Enable?",
				Default: false,
			}
			var response bool
			err := survey.AskOne(&confirm, &response)
			handleErr(err)

			if !response {
				return
			}

			viper.Set(key.AnilistEnable, response)
			err = viper.WriteConfig()
			if err != nil {
				switch err.(type) {
				case viper.ConfigFileNotFoundError:
					err = viper.SafeWriteConfig()
					handleErr(err)
				default:
					handleErr(err)
					log.Error(err)
				}
			}
		}

		if viper.GetString(key.AnilistID) == "" {
			input := survey.Input{
				Message: "Anilist client ID is not set. Please enter it:",
				Help:    "",
			}
			var response string
			err := survey.AskOne(&input, &response)
			handleErr(err)

			if response == "" {
				return
			}

			viper.Set(key.AnilistID, response)
			err = viper.WriteConfig()
			handleErr(err)
		}

		if viper.GetString(key.AnilistSecret) == "" {
			input := survey.Input{
				Message: "Anilist client secret is not set. Please enter it:",
				Help:    "",
			}
			var response string
			err := survey.AskOne(&input, &response)
			handleErr(err)

			if response == "" {
				return
			}

			viper.Set(key.AnilistSecret, response)
			err = viper.WriteConfig()
			handleErr(err)
		}

		if viper.GetString(key.AnilistCode) == "" {
			authURL := anilist.New().AuthURL()
			confirmOpenInBrowser := survey.Confirm{
				Message: "Open browser to authenticate with Anilist?",
				Default: false,
			}

			var openInBrowser bool
			err := survey.AskOne(&confirmOpenInBrowser, &openInBrowser)
			if err == nil && openInBrowser {
				err = open.Start(authURL)
			}

			if err != nil || !openInBrowser {
				fmt.Println("Please open the following URL in your browser:")
				fmt.Println(authURL)
			}

			input := survey.Input{
				Message: "Anilist code is not set. Please copy it from the link and paste in here:",
				Help:    "",
			}

			var response string
			err = survey.AskOne(&input, &response)
			handleErr(err)

			if response == "" {
				return
			}

			viper.Set(key.AnilistCode, response)
			err = viper.WriteConfig()
			handleErr(err)
		}

		fmt.Printf("%s Anilist integration was set up\n", icon.Get(icon.Success))
	},
}
