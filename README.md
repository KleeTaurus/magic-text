# magictext

Generate a high-quality summary of a given text input.

# TODO

- [x] Support goroutines to improve performance
- [x] Generate separate levels of summary and write to different files
- [x] Support explicitly set OpenAI APIKey
- [x] Generate parent summary with fixed number of chunks
- [x] Calculate tokens with [tiktoken-go](https://github.com/pkoukk/tiktoken-go)
- [x] Add retry logic for goroutines
- [ ] fix token regex performance issues
- [ ] remove video header and tail from rst file
- [ ] return the error instance, don't call log's fatal function directly
- [ ] Write project documentation