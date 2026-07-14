---
name: okf-server-creator
description: Guides the creation of a Go HTTP server and Cytoscape.js/Marked.js frontend to visualize and traverse OKF-based knowledge graphs.
---

# OKF Server Creator Instructions

Use this skill when you need to construct a web server and an interactive graph visualization to represent, explore, and traverse an Open Knowledge Format (OKF) bundle.

## Step 1: Read Technical Specifications
Use `load_skill_resource` to read the following reference files:
- `references/server-spec.md` for Go parser logic, JSON schemas, and routing.
- `references/cytoscape-template.html` for the Cytoscape.js and Marked.js HTML template.

## Step 2: Implement the Go Backend Parser
Create a Go package or utility (e.g. `parser.go`) that:
1. Walks the OKF bundle directory.
2. Identifies all `.md` files (ignoring reserved files like `log.md` unless requested).
3. Parses the YAML frontmatter of each file using a standard parser to extract `type`, `title`, `description`, `resource`, `tags`, and `timestamp`.
4. Scrapes the Markdown body for relative internal links (e.g., `[label](../category/target.md)`) to build the directed connection graph (edges).
5. Serves this unified graph representation as a JSON payload on `/api/graph`.

## Step 3: Implement the HTTP Server
Construct a basic Go web server that:
1. Registers the `/api/graph` JSON endpoint.
2. Serves the static HTML visualizer on the root path `/`.
3. Handles standard graceful shutdowns and CORS (if necessary).

## Step 4: Build the Cytoscape.js and Marked.js Frontend
Using the template, ensure the frontend:
1. Fetches the `/api/graph` data.
2. Initializes Cytoscape.js to draw the force-directed graph. Style nodes differently based on their `type` (e.g. Blue for tables, Green for datasets, Purple for metrics).
3. Renders the selected node's details in a sidebar:
   - Displays all frontmatter metadata.
   - Uses Marked.js to render the Markdown body inside the sidebar.
4. **Enables Traversal**: Intercepts click events on relative links within the sidebar markdown container, converting them to a graph selection update (updates the selected Cytoscape node and details panel, avoiding full-page reload).
