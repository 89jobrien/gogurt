package prompts

const RagPrompt = `
Answer the following question based on this context:
---
Context:
{{.context}}
---
Question: {{.question}}`