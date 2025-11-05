// Package ui implements the terminal user interface for the application.
package ui

import (
	"fmt"

	"github.com/ngundlach/alsa-profile-switcher/pactl"

	tea "github.com/charmbracelet/bubbletea"
)

type UIState struct {
	cards         []pactl.Card
	profileKeys   []string
	cursorPos     int
	err           error
	currentScreen screen
	selectedCard  pactl.Card
	selectedCardPos int
	loading       bool
}

type screen int

const (
	deviceScreen screen = iota
	profileScreen
	errorScreen
)

const (
	up   int = -1
	down int = 1
)

func InitialState() *UIState {
	return &UIState{currentScreen: deviceScreen}
}

func (ui *UIState) Init() tea.Cmd {
	return fetchCardDataCmd
}

func (ui *UIState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return ui.handleKeyMsg(msg.String())
	case dataFetchedMsg:
		ui.loading = false
		ui.cards = msg.cards
		if ui.currentScreen == profileScreen {
			ui.validateProfileScreen()
		}
		return ui, nil
	case errorMsg:
		ui.err = msg.err
		ui.currentScreen = errorScreen
		ui.loading = false
		return ui, nil
	case profileChangedMsg:
		ui.loading = false
		return ui, fetchCardDataCmd
	}
	return ui, nil
}
func (ui *UIState) validateProfileScreen() {
found := false;
for i, card := range ui.cards {
	if card.Index == ui.selectedCard.Index {
		ui.selectedCard = card
		ui.selectedCardPos = i
		found = true
		break
	}
}
if !found {
	ui.err = fmt.Errorf("selected device is no longer available")
	ui.currentScreen = errorScreen
}
}
func (ui *UIState) View() string {
	s := ""
	if ui.loading {
		s += "loading"
	} else {
		switch ui.currentScreen {
		case deviceScreen:
			s += ui.renderDeviceListScreen()
		case profileScreen:
			s += ui.renderProfileListScreen()
		case errorScreen:
			s += ui.renderErrorScreen()
		}
	}
	return s
}

func (ui *UIState) handleSelection() (tea.Model, tea.Cmd) {
	switch ui.currentScreen {
	case deviceScreen:
		ui.selectedCardPos = ui.cursorPos
		ui.selectedCard = ui.cards[ui.selectedCardPos]
		ui.currentScreen = profileScreen
		ui.cursorPos = 0
		ui.profileKeys = pactl.SortKeys(ui.selectedCard.Profiles)
	case profileScreen:
		ui.loading = true
		profileKey := ui.profileKeys[ui.cursorPos]
		return ui, setActiveProfileCmd(ui.selectedCard.Name, profileKey)
	}
	return ui, nil
}

func (ui *UIState) returnToDeviceScreen() {
	switch ui.currentScreen {
	case profileScreen:
		ui.cursorPos = ui.selectedCardPos
		ui.currentScreen = deviceScreen
		ui.profileKeys = nil
	}
}

func (ui *UIState) handleKeyMsg(msg string) (tea.Model, tea.Cmd) {
	if ui.loading {
		if msg == "q" || msg == "ctrl+c" {
			return ui, tea.Quit
		}
		return ui, nil
	}
	switch msg {
	case "q", "ctrl+c":
		return ui, tea.Quit
	case "up", "k":
		ui.cursorPos = ui.changeCursor(up)
		return ui, nil
	case "down", "j":
		ui.cursorPos = ui.changeCursor(down)
		return ui, nil
	case "enter", "right":
		return ui.handleSelection()
	case "backspace", "left":
		ui.returnToDeviceScreen()
		return ui, nil
	case "r":
		ui = &UIState{currentScreen: deviceScreen}
		return ui, fetchCardDataCmd
	}
	return ui, nil
}

func (ui *UIState) changeCursor(direction int) int {
	length := 0
	switch ui.currentScreen {
	case deviceScreen:
		length = len(ui.cards)
	case profileScreen:
		length = len(ui.cards[ui.selectedCardPos].Profiles)
	}
	newCursor := ui.cursorPos + direction
	if newCursor > length-1 || newCursor < 0 {
		return ui.cursorPos
	}
	return ui.cursorPos + direction
}

func (ui *UIState) renderProfileListScreen() string {
	s := "Select profile:\n\n"
	var cursor string
	for i, v := range ui.profileKeys {
		active := " "
		if ui.selectedCard.ActiveProfile == v {
			active = "*"
		}
		if i == ui.cursorPos {
			cursor = ">"
		} else {
			cursor = " "
		}
		s += fmt.Sprintf("[%s]%s %s\n", cursor, active, v)
	}
	s += "\nPress 'q' to quit. 'Backspace' to go back.\n"
	return s
}

func (ui *UIState) renderDeviceListScreen() string {
	s := "Select device:\n\n"
	var cursor string
	for i, v := range ui.cards {
		if i == ui.cursorPos {
			cursor = ">"
		} else {
			cursor = " "
		}
		s += fmt.Sprintf("[%s] %s (%s)\n", cursor, v.Properties.DeviceNick, v.Properties.DeviceProductName)
	}
	s += "\nPress 'q' to quit.\n"
	return s
}

func (ui *UIState) renderErrorScreen() string {
	s := "An error occured:\n\n"
	s += fmt.Sprintf("%v\n", ui.err)
	s += "\nPress 'q' top quit. 'r' to reload.\n"
	return s
}
