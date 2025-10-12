package ui

import (
	"fmt"

	"github.com/ngundlach/alsa-profile-switcher/pactl"

	tea "github.com/charmbracelet/bubbletea"
)

const up int = -1
const down int = 1

type UiState struct {
	cards         []pactl.Card
	profileKeys   []string
	cursorPos     int
	err           error
	currentScreen screen
	selectedCard  int
	loading       bool
}

type screen int

const (
	deviceScreen screen = iota
	profileScreen
	errorScreen
)

func (ui *UiState) Init() tea.Cmd {
	return fetchCardDataCmd
}

func (ui *UiState) handleSelection() (tea.Model, tea.Cmd) {
	switch ui.currentScreen {
	case deviceScreen:
		ui.selectedCard = ui.cursorPos
		ui.currentScreen = profileScreen
		ui.cursorPos = 0
		ui.profileKeys = pactl.SortKeys(ui.cards[ui.selectedCard].Profiles)
	case profileScreen:
		ui.loading = true
		profileKey := ui.profileKeys[ui.cursorPos]
		return ui, setActiveProfileCmd(ui.cards[ui.selectedCard].Name, profileKey)
	}
	return ui, nil
}
func (ui *UiState) returnToDeviceScreen() {
	switch ui.currentScreen {
	case profileScreen:
		ui.cursorPos = ui.selectedCard
		ui.currentScreen = deviceScreen
		ui.profileKeys = nil
	}
}
func (ui *UiState) handleKeyMsg(msg string) (tea.Model, tea.Cmd) {
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
		ui = &UiState{currentScreen: deviceScreen}
		return ui, fetchCardDataCmd
	}
	return ui, nil
}

func (ui *UiState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return ui.handleKeyMsg(msg.String())
	case dataFetchedMsg:
		ui.loading = false
		ui.cards = msg.cards
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

func (ui *UiState) changeCursor(change int) int {
	length := 0
	switch ui.currentScreen {
	case deviceScreen:
		length = len(ui.cards)
	case profileScreen:
		length = len(ui.cards[ui.selectedCard].Profiles)
	}
	newCursor := ui.cursorPos + change
	if newCursor > length-1 || newCursor < 0 {
		return ui.cursorPos
	}
	return ui.cursorPos + change
}

func (ui *UiState) renderProfileListScreen() string {
	s := "Select profile:\n\n"
	cursor := " "
	for i, v := range ui.profileKeys {
		active := " "
		if ui.cards[ui.selectedCard].ActiveProfile == v {
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

func (ui *UiState) renderDeviceListScreen() string {
	s := "Select device:\n\n"
	cursor := " "
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

func (ui *UiState) renderErrorScreen() string {
	s := "An error occured:\n\n"
	s += fmt.Sprintf("%v\n", ui.err)
	s += "\nPress 'q' top quit. 'r' to reload.\n"
	return s
}
func (ui *UiState) View() string {
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

func InitialState() *UiState {
	return &UiState{currentScreen: deviceScreen}
}
