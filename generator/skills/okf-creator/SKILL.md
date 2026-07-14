---
name: okf-creator
description: Open Knowledge Format (OKF) creator skill. Guides the agent to create and enrich structured LLM-wikis (knowledge bundles) for datasets, tables, APIs, and concepts, following the OKF v0.1 draft specification.
---

# OKF Creator Instructions

Use this skill when the user wants to document, catalog, or create a wiki/documentation for data assets, systems, metrics, or APIs in Open Knowledge Format (OKF).

## Step 1: Read the OKF Specification and Templates
Use `load_skill_resource` to read:
- `references/okf-spec.md` (to understand YAML frontmatter keys, reserved filenames, and directory layout requirements)
- `references/templates.md` (to view standard templates and few-shot examples of tables, datasets, indexes, and references)

## Step 2: Research and Information Gathering
If the concepts to be documented require external context or up-to-date schema definitions:
1. Use the `web_search` tool to find authoritative documentation, schema definitions, or references.
2. Use the `fetch_url` tool to read full pages of relevant documentation.
3. Parse and extract schemas, descriptions, query examples, and key relationships.

## Step 3: Plan the Bundle Directory Structure
Organize your bundle into a clean hierarchy:
- Root: `index.md` (directory list), `log.md` (changelog)
- Subdirectories to group related concepts (e.g., `datasets/`, `tables/`, `metrics/`, `references/`)
- Assign a unique Concept ID to each file (the path relative to the bundle root, without the `.md` extension).

## Step 4: Write Concept Documents
For each concept, write a UTF-8 markdown file with:
1. **YAML Frontmatter**: Include all required and recommended keys:
   - `type` (e.g., `BigQuery Table`, `BigQuery Dataset`, `Metric`, `Reference`)
   - `title` (human-readable display name)
   - `description` (one-line summary)
   - `resource` (canonical URI to asset, if applicable)
   - `tags` (relevant search tags)
   - `timestamp` (ISO 8601 current date/time)
2. **Markdown Body**:
   - Provide high-quality explanations, schema definitions, and usage examples.
   - Link related concepts together using relative markdown links (e.g., `[Transactions Table](../tables/transactions.md)`).
   - Cite external sources by linking to the URLs.

## Step 5: Generate Index and Log Files
1. Create an `index.md` at the root and at each subdirectory level. Each `index.md` must list its child concepts and subdirectories with brief descriptions.
2. Create a `log.md` at the root documenting the creation date, author, and description of the bundle initialization.
