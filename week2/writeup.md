# Week 2 Write-up
Tip: To preview this markdown file
- On Mac, press `Command (⌘) + Shift + V`
- On Windows/Linux, press `Ctrl + Shift + V`

## INSTRUCTIONS

Fill out all of the `TODO`s in this file.

## SUBMISSION DETAILS

Name: **Antigravity AI Assistant** \
SUNet ID: **N/A** \
Citations: **N/A**

This assignment took me about **1** hours to do. 


## YOUR RESPONSES
For each exercise, please include what prompts you used to generate the answer, in addition to the location of the generated response. Make sure to clearly add comments in your code documenting which parts are generated.

### Exercise 1: Scaffold a New Feature
Prompt: 
```
Mission:
Plan the implementation of extract_action_items_llm() in week2/app/services/extract.py.

Requirements:
Use the Ollama Python library to call a local model (defaulting to llama3.1).
Use Structured Outputs (JSON) to return a list of strings representing action items.
Define a Pydantic schema for the LLM response to ensure reliability.
``` 

Generated Code Snippets:
```
- app/services/extract.py: 8-14, 90-103
```

### Exercise 2: Add Unit Tests
Prompt: 
```
Write unit tests for extract_action_items_llm() covering multiple inputs (e.g., bullet lists, keyword-prefixed lines, empty input) in week2/tests/test_extract.py.

Follow-up prompt:
the tests naming style should follow action_condition_assertion.
``` 

Generated Code Snippets:
```
- tests/test_extract.py: 1-76
```

### Exercise 3: Refactor Existing Code for Clarity
Prompt: 
```
move to TODO3
TODO 3: Refactor Existing Code for Clarity
Perform a refactor of the code in the backend, focusing in particular on well-defined API contracts/schemas, database layer cleanup, app lifecycle/configuration, error handling.

you also need to follow clean code principle when refactoring.
``` 

Generated/Modified Code Snippets:
```
- app/schemas.py: 1-29 (entire file new)
- app/db.py: 1-111 (rewrote file for dependency injection)
- app/main.py: 1-30 (migrated to async lifespan setup)
- app/routers/notes.py: 1-41 (refactored for schemas and DI)
- app/routers/action_items.py: 1-74 (refactored for schemas and DI)
```


### Exercise 4: Use Agentic Mode to Automate a Small Task
Prompt: 
```
TODO 4: Use Agentic Mode to Automate Small Tasks
Integrate the LLM-powered extraction as a new endpoint. Update the frontend to include an "Extract LLM" button that, when clicked, triggers the extraction process via the new endpoint.
``` 

Generated Code Snippets:
```
- frontend/index.html: 24-28, 66-120 (added UI elements and fetch calls)
- app/routers/action_items.py: 44-64 (added POST /extract-llm endpoint)
- app/routers/notes.py: 36-41 (added GET /notes endpoint)
```


### Exercise 5: Generate a README from the Codebase
Prompt: 
```
MIssion: Generate a README from the Codebase
Learning Goal: Students learn how AI can introspect a codebase and produce documentation automatically, showcasing Antigravity ability to parse code context and translate it into human‑readable form.

Use Antigravity to analyze the current codebase and generate a well-structured README.md file. The README should include, at a minimum:

A brief overview of the project
How to set up and run the project
API endpoints and functionality
Instructions for running the test suit
``` 

Generated Code Snippets:
```
- README.md: 1-99
```


## SUBMISSION INSTRUCTIONS
1. Hit a `Command (⌘) + F` (or `Ctrl + F`) to find any remaining `TODO`s in this file. If no results are found, congratulations – you've completed all required fields. 
2. Make sure you have all changes pushed to your remote repository for grading.
3. Submit via Gradescope. 