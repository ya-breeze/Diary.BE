---
type: "always_apply"
---

# Development Guidelines

This document captures important development practices, patterns, and lessons learned from implementing features in the Diary.BE project.

## Database Operations

### Change Tracking
- When implementing database change tracking, always use atomic transactions to ensure both the main operation and change record creation succeed or fail together.

### Mobile Synchronization APIs
- For mobile synchronization APIs, use auto-incrementing IDs for change tracking and include hasMore/nextId fields for efficient pagination.

## Type Safety

### Integer Conversion
- When converting between uint and int32 types, always add bounds checking to prevent integer overflow and use #nosec annotations for intentionally safe operations.

## Task Management

### Feature Breakdown
- Break complex features into granular tasks representing ~20 minutes of work each, with clear acceptance criteria and proper dependency sequencing.

## Testing Strategy

### Test Pyramid
- Test pyramid approach: unit tests for models, integration tests for storage operations, end-to-end tests for complete workflows, and concurrency tests for atomic operations.

## Code Quality

### Linting
- Always run 'make lint' after completing any implementation task and fix all linting issues before considering the work complete.

## Architecture Patterns

### Atomic Operations
- All database operations that involve multiple related changes (e.g., item updates + change tracking) must be wrapped in transactions to maintain data consistency.

### API Design
- Design APIs with pagination from the start using consistent patterns (hasMore, nextId)
- Use OpenAPI-first development to ensure client-server consistency
- Include proper error handling and HTTP status codes

### Change Tracking Implementation
- Store complete item snapshots in change records to enable mobile apps to reconstruct state
- Include metadata fields for debugging and analytics
- Distinguish between operation types (create, update, delete) for different client-side handling

## Development Workflow

### Implementation Process
1. Start with data model design and database schema
2. Implement storage layer with proper indexing
3. Create API endpoints following OpenAPI specification
4. Write comprehensive tests (unit, integration, end-to-end)
5. Run linting and fix all issues
6. Verify all tests pass

### Testing Requirements
- Unit tests for all models and business logic
- Integration tests for database operations
- End-to-end tests for complete user workflows
- Concurrency tests for operations that must be atomic
- Edge case testing (empty states, invalid parameters, error conditions)

### Code Quality Standards
- All linting issues must be resolved before code completion
- Use proper error handling with meaningful messages
- Include security annotations (#nosec) for intentionally safe operations
- Follow Go best practices for type safety and memory management

## Mobile App Considerations

### Synchronization Design
- Implement incremental sync to reduce bandwidth usage
- Design for offline-first mobile applications
- Provide efficient change tracking with proper indexing
- Consider future conflict resolution strategies in the design

### API Efficiency
- Use auto-incrementing IDs for efficient change queries
- Implement proper pagination to handle large datasets
- Include metadata for debugging and analytics
- Design APIs to minimize round trips

## Documentation Requirements

- Maintain OpenAPI specifications as the source of truth for APIs
- Document architectural decisions and their rationale
- Include clear acceptance criteria for all tasks
- Keep development guidelines updated with lessons learned
