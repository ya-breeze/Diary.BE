# Diary.BE
A web-based personal diary application with a focus on privacy and simplicity.
Write, store, and manage your daily thoughts and experiences securely.

## Features
- 🔒 Secure authentication with JWT tokens
- 📝 Markdown support for rich text formatting
- 🌐 Web-based interface using Bootstrap
- 📱 Responsive design for mobile and desktop
- 💾 SQLite database for easy deployment
- 🐳 Docker support for containerized deployment

## Technical Overview
The application is built using:
- Go (backend)
- Bootstrap (frontend)
- SQLite (database)
- OpenAPI/Swagger (API documentation)

### Architecture
The application follows a layered architecture:
1. **Web Layer** - Handles HTTP requests and user interface (webapp/)
2. **API Layer** - RESTful API defined using OpenAPI spec (api/openapi.yaml)
3. **Service Layer** - Business logic implementation
4. **Data Layer** - Database operations using GORM

### Security
- Password hashing using bcrypt
- Session management with secure cookies
- JWT-based API authentication

## Getting Started
### Prerequisites
- Go 1.24+
- SQLite
- Docker (optional)

### Quick Start
1. Using Docker:
```bash
docker-compose up
```

2. Using Make:
```bash
make run
```

Default credentials:
- Username: test
- Password: test

## Development
- `make build` - Build the application
- `make test` - Run tests
- `make lint` - Run linters
- `make validate` - Validate code

## Configuration
Environment variables:
- `GB_USERS` - User credentials
- `GB_DBPATH` - Database path
- `GB_ASSETPATH` - Assets path
