package planner

const DDGSPlannerPrompt = `
Based on the user's goal, create a plan consisting of a sequence of tool calls.
Here are the available tools:

{{.tool_descriptions}}

Goal: {{.goal}}

Important Rules for the plan:
1. The 'serpapi_search' tool directly returns search results. The plan should only contain a single step, which is the search itself.
2. The plan MUST be a valid, flat JSON array of objects.
3. Each object must have a 'tool' and 'args' key.
4. All string values in the JSON MUST be simple, self-contained strings.
5. DO NOT include any comments or nested arrays.

Example of a valid plan:
[{"tool": "serpapi_search", "args": {"query": "What is the capital of New Jersey?"}}]`