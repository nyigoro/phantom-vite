package engine

import (
	"context"
	"time"
)

// Config represents the configuration for an automation engine
type Config struct {
	// Engine-specific settings
	Engine   string `json:"engine"`   // puppeteer, playwright, selenium
	Headless bool   `json:"headless"` // run in headless mode
	
	// Browser settings
	Viewport ViewportConfig `json:"viewport"` // browser viewport size
	Timeout  time.Duration  `json:"timeout"`  // default timeout for operations
	
	// Advanced settings
	UserAgent    string            `json:"user_agent,omitempty"`    // custom user agent
	ExtraHeaders map[string]string `json:"extra_headers,omitempty"` // additional HTTP headers
	Proxy        string            `json:"proxy,omitempty"`         // proxy server URL
	
	// Browser launch options
	ExecutablePath string   `json:"executable_path,omitempty"` // custom browser executable
	Args           []string `json:"args,omitempty"`            // additional browser arguments
	
	// Plugin configuration
	Plugins []PluginConfig `json:"plugins,omitempty"`
}

// ViewportConfig represents browser viewport dimensions
type ViewportConfig struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// PluginConfig represents a plugin configuration
type PluginConfig struct {
	Path    string                 `json:"path"`
	Name    string                 `json:"name,omitempty"`
	Enabled bool                   `json:"enabled"`
	Options map[string]interface{} `json:"options,omitempty"`
}

// ScreenshotOptions represents options for taking screenshots
type ScreenshotOptions struct {
	// Output settings
	Path    string `json:"path"`              // file path to save screenshot
	Format  string `json:"format,omitempty"`  // png, jpeg, webp
	Quality int    `json:"quality,omitempty"` // JPEG quality (0-100)
	
	// Capture settings
	FullPage bool `json:"full_page,omitempty"` // capture full scrollable page
	
	// Clipping rectangle (optional)
	Clip *ClipOptions `json:"clip,omitempty"`
	
	// Advanced options
	OmitBackground bool `json:"omit_background,omitempty"` // transparent background
}

// ClipOptions represents a rectangular area to clip from the screenshot
type ClipOptions struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// WaitOptions represents options for waiting operations
type WaitOptions struct {
	Timeout    time.Duration `json:"timeout,omitempty"`     // maximum wait time
	Visible    bool          `json:"visible,omitempty"`     // wait for element to be visible
	Hidden     bool          `json:"hidden,omitempty"`      // wait for element to be hidden
	Polling    time.Duration `json:"polling,omitempty"`     // polling interval
	RetryCount int           `json:"retry_count,omitempty"` // number of retries
}

// NavigationOptions represents options for page navigation
type NavigationOptions struct {
	Timeout         time.Duration `json:"timeout,omitempty"`          // navigation timeout
	WaitUntil       string        `json:"wait_until,omitempty"`       // load, domcontentloaded, networkidle0, networkidle2
	Referer         string        `json:"referer,omitempty"`          // HTTP referer header
	WaitForSelector string        `json:"wait_for_selector,omitempty"` // wait for specific selector after navigation
}

// ElementHandle represents a handle to a DOM element
type ElementHandle interface {
	Click() error
	Type(text string) error
	GetAttribute(name string) (string, error)
	GetProperty(name string) (interface{}, error)
	IsVisible() (bool, error)
	BoundingBox() (*BoundingBox, error)
}

// BoundingBox represents the bounding box of an element
type BoundingBox struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// Engine interface defines the contract for browser automation engines
type Engine interface {
	// Basic engine operations
	Name() string
	Initialize(config Config) error
	Close() error
	
	// Page operations
	Navigate(ctx context.Context, url string) (*Page, error)
	NewPage(ctx context.Context) (*Page, error)
	GetPages(ctx context.Context) ([]*Page, error)
	
	// Screenshot operations
	Screenshot(ctx context.Context, options ScreenshotOptions) error
	
	// Advanced operations
	SetUserAgent(userAgent string) error
	SetExtraHeaders(headers map[string]string) error
	SetViewport(viewport ViewportConfig) error
}

// Page interface defines operations that can be performed on a web page
type Page interface {
	// Basic page information
	Title() (string, error)
	URL() (string, error)
	Content() (string, error)
	
	// Navigation
	Navigate(url string, options *NavigationOptions) error
	Reload(options *NavigationOptions) error
	GoBack() error
	GoForward() error
	
	// Element operations
	QuerySelector(selector string) (ElementHandle, error)
	QuerySelectorAll(selector string) ([]ElementHandle, error)
	WaitForSelector(selector string, options *WaitOptions) (ElementHandle, error)
	
	// Script execution
	ExecuteScript(script string) (interface{}, error)
	ExecuteScriptAsync(script string) (interface{}, error)
	
	// Input operations
	Click(selector string) error
	Type(selector string, text string) error
	Fill(selector string, text string) error
	Select(selector string, values ...string) error
	
	// Screenshot operations
	Screenshot(options ScreenshotOptions) error
	
	// Waiting operations
	WaitForNavigation(options *NavigationOptions) error
	WaitForTimeout(timeout time.Duration) error
	WaitForFunction(pageFunction string, options *WaitOptions) error
	
	// Cookie operations
	GetCookies() ([]Cookie, error)
	SetCookies(cookies []Cookie) error
	ClearCookies() error
	
	// Advanced operations
	SetViewport(viewport ViewportConfig) error
	GetMetrics() (map[string]interface{}, error)
	EmulateDevice(device Device) error
	
	// Lifecycle
	Close() error
}

// Cookie represents an HTTP cookie
type Cookie struct {
	Name     string  `json:"name"`
	Value    string  `json:"value"`
	Domain   string  `json:"domain,omitempty"`
	Path     string  `json:"path,omitempty"`
	Expires  int64   `json:"expires,omitempty"` // Unix timestamp
	HTTPOnly bool    `json:"httpOnly,omitempty"`
	Secure   bool    `json:"secure,omitempty"`
	SameSite string  `json:"sameSite,omitempty"` // Strict, Lax, None
	Priority string  `json:"priority,omitempty"` // Low, Medium, High
}

// Device represents a device to emulate
type Device struct {
	Name            string         `json:"name"`
	UserAgent       string         `json:"userAgent"`
	Viewport        ViewportConfig `json:"viewport"`
	DeviceScaleFactor float64     `json:"deviceScaleFactor"`
	IsMobile        bool           `json:"isMobile"`
	HasTouch        bool           `json:"hasTouch"`
}

// EngineError represents an error that occurred during engine operations
type EngineError struct {
	Engine    string `json:"engine"`
	Operation string `json:"operation"`
	Message   string `json:"message"`
	Cause     error  `json:"cause,omitempty"`
}

func (e *EngineError) Error() string {
	if e.Cause != nil {
		return e.Engine + " " + e.Operation + ": " + e.Message + " (caused by: " + e.Cause.Error() + ")"
	}
	return e.Engine + " " + e.Operation + ": " + e.Message
}

func (e *EngineError) Unwrap() error {
	return e.Cause
}

// NewEngineError creates a new engine error
func NewEngineError(engine, operation, message string, cause error) *EngineError {
	return &EngineError{
		Engine:    engine,
		Operation: operation,
		Message:   message,
		Cause:     cause,
	}
}

// Common predefined devices for emulation
var (
	DeviceIPhone12 = Device{
		Name:              "iPhone 12",
		UserAgent:         "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/604.1",
		Viewport:          ViewportConfig{Width: 390, Height: 844},
		DeviceScaleFactor: 3,
		IsMobile:          true,
		HasTouch:          true,
	}
	
	DeviceIPhoneSE = Device{
		Name:              "iPhone SE",
		UserAgent:         "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/604.1",
		Viewport:          ViewportConfig{Width: 375, Height: 667},
		DeviceScaleFactor: 2,
		IsMobile:          true,
		HasTouch:          true,
	}
	
	DevicePixel5 = Device{
		Name:              "Pixel 5",
		UserAgent:         "Mozilla/5.0 (Linux; Android 11; Pixel 5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.91 Mobile Safari/537.36",
		Viewport:          ViewportConfig{Width: 393, Height: 851},
		DeviceScaleFactor: 3,
		IsMobile:          true,
		HasTouch:          true,
	}
	
	DeviceIPadPro = Device{
		Name:              "iPad Pro",
		UserAgent:         "Mozilla/5.0 (iPad; CPU OS 14_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/604.1",
		Viewport:          ViewportConfig{Width: 1024, Height: 1366},
		DeviceScaleFactor: 2,
		IsMobile:          false,
		HasTouch:          true,
	}
	
	DeviceDesktop = Device{
		Name:              "Desktop",
		UserAgent:         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		Viewport:          ViewportConfig{Width: 1920, Height: 1080},
		DeviceScaleFactor: 1,
		IsMobile:          false,
		HasTouch:          false,
	}
)

// DefaultConfig returns a default configuration
func DefaultConfig() Config {
	return Config{
		Engine:   "puppeteer",
		Headless: true,
		Viewport: ViewportConfig{Width: 1920, Height: 1080},
		Timeout:  30 * time.Second,
		Plugins:  []PluginConfig{},
	}
}