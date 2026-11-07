package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) renderStorageModal() string {
	var title, subtitle, placeholder string
	switch m.StorageModalType {
	case ModalMkdir:
		title = "Create New Directory"
		subtitle = "Folder will be created inside current directory"
		placeholder = "Enter folder name..."
	case ModalRename:
		title = "Rename Item"
		subtitle = "Enter new filename or directory name"
		placeholder = "Enter new name..."
	case ModalPush:
		title = "Upload File from PC to Phone"
		subtitle = "Enter absolute host file path or drag and drop here"
		placeholder = "/Users/username/file.mp4"
	default:
		return ""
	}

	displayInput := m.StorageModalInput
	if displayInput == "" {
		displayInput = StyleHelp.Render(placeholder)
	}

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorPrimary).
		Padding(1, 2).
		Width(64).
		Render(fmt.Sprintf(
			"%s\n%s\n\n%s\n\n%s",
			lipgloss.NewStyle().Bold(true).Foreground(ColorPrimary).Render(title),
			StyleHelp.Render(subtitle),
			lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(ColorSecondary).
				Padding(0, 1).
				Width(58).
				Render(displayInput),
			StyleHelp.Render("[Enter] Confirm   [Esc] Cancel"),
		))

	return lipgloss.Place(m.Width, m.Height-4, lipgloss.Center, lipgloss.Center, box)
}

func handleStorageModalKey(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		m.StorageModalOpen = false
		m.StorageModalInput = ""
		m.StorageModalType = ModalNone
		return m, nil

	case tea.KeyBackspace:
		if len(m.StorageModalInput) > 0 {
			m.StorageModalInput = m.StorageModalInput[:len(m.StorageModalInput)-1]
		}
		return m, nil

	case tea.KeyEnter:
		input := strings.TrimSpace(m.StorageModalInput)
		m.StorageModalOpen = false
		m.StorageModalInput = ""
		action := m.StorageModalType
		m.StorageModalType = ModalNone

		if input == "" {
			return m, nil
		}

		currentDir := m.StoragePath
		if currentDir == "" {
			currentDir = "/sdcard"
		}

		switch action {
		case ModalMkdir:
			m.Status = fmt.Sprintf("Creating directory %s...", input)
			return m, tea.Sequence(
				MkdirStorageItemCmd(m.Executor, currentDir, input),
				LoadStorageCmd(m.Executor, currentDir),
			)

		case ModalRename:
			if m.Storage == nil || len(m.Storage.Items) == 0 || m.StorageCursor >= len(m.Storage.Items) {
				return m, nil
			}
			targetItem := m.Storage.Items[m.StorageCursor]
			m.Status = fmt.Sprintf("Renaming to %s...", input)
			return m, tea.Sequence(
				RenameStorageItemCmd(m.Executor, targetItem.Path, input),
				LoadStorageCmd(m.Executor, currentDir),
			)

		case ModalPush:
			m.Status = fmt.Sprintf("Uploading %s to device...", input)
			return m, tea.Sequence(
				PushStorageItemCmd(m.Executor, input, currentDir),
				LoadStorageCmd(m.Executor, currentDir),
			)
		}
		return m, nil

	case tea.KeyRunes:
		m.StorageModalInput += string(msg.Runes)
		return m, nil

	default:
		return m, nil
	}
}
