package main

import (
		"context"
		"log"
		"os"
		"os/signal"
		"syscall"
		"time"

		"google.golang.org/grpc"
		"google.golang.org/grpc/credentials/insecure"
		stock "github.com/KingBean4903/StockTicker/stock"
)

func main() {
	
	conn, err := grpc.Dial(
					"localhost:50051",
					grpc.WithTransportCredentials(insecure.NewCredentials()),
					grpc.WithBlock(),
	)

	if err != nil {
				log.Fatalf("failed to connect: %v", err)
	}

	defer conn.Close()

	client := stock.NewStockTickerClient(conn)

	req := &stock.StockRequest{
					Symbols: []string{"AAPL", "GOOG", "MSFT"},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := client.Subscribe(ctx, req)
	if err != nil {
			log.Fatalf("Subscribe failed: %v", err)
	}

	log.Printf("Subscribed to updates for: %v", req.Symbols)
	log.Println("(Press Ctrl+C to stop)")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	for {
			
			select {
			case <-sigCh:
							log.Println("Shutting down ... ")
							cancel()
							return
			default: 
							resp, err := stream.Recv()
							if err != nil {
											log.Printf("Stream closed: %v", err)
											return
							}

							log.Printf("%s: $%.2f (%s)", 
													resp.Symbol,
												  resp.Price,
													time.Unix(resp.Timestamp, 0).Format("2006-01-02 15:04:05"),
								)
			}

	}

}
