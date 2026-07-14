---
name: adk-go-coder
description: Comprehensive guide for building agents, multi-agent systems, and graph workflows in Go using ADK v2.0. Contains references for bootstrap templates, workflow graphs, human-in-the-loop, resilience, and v1-to-v2 migration.
---

# Go ADK v2.0 Engineering Skill

This skill guides the agent in writing, refactoring, and testing Go applications built on Google Agent Development Kit (ADK) v2.0 (`google.golang.org/adk/v2`) and the new GenAI library (`google.golang.org/genai`).

## Step 1: Analyze Requirements
Determine which ADK v2.0 features are required:
*   **Simple Conversational Agent**: Single LLM agent with tools. (See `references/examples.md` for bootstrap template).
*   **Graph-based Workflow**: Complex control flow, conditional routing, loops, or parallel fanning. (See `references/workflow-graphs.md`).
*   **Human-in-the-Loop (HITL)**: Requires human validation, approval, or input mid-run. (See `references/hitl-resilience.md`).
*   **Migration from v1**: Upgrading existing Go agent code to ADK v2.0. (See `references/migration.md`).

## Step 2: Load the Relevant Reference
Use `load_skill_resource` to read the specific reference file(s) for the chosen pattern:
*   `references/examples.md` - Complete agent bootstrap code and tool registration.
*   `references/workflow-graphs.md` - Node definitions (Function, Emitting, Dynamic, Join, State-bound), routing, and loops.
*   `references/hitl-resilience.md` - Pausing workflows for human approval, handoff/re-entry modes, retry configurations, and isolation scopes.
*   `references/migration.md` - Signature changes, `agent.Context` usage, and event stream assertions in tests.

## Step 3: Implement Code
Follow these core Go ADK v2.0 rules:
1.  **Imports**: Always use `google.golang.org/adk/v2/...` module paths.
2.  **Context**: Use `agent.Context` for all tools, callbacks, and workflow node functions.
3.  **LLM Agent Modes**: Set agent modes (`Chat`, `Task`, `SingleTurn`) in `llmagent.Config` to auto-install necessary orchestrator tools.
4.  **Resilience**: Carry retry and timeout configurations directly inside `workflow.NodeConfig`.
5.  **State Plumbing**: Avoid manual session state retrieval; use `NewFunctionNodeFromState` with `state:"<key>"` struct tags.
