package ui

import (
	"github.com/ngundlach/alsa-profile-switcher/pactl"

	tea "github.com/charmbracelet/bubbletea"
)

func fetchCardDataCmd() tea.Msg {
	cards, err := pactl.FetchDeviceData()
	if err != nil {

		return errorMsg{err: err}
	}
	return dataFetchedMsg{cards}
}
func setActiveProfileCmd(card string, profile string) tea.Cmd {
	return func() tea.Msg {
		err := pactl.SetActiveProfile(card, profile)
		if err != nil {
			return errorMsg{err: err}
		} else {
			return profileChangedMsg{}
		}
	}
}
