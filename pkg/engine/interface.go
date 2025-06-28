// pkg/engine/interface.go
package engine

import "context"

type Engine interface {
    Name() string
    Initialize(config Config) error
    Navigate(ctx context.Context, url string) (*Page, error)
    Screenshot(ctx context.Context, options ScreenshotOptions) error
    Close() error
}

type Page interface {
    Title() (string, error)
    Content() (string, error)
    ExecuteScript(script string) (interface{}, error)
    WaitForSelector(selector string) error
}
