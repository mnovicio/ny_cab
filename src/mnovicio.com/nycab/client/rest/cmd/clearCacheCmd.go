package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(clearCacheCmd)
}

var clearCacheCmd = &cobra.Command{
	Use:   "clear-cache",
	Short: "Clears cached data on the server",
	Long:  `Clears cached data on the server`,
	Run: func(cmd *cobra.Command, args []string) {
		now := time.Now()
		log.Printf("clearCacheCmd REST started at %s", now)
		defer trackTime(now, "clearCacheCmd REST")
		server, _ := cmd.Flags().GetString("server")

		var body string

		// Call ClearCacheV1
		resp, err := http.Get(server + "/v1/cabtrips/clearcache")
		if err != nil {
			log.Fatalf("failed to call ClearCacheV1 method: %v", err)
		}
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			body = fmt.Sprintf("failed read ClearCacheV1 response body: %v", err)
		} else {
			body = string(bodyBytes)
		}
		log.Printf("ClearCacheV1 response: Code=%d, Body=%s\n\n", resp.StatusCode, body)
	},
}
