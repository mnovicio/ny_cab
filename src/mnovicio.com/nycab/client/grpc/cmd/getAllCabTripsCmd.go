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
	rootCmd.AddCommand(getAllCabTrips)
	getAllCabTrips.PersistentFlags().BoolP("ignore-cache", "", false, "Ignore cached data and force fetch DB")
}

var getAllCabTrips = &cobra.Command{
	Use:   "get-all-cab-trip-count",
	Short: "Prints all cab trips on record",
	Long:  `Prints all cab trips on record`,
	Run: func(cmd *cobra.Command, args []string) {
		now := time.Now()
		log.Printf("getAllCabTrips gRPC started at %s", now)
		defer trackTime(now, "getAllCabTrips gRPC")
		server, _ := cmd.Flags().GetString("server")
		ignoreCache, _ := cmd.Flags().GetBool("ignore-cache")

		log.Printf("Dialing gRPC server: %s", server)
		conn, err := grpc.Dial(server, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Unable to connect to NY CAB gRPC server at [%s]", server)
		}

		nyCabClient := pbsvc.NewNYCabServiceClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
		defer cancel()

		request := &pbsvc.GetAllCabTripsRequestV1{
			IgnoreCache: ignoreCache,
		}

		response, err := nyCabClient.GetAllCabTripCountPerDayV1(ctx, request)
		if err != nil {
			log.Fatalf("Failed calling GetAllCabTripCountPerDayV1 RPC from %s", server)
		}

		log.Printf("GetAllCabTripCountPerDayV1 response=[%+v]", response)
	},
}
