# Week 4 Write-up
Tip: To preview this markdown file
- On Mac, press `Command (⌘) + Shift + V`
- On Windows/Linux, press `Ctrl + Shift + V`

## INSTRUCTIONS

Fill out all of the `TODO`s in this file.

## SUBMISSION DETAILS

Name: **HungUnicorn** \
SUNet ID: **HUnicorn** \
Citations: **Claude Code Documentation (https://docs.anthropic.com/en/docs/claude-code)**

This assignment took me about **4** hours to do. 


## YOUR RESPONSES
### Automation #1
a. Design inspiration (e.g. cite the best-practices and/or sub-agents docs)
> My inspiration comes from the [Claude Code Best Practices](https://www.anthropic.com/engineering/claude-code-best-practices) guide, which highlights building common developer operations (formatting, linting, testing) directly into Claude to act as an intelligent CI/CD orchestrator on a developer's local machine.

b. Design of each automation, including goals, inputs/outputs, steps
> **Goals**: Run formatting, linting, and unit tests to ensure the repository is ready for a PR. 
> **Inputs**: None. 
> **Outputs**: A success message ("✅ Code is clean and tests passed.") or a summarized list of specific errors. 
> **Steps**: Run `make format`, run `make lint`, run `make test`, evaluate pass/fail criteria, and synthesize failures into actionable summaries.

c. How to run it (exact commands), expected outputs, and rollback/safety notes
> **Command**: `/verify`
> **Expected outputs**: If everything passes, it outputs the ready for submission message. If it fails, Claude gives a summarized explanation of the errors.
> **Rollback/safety**: Modifies formatting locally which is easily reversible with git. Linting and testing are read-only. Safe to run at any time.

d. Before vs. after (i.e. manual workflow vs. automated workflow)
> **Before**: A developer manually runs `make format`, `make lint`, and `make test` separately. If anything fails, they must manually scroll up and read through stack traces to figure out what broke.
> **After**: A single command runs the entire pipeline, and rather than hunting through stack traces, Claude synthesizes any errors and can automatically proceed to fix them.

e. How you used the automation to enhance the starter application
> I used `/verify` extensively when developing the backend tests (`test_notes_search.py`, `test_ping.py`, `test_validation.py`) and expanding the backend notes endpoints. It allowed me to rapidly format my code and ensure that nothing broke, enforcing CI checks before committing.


### Automation #2
a. Design inspiration (e.g. cite the best-practices and/or sub-agents docs)
> This automation was inspired by the [SubAgents documentation](https://docs.anthropic.com/en/docs/claude-code/sub-agents) and "Explain code" examples, turning Claude into an automated code-tour guide that consistently breaks down module logic and architecture.

b. Design of each automation, including goals, inputs/outputs, steps
> **Goals**: Explain a given file's logic, its dependencies, and its role within the larger FastAPI architecture.
> **Inputs**: A file path argument (`$PATH`).
> **Outputs**: A structured markdown breakdown of the file along with a one-sentence summary for a PR description.
> **Steps**: Read contents of `$PATH`, list primary functions/classes, identify external dependencies, explain interactions with the FastAPI lifecycle, and provide the PR summary.

c. How to run it (exact commands), expected outputs, and rollback/safety notes
> **Command**: `/walkthrough backend/app/routers/notes.py`
> **Expected outputs**: A concise, consistent markdown breakdown explaining the module, responsibilities, and external dependencies.
> **Rollback/safety**: Completely read-only and safe to run. It does not modify any code or state.

d. Before vs. after (i.e. manual workflow vs. automated workflow)
> **Before**: Onboarding developers have to manually trace code to understand module dependencies, or ask open-ended questions that yield wildly different documentation styles.
> **After**: Devs can use a standardized prompt that guarantees a consistent architectural summary and PR description for any file in the repository.

e. How you used the automation to enhance the starter application
> I used `/walkthrough` on core files like `backend/app/routers/notes.py` to quickly understand how the FastAPI lifecycle was currently structured. It made it easier to extend the application with searching capabilities, and I also used it to generate sections for the comprehensive `API.md` file.


### *(Optional) Automation #3*
*If you choose to build additional automations, feel free to detail them here!*

a. Design inspiration (e.g. cite the best-practices and/or sub-agents docs)
> N/A

b. Design of each automation, including goals, inputs/outputs, steps
> N/A

c. How to run it (exact commands), expected outputs, and rollback/safety notes
> N/A

d. Before vs. after (i.e. manual workflow vs. automated workflow)
> N/A

e. How you used the automation to enhance the starter application
> N/A
