# Go Backend Specification for OKF Graph Server

This reference details the Go architecture required to parse an OKF bundle and serve it over HTTP.

---

## 1. Graph Data Models
To feed Cytoscape.js, the backend must expose a unified JSON graph structure.

### Node Structure
```go
type NodeData struct {
	ID          string   `json:"id"`          // Concept ID (e.g. "tables/users")
	Title       string   `json:"title"`       // Display title
	Type        string   `json:"type"`        // Frontmatter type
	Description string   `json:"description"` // Short description
	Resource    string   `json:"resource"`    // Canonical URI
	Tags        []string `json:"tags"`        // Search/filter tags
	Body        string   `json:"body"`        // Markdown content (excluding frontmatter)
}

type Node struct {
	Data NodeData `json:"data"`
}
```

### Edge Structure
```go
type EdgeData struct {
	ID     string `json:"id"`     // Unique edge identifier
	Source string `json:"source"` // Source node ID
	Target string `json:"target"` // Target node ID
}

type Edge struct {
	Data EdgeData `json:"data"`
}
```

### Response Payload
```go
type GraphResponse struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}
```

---

## 2. Parsing Algorithm
To construct the graph, the server must traverse the directory structure:

1. **Extract Frontmatter**: Delimit the start and end of YAML block via `---` lines, then unmarshal using a YAML library (e.g. `gopkg.in/yaml.v3`).
2. **Extract Markdown Body**: Keep everything after the second `---` delimiter.
3. **Extract Outbound Links**: Use a regular expression to parse relative markdown links:
   ```go
   var linkRegex = regexp.MustCompile(`\[[^\]]*\]\(([^)]+\.md)\)`)
   ```
   For each match, clean the relative path using `path.Clean()` and resolve it against the source file's directory to find the target Concept ID.
   *Example*: A link `[Users](../tables/users.md)` inside `datasets/ga4.md` resolves to `tables/users`.
4. **Build Edges**: Append a directed edge from the source Concept ID to the resolved target Concept ID.

---

## 3. Go HTTP Endpoint
Create a handler to deliver the graph JSON:

```go
func GraphHandler(respWriter http.ResponseWriter, req *http.Request) {
	// Parse the OKF bundle
	graph, err := ParseOKFDirectory("./wiki")
	if err != nil {
		http.Error(respWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	respWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(respWriter).Encode(graph)
}
```
