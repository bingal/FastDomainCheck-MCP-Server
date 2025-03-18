# FastDomainCheck MCP Server

A Model Context Protocol for checking domain name registration status in bulk.

## Features

- Bulk domain registration status checking
- Dual verification using WHOIS and DNS
- Support for IDN (Internationalized Domain Names)
- Concise output format
- Built-in input validation and error handling

## Tool Documentation

### check_domains

Check registration status for multiple domain names.

#### Input Format

```json
{
  "domains": ["example.com", "test.com"]
}
```

Parameters:
- `domains`: Array of strings containing domain names to check
  - Maximum length of 255 characters per domain
  - Maximum 50 domains per request
  - No empty domain names allowed

#### Output Format

```json
{
  "results": {
    "example.com": {
      "registered": true
    },
    "test.com": {
      "registered": false
    }
  }
}
```

Response Fields:
- `results`: Object with domain names as keys and check results as values
  - `registered`: Boolean
    - `true`: Domain is registered
    - `false`: Domain is available

#### Error Handling

The tool will return an error in the following cases:
1. Empty domains list
2. More than 50 domains in request
3. Empty domain name
4. Domain name exceeding 255 characters
5. Result serialization failure

Error Response Format:
```json
{
  "error": "Error: domains list cannot be empty"
}
```

#### Usage Examples

Check multiple domains:
> Request
```json
{
  "domains": ["example.com", "test123456.com"]
}
```

> Response
```json
{
  "results": {
    "example.com": {
      "registered": true
    },
    "test123456.com": {
      "registered": false
    }
  }
}
```


## Performance Considerations

1. Domain checks are executed sequentially, taking approximately 0.3-1 second per domain
2. Maximum 50 domains per request to prevent resource exhaustion
3. WHOIS query timeout set to 10 seconds
4. DNS query fallback when WHOIS query fails

## Error Handling Strategy

1. Input Validation: Comprehensive validation before processing
2. Dual Verification: WHOIS primary, DNS fallback
3. Timeout Management: Reasonable timeouts for all network operations
4. Detailed Error Messages: Clear error descriptions for troubleshooting

## Usage

### Download Binary

Download the binary file from the release page.
https://github.com/bingal/FastDomainCheck-MCP-Server/releases

### For Mac/Linux
```bash
chmod +x FastDomainCheck-MCP-Server
```

### MCP Server Settings

#### Configuring FastDomainCheck MCP in Claude Deskto
Modify your claude-desktop-config.json file as shown below

> Mac/Linux
```json
{
  "mcpServers": {
    "fastdomaincheck": {
      "command": "/path/to/FastDomainCheck-MCP-Server",
      "args": []
    }
  }
}
```

> Windows
```json
{
  "mcpServers": {
    "fastdomaincheck": {
      "command": "path/to/FastDomainCheck-MCP-Server.exe",
      "args": []
    }
  }
}
```

## Development Guide

### Requirements

- Go 1.16 or higher
- Network connectivity (for WHOIS and DNS queries)
- Dependencies:
  - github.com/metoro-io/mcp-golang
  - Other project-specific dependencies


### Build

```bash
# Linux
go build -o FastDomainCheck-MCP-Server

# Windows
go build -o FastDomainCheck-MCP-Server.exe
```
