# Go Elasticsearch API

This project is a Go-based search API that leverages Elasticsearch to provide both semantic and keyword-based search functionalities. It's designed to search through a product catalog.

## Features

- **Semantic Search**: Utilizes OpenAI embeddings to understand the meaning behind search queries for more relevant results.
- **Keyword Search**: Falls back to traditional keyword matching if semantic search is not possible or yields no results.
- **Filtering**: Supports filtering search results by tags.
- **Sorting**: Allows sorting results by popularity.

## Project Structure

```
.env.example         # Example environment file
.gitignore           # Git ignore file
api/
  handler.go         # HTTP request handlers
data/
  products.json      # Sample product data
elastic/
  client.go          # Elasticsearch client initialization
  index.go           # Elasticsearch index creation and data insertion
  search.go          # Elasticsearch search logic (Note: current search logic is in api/handler.go)
embedding/
  embedding.go       # OpenAI embedding generation
go.mod               # Go module file
go.sum               # Go module checksums
main.go              # Main application entry point
README.md            # This file
Dockerfile           # Docker configuration for building and running the application
```

## Prerequisites

- Go (version 1.24.1 or higher recommended, as per `go.mod`)
- Elasticsearch (running on `http://localhost:9200`)
- An OpenAI API Key

## Running the Application

### Locally

1.  **Set up Environment Variables**:
    Ensure you have an `OPENAI_API_KEY` environment variable set. You can create a `.env` file in the project root by copying `.env.example` and editing it:
    ```bash
    cp .env.example .env
    # Then, edit .env to add your OPENAI_API_KEY
    # Example .env content:
    # OPENAI_API_KEY="your_actual_openai_api_key"
    ```
    The application loads this at startup.
    Also, ensure Elasticsearch is running and accessible (default: `http://localhost:9200`).

2.  **Run the application**:
    ```bash
    go run main.go
    ```
    The server will start, typically on port 8080 (as defined in `main.go` and exposed in `Dockerfile`).

### Using Docker

1.  **Build the Docker image**:
    From the project root directory:
    ```bash
    docker build -t go-elastic-search .
    ```

2.  **Run the Docker container**:
    Replace `"your_actual_openai_api_key"` with your actual OpenAI API key.
    ```bash
    docker run -p 8080:8080 -e OPENAI_API_KEY="your_actual_openai_api_key" go-elastic-search
    ```
    This command:
    - Runs the container (by default in the foreground, add `-d` for detached mode).
    - Maps port 8080 of the container to port 8080 on your host machine.
    - Passes the `OPENAI_API_KEY` as an environment variable to the application running inside the container.
    The application inside the container will attempt to connect to Elasticsearch. If Elasticsearch is running on your host machine, you might need to use `host.docker.internal:9200` (on Docker Desktop for Mac/Windows) or your host's IP address as the Elasticsearch address for the application. By default, the application uses `http://localhost:9200`.

## Setup

1.  **Clone the repository (if you haven't already):**
    ```bash
    git clone <repository-url>
    cd go-elastic-search
    ```

2.  **Set up environment variables:**
    Create a `.env` file in the root of the project by copying `.env.example`:
    ```bash
    cp .env.example .env
    ```
    Open the `.env` file and add your OpenAI API key:
    ```
    OPENAI_API_KEY=your_openai_api_key_here
    ```

3.  **Install dependencies:**
    Go modules should handle dependencies automatically when you build or run the project. If you need to explicitly download them:
    ```bash
    go mod tidy
    ```

4.  **Ensure Elasticsearch is running:**
    Make sure your Elasticsearch instance is accessible at `http://localhost:9200`.

## Running the Application

1.  **Start the Go application:**
    ```bash
    go run main.go
    ```
    The server will start, and you should see a log message:
    ```
    INFO main.go:20 Server running at :8080
    ```
    The application will automatically:
    - Initialize the Elasticsearch client.
    - Create the `products` index if it doesn't exist.
    - Bulk insert product data from `data/products.json` into the `products` index.

## API Usage

Once the server is running, you can use the search API endpoint:

**Endpoint:** `GET /search`

**Query Parameters:**

-   `q` (required): The search query string.
-   `tag` (optional): Filter results by a specific tag.
-   `sort` (optional): Sort results. Currently supports `popularity`.

**Example Request:**

To search for "modern smartphone" and sort by popularity:
```
http://localhost:8080/search?q=modern%20smartphone&sort=popularity
```

To search for products tagged "electronics":
```
http://localhost:8080/search?q=gadget&tag=electronics
```

**Example Response (JSON):**

The API will return a JSON response containing the search results from Elasticsearch.

## How it Works

1.  When a request hits the `/search` endpoint, the `SearchHandler` in `api/handler.go` is invoked.
2.  It attempts to generate an embedding vector for the query using `embedding.GetEmbedding`.
3.  If embedding generation is successful, a semantic (KNN) search is performed against the `embedding` field in the `products` index.
4.  If embedding fails or the semantic search yields no results, the system falls back to a traditional keyword-based `multi_match` query against the `name`, `description`, and `tags` fields.
5.  Filters (by tag) and sorting (by popularity) are applied to the Elasticsearch query as specified.
6.  The search results are returned as a JSON response.

## Key Files

-   `main.go`: Initializes the Elasticsearch client, creates the index, inserts data, and starts the HTTP server.
-   `elastic/client.go`: Contains the Elasticsearch client initialization logic.
-   `elastic/index.go`: Handles the creation of the `products` index and bulk insertion of product data from `data/products.json`.
-   `api/handler.go`: Contains the `SearchHandler` function which processes search requests, interacts with the embedding service and Elasticsearch, and handles fallback logic.
-   `embedding/embedding.go`: Responsible for generating text embeddings using the OpenAI API.
-   `go.mod`: Defines the project's module and its dependencies.
-   `Dockerfile`: Contains instructions to build and run the application in a Docker container.