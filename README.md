# Assignments for CS146S: The Modern Software Developer

This is the home of the assignments for [CS146S: The Modern Software Developer](https://themodernsoftware.dev), taught at Stanford University fall 2025.

## Repo Setup
These steps work with Python 3.12.

1. Install Anaconda
   - Download and install: [Anaconda Individual Edition](https://www.anaconda.com/download)
   - Open a new terminal so `conda` is on your `PATH`.

2. Create and activate a Conda environment (Python 3.12)
   ```bash
   conda create -n cs146s python=3.12 -y
   conda activate cs146s
   ```

3. Install Poetry
   ```bash
   curl -sSL https://install.python-poetry.org | python -
   ```

4. Install project dependencies with Poetry (inside the activated Conda env)
   From the repository root:
   ```bash
   poetry install --no-interaction
   ```

## Week 1 — Prompting Techniques
- **K-Shot Prompting**: Overcoming tokenization limits for character-level tasks.
- **Zero-Shot CoT**: Using mathematical theorems (Euler's Totient) to solve large exponents.
- **Tool Calling**: Forcing models to output strictly parsable JSON for automated execution.
- **Majority Voting**: Running multiple iterations and taking the most common answer to reduce errors.
- **RAG**: Providing external API documentation to the model to generate accurate code.
- **Reflexion**: Implementing a loop where the model corrects its code based on test failures.

## Week 2 — AI IDE
Using Antigravity to evolve a basic note-taking app with programmatic heuristic extraction into a modern, robust, and AI-powered full-stack application.

## Week 3 — Build a Custom MCP Server
A Model Context Protocol (MCP) server written in Go that wraps the [Open-Meteo](https://open-meteo.com/) weather API. It runs locally over **STDIO transport** and integrates with any MCP-compatible client.

## Week 4 - The Autonomous Coding Agent
 Build automations with custom slash commands-`.claude/skills/*.md`, `CLAUDE.md` files for repository or context guidance, Claude SubAgents (role-specialized agents working together) that improves developer workflow.

## Week 6 - Scan and Fix Vulnerabilities with Semgrep
Run static analysis against the provided app using **Semgrep**. Triage findings and remediate security issues.
