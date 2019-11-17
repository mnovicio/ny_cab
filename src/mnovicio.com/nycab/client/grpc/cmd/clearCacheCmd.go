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
	rootCmd.AddCommand(clearCacheCmd)
}

var clearCacheCmd = &cobra.Command{
	Use:   "clear-cache",
	Short: "Clears cached data on the server",
	Long:  `Clears cached data on the server`,
	Run: func(cmd *cobra.Command, args []string) {
		now := time.Now()
		log.Printf("clearCacheCmd gRPC started at %s", now)
		defer trackTime(now, "clearCacheCmd gRPC")
		server, _ := cmd.Flags().GetString("server")

		log.Printf("Dialing gRPC server: %s", server)
		conn, err := grpc.Dial(server, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Unable to connect to NY CAB gRPC server at [%s]", server)
		}

		nyCabClient := pbsvc.NewNYCabServiceClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
		defer cancel()

		request := &pbsvc.ClearCacheRequestV1{}

		response, err := nyCabClient.ClearCacheV1(ctx, request)
		if err != nil {
			log.Fatalf("Failed calling ClearCacheV1 RPC from %s", server)
		}

		log.Printf("ClearCacheV1 response=[%+v]", response)
	},
}
