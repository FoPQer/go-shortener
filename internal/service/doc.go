// Package service contains application-level business logic for the URL shortener.
//
// The package orchestrates repository interactions and provides use-case oriented
// services used by HTTP handlers.
//
// Main responsibilities:
//   - URL lifecycle operations (create, resolve, batch create, user-scoped delete)
//   - JSON input/output transformations for API payloads
//   - User retrieval and creation operations
//   - Runtime configuration access for service-level helpers
//
// Service methods are designed to remain transport-agnostic and operate on domain
// models from internal/model and repository contracts from internal/repository.
package service
