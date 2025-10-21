package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	inventoryV1 "shared/pkg/proto/inventory/v1"
)

const grpcPort = 50051

type InventoryService struct {
	inventoryV1.UnimplementedInventoryServiceServer
	mu    sync.RWMutex
	parts map[string]*inventoryV1.Part
}

func (s *InventoryService) GetPart(_ context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	if len(req.Uuid) == 0 {
		return nil, status.Error(codes.InvalidArgument, "uuid is empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	part, ok := s.parts[req.GetUuid()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "part with uuid %s not found", req.GetUuid())
	}

	return &inventoryV1.GetPartResponse{
		Part: part,
	}, nil
}

// Sort our Parts by Filter adn return list of them
func (s *InventoryService) ListParts(_ context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	filtredParts := FilterParts(s.parts, req.Filter)

	if len(filtredParts) == 0 {
		return nil, status.Errorf(codes.NotFound, "parts with this filter %s were not found", req.Filter)
	}

	return &inventoryV1.ListPartsResponse{
		Parts: filtredParts,
	}, nil
}

func FilterParts(parts map[string]*inventoryV1.Part, Filter *inventoryV1.PartsFilter) []*inventoryV1.Part {
	var PartsToReturn []*inventoryV1.Part
	if Filter == nil {
		PartsToReturn = make([]*inventoryV1.Part, 0, len(parts))
		for i := range parts {
			PartsToReturn = append(PartsToReturn, parts[i])
		}
		return PartsToReturn
	}

	for _, part := range parts {
		// Filter by UUID
		if len(Filter.Uuids) > 0 && !contains(Filter.Uuids, part.Uuid) {
			continue
		}
		if len(Filter.Names) > 0 && !contains(Filter.Names, part.Name) {
			continue
		}
		if len(Filter.Categories) > 0 {
			found := false
			for _, v := range Filter.Categories {
				if v == part.Category {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		if len(Filter.ManufacturerCountries) > 0 && !contains(Filter.ManufacturerCountries, part.Manufacturer.Country) {
			continue
		}
		if len(Filter.Tags) > 0 && !anyMatch(part.Tags, Filter.Tags) {
			continue
		}

		PartsToReturn = append(PartsToReturn, part)

	}
	return PartsToReturn
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// checks tags intersections empty or not
func anyMatch(tags, filter []string) bool {
	for _, tag := range tags {
		for _, f := range filter {
			if tag == f {
				return true
			}
		}
	}
	return false
}

func (s *InventoryService) initParts() {
	parts := generateParts()

	for _, part := range parts {
		s.parts[part.Uuid] = part
	}
}

func generateParts() []*inventoryV1.Part {
	names := []string{
		"Main Engine",
		"Reserve Engine",
		"Thruster",
		"Fuel Tank",
		"Left Wing",
		"Right Wing",
		"Window A",
		"Window B",
		"Control Module",
		"Stabilizer",
	}

	descriptions := []string{
		"Primary propulsion unit",
		"Backup propulsion unit",
		"Thruster for fine adjustments",
		"Main fuel tank",
		"Left aerodynamic wing",
		"Right aerodynamic wing",
		"Front viewing window",
		"Side viewing window",
		"Flight control module",
		"Stabilization fin",
	}

	var parts []*inventoryV1.Part
	for i := 0; i < gofakeit.Number(1, 50); i++ {
		idx := gofakeit.Number(0, len(names)-1)
		parts = append(parts, &inventoryV1.Part{
			Uuid:          uuid.NewString(),
			Name:          names[idx],
			Description:   descriptions[idx],
			Price:         roundTo(gofakeit.Float64Range(100, 10_000)),
			StockQuantity: int64(gofakeit.Number(1, 100)),
			Category:      inventoryV1.Category(gofakeit.Number(1, 4)), //nolint:gosec // safe: gofakeit.Number returns 1..4
			Dimensions:    generateDimensions(),
			Manufacturer:  generateManufacturer(),
			Tags:          generateTags(),
			Metadata:      generateMetadata(),
			CreatedAt:     timestamppb.Now(),
		})
	}

	return parts
}

func generateDimensions() *inventoryV1.Dimensions {
	return &inventoryV1.Dimensions{
		Length: roundTo(gofakeit.Float64Range(1, 1000)),
		Width:  roundTo(gofakeit.Float64Range(1, 1000)),
		Height: roundTo(gofakeit.Float64Range(1, 1000)),
		Weight: roundTo(gofakeit.Float64Range(1, 1000)),
	}
}

func generateManufacturer() *inventoryV1.Manufacturer {
	return &inventoryV1.Manufacturer{
		Name:    gofakeit.Name(),
		Country: gofakeit.Country(),
		Website: gofakeit.URL(),
	}
}

func generateTags() []string {
	var tags []string
	for i := 0; i < gofakeit.Number(1, 10); i++ {
		tags = append(tags, gofakeit.EmojiTag())
	}

	return tags
}

func generateMetadata() map[string]*inventoryV1.Value {
	metadata := make(map[string]*inventoryV1.Value)

	for i := 0; i < gofakeit.Number(1, 10); i++ {
		metadata[gofakeit.Word()] = generateMetadataValue()
	}

	return metadata
}

func generateMetadataValue() *inventoryV1.Value {
	switch gofakeit.Number(0, 3) {
	case 0:
		return &inventoryV1.Value{
			Kind: &inventoryV1.Value_StringValue{
				StringValue: gofakeit.Word(),
			},
		}

	case 1:
		return &inventoryV1.Value{
			Kind: &inventoryV1.Value_Int64Value{
				Int64Value: int64(gofakeit.Number(1, 100)),
			},
		}

	case 2:
		return &inventoryV1.Value{
			Kind: &inventoryV1.Value_DoubleValue{
				DoubleValue: roundTo(gofakeit.Float64Range(1, 100)),
			},
		}

	case 3:
		return &inventoryV1.Value{
			Kind: &inventoryV1.Value_BoolValue{
				BoolValue: gofakeit.Bool(),
			},
		}

	default:
		return nil
	}
}

func roundTo(x float64) float64 {
	return math.Round(x*100) / 100
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()

	s := grpc.NewServer()

	service := &InventoryService{
		parts: make(map[string]*inventoryV1.Part),
	}

	inventoryV1.RegisterInventoryServiceServer(s, service)

	service.initParts()
	// Включаем рефлексию для отладки
	reflection.Register(s)

	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}
