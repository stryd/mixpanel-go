package mixpanel

type Option interface {
	Apply(*Config)
}

type ApiUrlOption struct {
	value string
}

func (o ApiUrlOption) Apply(c *Config) {
	c.ApiUrl = o.value
}

func WithApiUrl(value string) Option {
	return ApiUrlOption{value: value}
}

type TokenOption struct {
	value string
}

func (o TokenOption) Apply(c *Config) {
	c.Token = o.value
}

func WithToken(value string) Option {
	return TokenOption{value: value}
}

type SecretOption struct {
	value string
}

func (o SecretOption) Apply(c *Config) {
	c.Secret = o.value
}

func WithSecret(value string) Option {
	return SecretOption{value: value}
}

type ProjectIDOption struct {
	value string
}

func (o ProjectIDOption) Apply(c *Config) {
	c.ProjectID = o.value
}

func WithProjectID(value string) Option {
	return ProjectIDOption{value: value}
}
