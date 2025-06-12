package main


import (
//	"context"
	"log"
	"math/rand"
	"net"
//	"sync"
	"time"
	"google.golang.org/grpc"
//	"google.golang.org/grpc/credentials/insecure"
	stock "github.com/KingBean4903/StockTicker/stock"
)

type stockTickerServer struct {
	stock.UnimplementedStockTickerServer  
}


func (s *stockTickerServer) Subscribe(req *stock.StockRequest, stream stock.StockTicker_SubscribeServer) error {

	for  {
		
		for _, symbol := range req.Symbols {
				
			resp := &stock.StockResponse{
					Symbol: symbol,
					Price: generateFakePrice(symbol),
					Timestamp: time.Now().Unix(),
			}

			if err := stream.Send(resp); err != nil {
					log.Printf("Client disconnected: %v", err)
					return err
			}	
		}

		time.Sleep(1 * time.Second)

	}

}

func generateFakePrice(symbol string) float64 {
		basePrices := map[string]float64{
					"AAPL" : 170.0,
					"GOOG" : 148.0,
					"MSFT" : 410.0,
		}

		base := basePrices[symbol]
		fluctuation := rand.Float64() * 2.0 - 1.0
		return base + fluctuation
}



func main() {

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
			log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	stock.RegisterStockTickerServer(s, &stockTickerServer{})

	log.Printf("Server started at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to start server: %v", err)
	}

}




