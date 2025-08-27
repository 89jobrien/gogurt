package rag

const BasicRagPrompt = `
Answer the following question based on this context:
---
Context:
{{.context}}
---
Question: {{.question}}`
