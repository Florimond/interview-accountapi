package provider

import "fmt"

type restProvider struct {
	resourceURL string
}

// Path returns the path for the resource
func (p *restProvider) MakeURL(options ...Option) string {
	for _, o := range options{
		switch {
		
		}
	}


	return trim(p.resourceURL) + formatOptions(options)
}



// NewRESTProvider creates a new REST provider
func NewRESTProvider(format string, args ...interface{}) Provider {
	return &restProvider{
		resourceURL: fmt.Sprintf(format, args...),
	}
}


// formatOptions formats a set of URL options
func formatOptions(options []Option) string {
	opts, hasOpts := "", false
	if options != nil && len(options) > 0 {
		for _, option := range options {
			if !hasOpts {
				hasOpts = true
				opts += "?"
			} else {
				opts += "&"
			}
			opts += string(option)
		}
	}
	return opts
}

// Trim removes both suffix and prefix
func trim(v string) string {
	return strings.TrimSuffix(strings.TrimPrefix(v, "/"), "/")
}


func For(userID string) Provider {
	return provider.NewRestProvider("user/%s/publickeys/", userID)
}

p := publickeys.For("florimond")


c.FindByID(ctx, p, "key1")