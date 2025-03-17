package checker

import (
	"fmt"
	"net"
	"strings"
	"time"
	"unicode"

	"github.com/bingal/FastDomainCheck-MCP-Server/config"
)

// DomainStatus represents the status of a domain name
type DomainStatus struct {
	Status     string            `json:"status"`     // registered, unregistered, failed
	Reason     string            `json:"reason"`     // Brief explanation of the status
	Details    string            `json:"details"`    // Detailed information
	TimingInfo map[string]string `json:"timingInfo"` // Query timing information
}

// DomainChecker handles domain registration status checks
type DomainChecker struct {
	config *config.Config
}

// NewDomainChecker creates a new domain checker instance
func NewDomainChecker(cfg *config.Config) *DomainChecker {
	return &DomainChecker{
		config: cfg,
	}
}

// getTLD extracts the top-level domain from a domain name
func (dc *DomainChecker) getTLD(domain string) string {
	parts := strings.Split(strings.ToLower(domain), ".")

	// Handle IDN (Internationalized Domain Names)
	for _, part := range parts {
		if containsChinese(part) {
			if _, ok := dc.config.WhoisServers[part]; ok {
				return part
			}
			return ""
		}
	}

	// Handle compound TLDs (e.g., .com.cn)
	if len(parts) >= 3 {
		compoundTLD := fmt.Sprintf("%s.%s", parts[len(parts)-2], parts[len(parts)-1])
		if _, ok := dc.config.WhoisServers[compoundTLD]; ok {
			return compoundTLD
		}
	}

	// Check regular TLD
	tld := parts[len(parts)-1]
	if _, ok := dc.config.WhoisServers[tld]; ok {
		return tld
	}
	return ""
}

// containsChinese checks if a string contains Chinese characters
func containsChinese(s string) bool {
	for _, r := range s {
		if unicode.Is(unicode.Han, r) {
			return true
		}
	}
	return false
}

// queryWhois performs a WHOIS query for the domain
func (dc *DomainChecker) queryWhois(domain string) (string, error) {
	tld := dc.getTLD(domain)
	whoisServer := dc.config.WhoisServers[tld]
	if whoisServer == "" {
		whoisServer = "whois.iana.org"
	}

	conn, err := net.DialTimeout("tcp", whoisServer+":43", 10*time.Second)
	if err != nil {
		return "", fmt.Errorf("failed to connect to WHOIS server: %v", err)
	}
	defer conn.Close()

	conn.Write([]byte(domain + "\r\n"))
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))

	var response strings.Builder
	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			break
		}
		response.Write(buf[:n])
	}

	return response.String(), nil
}

// checkDNSRecord checks if a domain is registered using DNS lookup
func (dc *DomainChecker) checkDNSRecord(domain string) (bool, error) {
	// Use net.LookupHost instead of dns package
	_, err := net.LookupHost(domain)
	if err != nil {
		if dnsErr, ok := err.(*net.DNSError); ok && dnsErr.IsNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// CheckDomain checks a single domain and returns its status
func (dc *DomainChecker) CheckDomain(domain string) DomainStatus {
	timingInfo := make(map[string]string)

	// Check if TLD is supported
	tld := dc.getTLD(domain)
	if tld == "" {
		return DomainStatus{
			Status:     "failed",
			Reason:     "Unsupported TLD",
			Details:    fmt.Sprintf("Unsupported top-level domain: %s", strings.Split(domain, ".")[len(strings.Split(domain, "."))-1]),
			TimingInfo: timingInfo,
		}
	}

	// Step 1: WHOIS query
	whoisStartTime := time.Now()
	whoisResponse, err := dc.queryWhois(domain)
	whoisTime := time.Since(whoisStartTime)
	timingInfo["whois_query"] = fmt.Sprintf("%.2fs", whoisTime.Seconds())

	if err != nil {
		// If WHOIS fails, try DNS lookup
		dnsStartTime := time.Now()
		dnsResult, dnsErr := dc.checkDNSRecord(domain)
		dnsTime := time.Since(dnsStartTime)
		timingInfo["dns_query"] = fmt.Sprintf("%.2fs", dnsTime.Seconds())

		if dnsErr == nil {
			if dnsResult {
				return DomainStatus{
					Status:     "registered",
					Reason:     "DNS resolution successful",
					Details:    "Domain is registered (confirmed by DNS lookup)",
					TimingInfo: timingInfo,
				}
			}
			return DomainStatus{
				Status:     "unregistered",
				Reason:     "No DNS records found",
				Details:    "Domain is not registered (confirmed by DNS lookup)",
				TimingInfo: timingInfo,
			}
		}
		return DomainStatus{
			Status:     "failed",
			Reason:     "Both WHOIS and DNS queries failed",
			Details:    fmt.Sprintf("WHOIS error: %v, DNS error: %v", err, dnsErr),
			TimingInfo: timingInfo,
		}
	}

	// Check for unregistered domain patterns
	notFoundPatterns := dc.config.NotFoundPattern[tld]
	if len(notFoundPatterns) == 0 {
		notFoundPatterns = dc.config.NotFoundPattern["default"]
	}

	for _, pattern := range notFoundPatterns {
		if strings.Contains(whoisResponse, pattern) {
			return DomainStatus{
				Status:     "unregistered",
				Reason:     "Domain available according to WHOIS",
				Details:    "Domain is not registered (confirmed by WHOIS)",
				TimingInfo: timingInfo,
			}
		}
	}

	// Default to registered if no unregistered patterns found
	return DomainStatus{
		Status:     "registered",
		Reason:     "Domain registered according to WHOIS",
		Details:    "Domain is registered (confirmed by WHOIS)",
		TimingInfo: timingInfo,
	}
}

// CheckDomains performs bulk domain status checks
func (dc *DomainChecker) CheckDomains(domains []string) []DomainStatus {
	results := make([]DomainStatus, len(domains))
	for i, domain := range domains {
		results[i] = dc.CheckDomain(domain)
	}
	return results
}
