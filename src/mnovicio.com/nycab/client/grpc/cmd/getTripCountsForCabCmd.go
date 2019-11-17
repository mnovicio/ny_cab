package cmd

import (
	"context"
	"log"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	pbsvc "mnovicio.com/nycab/protocol/rpc"
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
Example: ./ny_cab_client_grpc get-trip-counts-for-cab --cab-ids="cab1,cab2" --pickup-date="2013-12-01" --ignore-cache=true`,
	Run: func(cmd *cobra.Command, args []string) {
		now := time.Now()
		log.Printf("getTripCountsForCab gRPC started at %s", now)
		defer trackTime(now, "getTripCountsForCab gRPC")
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

		log.Printf("Dialing gRPC server: %s", server)
		conn, err := grpc.Dial(server, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Unable to connect to NY CAB gRPC server at [%s]", server)
		}

		nyCabClient := pbsvc.NewNYCabServiceClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
		defer cancel()

		request := &pbsvc.GetTripCountsForCabIDsRequestV1{
			CabIds:      cabIds,
			PickupDate:  pickUpDate,
			IgnoreCache: ignoreCache,
		}

		response, err := nyCabClient.GetTripCountsForCabIDsV1(ctx, request)
		if err != nil {
			log.Fatalf("Failed calling GetTripCountsForCabIDsV1 RPC from %s", server)
		}

		log.Printf("GetTripCountsForCabIDsV1 response=[%+v]", response)

	},
}
