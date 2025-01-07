# Modular APIs

- **Idea:** A user should be able to easily choose which KG and weather API they'd like to use, and a dev should be able to easily add one.
- **How:** Along the lines of:
```
type KnowledgeGraphProvider struct {
    // OpenAI GPT-4o
	Name         string
    // gpt-4o
    ID           string
    // openai.com
	SignUpLink   string
	NeedsPayment bool
    // completions endpoint
	APIEndpoint  string
	// streaming capable
	IGCapable    bool
    // smart enough for LLM commands?
    Smart        bool
    KG           KGFuncs
}

type KGFuncs interface {
    Process(input string, ig bool, output chan string, end chan bool) error
    Test() bool
}
```
- This selection will be in (the equivalent of) Bot Settings, as this selection will be per-bot rather than wire-pod wide. Though, I will probably allow for a default option.