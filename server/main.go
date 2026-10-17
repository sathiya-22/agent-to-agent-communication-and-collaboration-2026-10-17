package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
	pb "agent-tool-registry/registry" // Adjust module path if needed
)

// toolRegistryServer implements the ToolRegistry service.
type toolRegistryServer struct {
	pb.UnimplementedToolRegistryServer
	mu    sync.Mutex
	tools map[string]*pb.ToolDefinition
}

func newToolRegistryServer() *toolRegistryServer {
	return &toolRegistryServer{
		tools: make(map[string]*pb.ToolDefinition),
	}
}

// RegisterTool implements pb.ToolRegistryServer.
func (s *toolRegistryServer) RegisterTool(ctx context.Context, in *pb.ToolDefinition) (*pb.RegisterToolResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tools[in.Name]; exists {
		return &pb.RegisterToolResponse{Success: false, Message: fmt.Sprintf("Tool '%s' already registered", in.Name)}, nil
	}

	s.tools[in.Name] = in
	log.Printf("Registered tool: %s", in.Name)
	return &pb.RegisterToolResponse{Success: true, Message: fmt.Sprintf("Tool '%s' registered successfully", in.Name)}, nil
}

// ListTools implements pb.ToolRegistryServer.
func (s *toolRegistryServer) ListTools(ctx context.Context, in *pb.ListToolsRequest) (*pb.ListToolsResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var toolList []*pb.ToolDefinition
	for _, tool := range s.tools {
		// In a real system, we might filter by agent_id, but for this prototype, we list all.
		toolList = append(toolList, tool)
	}
	return &pb.ListToolsResponse{Tools: toolList}, nil
}

// InvokeTool implements pb.ToolRegistryServer.
func (s *toolRegistryServer) InvokeTool(ctx context.Context, in *pb.ToolInvocationRequest) (*pb.ToolInvocationResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	tool, exists := s.tools[in.Tool_Name]
	if !exists {
		return &pb.ToolInvocationResponse{Success: false, Error: fmt.Sprintf("Tool '%s' not found", in.Tool_Name)}, nil
	}

	log.Printf("Invoking tool '%s' with arguments: %v", in.Tool_Name, in.Arguments)

	// --- FIXTURE DATA / MOCK TOOL EXECUTION ---
	// In a real multi-agent system, this would involve routing the invocation
	// to the actual agent that registered the tool, waiting for its response,
	// and then returning it. For this prototype, we simulate the execution.
	switch in.Tool_Name {
	case "CalculateSum":
		aStr, okA := in.Arguments["a"]
		bStr, okB := in.Arguments["b"]
		if !okA || !okB {
			return &pb.ToolInvocationResponse{Success: false, Error: "Missing 'a' or 'b' arguments for CalculateSum"}, nil
		}
		var a, b int
		_, errA := fmt.Sscanf(aStr, "%d", &a)
		_, errB := fmt.Sscanf(bStr, "%d", &b)
		if errA != nil || errB != nil {
			return &pb.ToolInvocationResponse{Success: false, Error: "Invalid integer arguments for CalculateSum"}, nil
		}
		sum := a + b
		return &pb.ToolInvocationResponse{Success: true, Result: fmt.Sprintf("%d", sum)}, nil
	case "FetchWeatherData":
		location, ok := in.Arguments["location"]
		if !ok {
			return &pb.ToolInvocationResponse{Success: false, Error: "Missing 'location' argument for FetchWeatherData"}, nil
		}
		// Mock weather data
		weather := fmt.Sprintf("Weather in %s: Sunny, 25°C", location)
		return &pb.ToolInvocationResponse{Success: true, Result: weather}, nil
	default:
		return &pb.ToolInvocationResponse{Success: false, Error: fmt.Sprintf("Tool '%s' has no mock implementation", in.Tool_Name)}, nil
	}
	// --- END FIXTURE DATA ---
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterToolRegistryServer(s, newToolRegistryServer())
	log.Printf("Tool Registry Service started on %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
