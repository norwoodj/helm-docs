package helm

func ParseComment(valueKey string, commentLines []string) (string, ChartValueDescription) {
	var c ChartValueDescription

	match := valuesDescriptionRegex.FindStringSubmatch(commentLines[0])
	if match[1] != "" {
		valueKey = match[1]
	}

	c.Description = match[2]

	for _, line := range commentLines[1:] {
		defaultCommentMatch := defaultValueRegex.FindStringSubmatch(line)

		if len(defaultCommentMatch) > 1 {
			c.Default = defaultCommentMatch[1]
			continue
		}

		commentContinuationMatch := commentContinuationRegex.FindStringSubmatch(line)

		if len(commentContinuationMatch) > 1 {
			c.Description += " " + commentContinuationMatch[1]
			continue
		}
	}

	return valueKey, c
}
