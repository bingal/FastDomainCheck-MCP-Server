package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/bingal/FastDomainCheck-MCP-Server/checker"
	"github.com/bingal/FastDomainCheck-MCP-Server/config"
	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
)

// CheckDomainsRequest represents a bulk domain check request
type CheckDomainsRequest struct {
	Domains []string `json:"domains" jsonschema:"required,description=List of domains to check"`
}

// SimpleDomainStatus represents the simplified domain status
type SimpleDomainStatus struct {
	Registered bool `json:"registered"` // true: registered, false: available
}

// SimpleCheckDomainsResponse represents the simplified bulk domain check response
type SimpleCheckDomainsResponse struct {
	Results map[string]SimpleDomainStatus `json:"results"` // Domain check results
}

func main() {
	// Parse command line flags
	skipHealthCheck := flag.Bool("skip-health-check", false, "Skip starting the health check server")
	flag.Parse()

	// Create configuration and domain checker
	cfg := config.NewConfig()
	domainChecker := checker.NewDomainChecker(cfg)

	// Start health check server if not skipped
	if !*skipHealthCheck {
		go func() {
			log.Println("Starting health check server on :8080")
			http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
			})
			if err := http.ListenAndServe(":8080", nil); err != nil {
				log.Printf("Health check server error: %v", err)
			}
		}()
	}

	// Create STDIO transport server
	server := mcp_golang.NewServer(stdio.NewStdioServerTransport())

	// Register bulk domain check tool
	err := server.RegisterTool("check_domains",
		`Check if multiple domain names are registered.

Input: A list of domain names to check (e.g. ["example.com", "test.com"])
Output: JSON object containing registration status of each domain:
{
  "results": {
    "example.com": {
      "registered": true
    },
    "test.com": {
      "registered": false
    }
  }
}`,
		func(arguments CheckDomainsRequest) (*mcp_golang.ToolResponse, error) {
			// Input validation
			if len(arguments.Domains) == 0 {
				return nil, fmt.Errorf("domains list cannot be empty")
			}
			if len(arguments.Domains) > 50 {
				return nil, fmt.Errorf("maximum 50 domains allowed per request")
			}

			// Check domain format
			for _, domain := range arguments.Domains {
				if len(domain) == 0 {
					return nil, fmt.Errorf("domain name cannot be empty")
				}
				if len(domain) > 255 {
					return nil, fmt.Errorf("domain name cannot exceed 255 characters: %s", domain)
				}
			}

			// Perform domain check
			results := domainChecker.CheckDomains(arguments.Domains)

			// Convert results to simplified format
			response := SimpleCheckDomainsResponse{
				Results: make(map[string]SimpleDomainStatus),
			}
			for i, domain := range arguments.Domains {
				response.Results[domain] = SimpleDomainStatus{
					Registered: results[i].Status == "registered",
				}
			}

			// Serialize results
			jsonData, err := json.Marshal(response)
			if err != nil {
				return nil, fmt.Errorf("failed to serialize results: %v", err)
			}

			return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(string(jsonData))), nil
		})
	if err != nil {
		log.Fatalf("failed to register check_domains tool: %v", err)
	}

	if err := server.Serve(); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}

	done := make(chan struct{})
	<-done
}
