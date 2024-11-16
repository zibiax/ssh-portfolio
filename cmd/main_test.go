package main

import (
	"strings"
	"testing"
	tea "github.com/charmbracelet/bubbletea"
    "os"
    "time"
)

// TestInitialModel tests the initial model creation
func TestInitialModel(t *testing.T) {
	fingerprint := "test-fingerprint"
	model := initialModel(fingerprint)

	// Test initial values
	if model.currentPage != "main" {
		t.Errorf("Expected initial page to be 'main', got %s", model.currentPage)
	}

	if model.cursor != 0 {
		t.Errorf("Expected initial cursor to be 0, got %d", model.cursor)
	}

	if len(model.choices) != 4 {
		t.Errorf("Expected 4 choices, got %d", len(model.choices))
	}

	if len(model.projects) != 4 {
		t.Errorf("Expected 4 projects, got %d", len(model.projects))
	}

	if model.showLink != false {
		t.Errorf("Expected showLink to be false")
	}

	if model.userFingerprint != fingerprint {
		t.Errorf("Expected fingerprint to be %s, got %s", fingerprint, model.userFingerprint)
	}
}

// TestGetKeyFingerPrint tests the key fingerprint generation
func TestGetKeyFingerPrint(t *testing.T) {
	// Test with nil key
	nilFingerprint := getKeyFingerPrint(nil)
	if nilFingerprint != "unknown-key" {
		t.Errorf("Expected 'unknown-key' for nil key, got %s", nilFingerprint)
	}

	// Test with non-nil case using a fixed test string
	testFingerprint := "test-fingerprint"
	model := initialModel(testFingerprint)
	if model.userFingerprint != testFingerprint {
		t.Errorf("Expected fingerprint %s, got %s", testFingerprint, model.userFingerprint)
	}
}

// TestModelUpdate tests the model's update function
func TestModelUpdate(t *testing.T) {
	model := initialModel("test")

	// Test navigation
	testCases := []struct {
		name     string
		msg      tea.Msg
		wantPage string
		wantCursor int
	}{
		{
			name: "Move cursor down",
			msg: tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")},
			wantPage: "main",
			wantCursor: 1,
		},
		{
			name: "Move cursor up",
			msg: tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")},
			wantPage: "main",
			wantCursor: 0,
		},
		{
			name: "Select menu item",
			msg: tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(" ")},
			wantPage: "aboutme",
			wantCursor: 0,
		},
		{
			name: "Return to main",
			msg: tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")},
			wantPage: "main",
			wantCursor: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			newModel, _ := model.Update(tc.msg)
			m := newModel.(Model)
			
			if m.currentPage != tc.wantPage {
				t.Errorf("Expected page %s, got %s", tc.wantPage, m.currentPage)
			}
			
			if m.cursor != tc.wantCursor {
				t.Errorf("Expected cursor %d, got %d", tc.wantCursor, m.cursor)
			}
		})
	}
}

// TestModelView tests the view rendering
func TestModelView(t *testing.T) {
	model := initialModel("test-fingerprint")

	// Test main page view
	view := model.View()
	if !strings.Contains(view, "Welcome to Martin's Portfolio!") {
		t.Error("Main page view should contain welcome message")
	}

	// Test projects page view
	model.currentPage = "projects"
	view = model.View()
	if !strings.Contains(view, "=== Projects ===") {
		t.Error("Projects page view should contain projects header")
	}
	
	// Verify all projects are rendered
	for _, project := range model.projects {
		if !strings.Contains(view, project.name) {
			t.Errorf("Projects view should contain project %s", project.name)
		}
	}

	// Test about page view
	model.currentPage = "aboutme"
	view = model.View()
	if !strings.Contains(view, "=== About Me ===") {
		t.Error("About page view should contain about header")
	}

	// Test skills page view
	model.currentPage = "skills"
	view = model.View()
	if !strings.Contains(view, "=== Skills ===") {
		t.Error("Skills page view should contain skills header")
	}

	// Test contact page view
	model.currentPage = "contact"
	view = model.View()
	if !strings.Contains(view, "=== Contact ===") {
		t.Error("Contact page view should contain contact header")
	}
}

// TestLogAccess tests the logging functionality
func TestLogAccess(t *testing.T) {
    // Save the current working directory
    originalCwd, err := os.Getwd()
    if err != nil {
        t.Fatalf("Failed to get current working directory: %v", err)
    }

    // Create temporary directory
    tmpDir, err := os.MkdirTemp("", "test-logs")
    if err != nil {
        t.Fatalf("Failed to create temp directory: %v", err)
    }
    defer os.RemoveAll(tmpDir) // Clean up after test
    
    // Change to temp directory
    if err := os.Chdir(tmpDir); err != nil {
        t.Fatalf("Failed to change directory: %v", err)
    }
    // Ensure we change back to original directory after test
    defer os.Chdir(originalCwd)

    // Test logging
    fingerprint := "test-fingerprint"
    status := "connected"
    logAccess(fingerprint, status)

    // Read log file (it will be created as "access.log" in temp directory)
    content, err := os.ReadFile("access.log")
    if err != nil {
        t.Fatalf("Failed to read log file: %v", err)
    }

    logContent := string(content)
    if !strings.Contains(logContent, fingerprint) {
        t.Error("Log should contain fingerprint")
    }
    if !strings.Contains(logContent, status) {
        t.Error("Log should contain status")
    }

    // Verify timestamp format
    if !strings.Contains(logContent, time.Now().Format(time.RFC3339)[:10]) {
        t.Error("Log should contain today's date")
    }
}
