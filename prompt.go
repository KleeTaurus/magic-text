package magictext

const (
	generateSummaryPrompt = "`reset` " +
		"`no quotes` " +
		"`no explanations` " +
		"`no prompt` " +
		"`no self-reference` " +
		"`no apologies` " +
		"`no filler` " +
		"`just answer` " + `
I will give you text content, you will rewrite it and output that in a short 
summarized version of my text. Keep the meaning the same. Ensure that the 
revised content has significantly fewer characters than the original text, 
and no more than 150 Chinese words, the fewer the better.
Only give me the output and nothing else. Now, using the concepts above, 
summarize the following text. Respond in Chinese language:
`

	generateSummaryPromptWithTopics = "`reset` " +
		"`no quotes` " +
		"`no explanations` " +
		"`no prompt` " +
		"`no self-reference` " +
		"`no apologies` " +
		"`no filler` " +
		"`just answer` " + `
I will give you text content, you will rewrite it and output that in a short 
summarized version of my text. Keep the meaning the same. Ensure that the 
revised content has significantly fewer characters than the original text, 
and no more than 150 Chinese words, the fewer the better.
` +
		`When generating text summaries, expand around the following topics as
much as possible: %s` +
		`Only give me the output and nothing else. Now, using the concepts above, 
summarize the following text. Respond in Chinese language:
`

	generateTitlePrompt = ""

	extractNounsPrompt = ""
)
