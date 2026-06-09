# Go Job Queue

A lightweight background job queue written in Go with Redis. The project provides an HTTP API for submitting jobs and checking their status, while a background worker processes jobs asynchronously.

## Features

* HTTP API for creating jobs
* Redis-backed queue
* Background worker for processing jobs
* Job status tracking
* Automatic retry support (up to 3 attempts)
* Graceful shutdown support
* Example job types:

  * `email`
  * `resize_image`

---

## Tech Stack

* Go
* Redis
* HTTP REST API
* UUIDs for job tracking
* Goroutines and channels
* Graceful shutdown with `SIGINT` and `SIGTERM` handling

---

## Prerequisites

### Go

This project requires **Go 1.25+**.

Verify your installation:

```bash
go version
```

Install Go if needed:

https://go.dev/dl/

### Redis

A Redis server must be running locally on the default port:

```
localhost:6379
```

Start Redis with Docker:

```bash
docker run -p 6379:6379 redis
```

Or start your local Redis server:

```bash
redis-server
```

Verify Redis is running:

```bash
redis-cli ping
```

Expected output:

```
PONG
```

---

## How to Run

Clone the repository:

```bash
git clone https://github.com/yfaheid/go-job-queue.git
cd go-job-queue
```

Install dependencies:

```bash
go mod tidy
```

Start Redis if it is not already running.

Run the application:

```bash
go run .
```

The application starts:

* An HTTP API server on `localhost:8080`
* A background worker for processing jobs

Example startup output:

```
Worker started, waiting for jobs...
API listening on :8080
```

---

## API Endpoints

### Enqueue a Job

**POST** `/jobs`

Request body:

```json
{
  "type": "email",
  "payload": "hello@example.com"
}
```

Example:

```bash
curl -X POST http://localhost:8080/jobs \
  -H "Content-Type: application/json" \
  -d '{
    "type":"email",
    "payload":"hello@example.com"
  }'
```

Example response:

```json
{
  "status": "enqueued"
}
```

---

### Check Job Status

Each submitted job receives a unique ID.

Retrieve a specific job:

```
GET /jobs/{job_id}
```

Example:

```bash
curl http://localhost:8080/jobs/<JOB_ID>
```

Example response:

```json
{
  "id": "7c1f8b77-6a13-42d2-a2d3-7b0ec3d5a4d1",
  "type": "email",
  "status": "completed",
  "attempts": "1",
  "max_attempts": "3"
}
```

#### Job Status Values

| Status    | Description                    |
| --------- | ------------------------------ |
| pending   | Waiting to be processed        |
| running   | Currently being processed      |
| completed | Successfully finished          |
| failed    | Maximum retry attempts reached |

---

### List Recent Jobs

Retrieve recently submitted jobs:

```bash
curl http://localhost:8080/jobs
```

Example response:

```json
[
  {
    "id": "...",
    "type": "email",
    "status": "completed",
    "attempts": "1",
    "max_attempts": "3"
  }
]
```

---

## Example Workflow

### 1. Start Redis

```bash
docker run -p 6379:6379 redis
```

### 2. Start the application

```bash
go run .
```

### 3. Submit a job

```bash
curl -X POST http://localhost:8080/jobs \
  -H "Content-Type: application/json" \
  -d '{
    "type":"resize_image",
    "payload":"image.jpg"
  }'
```

### 4. List jobs

```bash
curl http://localhost:8080/jobs
```

Copy the returned job ID.

### 5. Check job status

```bash
curl http://localhost:8080/jobs/<JOB_ID>
```

---

## Supported Job Types

Currently implemented:

| Type         | Description                 |
| ------------ | --------------------------- |
| email        | Simulates sending an email  |
| resize_image | Simulates resizing an image |

Additional job types can be added by extending the worker's processing logic.

---

## Architecture

```
Client
  |
  v
HTTP API
  |
  v
Redis Queue
  |
  v
Background Worker
  |
  v
Job Processing
  |
  v
Job Status Stored in Redis
```

---

## Retry Behavior

* New jobs start with:

  * Status: `pending`
  * Attempts: `0`
  * Maximum attempts: `3`
* Failed jobs are automatically retried.
* After the final retry attempt, the job is marked as `failed`.

---

## Graceful Shutdown

The application handles `SIGINT` and `SIGTERM` signals gracefully.

When a shutdown signal is received (for example, by pressing `Ctrl+C` or when a Docker container is stopped), the application:

* Stops accepting new work.
* Allows in-flight jobs to finish processing.
* Shuts down cleanly without interrupting active jobs.

This behavior helps prevent partially processed jobs and makes the service suitable for production deployments and containerized environments such as Docker or Kubernetes.

Examples:

```bash
docker stop <container>
```

or simply:

```
Ctrl+C
```

will trigger a graceful shutdown sequence instead of terminating the process immediately.

---

## Notes

* Redis must be running before starting the application.
* The API server and background worker run within the same process.
* Job metadata is stored in Redis and can be queried through the API.
* Recent job IDs are maintained for listing recently submitted jobs.
* The application supports graceful shutdown to ensure active jobs complete before exit.
* The project is designed as a simple foundation for building more advanced background processing systems in Go.
