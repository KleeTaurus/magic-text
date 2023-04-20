# magictext

Generate a high-quality summary of a given text input.

# TODO

- [ ] return the error instance, don't call log's fatal function directly
- [ ] Optimize logging output
- [x] Support goroutines to improve performance
- [ ] Write project documentation
- [x] Generate separate levels of summary and write to different files
- [ ] Support explicitly set OpenAI APIKey
- [x] Generate parent summary with fixed number of chunks
- [x] Calculate tokens with [tiktoken-go](https://github.com/pkoukk/tiktoken-go)
- [x] Add retry logic for goroutines
