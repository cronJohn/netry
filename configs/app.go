package configs

// Preconfigured nmap scan modes
var scanModeMap = map[string][]string{
	"discover": {"-sn"},
	"default":  {""},
	"full":     {"-A"},
	// NOTE: NetworkInterfaceMode could be added for --iflist flag
}

const (
	VersionLevel         int  = 7
	GetOSInfo            bool = false
	GetDefaultScriptInfo bool = false
	GetTracerouteInfo    bool = false
)
