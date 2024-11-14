package main

import (
	_"crypto/sha256"
	_"encoding/base64"
	_"fmt"
	_"log"
	_"os"
	_"strings"
	_"time"

	// tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	_"github.com/charmbracelet/wish"
	// bm "github.com/charmbracelet/wish/bubbletea"
	_"github.com/gliderlabs/ssh"
)

type Project struct {
    name string
    description string
    techStack string
    link string
}

type Model struct {
    currentPage string
    choices []string
    cursor int
    projects []Project
    showLink bool
    userFingerprint string
}

var (
    linkStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("86")).Underline(true)

    headerStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
)

