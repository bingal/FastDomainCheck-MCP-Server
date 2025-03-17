package config

// Config contains all configuration settings
type Config struct {
	WhoisServers    map[string]string   // WHOIS server mappings for different TLDs
	NotFoundPattern map[string][]string // Patterns indicating domain availability
}

// NewConfig creates a new configuration instance
func NewConfig() *Config {
	return &Config{
		WhoisServers: map[string]string{
			// Generic TLDs
			"com":  "whois.verisign-grs.com",
			"net":  "whois.verisign-grs.com",
			"org":  "whois.pir.org",
			"info": "whois.nic.info",
			"biz":  "whois.nic.biz",
			"mobi": "whois.nic.mobi",
			"edu":  "whois.educause.edu",
			"gov":  "whois.dotgov.gov",
			"mil":  "whois.nic.mil",
			"int":  "whois.iana.org",

			// Country and Region TLDs
			"cn": "whois.cnnic.cn",
			"hk": "whois.hkirc.hk",
			"tw": "whois.twnic.net.tw",

			// Popular New TLDs
			"io":     "whois.nic.io",
			"ai":     "whois.nic.ai",
			"me":     "whois.nic.me",
			"cc":     "whois.nic.cc",
			"tv":     "whois.nic.tv",
			"co":     "whois.nic.co",
			"xyz":    "whois.nic.xyz",
			"top":    "whois.nic.top",
			"vip":    "whois.nic.vip",
			"club":   "whois.nic.club",
			"shop":   "whois.nic.shop",
			"site":   "whois.nic.site",
			"wang":   "whois.nic.wang",
			"xin":    "whois.nic.xin",
			"app":    "whois.nic.google",
			"dev":    "whois.nic.google",
			"cloud":  "whois.nic.cloud",
			"online": "whois.nic.online",
			"store":  "whois.nic.store",

			// Chinese Regional TLDs
			"com.cn": "whois.cnnic.cn",
			"net.cn": "whois.cnnic.cn",
			"org.cn": "whois.cnnic.cn",
			"gov.cn": "whois.cnnic.cn",

			// Chinese IDN TLDs
			"中国":  "whois.cnnic.cn",
			"公司":  "whois.cnnic.cn",
			"网络":  "whois.cnnic.cn",
			"商城":  "whois.cnnic.cn",
			"网店":  "whois.cnnic.cn",
			"中文网": "whois.cnnic.cn",
		},
		NotFoundPattern: map[string][]string{
			"default": {"Domain not found", "No match for", "NOT FOUND", "No Data Found", "No entries found"},
			"cn":      {"no matching record", "No matching record"},
			"com":     {"No match for", "NOT FOUND", "No Data Found"},
			"net":     {"No match for", "NOT FOUND", "No Data Found"},
		},
	}
}
