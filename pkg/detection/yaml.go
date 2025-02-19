package detection

// GenerateYAML generates a nexlayer.yaml based on detected project info
func GenerateYAML(info *ProjectInfo) (string, error) {
	return GenerateYAMLFromTemplate(info)
}
