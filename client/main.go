package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "agent-tool-registry/registry" // Adjust module path if needed
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewToolRegistryClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// --- Agent A: Registers CalculateSum tool ---
	log.Println("\n--- Agent A registering CalculateSum ---")
	registerRespA, err := c.RegisterTool(ctx, &pb.ToolDefinition{
		Name:        "CalculateSum",
		Description: "Calculates the sum of two integers.",
		Parameters: map[string]*pb.ToolParameter{
			"a": {Type: pb.ToolParameter_INT, Required: true, Description: "First integer"},
			"b": {Type: pb.ToolParameter_INT, Required: true, Description: "Second integer"},
		},
		Return_Value: &pb.ToolParameter{Type: pb.ToolParameter_INT, Description: "The sum of a and b"},
	})
	if err != nil {
		log.Fatalf("could not register tool A: %v", err)
	}
	log.Printf("RegisterTool (Agent A) response: %v", registerRespA.Message)

	// --- Agent B: Registers FetchWeatherData tool ---
	log.Println("\n--- Agent B registering FetchWeatherData ---")
	registerRespB, err := c.RegisterTool(ctx, &pb.ToolDefinition{
		Name:        "FetchWeatherData",
		Description: "Fetches current weather data for a given location.",
		Parameters: map[string]*pb.ToolParameter{
			"location": {Type: pb.ToolParameter_STRING, Required: true, Description: "City name (e.g., 'London')"},
		},
		Return_Value: &pb.ToolParameter{Type: pb.ToolParameter_STRING, Description: "Current weather conditions"},
	})
	if err != nil {
		log.Fatalf("could not register tool B: %v", err)
	}
	log.Printf("RegisterTool (Agent B) response: %v", registerRespB.Message)

	// Give the server a moment to process (not strictly necessary with gRPC, but good for sequential logs)
	time.Sleep(50 * time.Millisecond)

	// --- Agent C: Discovers tools ---
	log.Println("\n--- Agent C discovering tools ---")
	listResp, err := c.ListTools(ctx, &pb.ListToolsRequest{})
	if err != nil {
		log.Fatalf("could not list tools: %v", err)
	}
	log.Println("Available tools:")
	for _, tool := range listResp.Tools {
		log.Printf("- %s: %s", tool.Name, tool.Description)
		for paramName, param := range tool.Parameters {
			log.Printf("  - Param: %s (Type: %s, Required: %t)", paramName, param.Type.String(), param.Required)
		}
	}

	// --- Agent C: Invokes CalculateSum tool ---
	log.Println("\n--- Agent C invoking CalculateSum ---")
	invokeSumResp, err := c.InvokeTool(ctx, &pb.ToolInvocationRequest{
		Tool_Name: "CalculateSum",
		Arguments: map[string]string{"a": "10", "b": "25"},
	})
	if err != nil {
		log.Fatalf("could not invoke CalculateSum: %v", err)
	}
	if invokeSumResp.Success {
		log.Printf("CalculateSum result: %s", invokeSumResp.Result)
	} else {
		log.Printf("CalculateSum failed: %s", invokeSumResp.Error)
	}

	// Give the server a moment
	time.Sleep(50 * time.Millisecond)

	// --- Agent C: Invokes FetchWeatherData tool ---
	log.Println("\n--- Agent C invoking FetchWeatherData ---")
	invokeWeatherResp, err := c.InvokeTool(ctx, &pb.ToolInvocationRequest{
		Tool_Name: "FetchWeatherData",
		Arguments: map[string]string{"location": "London"},
	})
	if err != nil {
		log.Fatalf("could not invoke FetchWeatherData: %v", err)
	}
	if invokeWeatherResp.Success {
		log.Printf("FetchWeatherData result: %s", invokeWeatherResp.Result)
	} else {
		log.Printf("FetchWeatherData failed: %s", invokeWeatherResp.Error)
	}

	// --- Agent C: Invokes a non-existent tool (error case) ---
	log.Println("\n--- Agent C invoking a non-existent tool ---")
	invokeNonExistentResp, err := c.InvokeTool(ctx, &pb.ToolInvocationRequest{
		Tool_Name: "NonExistentTool",
		Arguments: map[string]string{},
	})
	if err != nil {
		log.Printf("Invocation error (expected): %v", err) // gRPC error, not application error
	} else {
		if !invokeNonExistentResp.Success {
			log.Printf("NonExistentTool failed (expected): %s", invokeNonExistentResp.Error)
		} else {
			log.Printf("NonExistentTool unexpectedly succeeded: %s", invokeNonExistentResp.Result)
		}
	}
}
