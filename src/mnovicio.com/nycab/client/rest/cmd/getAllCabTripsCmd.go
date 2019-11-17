package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(getAllCabTrips)
	getAllCabTrips.PersistentFlags().BoolP("ignore-cache", "", false, "Ignore cached data and force fetch DB")
}

var getAllCabTrips = &cobra.Command{
	Use:   "get-all-cab-trip-count",
	Short: "Prints all cab trips on record",
	Long:  `Prints all cab trips on record`,
	Run: func(cmd *cobra.Command, args []string) {
		now := time.Now()
		log.Printf("getAllCabTrips REST started at %s", now)
		defer trackTime(now, "getAllCabTrips REST")
		server, _ := cmd.Flags().GetString("server")
		ignoreCache, _ := cmd.Flags().GetBool("ignore-cache")

		var body string

		// Call GetAllCabTripCountPerDayV1
		resp, err := http.Post(server+"/v1/cabtrips", "application/json", strings.NewReader(fmt.Sprintf(`
			{
				"ignore_cache": %t
			}
		`, ignoreCache)))
		if err != nil {
			log.Fatalf("failed to call GetAllCabTripCountPerDayV1 method: %v", err)
		}
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			body = fmt.Sprintf("failed read GetAllCabTripCountPerDayV1 response body: %v", err)
		} else {
			body = string(bodyBytes)
		}
		log.Printf("GetAllCabTripCountPerDayV1 response: Code=%d, Body=%s\n\n", resp.StatusCode, body)
	},
}
