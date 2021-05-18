package helm

func ParseComment(commentLines []string) (string, ChartValueDescription) {
	var valueKey string
	var c ChartValueDescription

	match := valuesDescriptionRegex.FindStringSubmatch(commentLines[0])
	if match[1] != "" {
		valueKey = match[1]
	}

	c.Description = match[2]

	valueTypeMatch := valueTypeRegex.FindStringSubmatch(c.Description)
	if len(valueTypeMatch) > 0 && valueTypeMatch[1] != "" {
		c.ValueType = valueTypeMatch[1]
		c.Description = valueTypeMatch[2]
	}

	for _, line := range commentLines[1:] {
		defaultCommentMatch := defaultValueRegex.FindStringSubmatch(line)
		notationTypeCommentMatch := valueNotationTypeRegex.FindStringSubmatch(line)

		if len(defaultCommentMatch) > 1 {
			c.Default = defaultCommentMatch[1]
			continue
		}

		if len(notationTypeCommentMatch) > 1 {
			c.NotationType = notationTypeCommentMatch[1]
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
