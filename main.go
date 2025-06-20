package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
)

func main() {
	loadDotEnv()

	// Initialize server
	// mux := http.NewServeMux()

	// Basic health check endpoint
	// mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
	// 	w.WriteHeader(http.StatusOK)
	// 	fmt.Fprint(w, "OK")
	// })

	// mux.HandleFunc("/sample", func(w http.ResponseWriter, r *http.Request) {
	// 	w.WriteHeader(http.StatusOK)
	// 	fmt.Fprint(w, "OK")
	// })

	// port := ":8080"
	// fmt.Printf("Server starting on port %s\n", port)

	// if err := http.ListenAndServe(port, mux); err != nil {
	// 	log.Fatal("Server failed to start:", err)
	// }

	runSample("http://localhost:8080/mcp")
}

func runSample(httpURL string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	httpTransport, err := transport.NewStreamableHTTP(httpURL)
	if err != nil {
		log.Fatalf("Failed to create HTTP transport: %v", err)
	}

	// Create client with the transport
	mcpClient := client.NewClient(httpTransport)
	// Start the client
	if err := mcpClient.Start(ctx); err != nil {
		log.Fatalf("Failed to start client: %v", err)
	}

	// Set up notification handler
	mcpClient.OnNotification(func(notification mcp.JSONRPCNotification) {
		fmt.Printf("Received notification: %s\n", notification.Method)
	})

	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "MCP-Go Simple Client Example",
		Version: "1.0.0",
	}
	initRequest.Params.Capabilities = mcp.ClientCapabilities{}

	serverInfo, err := mcpClient.Initialize(ctx, initRequest)
	if err != nil {
		log.Fatalf("Failed to initialize: %v", err)
	}

	// Display server information
	fmt.Printf("Connected to server: %s (version %s)\n",
		serverInfo.ServerInfo.Name,
		serverInfo.ServerInfo.Version)
	fmt.Printf("Server capabilities: %+v\n", serverInfo.Capabilities)

	var tools []mcp.Tool
	// List available tools if the server supports them
	if serverInfo.Capabilities.Tools != nil {
		fmt.Println("Fetching available tools...")
		toolsRequest := mcp.ListToolsRequest{}
		toolsResult, err := mcpClient.ListTools(ctx, toolsRequest)
		if err != nil {
			log.Printf("Failed to list tools: %v", err)
		} else {
			fmt.Printf("Server has %d tools available\n", len(toolsResult.Tools))
			for i, tool := range toolsResult.Tools {
				fmt.Printf("  %d. %s - %s - %s\n", i+1, tool.Name, tool.Description, tool.InputSchema)
			}
			tools = toolsResult.Tools
		}
	}

	// List available resources if the server supports them
	if serverInfo.Capabilities.Resources != nil {
		fmt.Println("Fetching available resources...")
		resourcesRequest := mcp.ListResourcesRequest{}
		resourcesResult, err := mcpClient.ListResources(ctx, resourcesRequest)
		if err != nil {
			log.Printf("Failed to list resources: %v", err)
		} else {
			fmt.Printf("Server has %d resources available\n", len(resourcesResult.Resources))
			for i, resource := range resourcesResult.Resources {
				fmt.Printf("  %d. %s - %s\n", i+1, resource.URI, resource.Name)
			}
			// resources = resourcesResult.Resources
		}
	}

	client := anthropic.NewClient(
		option.WithAPIKey(os.Getenv("ANTHROPIC_API_KEY")), // defaults to os.LookupEnv("ANTHROPIC_API_KEY")

	)

	messages := []anthropic.MessageParam{
		{
			Role: "user",
			Content: []anthropic.ContentBlockParamUnion{
				{
					OfText: &anthropic.TextBlockParam{
						Text: "Register a new customer",
					},
				},
			},
		},
	}

	// https://docs.anthropic.com/en/api/messages#auto
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		response, err := client.Messages.New(
			ctx,
			anthropic.MessageNewParams{
				//  "claude-3-5-sonnet-20240620"
				Model:     anthropic.ModelClaude3_7SonnetLatest,
				MaxTokens: 1024,
				Messages:  messages,
				System: []anthropic.TextBlockParam{
					{
						Text: "You are a helpful assistant that can use the tools provided to you. To manage a customer database",
					},
				},

				// Convert mcp.Tool from mcp client response to anthropic.ToolParam
				Tools: convertMcpToolToAnthropicTool(tools),
			},
		)
		if err != nil {
			log.Fatalf("Failed to send message: %v", err)
		}

		responseMessage := anthropic.MessageParam{
			Role:    "assistant",
			Content: []anthropic.ContentBlockParamUnion{},
		}
		toolResults := []anthropic.ToolResultBlockParam{}
		for _, content := range response.Content {
			if content.Type == "tool_use" {
				responseMessage.Content = append(responseMessage.Content, anthropic.ContentBlockParamUnion{
					OfToolUse: &anthropic.ToolUseBlockParam{
						ID:    content.ID,
						Name:  content.Name,
						Input: content.Input,
					},
				})
				// TODO: call tool
				callTool(content.Name, content.Input, tools, mcpClient)
				toolResults = append(toolResults, anthropic.ToolResultBlockParam{
					ToolUseID: content.ID,
					Content: []anthropic.ToolResultBlockParamContentUnion{
						{OfText: &anthropic.TextBlockParam{Text: "registered"}},
					},
				})
			}

			if content.Type == "text" {
				fmt.Printf("\033[94m%s\033[0m\n", content.Text)
				responseMessage.Content = append(responseMessage.Content, anthropic.ContentBlockParamUnion{
					OfText: &anthropic.TextBlockParam{Text: content.Text},
				})
			}
		}

		messages = append(messages, responseMessage)

		userInput := ""
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		userInput = strings.TrimSpace(text)

		if userInput == "exit" {
			break
		}

		userMessage := anthropic.MessageParam{
			Role:    "user",
			Content: []anthropic.ContentBlockParamUnion{},
		}
		for _, toolResult := range toolResults {
			userMessage.Content = append(userMessage.Content, anthropic.ContentBlockParamUnion{
				OfToolResult: &toolResult,
			})
		}

		userMessage.Content = append(userMessage.Content, anthropic.ContentBlockParamUnion{
			OfText: &anthropic.TextBlockParam{Text: userInput},
		})
		messages = append(messages, userMessage)
	}

	mcpClient.Close()
}

func callTool(name string, input []byte, tools []mcp.Tool, mcpClient *client.Client) {
	fmt.Printf("\033[33mTool use:%s - %s\033[0m\n", name, string(input))
	for _, tool := range tools {
		if tool.Name == name {
			fmt.Printf("Registering customer: %s\n", string(input))
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			var arguments map[string]any
			err := json.Unmarshal(input, &arguments)
			if err != nil {
				fmt.Printf("Error unmarshalling input: %v\n", err)
				return
			}
			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name:      tool.Name,
					Arguments: arguments,
				},
			}

			toolResult, err := mcpClient.CallTool(ctx, request)
			if err != nil {
				fmt.Printf("Error calling tool: %v\n", err)
			}
			fmt.Printf("Tool result: %+v\n", toolResult)
			return
		}
	}
}

func convertMcpToolToAnthropicTool(mcpTools []mcp.Tool) []anthropic.ToolUnionParam {
	anthropicTools := make([]anthropic.ToolParam, len(mcpTools))
	for i, mcpTool := range mcpTools {
		anthropicTools[i] = anthropic.ToolParam{
			Name:        mcpTool.Name,
			Description: anthropic.String(mcpTool.Description),
			InputSchema: anthropic.ToolInputSchemaParam{
				// Type: anthropic.Object,
				Properties: mcpTool.InputSchema.Properties,
				Required:   mcpTool.InputSchema.Required,
			},
		}
	}
	return convertToolsToToolUnionParam(anthropicTools)
}

func convertToolsToToolUnionParam(tools []anthropic.ToolParam) []anthropic.ToolUnionParam {
	toolUnionParams := make([]anthropic.ToolUnionParam, len(tools))
	for i, tool := range tools {
		toolUnionParams[i] = anthropic.ToolUnionParam{
			OfTool: &tool,
		}
	}
	return toolUnionParams
}

func loadDotEnv() {
	// Load .env file without using external packages
	envFile, err := os.Open(".env")
	if err != nil {
		log.Fatalf("Error opening .env file: %v", err)
	}
	defer envFile.Close()

	scanner := bufio.NewScanner(envFile)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		os.Setenv(parts[0], parts[1])
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading .env file: %v", err)
	}
}
