# Diary.BE
Diary.BE is a personal diary web application built in Go with SQLite database, Bootstrap frontend, and OpenAPI-generated server/client code. The application provides secure diary entry management with JWT authentication, file asset management, and a responsive web interface.

Always reference these instructions first and fallback to search or bash commands only when you encounter unexpected information that does not match the info here.

## Working Effectively

### Bootstrap, Build, and Test the Repository
- **CRITICAL**: All commands below MUST be run from the repository root directory
- Install required tools: `go install mvdan.cc/gofumpt@v0.7.0`
- Build: `make build` -- takes 1-2 seconds after initial setup. NEVER CANCEL. Set timeout to 3+ minutes for first run.
- Test: `make test` -- takes 4-5 seconds. NEVER CANCEL. Set timeout to 2+ minutes.
- Lint: `make lint` -- takes ~10 second after setup. NEVER CANCEL. Set timeout to 8+ minutes for first run.
- Format: `/home/runner/go/bin/gofumpt -l -w .` -- takes <1 second (use full path or ensure GOPATH/bin is in PATH).
- Validate OpenAPI: `HOST_PWD=$(pwd) make validate` -- takes 1-2 seconds (requires HOST_PWD environment variable).

### Run the Application
- **ALWAYS build first**: `make build`
- **Local run**: `make run` (starts server on port 8080 with test user)
- **Manual run**: 
  ```bash
  GB_USERS=test@test.com:JDJhJDEwJC9sVWJpTlBYVlZvcU9ZNUxIZmhqYi4vUnRuVkJNaEw4MTQ2VUdFSXRDeE9Ib0ZoVkRLR3pl \
  GB_DISABLEIMPORTERS=true \
  GB_DBPATH=$(pwd)/diary.db \
  GB_ASSETPATH=$(pwd)/diary-assets \
  ./bin/diary server
  ```
- **Docker run**: `docker compose up --build` -- FAILS in sandboxed environments due to Alpine package repo restrictions. Document as "does not work in restricted environments".

### Code Generation (Optional)
- **CRITICAL**: `HOST_PWD=$(pwd) make generate` -- currently has issues with OpenAPI code generation. May fail with model generation errors.
- Only run if modifying `api/openapi.yaml`
- Generated code is already present in `pkg/generated/`

## Validation

### Pre-commit Validation
- **ALWAYS run these before committing**:
  ```bash
  make build
  make test  
  /home/runner/go/bin/gofumpt -l -w .  # or ensure GOPATH/bin is in PATH
  make lint
  ```
- **Expected CI behavior**: CI will fail if lint errors exist or tests fail

## Common Issues and Workarounds

### Configuration
Environment variables (see `make run` for examples):
- **`GB_USERS`**: Base64-encoded user credentials (format: `username:bcrypt_hash`)
- **`GB_DBPATH`**: SQLite database file path
- **`GB_ASSETPATH`**: Directory for user-uploaded assets
- **`GB_DISABLEIMPORTERS`**: Disable background data importers

### Test User Credentials
- Username: `test@test.com`
- Password: `test`
- Encoded: `test@test.com:JDJhJDEwJC9sVWJpTlBYVlZvcU9ZNUxIZmhqYi4vUnRuVkJNaEw4MTQ2VUdFSXRDeE9Ib0ZoVkRLR3pl`

## Development Workflow

### Making Changes
1. **Always build and test existing code first**: `make build && make test`
2. **If build fails with missing generated packages**: Run `git restore pkg/generated/` to restore the working generated code
3. Make your changes
4. **Test incrementally**: `make build && make test` after each change
5. **Manual validation**: Start server and test affected functionality
6. **Lint and format**: Run `make lint`
7. **Final validation**: Complete authentication and web interface tests

### Adding New Features
- Update `api/openapi.yaml` for API changes
- Run `HOST_PWD=$(pwd) make generate` to regenerate client/server code
- Add tests in `test/flows/` for new functionality
- Follow existing patterns in `pkg/server/api/` for handlers
- Update database models in `pkg/database/models/` if needed

### Common File Patterns
After changing API contracts:
- Run `HOST_PWD=$(pwd) make generate` to regenerate client/server code
- Check `pkg/server/api/` for handler implementations
- Update `test/flows/` tests for new functionality

Remember: This codebase uses Ginkgo/Gomega for testing, GORM for database operations, and Gorilla for HTTP routing.
