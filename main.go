package main

import (
	"crypto/sha256"
	"encoding/base64"
    "fmt"
    "log"
	"os"
   _ "strings" 
    "time"

	// tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	_ "github.com/charmbracelet/wish"

	// bm "github.com/charmbracelet/wish/bubbletea"
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

func initialModel(fingerprint string) Model {
    return Model{
        currentPage: "main",

        choices: []string{"About Me", "Projects", "Skills", "Contact"},

        cursor: 0,

        projects: []Project{
            {
            name: "Cliphive",
            description: "Working on a Pastebin website that uses golang as it's backend",
            techStack: "Go, and a little JS",
            link: "https://github.com/zibiax/cliphive",
        },
        {
            name: "Flask Portfolio website",
            description: "This is my portoflio website that uses python-flask, the url to it is https://evenbom.se",
            techStack: "Python-flask, Sqlite, Docker",
            link: "https://github.com/zibiax/flask-portfolio",
        },
        {
            name: "Draugen",
            description: "A group project from when I studied VR-developer, this is just a unity game thats quite simple. It also supports multiplayer",
            techStack: "Unity, C#",
            link: "https://github.com/brinobre/Draugen",
        },
        {
            name: "DraftGPT",
            description: "A project that takes advantage of openAI api, to diagnose error outputs from servers",
            techStack: "Python",
            link: "https://github.com/Sanjivani-S/DraftGPT",
        },
    },
    showLink: false,
    userFingerprint: fingerprint,
    }
}

func getKeyFingerPrint(key ssh.PublicKey) string {
    hash := sha256.Sum256(key.Marshal())
    return base64.StdEncoding.EncodeToString(hash[:])
}

func logAccess(fingerprint, status string) {
    f, err := os.OpenFile("access.log", 
        os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    
    if err != nil {
        log.Printf("Error opening log file: %v", err)
        return
    }

    defer f.Close()

    timestamp := time.Now().Format(time.RFC3339)

    logEntry := fmt.Sprintf("%s - Fingerprint: %s - Status: %s\n",
        timestamp, fingerprint, status)

    if _, err := f.WriteString(logEntry); err != nil {
        log.Printf("Error writing to log: %v", err)
    }
}


