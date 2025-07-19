# GoPodcasts

GoPodcasts is a backend service written in Go designed to fetch, parse, and locally store podcast feeds. It acts as a personal, self-hostable podcast aggregation engine, providing a robust foundation for building a full-featured podcast listening application.

The core of the service is a data sync pipeline that ingests RSS feeds, stores podcast and episode metadata in a local SQLite database, and gracefully handles duplicates to ensure data integrity.

## Current Status & Features

This project is currently in the **backend engine development phase**. The core data persistence and sync logic is complete and tested.

-   **RSS Feed Fetching:** Reliably fetches podcast data from a given RSS feed URL.
-   **XML Parsing:** Parses complex podcast RSS XML into structured Go objects.
-   **SQLite Database Persistence:**
    -   Creates a local `data.db` file with tables for podcasts and episodes.
    -   Uses foreign keys to maintain the relationship between a podcast and its episodes.
-   **"Upsert" Logic:**
    -   When adding a podcast, it uses `INSERT OR IGNORE` to prevent duplicate entries based on the feed URL.
    -   If a podcast is new, it's inserted; if it already exists, its existing ID is retrieved.
    -   When adding episodes, it uses `INSERT OR IGNORE` to prevent duplicate entries based on the episode's unique enclosure URL.
-   **Robust Error Handling:** The entire pipeline includes contextual error handling to manage network failures, server errors, and malformed data.
-   **Comprehensive Unit Tests:** The core fetching logic is covered by a suite of unit tests that mock HTTP responses to validate both success and failure scenarios.

## Technology Stack

-   **Language:** Go
-   **Database:** SQLite 3
-   **Key Libraries:**
    -   `database/sql` (for standard database interaction)
    -   `encoding/xml` (for RSS parsing)
    -   `net/http` (for fetching feeds)
    -   `github.com/mattn/go-sqlite3` (as the database driver)

## How to Run (Current Version)

The application currently runs as a command-line script that fetches a single, hardcoded podcast feed and saves it to the database.

1.  **Prerequisites:**
    -   Go (version 1.18+ recommended)
    -   A C compiler (required by the `go-sqlite3` driver, e.g., `gcc` on Linux/macOS, `MinGW` on Windows)

2.  **Install Dependencies:**
    ```sh
    go get github.com/mattn/go-sqlite3
    ```

3.  **Run the application:**
    ```sh
    go run main.go
    ```
    Upon successful execution, a `data.db` file will be created in the project root containing the fetched podcast and episode data.

## Project Roadmap

The core data engine is complete. The next phases will focus on building out the application into a full service.

-   [ ] **Phase 1: Project Restructuring**
    -   [x] Reorganize code into logical packages (`storage`, `fetching`, `api`).
-   [ ] **Phase 2: API Development**
    -   [x] Implement an HTTP server.
    -   [ ] Create RESTful API endpoints (`GET /podcasts`, `GET /podcasts/{id}/episodes`).
    -   [x] Create a `POST /podcasts` endpoint to add new podcasts.
-   [ ] **Phase 3: Frontend Client**
    -   [x] Develop a simple web client (initially vanilla, later React) to consume the API.
-   [ ] **Phase 4: Advanced Features**
    -   [ ] Implement user authentication.
    -   [ ] Create a background worker to periodically refresh all subscribed feeds.
    -   [ ] Add functionality to download and cache audio files for offline listening.

## Other TODO

- [ ] Create `GET /podcasts/{id}` endpoint to get more than just metadata
- [ ] Create `GET /podcasts/{id}/episodes` endpoint to get all episodes
  - [ ] Create `GET /podcasts/{id}/episodes/{id}` to get data about a single episode