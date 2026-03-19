package tests

import (
	"testing"

	"GoClaw/RAG/config"
	"GoClaw/RAG/repository"
)

// Mock configuration for testing
func mockConfig() *config.Configuration {
	return &config.Configuration{
		OpenAIAPIKey:    "test-key", // This will cause errors in real API calls, but we'll mock the embedding function in tests
		EmbeddingModel:  "text-embedding-ada-002",
		ChatModel:       "gpt-3.5-turbo",
		ShortTermWindow: 5,
		TopK:            3,
		UseWeaviate:     false,
		WeaviateURL:     "http://localhost:8080",
		WeaviateAPIKey:  "",
	}
}

// Test InMemoryVectorStore basic functionality
func TestInMemoryVectorStore(t *testing.T) {
	// Create a mock config
	cfg := mockConfig()

	// Create store
	store := repository.NewInMemoryVectorStore()

	// Test AddDocument (this will fail because we don't have a real API key, but we can check the error)
	err := store.AddDocument("test1", "This is a test document", cfg)
	if err == nil {
		t.Error("Expected error when adding document without valid OpenAI API key, got nil")
	}

	// Test Retrieve with empty store
	results := store.Retrieve("test query", cfg)
	if results != nil && len(results) > 0 {
		t.Error("Expected empty results from empty store, got:", results)
	}

	// Test that store implements VectorStore interface
	var _ repository.VectorStore = store
}
