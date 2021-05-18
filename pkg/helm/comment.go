package helm

func ParseComment(commentLines []string) (string, ChartValueDescription) {
	var valueKey string
	var c ChartValueDescription
	var docStartIdx int

	for i := range commentLines {
		match := valuesDescriptionRegex.FindStringSubmatch(commentLines[i])
		if len(match) < 3 {
			continue
		}

		valueKey = match[1]
		c.Description = match[2]
		docStartIdx = i
		break
	}

	for _, line := range commentLines[docStartIdx+ 1:] {
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
