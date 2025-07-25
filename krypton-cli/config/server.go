package config

const (
	Default = "default"
)

func (s *Server) getAddress(module, server string) string {
	if v, ok := s.Addresses[server]; ok {
		return v
	} else if v, ok := s.Addresses[Default]; ok {
		return v
	}
	logger.Debug("Error looking up address: ", module, server, s.Addresses)
	return ""
}
