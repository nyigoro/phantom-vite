// pkg/plugin/manager.go
type PluginManager struct {
    plugins map[string]Plugin
}

type Plugin interface {
    Name() string
    Version() string
    Execute(ctx context.Context, engine Engine, args []string) error
    Dependencies() []string
}
