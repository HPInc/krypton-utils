package config

const (
	DefaultProfile = "ci"
)

// apply the current selected profile
func (c *Config) applyProfile() {
	p := c.getProfileByName(c.ProfileName)
	if p == nil {
		logger.Fatalf("profile %s not found", c.ProfileName)
	}
	c.CurrentProfile = p
	p.apply()
}

// look up profile by name
func (c *Config) getProfileByName(profile string) *Profile {
	for _, p := range c.Profiles {
		if p.Name == profile {
			return &p
		}
	}
	return nil
}
