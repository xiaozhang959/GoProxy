package config

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"
	"sync"
)

const DefaultPassword = "goproxy"

func dataDir() string {
	if d := os.Getenv("DATA_DIR"); d != "" {
		os.MkdirAll(d, 0755)
		return d + "/"
	}
	return ""
}

func ConfigFile() string { return dataDir() + "config.json" }

// ProxySource 定义一个可配置的公开代理抓取源。
type ProxySource struct {
	URL      string `json:"url"`
	Protocol string `json:"protocol"` // http 或 socks5
}

var defaultFastFetchSources = []ProxySource{
	// ProxyScraper - 每30分钟更新
	{"https://raw.githubusercontent.com/ProxyScraper/ProxyScraper/main/http.txt", "http"},
	{"https://raw.githubusercontent.com/ProxyScraper/ProxyScraper/main/socks4.txt", "socks5"},
	{"https://raw.githubusercontent.com/ProxyScraper/ProxyScraper/main/socks5.txt", "socks5"},
	// monosans - 每小时更新
	{"https://raw.githubusercontent.com/monosans/proxy-list/main/proxies/http.txt", "http"},
	// prxchk - 频繁更新
	{"https://raw.githubusercontent.com/prxchk/proxy-list/main/http.txt", "http"},
	{"https://raw.githubusercontent.com/prxchk/proxy-list/main/socks5.txt", "socks5"},
	{"https://raw.githubusercontent.com/prxchk/proxy-list/main/socks4.txt", "socks5"},
	// sunny9577 - 自动抓取更新
	{"https://cdn.jsdelivr.net/gh/sunny9577/proxy-scraper/generated/http_proxies.txt", "http"},
	{"https://cdn.jsdelivr.net/gh/sunny9577/proxy-scraper/generated/socks5_proxies.txt", "socks5"},
	{"https://cdn.jsdelivr.net/gh/sunny9577/proxy-scraper/generated/socks4_proxies.txt", "socks5"},
}

var defaultSlowFetchSources = []ProxySource{
	// TheSpeedX - 每天更新，量大
	{"https://raw.githubusercontent.com/TheSpeedX/SOCKS-List/master/http.txt", "http"},
	{"https://raw.githubusercontent.com/TheSpeedX/SOCKS-List/master/socks4.txt", "socks5"},
	{"https://raw.githubusercontent.com/TheSpeedX/SOCKS-List/master/socks5.txt", "socks5"},
	// monosans SOCKS
	{"https://raw.githubusercontent.com/monosans/proxy-list/main/proxies/socks4.txt", "socks5"},
	{"https://raw.githubusercontent.com/monosans/proxy-list/main/proxies/socks5.txt", "socks5"},
	// databay-labs - 备用源
	{"https://cdn.jsdelivr.net/gh/databay-labs/free-proxy-list/http.txt", "http"},
	{"https://cdn.jsdelivr.net/gh/databay-labs/free-proxy-list/socks5.txt", "socks5"},
	// Anonym0usWork1221 - 量大质量尚可
	{"https://cdn.jsdelivr.net/gh/Anonym0usWork1221/Free-Proxies/proxy_files/http_proxies.txt", "http"},
	{"https://cdn.jsdelivr.net/gh/Anonym0usWork1221/Free-Proxies/proxy_files/socks5_proxies.txt", "socks5"},
	{"https://cdn.jsdelivr.net/gh/Anonym0usWork1221/Free-Proxies/proxy_files/socks4_proxies.txt", "socks5"},
	// ALIILAPRO
	{"https://cdn.jsdelivr.net/gh/ALIILAPRO/Proxy/http.txt", "http"},
	{"https://cdn.jsdelivr.net/gh/ALIILAPRO/Proxy/socks4.txt", "socks5"},
	// vakhov/fresh-proxy-list
	{"https://cdn.jsdelivr.net/gh/vakhov/fresh-proxy-list/http.txt", "http"},
	{"https://cdn.jsdelivr.net/gh/vakhov/fresh-proxy-list/socks5.txt", "socks5"},
	{"https://cdn.jsdelivr.net/gh/vakhov/fresh-proxy-list/socks4.txt", "socks5"},
	// Zaeem20
	{"https://cdn.jsdelivr.net/gh/Zaeem20/FREE_PROXIES_LIST/http.txt", "http"},
	{"https://cdn.jsdelivr.net/gh/Zaeem20/FREE_PROXIES_LIST/socks4.txt", "socks5"},
	// hookzof - socks5 专项
	{"https://cdn.jsdelivr.net/gh/hookzof/socks5_list/proxy.txt", "socks5"},
	// proxy4parsing
	{"https://cdn.jsdelivr.net/gh/proxy4parsing/proxy-list/http.txt", "http"},
	{"https://cdn.jsdelivr.net/gh/proxy4parsing/proxy-list/socks5.txt", "socks5"},
}

func cloneProxySources(sources []ProxySource) []ProxySource {
	clone := make([]ProxySource, len(sources))
	copy(clone, sources)
	return clone
}

func DefaultFastFetchSources() []ProxySource {
	return cloneProxySources(defaultFastFetchSources)
}

func DefaultSlowFetchSources() []ProxySource {
	return cloneProxySources(defaultSlowFetchSources)
}

// NormalizeProxySources 清理配置源，兼容 socks4 源按 socks5 代理处理。
func NormalizeProxySources(sources []ProxySource) []ProxySource {
	normalized, _ := normalizeProxySources(sources, false)
	return normalized
}

// ValidateProxySources 校验并规范化页面提交的源配置。
func ValidateProxySources(sources []ProxySource) ([]ProxySource, error) {
	return normalizeProxySources(sources, true)
}

func normalizeProxySources(sources []ProxySource, strict bool) ([]ProxySource, error) {
	result := make([]ProxySource, 0, len(sources))
	seen := make(map[string]struct{})

	for _, source := range sources {
		sourceURL := strings.TrimSpace(source.URL)
		protocol := strings.ToLower(strings.TrimSpace(source.Protocol))
		if sourceURL == "" && protocol == "" {
			continue
		}
		if protocol == "socks4" {
			protocol = "socks5"
		}
		if protocol != "http" && protocol != "socks5" {
			if strict {
				return nil, fmt.Errorf("invalid source protocol: %s", source.Protocol)
			}
			continue
		}
		if sourceURL == "" {
			if strict {
				return nil, fmt.Errorf("empty source url")
			}
			continue
		}
		if strict {
			parsed, err := url.Parse(sourceURL)
			if err != nil || parsed.Scheme == "" || parsed.Host == "" {
				return nil, fmt.Errorf("invalid source url: %s", sourceURL)
			}
			if parsed.Scheme != "http" && parsed.Scheme != "https" {
				return nil, fmt.Errorf("unsupported source url scheme: %s", parsed.Scheme)
			}
		}

		key := protocol + " " + sourceURL
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, ProxySource{URL: sourceURL, Protocol: protocol})
	}

	return result, nil
}

type Config struct {
	// WebUI 端口
	WebUIPort string

	// WebUI 密码 SHA256 哈希
	WebUIPasswordHash string

	// 代理列表文本接口访问令牌
	ProxyListToken string

	// 代理池本地监听端口（随机轮换模式）
	ProxyPort string

	// 稳定代理端口（最低延迟模式）
	StableProxyPort string

	// SOCKS5 服务端口（随机轮换模式）
	SOCKS5Port string

	// 稳定 SOCKS5 端口（最低延迟模式）
	StableSOCKS5Port string

	// 代理服务认证配置
	ProxyAuthEnabled      bool   // 是否启用代理认证（默认 false）
	ProxyAuthUsername     string // 代理认证用户名（默认 "proxy"）
	ProxyAuthPassword     string // 代理认证密码明文（用于 SOCKS5）
	ProxyAuthPasswordHash string // 代理认证密码 SHA256 哈希（用于 HTTP）

	// 地理过滤配置
	BlockedCountries []string // 屏蔽的国家代码列表（如 ["CN", "RU"]，默认 ["CN"]）
	AllowedCountries []string // 允许的国家代码列表（白名单，非空时优先于黑名单）

	// SQLite 数据库路径
	DBPath string

	// ========== 池子容量配置 ==========
	PoolMaxSize        int     // 代理池总容量（默认100）
	PoolHTTPRatio      float64 // HTTP协议占比（默认0.5）
	PoolMinPerProtocol int     // 每协议最小保证（默认10）

	// ========== 延迟标准配置 ==========
	MaxLatencyMs          int // 标准模式最大延迟（默认2000ms）
	MaxLatencyEmergency   int // 紧急模式放宽延迟（默认3000ms）
	MaxLatencyHealthy     int // 健康模式严格延迟（默认1500ms）
	MaxLatencyDegradation int // 降级模式超宽松延迟（默认5000ms）

	// ========== 验证配置 ==========
	ValidateConcurrency int    // 验证并发数（默认300）
	ValidateTimeout     int    // 验证超时（秒）（默认8）
	ValidateURL         string // 验证目标 URL

	// ========== 健康检查配置 ==========
	HealthCheckInterval    int // 状态监控间隔（分钟）（默认5）
	HealthCheckBatchSize   int // 每批验证数量（默认20）
	HealthCheckConcurrency int // 批次内并发数（默认50）

	// ========== 优化配置 ==========
	OptimizeInterval    int     // 优化轮换间隔（分钟）（默认30）
	OptimizeConcurrency int     // 优化时并发数（默认100）
	ReplaceThreshold    float64 // 替换阈值（默认0.7，新代理需快30%）

	// ========== IP查询配置 ==========
	IPQueryRateLimit int // IP查询限流（次/秒）（默认10）

	// ========== 源管理配置 ==========
	SourceFailThreshold    int           // 源降级阈值（默认3）
	SourceDisableThreshold int           // 源禁用阈值（默认5）
	SourceCooldownMinutes  int           // 源禁用冷却时间（默认30）
	FetchFastSources       []ProxySource // 快速更新源（补充/紧急抓取）
	FetchSlowSources       []ProxySource // 慢速更新源（优化轮换抓取）

	// ========== 自定义订阅代理配置 ==========
	CustomProxyMode       string // 代理使用模式：mixed / custom_only / free_only（默认 mixed）
	CustomPriority        bool   // 混用模式下订阅代理优先（默认 true）
	CustomFreePriority    bool   // 混用模式下免费代理优先（默认 false）
	CustomProbeInterval   int    // 禁用代理探测唤醒间隔（分钟，默认 10）
	CustomRefreshInterval int    // 默认订阅刷新间隔（分钟，默认 60）
	SingBoxPath           string // sing-box 二进制路径（默认 "sing-box"）
	SingBoxBasePort       int    // sing-box 本地端口起始（默认 20000）

	// ========== 兼容旧配置 ==========
	MaxResponseMs int // 已废弃，使用 MaxLatencyMs 替代
	MaxFailCount  int // 代理失败次数阈值
	MaxRetry      int // 请求失败后的重试次数
	FetchInterval int // 已废弃，由智能抓取器管理
	CheckInterval int // 已废弃，由 HealthCheckInterval 替代

}

var (
	globalCfg *Config
	cfgMu     sync.RWMutex
)

func passwordHash(plain string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(plain)))
}

func newUUID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Sprintf("generate proxy list token: %v", err))
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

func DefaultConfig() *Config {
	// 优先从环境变量 WEBUI_PASSWORD 读取密码，未设置时使用默认密码
	password := os.Getenv("WEBUI_PASSWORD")
	if password == "" {
		password = DefaultPassword
	}
	proxyListToken := strings.TrimSpace(os.Getenv("PROXY_LIST_TOKEN"))
	if proxyListToken == "" {
		proxyListToken = newUUID()
	}

	// 读取代理认证配置
	proxyAuthEnabled := os.Getenv("PROXY_AUTH_ENABLED") == "true"
	proxyAuthUsername := os.Getenv("PROXY_AUTH_USERNAME")
	if proxyAuthUsername == "" {
		proxyAuthUsername = "proxy"
	}
	proxyAuthPassword := os.Getenv("PROXY_AUTH_PASSWORD")
	proxyAuthHash := ""
	if proxyAuthPassword != "" {
		proxyAuthHash = passwordHash(proxyAuthPassword)
	}

	// 读取地理过滤配置
	blockedCountries := []string{"CN"} // 默认屏蔽中国大陆
	if blockedEnv := os.Getenv("BLOCKED_COUNTRIES"); blockedEnv != "" {
		// 支持逗号分隔的国家代码，如 "CN,RU,KP"
		countries := strings.Split(blockedEnv, ",")
		blockedCountries = make([]string, 0, len(countries))
		for _, c := range countries {
			c = strings.TrimSpace(strings.ToUpper(c))
			if c != "" {
				blockedCountries = append(blockedCountries, c)
			}
		}
	}

	// 读取白名单配置（白名单非空时优先于黑名单）
	var allowedCountries []string
	if allowedEnv := os.Getenv("ALLOWED_COUNTRIES"); allowedEnv != "" {
		countries := strings.Split(allowedEnv, ",")
		allowedCountries = make([]string, 0, len(countries))
		for _, c := range countries {
			c = strings.TrimSpace(strings.ToUpper(c))
			if c != "" {
				allowedCountries = append(allowedCountries, c)
			}
		}
	}

	// 读取订阅代理配置
	customProxyMode := os.Getenv("CUSTOM_PROXY_MODE")
	if customProxyMode == "" {
		customProxyMode = "mixed"
	}
	singBoxPath := os.Getenv("SINGBOX_PATH")
	if singBoxPath == "" {
		singBoxPath = "sing-box"
	}

	return &Config{
		// 基础服务配置
		WebUIPort:         ":7778",
		WebUIPasswordHash: passwordHash(password),
		ProxyListToken:    proxyListToken,
		ProxyPort:         ":7777",
		StableProxyPort:   ":7776",
		SOCKS5Port:        ":7779",
		StableSOCKS5Port:  ":7780",
		DBPath:            dataDir() + "proxy.db",

		// 代理认证配置
		ProxyAuthEnabled:      proxyAuthEnabled,
		ProxyAuthUsername:     proxyAuthUsername,
		ProxyAuthPassword:     proxyAuthPassword,
		ProxyAuthPasswordHash: proxyAuthHash,

		// 地理过滤配置
		BlockedCountries: blockedCountries,
		AllowedCountries: allowedCountries,

		// 池子容量配置
		PoolMaxSize:        100, // 总容量
		PoolHTTPRatio:      0.3, // HTTP占30%
		PoolMinPerProtocol: 10,  // 每协议最少10个

		// 延迟标准配置
		MaxLatencyMs:          2500, // 标准2.5秒
		MaxLatencyEmergency:   4000, // 紧急4秒
		MaxLatencyHealthy:     2000, // 健康2秒
		MaxLatencyDegradation: 5000, // 降级5秒

		// 验证配置
		ValidateConcurrency: 300,
		ValidateTimeout:     10, // 从8秒增加到10秒
		ValidateURL:         "http://www.gstatic.com/generate_204",

		// 健康检查配置
		HealthCheckInterval:    5,  // 5分钟
		HealthCheckBatchSize:   20, // 每批20个
		HealthCheckConcurrency: 50, // 批次并发50

		// 优化配置
		OptimizeInterval:    30,  // 30分钟
		OptimizeConcurrency: 100, // 并发100
		ReplaceThreshold:    0.7, // 新代理需快30%

		// IP查询配置
		IPQueryRateLimit: 10, // 10次/秒

		// 源管理配置
		SourceFailThreshold:    3,  // 失败3次降级
		SourceDisableThreshold: 5,  // 失败5次禁用
		SourceCooldownMinutes:  30, // 禁用30分钟
		FetchFastSources:       DefaultFastFetchSources(),
		FetchSlowSources:       DefaultSlowFetchSources(),

		// 自定义订阅代理配置
		CustomProxyMode:       customProxyMode,
		CustomPriority:        true,
		CustomProbeInterval:   10,
		CustomRefreshInterval: 60,
		SingBoxPath:           singBoxPath,
		SingBoxBasePort:       20000,

		// 兼容旧配置
		MaxResponseMs: 5000,
		MaxFailCount:  3,
		MaxRetry:      3,
		FetchInterval: 30,
		CheckInterval: 10,
	}
}

// Load 从文件加载配置，文件不存在则用默认值
func Load() *Config {
	cfg := DefaultConfig()
	data, err := os.ReadFile(ConfigFile())
	if err == nil {
		var saved savedConfig
		if json.Unmarshal(data, &saved) == nil {
			// 池子配置
			if saved.PoolMaxSize > 0 {
				cfg.PoolMaxSize = saved.PoolMaxSize
			}
			if saved.PoolHTTPRatio > 0 && saved.PoolHTTPRatio <= 1 {
				cfg.PoolHTTPRatio = saved.PoolHTTPRatio
			}
			if saved.PoolMinPerProtocol > 0 {
				cfg.PoolMinPerProtocol = saved.PoolMinPerProtocol
			}
			if strings.TrimSpace(saved.ProxyListToken) != "" {
				cfg.ProxyListToken = strings.TrimSpace(saved.ProxyListToken)
			}

			// 延迟配置
			if saved.MaxLatencyMs > 0 {
				cfg.MaxLatencyMs = saved.MaxLatencyMs
			}
			if saved.MaxLatencyEmergency > 0 {
				cfg.MaxLatencyEmergency = saved.MaxLatencyEmergency
			}
			if saved.MaxLatencyHealthy > 0 {
				cfg.MaxLatencyHealthy = saved.MaxLatencyHealthy
			}

			// 验证配置
			if saved.ValidateConcurrency > 0 {
				cfg.ValidateConcurrency = saved.ValidateConcurrency
			}
			if saved.ValidateTimeout > 0 {
				cfg.ValidateTimeout = saved.ValidateTimeout
			}

			// 健康检查配置
			if saved.HealthCheckInterval > 0 {
				cfg.HealthCheckInterval = saved.HealthCheckInterval
			}
			if saved.HealthCheckBatchSize > 0 {
				cfg.HealthCheckBatchSize = saved.HealthCheckBatchSize
			}

			// 优化配置
			if saved.OptimizeInterval > 0 {
				cfg.OptimizeInterval = saved.OptimizeInterval
			}
			if saved.ReplaceThreshold > 0 && saved.ReplaceThreshold <= 1 {
				cfg.ReplaceThreshold = saved.ReplaceThreshold
			}

			// 兼容旧配置
			if saved.FetchInterval > 0 {
				cfg.FetchInterval = saved.FetchInterval
			}
			if saved.CheckInterval > 0 {
				cfg.CheckInterval = saved.CheckInterval
			}

			// 抓取源配置
			if saved.FetchFastSources != nil {
				if sources := NormalizeProxySources(saved.FetchFastSources); len(sources) > 0 {
					cfg.FetchFastSources = sources
				}
			}
			if saved.FetchSlowSources != nil {
				if sources := NormalizeProxySources(saved.FetchSlowSources); len(sources) > 0 {
					cfg.FetchSlowSources = sources
				}
			}

			// 地理过滤配置（config.json 优先于环境变量）
			if saved.BlockedCountries != nil {
				cfg.BlockedCountries = saved.BlockedCountries
			}
			if saved.AllowedCountries != nil {
				cfg.AllowedCountries = saved.AllowedCountries
			}

			// 自定义订阅代理配置
			if saved.CustomProxyMode != "" {
				cfg.CustomProxyMode = saved.CustomProxyMode
			}
			if saved.CustomPriority != nil {
				cfg.CustomPriority = *saved.CustomPriority
			}
			if saved.CustomFreePriority != nil {
				cfg.CustomFreePriority = *saved.CustomFreePriority
			}
			if saved.CustomProbeInterval > 0 {
				cfg.CustomProbeInterval = saved.CustomProbeInterval
			}
			if saved.CustomRefreshInterval > 0 {
				cfg.CustomRefreshInterval = saved.CustomRefreshInterval
			}
			if saved.SingBoxPath != "" {
				cfg.SingBoxPath = saved.SingBoxPath
			}
			if saved.SingBoxBasePort > 0 {
				cfg.SingBoxBasePort = saved.SingBoxBasePort
			}
		}
	}
	cfgMu.Lock()
	globalCfg = cfg
	cfgMu.Unlock()
	return cfg
}

// Get 获取当前配置
func Get() *Config {
	cfgMu.RLock()
	defer cfgMu.RUnlock()
	return globalCfg
}

// savedConfig 持久化可调整的字段
type savedConfig struct {
	// 文本代理列表接口
	ProxyListToken string `json:"proxy_list_token,omitempty"`

	// 池子配置
	PoolMaxSize        int     `json:"pool_max_size"`
	PoolHTTPRatio      float64 `json:"pool_http_ratio"`
	PoolMinPerProtocol int     `json:"pool_min_per_protocol"`

	// 延迟配置
	MaxLatencyMs        int `json:"max_latency_ms"`
	MaxLatencyEmergency int `json:"max_latency_emergency"`
	MaxLatencyHealthy   int `json:"max_latency_healthy"`

	// 验证配置
	ValidateConcurrency int `json:"validate_concurrency"`
	ValidateTimeout     int `json:"validate_timeout"`

	// 健康检查配置
	HealthCheckInterval  int `json:"health_check_interval"`
	HealthCheckBatchSize int `json:"health_check_batch_size"`

	// 优化配置
	OptimizeInterval int     `json:"optimize_interval"`
	ReplaceThreshold float64 `json:"replace_threshold"`

	// 地理过滤配置
	BlockedCountries []string `json:"blocked_countries,omitempty"`
	AllowedCountries []string `json:"allowed_countries,omitempty"`

	// 自定义订阅代理配置
	CustomProxyMode       string `json:"custom_proxy_mode,omitempty"`
	CustomPriority        *bool  `json:"custom_priority,omitempty"`
	CustomFreePriority    *bool  `json:"custom_free_priority,omitempty"`
	CustomProbeInterval   int    `json:"custom_probe_interval,omitempty"`
	CustomRefreshInterval int    `json:"custom_refresh_interval,omitempty"`
	SingBoxPath           string `json:"singbox_path,omitempty"`
	SingBoxBasePort       int    `json:"singbox_base_port,omitempty"`

	// 抓取源配置
	FetchFastSources []ProxySource `json:"fetch_fast_sources,omitempty"`
	FetchSlowSources []ProxySource `json:"fetch_slow_sources,omitempty"`

	// 兼容旧配置
	FetchInterval int `json:"fetch_interval,omitempty"`
	CheckInterval int `json:"check_interval,omitempty"`
}

// Save 保存配置到文件，并更新内存配置
func Save(cfg *Config) error {
	cfgMu.Lock()
	*globalCfg = *cfg
	cfgMu.Unlock()

	customPriority := cfg.CustomPriority
	customFreePriority := cfg.CustomFreePriority
	data, err := json.MarshalIndent(savedConfig{
		ProxyListToken:        strings.TrimSpace(cfg.ProxyListToken),
		PoolMaxSize:           cfg.PoolMaxSize,
		PoolHTTPRatio:         cfg.PoolHTTPRatio,
		PoolMinPerProtocol:    cfg.PoolMinPerProtocol,
		MaxLatencyMs:          cfg.MaxLatencyMs,
		MaxLatencyEmergency:   cfg.MaxLatencyEmergency,
		MaxLatencyHealthy:     cfg.MaxLatencyHealthy,
		ValidateConcurrency:   cfg.ValidateConcurrency,
		ValidateTimeout:       cfg.ValidateTimeout,
		HealthCheckInterval:   cfg.HealthCheckInterval,
		HealthCheckBatchSize:  cfg.HealthCheckBatchSize,
		OptimizeInterval:      cfg.OptimizeInterval,
		ReplaceThreshold:      cfg.ReplaceThreshold,
		BlockedCountries:      cfg.BlockedCountries,
		AllowedCountries:      cfg.AllowedCountries,
		CustomProxyMode:       cfg.CustomProxyMode,
		CustomPriority:        &customPriority,
		CustomFreePriority:    &customFreePriority,
		CustomProbeInterval:   cfg.CustomProbeInterval,
		CustomRefreshInterval: cfg.CustomRefreshInterval,
		SingBoxPath:           cfg.SingBoxPath,
		SingBoxBasePort:       cfg.SingBoxBasePort,
		FetchFastSources:      NormalizeProxySources(cfg.FetchFastSources),
		FetchSlowSources:      NormalizeProxySources(cfg.FetchSlowSources),
		FetchInterval:         cfg.FetchInterval,
		CheckInterval:         cfg.CheckInterval,
	}, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(ConfigFile(), data, 0644)
}

// CalculateSlots 根据配置计算各协议的槽位数
func (c *Config) CalculateSlots() (httpSlots, socks5Slots int) {
	httpSlots = int(float64(c.PoolMaxSize) * c.PoolHTTPRatio)
	socks5Slots = c.PoolMaxSize - httpSlots

	// 保证最小值
	if httpSlots < c.PoolMinPerProtocol {
		httpSlots = c.PoolMinPerProtocol
	}
	if socks5Slots < c.PoolMinPerProtocol {
		socks5Slots = c.PoolMinPerProtocol
	}

	return
}

// GetLatencyThreshold 根据池子状态返回合适的延迟阈值
func (c *Config) GetLatencyThreshold(poolStatus string) int {
	switch poolStatus {
	case "emergency":
		return c.MaxLatencyEmergency
	case "critical":
		return c.MaxLatencyEmergency
	case "warning":
		// warning状态使用紧急标准，加快填充速度
		return c.MaxLatencyEmergency
	case "healthy":
		return c.MaxLatencyHealthy
	default:
		return c.MaxLatencyMs
	}
}
