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
	rootCmd.AddCommand(getTripCountsForCab)
	getTripCountsForCab.PersistentFlags().StringSliceP("cab-ids", "", []string{"D7D598CD99978BD012A87A76A7C891B7", "42D815590CE3A33F3A23DBF145EE66E3"}, "list of cab IDs to fetch")
	getTripCountsForCab.PersistentFlags().StringP("pickup-date", "", "2013-12-01", "pickup date")
	getTripCountsForCab.PersistentFlags().BoolP("ignore-cache", "", false, "Ignore cached data and force fetch DB")
}

var getTripCountsForCab = &cobra.Command{
	Use:   "get-trip-counts-for-cab",
	Short: "Prints cab trip count on given pickup date",
	Long: `Prints cab trip count on given pickup date
Example: ./ny_cab_client_rest get-trip-counts-for-cab --cab-ids="cab1,cab2" --pickup-date="2013-12-01" --ignore-cache=true`,
	Run: func(cmd *cobra.Command, args []string) {
		now := time.Now()
		log.Printf("getTripCountsForCab REST started at %s", now)
		defer trackTime(now, "getTripCountsForCab REST")
		server, _ := cmd.Flags().GetString("server")
		cabIds, _ := cmd.Flags().GetStringSlice("cab-ids")

		if len(cabIds) <= 0 {
			log.Fatal("empty cab ID list")
		}
		pickUpDate, _ := cmd.Flags().GetString("pickup-date")
		if pickUpDate == "" {
			log.Fatal("missing pickup-date")
		}

		ignoreCache, _ := cmd.Flags().GetBool("ignore-cache")

		var body string

		cbIDs := fmt.Sprintf("\"%s\"", strings.Join(cabIds, "\", \""))

		// Call GetTripCountsForCabIDsV1
		bodyRequest := fmt.Sprintf(`
		{
			"cab_ids": [%s],
			"pickup_date": "%s",
			"ignore_cache": %t
		}`, cbIDs, pickUpDate, ignoreCache)
		log.Println("body request: ", bodyRequest)
		resp, err := http.Post(server+"/v1/cabtrips/bypickupdate", "application/json", strings.NewReader(bodyRequest))
		if err != nil {
			log.Fatalf("failed to call GetTripCountsForCabIDsV1 method: %v", err)
		}
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			body = fmt.Sprintf("failed read GetTripCountsForCabIDsV1 response body: %v", err)
		} else {
			body = string(bodyBytes)
		}
		log.Printf("GetTripCountsForCabIDsV1 response: Code=%d, Body=%s\n\n", resp.StatusCode, body)
	},
}
