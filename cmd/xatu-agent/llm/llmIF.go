package llm

type LLMClientIF interface {
	Call(prompt string) (string, error)
}
