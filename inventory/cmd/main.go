package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	partAPI "inventory/internal/api/part/v1"
	partRepository "inventory/internal/repository/part"
	partService "inventory/internal/service/part"
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

	repo := partRepository.NewRepository()
	service := partService.NewService(repo)
	api := partAPI.NewAPI(service)

	inventoryV1.RegisterInventoryServiceServer(s, api)

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
