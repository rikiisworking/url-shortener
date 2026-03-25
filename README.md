# URL Shortener Backend

A high-performance URL shortener service built with **Go**, **Fiber**, **PostgreSQL** (pgx), and **Redis**.  

Designed as a production-grade backend example to showcase clean architecture, fast redirects, caching strategies, and modern Go practices.

![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go)
![Fiber](https://img.shields.io/badge/Fiber-v3-00AEEF?style=for-the-badge)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-4169E1?style=for-the-badge&logo=postgresql)
![Redis](https://img.shields.io/badge/Redis-DC382D?style=for-the-badge&logo=redis)
![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker)

## ✨ Features

- **Ultra-fast redirects** — Redis cache-first lookup (< 10ms typical)
- **Collision-resistant short codes** — Base62 random generation with retry logic
- **Click analytics** — Atomic increment in both Redis and PostgreSQL
- **Production-ready** — Connection pooling (pgxpool), graceful error handling, structured logging
- **Dockerized** — Easy local setup with PostgreSQL + Redis
- **Clean Architecture** — Handler → Service → Storage separation
- **Input validation** — Using go-playground/validator
- **Health check endpoint**

## 🚀 Quick Start

### 1. Clone the repository
```bash
git clone https://github.com/rikiisworking/url-shortener.git
cd url-shortener
```

### 2. Copy environment file
```bash
cp .env.example .env
```

### 3. Start databases with Docker
```bash
docker-compose up -d
```

### 4. Run the server
```bash
go run cmd/server/main.go
```

Server will start at `http://localhost:8080`

## 📡 API Endpoints

| Method | Endpoint              | Description                          | Request Body                  | Response |
|--------|-----------------------|--------------------------------------|-------------------------------|----------|
| POST   | `/api/shorten`        | Create a shortened URL               | `{"url": "https://..."}`      | JSON with short_url |
| GET    | `/:shortCode`         | Redirect to original URL (301)       | -                             | Redirect |
| GET    | `/health`             | Health check                         | -                             | `{"status": "healthy"}` |

### Example: Shorten a URL
```bash
curl -X POST http://localhost:8080/api/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://github.com"}'
```

Response:
```json
{
  "short_code": "aB3cD4e",
  "short_url": "http://localhost:8080/aB3cD4e",
  "original_url": "https://github.com"
}
```

## 🏗️ Architecture

```
Handler (Fiber) → Service (Business Logic) → Repository (pgx + Redis)
```

- **Why Redis first?** Redirects are the hot path. Caching dramatically improves performance and reduces database load.
- **Why pgx instead of GORM?** Lightweight, high-performance, full control over SQL, better for learning and production Go services.
- **Why Fiber?** One of the fastest web frameworks in Go (built on fasthttp), clean Express-like API.

## 🛠️ Tech Stack

- **Language**: Go 1.23+
- **Framework**: [Fiber v3](https://gofiber.io/)
- **Database**: PostgreSQL + [pgx/v5](https://github.com/jackc/pgx) (with pgxpool)
- **Cache**: Redis (go-redis/v9)
- **Validation**: go-playground/validator
- **Config**: Environment variables + godotenv
- **Container**: Docker + docker-compose

## 📁 Project Structure

```bash
.
├── cmd/server/main.go                 # Application entrypoint
├── internal/
│   ├── handler/                       # HTTP handlers
│   ├── service/                       # Business logic
│   ├── storage/                    # Data access (pgx + Redis)
│   ├── model/                         # Domain models
│   ├── config/                        # Configuration
│   └── util/                          # Helpers (short code generator)
├── docker-compose.yml
├── .env.example
└── README.md
```

## 🧪 Testing the Service

1. Shorten a URL using the POST endpoint
2. Open the returned short URL in your browser → should redirect instantly
3. Check Redis and PostgreSQL to verify caching and click count

## 📈 Future Enhancements (Planned / Ideas)

- Rate limiting middleware
- Custom alias support
- URL expiration & deletion
- Analytics dashboard endpoint
- OpenAPI/Swagger documentation
- Unit + integration tests
- Graceful shutdown

## 📄 License

MIT License — feel free to use this code in your own projects or portfolio.

---