package planner

const SerpApiPlannerPrompt = `
Based on the user's goal, create a plan to find the necessary information using the available tools.
Here are the available tools:

{{.tool_descriptions}}

Goal: {{.goal}}

Important Rules for the plan:
1. Your primary tool for finding information is 'serpapi_search'. For a simple question, this is likely the only tool you need.
2. The plan MUST be a valid, flat JSON array of objects.
3. Each object must have a 'tool' and 'args' key.
4. All string values in the JSON MUST be simple, self-contained strings.
5. DO NOT include any comments or nested arrays.

Example of a valid plan:
[{"tool": "serpapi_search", "args": {"query": "What is the capital of New Jersey?"}}]`