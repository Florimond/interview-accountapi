package contracts

// Document represents a generic document interface
type Document interface {
}

// Provider represents a contract for a provider of documents
type Provider interface {
	Path() string // Path returns a URI path
}

type staticProvider struct {
	path string
}

// Path returns the path for the resource
func (p *staticProvider) Path() string {
	return p.path
}

// NewStaticProvider creates a new static provider
func NewStaticProvider(path string) Provider {
	return &staticProvider{
		path: path,
	}
}
