package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"
	_ "strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
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

func (m Model) Init() tea.Cmd {
    return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c", "q":
            if m.currentPage != "main" {
                m.currentPage = "main"
                m.showLink = false
                return m, nil
            }
            return m, tea.Quit

        case "up", "k":
            if m.cursor > 0 {
                m.cursor--
            }

        case "down", "j":
            if m.currentPage == "projects" {
                if m.cursor < len(m.projects)-1 {
                    m.cursor++
            }
        } else {
            if m.cursor < len(m.choices)-1 {
                m.cursor++
            }
        }

    case "enter", " ":
        if m.currentPage == "main" {
            m.currentPage = strings.ToLower(strings.Replace(m.choices[m.cursor], " ", "", -1))
            m.cursor = 0
        } else if m.currentPage == "projects" {
            m.showLink = !m.showLink
        }
    }
}
    return m, nil
}

func (m Model) View() string {
    var s strings.Builder

    s.WriteString(headerStyle.Render("Connected with key fingerprint: \n"))
    s.WriteString(fmt.Sprintf("%s\n\n", m.userFingerprint))

    switch m.currentPage {
    case "main":
        s.WriteString(headerStyle.Render("Welcome to Martin's Portfolio!\n\n"))
        s.WriteString("Use ↑/↓ arrows to navigate, or j/k, enter to select, q to quit\n\n")

        for i, choice := range m.choices {
            cursor := " "
            if m.cursor == i {
                cursor = ">"
            }
            s.WriteString(fmt.Sprintf("%s %s\n", cursor, choice))
        }
    case "projects":
        s.WriteString(headerStyle.Render("=== Projects ===\n\n"))
        s.WriteString("Use↑/↓ to navigate, enter to show/hide links, q to go back\n\n")

        for i, project := range m.projects {
            cursor := " "
            if m.cursor == i {
                cursor = ">"
            }
            s.WriteString(fmt.Sprintf("%s %s\n", cursor, project.name))
            s.WriteString(fmt.Sprintf("   Description: %s\n", project.description))
            s.WriteString(fmt.Sprintf("   Tech Stack: %s\n", project.techStack))
            if m.showLink && m.cursor == i {
                s.WriteString(fmt.Sprintf("   Link: %s\n", linkStyle.Render(project.link)))
            }
            s.WriteString("\n")
        }
    
    case "aboutme":
        s.WriteString(headerStyle.Render("=== About Me ===\n"))
        s.WriteString("Software Developer etc.")
        s.WriteString("\nPress q to go back")

    case "skills":
        s.WriteString(headerStyle.Render("=== Skills ===\n"))
        s.WriteString("Languages: Python, Golang, Javascript, SQL etc\n")
        s.WriteString("Frameworks: Flask, Django\n")
        s.WriteString("Tools: Vim, Docker, Linux, Azure\n")
        s.WriteString("\nPress q to go back")

    case "contact":
        s.WriteString(headerStyle.Render("=== Contact ===\n"))
        s.WriteString("Email: martin.evenbom@gmail.com\n")
        s.WriteString("Github: github.com/zibiax\n")
        s.WriteString("\nPress q to go back")
    }

    return s.String()
}



