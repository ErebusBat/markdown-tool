package types

// ParseContext holds data collected during the parsing phase
type ParseContext struct {
	OriginalInput string
	DetectedType  ContentType
	Confidence    int
	Metadata      map[string]interface{}
}

// ContentType represents the type of content detected
type ContentType int

const (
	ContentTypeUnknown ContentType = iota
	ContentTypeURL
	ContentTypeGitHubURL
	ContentTypeGitHubLong
	ContentTypeJIRAURL
	ContentTypeJIRAComment
	ContentTypeNotionURL
	ContentTypeJIRAKey
	ContentTypeJIRAKeyWithDescription
	ContentTypePhone7Digit
	ContentTypePhone10Digit
	ContentTypePhone11Digit
)

// Parser interface for content detection and parsing
type Parser interface {
	Parse(input string) (*ParseContext, error)
	CanHandle(input string) bool
}

// Writer interface for output generation
type Writer interface {
	Write(ctx *ParseContext) (string, error)
	Vote(ctx *ParseContext) int
	GetName() string
}

// Config represents the application configuration
type Config struct {
	GitHub GitHubConfig `yaml:"github" mapstructure:"github"`
	JIRA   JIRAConfig   `yaml:"jira" mapstructure:"jira"`
	URL    URLConfig    `yaml:"url" mapstructure:"url"`
}

// GitHubConfig holds GitHub-specific configuration
type GitHubConfig struct {
	DefaultOrg  string            `yaml:"default_org" mapstructure:"default_org"`
	DefaultRepo string            `yaml:"default_repo" mapstructure:"default_repo"`
	Mappings    map[string]string `yaml:"mappings" mapstructure:"mappings"`
}

// JIRAConfig holds JIRA-specific configuration
type JIRAConfig struct {
	Domain   string   `yaml:"domain" mapstructure:"domain"`
	Projects []string `yaml:"projects" mapstructure:"projects"`
}

// URLConfig holds URL-specific configuration
type URLConfig struct {
	DomainMappings map[string]string `yaml:"domain_mappings" mapstructure:"domain_mappings"`
}
