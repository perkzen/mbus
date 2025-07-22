
# 🚌 Mbus – Marprom Bus Departures

## Overview

**Mbus** is a web application for searching and viewing Marprom bus departures. It provides a user-friendly interface to select departure and arrival stations, browse available routes, and view accurate, up-to-date schedules.

## ✨ Features

* Search for bus departures by station and date
* View detailed schedule and route information
* Fast, responsive UI with modern design

## 🛠 Tech Stack

**Frontend:**

* React
* TypeScript
* Shadcn UI
* TanStack Query & Router

**Backend:**

* Go
* PostgreSQL
* Redis
* OpenRouteService API

## 🚀 Getting Started

### Prerequisites

Ensure the following tools are installed on your machine:

* Node.js & [pnpm](https://pnpm.io)
* Go (version 1.23+)
* Docker & Docker Compose

---

### ⚙️ Example Environment Variables

#### Backend (`apps/bus-service/.env`)

```env
POSTGRES_URL=postgres://postgres:test@db:5432/test
REDIS_PASSWORD=your_redis_password
REDIS_ADDR=redis:6379
ORS_API_KEY=your_openrouteservice_api_key
```

#### Frontend (`apps/web/.env`)

```env
VITE_API_URL=http://backend:8080/api
```

---

### 🐳 Running with Docker Compose

Ensure your `.env` files are set up, then run:

```bash
docker compose up -d
```

This will start the backend, frontend, and required services like PostgreSQL and Redis.

---

### 💻 Running Locally (Without Docker)

You can also run the frontend and backend directly:

#### Frontend

```bash
cd apps/web
pnpm install
pnpm run dev
```

#### Backend

Make sure PostgreSQL and Redis are running via Docker Compose (`docker-compose.yml`).

```bash
cd apps/bus-service
go mod tidy
go run cmd/server/main.go
```

Alternatively, you can use:

* **Makefile:**

  ```bash
  make serve
  ```

* **[Air](https://github.com/cosmtrek/air)** for hot reload:

  ```bash
  air
  ```


## 📦 Scripts & Tooling

* **Migrations:** Managed using the [`goose`](https://github.com/pressly/goose) CLI with SQL-based migration files located in the `migrations/` directory.

* **Makefile:** Provides convenient shortcuts for common development tasks, including:

    * ✅ Starting the Go backend server
    * 🗃️ Running, reverting, and creating database migrations
    * 🌱 Seeding the database with test or production data
    * 🧹 Truncating all tables (useful for resetting the DB)
    * 🔍 Scraping live data from the Marprom website
    * 📄 Generating Swagger API documentation

You can view all available commands by running:

```bash
make help
```

> The `Makefile` also supports environment variable loading via `.env`, so everything works out of the box.


## To-Do
* [ ] Add unit tests for critical components
* [ ] translation for english
* [ ] Add support for holiday schedules
* [ ] Add a dark mode 
* [ ] Move from postgres to sqlite
* [ ] Add a feature to save favorite routes in local storage

