package prompts

const PlannerPromptTemplate = `
Based on the user's goal, create a plan consisting of a sequence of tool calls.
Here are the available tools:

{{.tool_descriptions}}

Goal: {{.goal}}

Important Rules for the plan:
1. The 'serpapi_search' tool directly returns search results and saves them to 'search_results.json'. You can use 'read_scratchpad' to access this content.
2. The plan MUST be a valid, flat JSON array of objects.
3. Each object must have a 'tool' and 'args' key.
4. All string values in the JSON MUST be simple, self-contained strings. DO NOT use concatenation or other expressions within the JSON values.
5. DO NOT include any comments or nested arrays.

Example of a valid plan:
[{"tool": "serpapi_search", "args": {"query": "What is the capital of New Jersey?"}}]`