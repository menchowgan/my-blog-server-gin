# My Blog Server Architecture

This diagram illustrates the architecture of the Go Gin web server, showing the flow of a typical API request and the interaction between different components.

```mermaid
graph TD
    subgraph "External"
        Client([<fa:fa-user> User/Client])
        MySQL_DB[(<fa:fa-database> MySQL Database)]
        Redis_Cache[(<fa:fa-bolt> Redis Cache)]
    end

    subgraph "My Blog Server (Go Application)"
        style GoApp fill:#f9f9f9,stroke:#333,stroke-width:2px

        A[<b>Gin Router</b><br>/router/index.go]
        B{<b>Middlewares</b><br>/middlewares/*<br>Logging, Auth, etc.}
        C[<b>API Controllers</b><br>/api/*/*<br>Business Logic]
        D[<b>View Layer (DTO)</b><br>/view/*/*<br>Data Transfer Objects]
        E[<b>Database Layer</b><br>/db/*<br>GORM CRUD Operations]
        F[<b>Redis Layer</b><br>/redis/*<br>Caching Logic]
        G[<b>Response Helpers</b><br>/response/*<br>Standard JSON Output]
        H[<b>Models (Entities)</b><br>/model/*<br>DB Table Structs]
        I[<b>Config</b><br>/config/index.go<br>Loaded from .env]

        A --> B
        B --> C
        C -- Uses --> D
        C -- Calls --> E
        C -- Calls --> F
        C -- Formats with --> G
        E -- Uses --> H
        E -- Connects to --> MySQL_DB
        F -- Connects to --> Redis_Cache
        I -- Injected into --> E & F & C
    end

    Client -- HTTP Request --> A
    G -- HTTP Response (JSON) --> Client

```

### How to View This Diagram

*   **GitHub/GitLab**: If you commit this file to a GitHub or GitLab repository, it will be automatically rendered as a visual diagram.
*   **Online Editors**: You can copy and paste the content of the `mermaid` block into an online editor like [Mermaid Live](https://mermaid.live).
*   **IDE Extensions**: IDEs like VS Code have extensions (e.g., "Markdown Preview Mermaid Support") that can render the diagram directly in your editor.

### Component Descriptions

*   **User/Client**: The end-user or a frontend application making API calls.
*   **Gin Router**: The entry point for all requests. It matches the URL to the appropriate controller handler.
*   **Middlewares**: Functions that execute before the main controller logic, used for cross-cutting concerns like logging, authentication, and error recovery.
*   **API Controllers**: The core of the application, containing the business logic for each API endpoint.
*   **View Layer (DTOs)**: Defines the structure of data sent to and from the API, decoupling the external API from the internal database models.
*   **Database Layer**: Handles all communication with the MySQL database, performing Create, Read, Update, Delete (CRUD) operations.
*   **Redis Layer**: Manages caching logic, communicating with the Redis server to store or retrieve frequently accessed data.
*   **Models (Entities)**: Go structs that map directly to your database tables.
*   **Response Helpers**: Utility functions to create consistent JSON success or error responses.
*   **Config**: Loads and provides access to application configuration (e.g., database credentials, server ports) from environment variables (`.env` file).
