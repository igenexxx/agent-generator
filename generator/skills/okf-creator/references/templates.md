# OKF Templates and Few-Shots Reference

This file contains standard templates and few-shot examples of various OKF document types, derived from real-world datasets like Google Analytics 4 (GA4).

---

## 1. Root Directory index.md Template
Used to index the top-level categories/directories.

### Example (Root index.md from GA4 Bundle)
```markdown
# Subdirectories

* [datasets](datasets/index.md) - A sample of obfuscated Google Analytics BigQuery event export data for three months from the Google Merchandise Store is available as a public dataset in BigQuery.
* [references](references/index.md) - This directory contains specifications, metrics, and references for Google Analytics 4.
```

---

## 2. Category index.md Template
Used to index all concepts inside a specific subdirectory.

### Example (datasets/index.md from GA4 Bundle)
```markdown
# BigQuery Dataset

* [BigQuery sample dataset for Google Analytics ecommerce web implementation](ga4_obfuscated_sample_ecommerce.md) - A sample of obfuscated Google Analytics BigQuery event export data.
```

---

## 3. BigQuery Table Concept Template
Used to document table structures, schemas, metrics, and query examples.

### Example (tables/events_.md from GA4 Bundle)
```markdown
---
type: BigQuery Table
resource: https://bigquery.googleapis.com/v2/projects/bigquery-public-data/datasets/ga4_obfuscated_sample_ecommerce/tables/events_*
title: Events table (Google Analytics BigQuery Export)
description: Contains Google Analytics event export data from the `ga4_obfuscated_sample_ecommerce` dataset.
tags:
- events
- Google Analytics
- BigQuery
- ecommerce
- schema
- basic queries
timestamp: '2026-05-28T22:53:05+00:00'
---

# Overview
The `events_` table is a sharded BigQuery table containing Google Analytics event export data from the `ga4_obfuscated_sample_ecommerce` dataset.

# Metrics
- [Event Count](../references/metrics/event_count.md) — Total number of events.
- [User Count](../references/metrics/user_count.md) — Total number of unique users.

# Schema
The `events_YYYYMMDD` table, created daily, contains the following fields:

## event
The event fields contain information that uniquely identifies an event.

- `batch_event_index` (INTEGER): A number indicating the sequential order of each event within a batch based on their order of occurrence on the device.
- `event_date` (STRING): The date when the event was logged (YYYYMMDD format).
- `event_name` (STRING): The name of the event.

### event_params RECORD
The `event_params` RECORD is repeated for each key that is associated with an event.

- `event_params.key` (STRING): The name of the event parameter.
- `event_params.value` (RECORD): A record containing the event parameter's value.
    - `event_params.value.string_value` (STRING): String representation of the parameter value.
```

---

## 4. Reference / Metric Concept Template
Used to define metrics, formulas, or standard definitions.

### Example (references/metrics/event_count.md from GA4 Bundle)
```markdown
---
type: Reference
resource: https://developers.google.com/analytics/bigquery/web-ecommerce-demo-dataset
title: Event Count
description: Total number of events.
tags:
- metric
timestamp: '2026-05-28T22:50:07+00:00'
---

# Event Count
Total number of events triggered by users.

## Calculation Formula
```sql
SELECT COUNT(*) as event_count
FROM `bigquery-public-data.ga4_obfuscated_sample_ecommerce.events_*`
```
```

---

## 5. Changelog log.md Template
Used to log changes in a chronological table.

### Example
```markdown
# Change Log

| Timestamp | Author | Description |
|-----------|--------|-------------|
| 2026-07-12T11:45:00Z | LLM Agent | Initialized OKF Bundle structure. |
| 2026-07-12T11:50:00Z | LLM Agent | Added tables/events_.md schema details. |
```
