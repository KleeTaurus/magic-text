package magictext

const (
	DefaultPrompt = "`reset` `no quotes` `no explanations` `no prompt` `no self-reference` `no apologies` `no filler` `just answer` "

	GenerateSummaryPrompt = DefaultPrompt + `
I will give you text content, you will rewrite it and output that in a short 
summarized version of my text. Keep the meaning the same. Ensure that the 
revised content has fewer characters than the original text and no more than 150 
Chinese words.

Only give me the output and nothing else. Now, using the concepts above, 
summarize the following text. Respond in Chinese language.

[text]

%s

[output]
`

	GenerateSummaryPromptWithTopic = DefaultPrompt + `
I will give you text content, you will rewrite it and output that in a short 
summarized version of my text. Keep the meaning the same. Ensure that the 
revised content has fewer characters than the original text and no more than 150 
Chinese words.

When generating text summaries, expand around the following topics as much as 
possible: %s` + `

Only give me the output and nothing else. Now, using the concepts above, 
summarize the following text. Respond in Chinese language.

[text]

%s

[output]
`

	GenerateTitlePrompt = DefaultPrompt + `
I want you to act as a video title generator. I will provide you with a topic 
or theme, and you need to generate only one creative and catchy video title for 
that topic. Your title should be short, attention-grabbing, and descriptive of 
the content of the video. Avoid using cliches or generic titles. The number of 
characters in the title should not exceed 25 Chinese characters. Respond in 
Chinese language.

[text]

%s
	
[output]
`

	ExtractNounsPrompt = DefaultPrompt + `
Find all user names, company names, product names, course names, and book names 
from the following text, and output them in the json format. Respond in Chinese 
language.

[output format]
{
	"usernames": [],
	"company_names": [],
	"product_names": [],
	"course_names": [],
	"book_names": [],
}

[text]

%s

[output]
`
)
