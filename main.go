package main

import (
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	spaceshipWidth  = 5
	spaceshipHeight = 3
)

func main() {
	rand.Seed(time.Now().UnixNano())
	p := tea.NewProgram(model{
		pos:     0,
		columns: 0,
		rows:    0,
		points:  0,
		bullets: []bullet{},
	}, tea.WithAltScreen(), tea.WithMouseAllMotion())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type bullet struct {
	xPos int
	yPos int
}

type enemy struct {
	xPos int
	yPos int
}

type model struct {
	pos     int
	columns int
	rows    int
	points  int
	bullets []bullet
	enemies []enemy
	beat    int
}

func (m model) Init() tea.Cmd {
	return tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tickMsg:
		m.beat++
		if m.beat > 100 {
			m.beat = 0
		}

		if m.beat%5 == 0 {
			newBullets := make([]bullet, len(m.bullets))
			offScreenCount := 0
			for i, b := range m.bullets {
				b.yPos--
				newBullets[i] = b

				if b.yPos <= 0 {
					offScreenCount++
				}
			}

			newBullets = newBullets[offScreenCount:]

			veryNewBullets := []bullet{}
			seen := make(map[bullet]struct{})
			for _, b := range newBullets {
				if _, ok := seen[b]; ok {
					continue
				}
				seen[b] = struct{}{}
				veryNewBullets = append(veryNewBullets, b)
			}

			m.bullets = veryNewBullets

			veryNewEnemies := []enemy{}
			for _, e := range m.enemies {
				hit := false
				for _, b := range m.bullets {
					if b.xPos == e.xPos && b.yPos == e.yPos || (b.xPos == e.xPos && b.yPos == e.yPos-1) {
						m.points++
						hit = true
						break
					}
				}
				if !hit {
					veryNewEnemies = append(veryNewEnemies, e)
				}
			}

			m.enemies = veryNewEnemies
		}

		if m.beat%50 == 0 {
			randomNumber := rand.Intn(100)
			if randomNumber > 90 {
				e := enemy{xPos: rand.Intn(m.columns), yPos: 0}
				if e.xPos > m.columns-spaceshipWidth {
					e.xPos = m.columns - spaceshipWidth
				}
				if e.xPos < spaceshipWidth {
					e.xPos = spaceshipWidth
				}
				m.enemies = append(m.enemies, e)
			}

			newEnemies := make([]enemy, len(m.enemies))
			offScreenCount := 0
			for i, e := range m.enemies {
				e.yPos++
				newEnemies[i] = e

				if e.yPos == m.rows {
					offScreenCount++
				}
			}

			m.enemies = newEnemies[offScreenCount:]

		}

		return m, tick

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeySpace:
			m.bullets = append(m.bullets, bullet{xPos: m.pos + 2, yPos: m.rows - spaceshipHeight - 1})

		case tea.KeyLeft:
			if m.pos > 1 {
				m.pos--
			}

		case tea.KeyRight:
			if m.pos < m.columns-spaceshipWidth {
				m.pos++
			}

		case tea.KeyRunes:
			switch string(msg.Runes) {
			case "q":
				return m, tea.Quit
			}
		}

	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseMotion:
			if msg.X > m.pos {
				m.pos++
			} else if msg.X < m.pos {
				m.pos--
			}
		case tea.MouseLeft:
			m.bullets = append(m.bullets, bullet{xPos: m.pos + 2, yPos: m.rows - spaceshipHeight - 1})
		}

	case tea.WindowSizeMsg:
		m.columns = msg.Width
		m.rows = msg.Height
	}

	return m, nil
}

func (m model) View() string {
	if m.rows == 0 {
		return ""
	}

	board := make([][]string, m.rows)

	// draw empty board
	for i := 0; i < m.rows; i++ {
		line := strings.Split(strings.Repeat(" ", m.columns), "")
		board[i] = line
	}

	// draw bullets
	for _, b := range m.bullets {
		board[b.yPos][b.xPos] = "*"
	}

	// draw enemies
	for _, e := range m.enemies {
		board[e.yPos][e.xPos] = "▓"
	}

	board[1][1] = strconv.Itoa(m.points)

	// draw spaceship
	board[m.rows-3][m.pos+2] = "_"
	board[m.rows-2][m.pos+1] = "/"
	board[m.rows-2][m.pos+3] = "\\"
	board[m.rows-1][m.pos] = "/"
	board[m.rows-1][m.pos+1] = "_"
	board[m.rows-1][m.pos+2] = "_"
	board[m.rows-1][m.pos+3] = "_"
	board[m.rows-1][m.pos+4] = "\\"

	var boardString string
	for i, line := range board {
		boardString += strings.Join(line, "")
		if i < len(board)-1 {
			boardString += "\n"
		}
	}

	return boardString
}

// Messages are events that we respond to in our Update function. This
// particular one indicates that the timer has ticked.
type tickMsg time.Time

func tick() tea.Msg {
	time.Sleep(time.Millisecond * 10)
	return tickMsg{}
}
