# Project Health Check
Intent: Runs formatting, linting, and tests to ensure the repo is PR-ready.

Steps:
1. Run `make format` to clean up code style.
2. Run `make lint` to check for logic errors or PEP 8 violations.
3. Run `make test` to execute the full pytest suite.
4. If everything passes, output: "✅ Code is clean and tests passed. Ready for submission."
5. If any step fails, summarize the specific errors that need fixing.
