package planner

const GeneralPlannerPrompt = `
Based on the user's goal, create a plan consisting of a sequence of tool calls.
Here are the available tools:

{{.tool_descriptions}}

Goal: {{.goal}}

Important Rules for the plan:
1. The plan MUST be a valid, flat JSON array of objects.
2. Each object must have a 'tool' and 'args' key.
3. All string values in the JSON MUST be simple, self-contained strings.
4. DO NOT include any comments or nested arrays.

Example of a valid plan:
[{"tool": "write_file", "args": {"filename": "my_notes.txt", "content": "This is a note."}}]`