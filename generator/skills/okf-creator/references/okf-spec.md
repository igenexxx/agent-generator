# Open Knowledge Format (OKF) Specification Reference

**Version 0.1 — Draft Summary**

OKF is a universal, vendor-neutral format for representing knowledge (metadata, context, schema definitions, queries) as plain markdown files with YAML frontmatter, organized in a directory hierarchy.

## 1. Key Goals
- **Human- and Agent-readable**: Markdown structure can be read via normal tools (e.g. `cat`) or directly by LLMs.
- **Version Control Friendly**: Pure text files allow standard git diffs, pull requests, and blame history.
- **Extensible**: Small set of required fields; producers can add custom YAML frontmatter or Markdown sections.
- **Graph-shaped**: Markdown links between files create semantic relationships beyond the filesystem tree.

## 2. Directory Layout and Reserved Files
A knowledge bundle is structured as follows:

```
my-bundle/
├── index.md                      # Required/Recommended: Directory catalog for progressive disclosure.
├── log.md                        # Optional: Chronological history of updates.
└── <category-subdirectory>/      # Subdirectories for logical grouping (e.g. tables, metrics, datasets).
    ├── index.md                  # Catalog of files in this subdirectory.
    └── <concept-document>.md     # A concept/asset document (e.g. users.md).
```

### Reserved Filenames
- `index.md`: Used for listing the contents of a directory. It lists subdirectories and concepts with descriptions.
- `log.md`: Used for listing the history of bundle updates.

## 3. Concept Document Format
Each document is a UTF-8 markdown file consisting of:
1. **YAML Frontmatter**: Delimited by `---` lines.
2. **Markdown Body**: Detailed prose, tables, schemas, or query samples.

### YAML Frontmatter Fields
- `type` (REQUIRED): A short string category, e.g. `BigQuery Table`, `BigQuery Dataset`, `Metric`, `Reference`, `API Endpoint`.
- `title` (RECOMMENDED): Human-readable display name.
- `description` (RECOMMENDED): One-line summary.
- `resource` (RECOMMENDED): Canonical URI for the underlying asset (e.g. Google Cloud resource link).
- `tags` (RECOMMENDED): Array of tags for categorization/filtering.
- `timestamp` (RECOMMENDED): ISO 8601 datetime of the last update.

## 4. Link & Relationship Rules
- **Internal Links**: Link from one concept to another using relative markdown links (e.g. `[Users Table](../tables/users.md)`). Standard parser/viewer tools rewire these paths.
- **External Citations**: Standard HTTP/HTTPS links referencing external documentation, web pages, or specs.
