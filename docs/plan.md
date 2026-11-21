# News-Terminal Architecture & Implementation Plan

This document outlines the detailed 3-stage plan for structuring the `news-terminal` application into a robust concurrent architecture. It covers the server layer, initial startup synchronization, source-isolated background workers with job status tracking and task queues, resilient LLM pipelines, and the TUI interface, along with technical patterns, file structures, protocols, and testing strategies.

## Stage 1: Server and Initial Startup Sync

### 1.1 Server Layer & File Structure
- **File Structure:** Create `internal/server/server.go` to house the query service wrapping `db.DB`.
- **Pattern & Protocol:** Implement a clean service pattern acting as the direct protocol boundary for the TUI. Expose thread-safe methods such as `GetSources()`, `GetStoriesBySource(sourceID int64, limit, offset int64)`, and `GetStoryDetail(storyID int64)`.
- **Testing:** Add unit tests in `internal/server/server_test.go` using a temporary SQLite database to verify queries.

### 1.2 Initial Startup Sync
- **Implementation:** Inside `main.go`, check if any sources in the database lack a `last_refreshed_at` timestamp.
- **Synchronous Execution:** If empty, synchronously fetch from all sources (Hacker News, Lobsters, TLDR), run resilient LLM batch processing, persist to the database, and set initial refresh times before starting background routines.
- **Testing:** Verify startup sync behavior with integration tests.

## Stage 2: Background Worker Goroutines, Job Queues, & Resilient Pipelines

### 2.1 Source-Isolated Workers & Job Tracking (`internal/worker/job.go`, `internal/worker/worker.go`)
- **File Structure:** Implement job state models (`StateIdle`, `StateRunning`, `StateSuccess`, `StateFailed`) and dedicated task queues per source.
- **Pattern & Protocol:** Use buffered Go channels or task dispatchers where each source (Hacker News, Lobsters, TLDR) has an isolated worker routine. 
- **Job Status & Retries:** Track execution statuses, last errors, and timestamps per source task. If a specific source sync fails, it records `StateFailed`, isolates the error, and schedules a backoff retry without impacting other sources.
- **Resilient LLM Pipeline:** Separate content fetching and LLM summarization/cleanup into dedicated pipeline stages with fallback error handling.
- **Testing:** Add unit tests in `internal/worker/worker_test.go` mocking scrapers, LLM providers, and verifying queue dispatch/retry logic.

## Stage 3: TUI Goroutine & Concurrency Orchestration

### 3.1 TUI Integration & File Structure
- **File Structure:** Organize TUI code under an `internal/tui/` directory (e.g., `internal/tui/app.go`).
- **Pattern & Protocol:** Run the TUI interface within its own goroutine or main application thread as required by the chosen TUI framework (e.g., BubbleTea). Communicate with the server layer via direct method calls or message passing channels.
- **Testing:** Use component testing and model unit tests for the TUI rendering and event handling logic.

### 3.2 Orchestration & Lifecycle Management
- **Main Wiring (`main.go`):** 
  1. Load config and initialize SQLite database.
  2. Run Stage 1 startup check and potential initial sync.
  3. Instantiate the Server and Source-Isolated Background Workers with Job Queues.
  4. Spawn the Background Workers and TUI inside separate goroutines.
- **Graceful Shutdown:** Implement signal handling (`SIGINT`, `SIGTERM`) using `os.Signal` and `context.WithCancel` to gracefully shut down the workers, server, and TUI without data loss or connection leaks.
