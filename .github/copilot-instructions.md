# Diary.BE - Personal Diary Application

Diary.BE is a Go-based web application that provides a secure personal diary with REST API, SQLite database, and Bootstrap frontend. The application supports diary entries, asset management, and JWT-based authentication.

**ALWAYS follow these instructions first and only search or use bash commands when the information here is incomplete or found to be in error.**

## Working Effectively

### Prerequisites and Environment Setup
- Go 1.24+ is required
- Docker for containerization and API generation
- SQLite for database operations

### Essential Go Tools Installation
These tools are REQUIRED for development and MUST be installed before starting:
```bash
# Install required Go tools
go install mvdan.cc/gofumpt@v0.8.0
go install golang.org/x/tools/cmd/goimports@latest
```

Add Go tools to PATH:
```bash
export PATH=$PATH:~/go/bin
```

### Bootstrap, Build, and Test Commands

#### Build the Application
```bash
make build
```
- **TIME ESTIMATE**: 1-2 minutes (first build includes dependency downloads)
- **NEVER CANCEL**: Wait for completion
- **OUTPUT**: Creates `bin/diary` executable (~22MB)

#### Run All Tests
```bash
make test
```
- **TIME ESTIMATE**: 15-20 seconds
- **NEVER CANCEL**: Set timeout to 30+ minutes for safety
- **VALIDATION**: All tests should pass, includes both unit tests and integration flow tests
- Uses Ginkgo test framework with both API service tests and end-to-end flow tests

#### Run Linting (CRITICAL for CI)
```bash
export PATH=$PATH:~/go/bin  # Required for gofumpt and goimports
make lint
```
- **TIME ESTIMATE**: 1-2 seconds (after tools are installed)
- **DEPENDENCIES**: Requires gofumpt and goimports in PATH
- **CI REQUIREMENT**: ALWAYS run before committing or CI will fail

#### Validate OpenAPI Specification
```bash
export HOST_PWD=$(pwd)
make validate
```
- **TIME ESTIMATE**: 5-10 seconds (first run downloads Docker image)
- **DEPENDENCIES**: Requires Docker and HOST_PWD environment variable
- **PURPOSE**: Validates api/openapi.yaml specification

#### Run the Application
```bash
make run
```
- **STARTUP TIME**: 2-3 seconds
- **DEFAULT PORT**: 8080
- **DEFAULT CREDENTIALS**: Username: `test`, Password: `test`
- **DATABASE**: Creates diary.db in project root
- **ASSETS**: Creates diary-assets/ directory for file uploads

#### Run with Docker Compose
```bash
docker compose up --build
```
- **BUILD TIME**: 3-5 minutes (NEVER CANCEL)
- **TIMEOUT**: Set timeout to 10+ minutes
- **PORT**: 8080 (configurable via PORT env var)

## Validation and Testing Scenarios

### MANDATORY Validation Steps
After making ANY changes, ALWAYS perform these validation steps:

1. **Build Validation**:
   ```bash
   make build
   ```

2. **Test Validation**:
   ```bash
   make test
   ```
   - MUST see "SUCCESS! Test Suite Passed"
   - Integration tests validate server startup and API functionality

3. **Linting Validation**:
   ```bash
   export PATH=$PATH:~/go/bin
   make lint
   ```
   - MUST pass or CI will fail

4. **Manual Application Testing**:
   Start the server and test core functionality:
   ```bash
   make run &
   sleep 3
   
   # Test web interface
   curl -s http://localhost:8080/ | head -5
   
   # Test API authentication
   curl -X POST -H "Content-Type: application/json" \
        -d '{"email":"test","password":"test"}' \
        http://localhost:8080/v1/authorize
   
   # Stop server
   pkill -f "diary server"
   ```
   - MUST return HTML for web interface
   - MUST return JWT token for API authentication

### End-to-End User Scenarios
When testing UI or API changes:

1. **Login Flow**: Test authentication with test user (test/test)
2. **Diary Entry**: Create, read, update diary entries
3. **Asset Upload**: Test file upload functionality
4. **API Access**: Validate all REST endpoints work with JWT token

## Code Generation and OpenAPI

### Generate API Client/Server Code
```bash
export HOST_PWD=$(pwd)
make generate
```
- **TIME ESTIMATE**: 30-60 seconds
- **DEPENDENCIES**: Docker, HOST_PWD environment variable
- **GENERATES**: Updates pkg/generated/goclient and pkg/generated/goserver
- **WARNING**: May fail with permission issues - generated code is typically pre-built

### Common Generation Issues
- Permission errors: Generated code is usually committed, regeneration not required for most changes
- Missing HOST_PWD: Always set `export HOST_PWD=$(pwd)` before running

## Project Structure and Navigation

### Key Directories
- `cmd/` - Main application and CLI commands
  - `cmd/main.go` - Application entry point
  - `cmd/commands/` - Server and user management commands
- `pkg/` - Core application packages
  - `pkg/server/` - HTTP server and API implementations
  - `pkg/auth/` - Authentication and JWT handling
  - `pkg/database/` - Database operations and models
  - `pkg/config/` - Configuration management
  - `pkg/generated/` - OpenAPI generated client/server code
- `api/` - OpenAPI specification
  - `api/openapi.yaml` - Complete API specification
- `webapp/` - Frontend templates and static files
  - `webapp/templates/` - HTML templates
  - `webapp/static/` - CSS, JS, and static assets
- `test/` - Test suites
  - `test/flows/` - Integration tests that start server and test complete workflows

### Important Files to Check After Changes
- When modifying API contracts: Always check `pkg/server/api/` implementations
- When changing authentication: Review `pkg/auth/` and test login scenarios
- When updating database models: Verify `pkg/database/` and run tests
- Configuration changes: Check `pkg/config/` and test with different env vars

## Configuration and Environment

### Environment Variables
- `GB_USERS` - User credentials (format: username:hashedpassword)
- `GB_DBPATH` - Database file path (default: ./diary.db)
- `GB_ASSETPATH` - Asset storage directory (default: ./diary-assets)
- `GB_DISABLEIMPORTERS` - Disable data importers (default: true)
- `HOST_PWD` - Required for Docker-based operations (set to $(pwd))

### Default Test Configuration
- Database: In-memory SQLite for tests, file-based for local development
- User: test/test (username/password)
- Port: 8080
- JWT: Auto-generated secret (warns if not set)

## Troubleshooting Common Issues

### Build Failures
- **Go version**: Requires Go 1.24+, check with `go version`
- **Dependencies**: Run `go mod download` if build fails with missing packages

### Test Failures
- **Ginkgo missing**: Install with `go install github.com/onsi/ginkgo/v2/ginkgo@latest`
- **Permission errors**: Ensure test can create temporary directories
- **Port conflicts**: Tests use random ports, but check for conflicts on 8080

### Linting Failures
- **Missing tools**: Install gofumpt and goimports (see Essential Go Tools section)
- **PATH issues**: Ensure `~/go/bin` is in PATH
- **Code formatting**: Run `gofumpt -l -w .` to fix formatting issues

### Runtime Issues
- **Database locked**: Stop any running instances before starting new one
- **Port in use**: Change port with environment variable or stop existing process
- **Asset directory**: Ensure GB_ASSETPATH directory is writable

### Docker Issues
- **Build timeouts**: Docker builds take 3-5 minutes, never cancel
- **Permission errors**: Ensure Docker daemon is running and accessible
- **Image pull failures**: Check internet connectivity for base image downloads

## Performance Expectations

### Timing Guidelines (NEVER CANCEL before these times)
- **Initial build**: 60-90 seconds (includes dependency downloads)
- **Subsequent builds**: 10-30 seconds
- **Test suite**: 15-20 seconds (NEVER CANCEL - set 30+ minute timeout)
- **Linting**: 1-2 seconds (after tools installed)
- **Docker build**: 3-5 minutes (NEVER CANCEL - set 10+ minute timeout)
- **Server startup**: 2-3 seconds
- **API response**: <100ms for most endpoints

### Resource Usage
- **Binary size**: ~22MB (Go binary with SQLite)
- **Memory usage**: ~50MB typical runtime
- **Database**: SQLite file grows with diary entries and assets
- **Asset storage**: Depends on uploaded files

## Architecture Overview

### Technology Stack
- **Backend**: Go 1.24+ with Gorilla Mux router
- **Database**: SQLite with GORM ORM
- **Frontend**: Bootstrap with Go templates
- **API**: REST API defined by OpenAPI 3.0 specification
- **Authentication**: JWT tokens with bcrypt password hashing
- **Testing**: Ginkgo/Gomega test framework

### Request Flow
1. HTTP requests → Gorilla Mux router
2. Authentication middleware → JWT validation
3. API handlers → Business logic in pkg/server/
4. Database operations → GORM → SQLite
5. Response rendering → JSON API or HTML templates

### Security Features
- bcrypt password hashing
- JWT token authentication
- Path traversal protection for assets
- Secure session management
- CORS support for API access