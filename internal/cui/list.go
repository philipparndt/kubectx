package cui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"log"
	"sort"
	"strings"
)

type LabelProvider[T any] interface {
	Label(T) string
}

func Select[T any](items []T, labelProvider LabelProvider[T]) *T {
	// Initialize screen
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("Error creating screen: %v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("Error initializing screen: %v", err)
	}
	defer screen.Fini()

	// Sort by label
	sorted := make([]T, len(items))
	copy(sorted, items)
	sort.Slice(sorted, func(i, j int) bool {
		return labelProvider.Label(sorted[i]) < labelProvider.Label(sorted[j])
	})

	// Customize styles for fzf-like appearance
	normalStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite)
	highlightStyle := tcell.StyleDefault.Background(tcell.ColorLightCyan).Foreground(tcell.ColorBlack)

	// Original list of options
	query := ""
	filteredItems := items
	selectedIndex := 0
	pageSize := 5 // Number of items to jump with PageUp/PageDown

	// Event handling loop
	for {
		screen.Clear()

		// Display the filter query
		printString(screen, 1, 0, fmt.Sprintf("> %s", query), normalStyle)

		// Filter the list based on the query
		filteredItems = filterItems(sorted, labelProvider, query)

		// Draw the filtered list
		for i, item := range filteredItems {
			style := normalStyle
			if i == selectedIndex {
				// Highlight the selected item
				style = highlightStyle
			}
			printString(screen, 1, i+2, labelProvider.Label(item), style) // +2 to leave space for the filter input
		}

		screen.Show()

		// Wait for a key event
		ev := screen.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyRune:
				// Handle `j` for down and `k` for up
				switch ev.Rune() {
				case 'j':
					if selectedIndex < len(filteredItems)-1 {
						selectedIndex++
					}
				case 'k':
					if selectedIndex > 0 {
						selectedIndex--
					}
				default:
					// Add typed characters to the query
					query += string(ev.Rune())
					selectedIndex = 0 // Reset selection after filtering
				}
			case tcell.KeyUp:
				// Move up in the filtered list
				if selectedIndex > 0 {
					selectedIndex--
				}
			case tcell.KeyDown:
				// Move down in the filtered list
				if selectedIndex < len(filteredItems)-1 {
					selectedIndex++
				}
			case tcell.KeyEnter:
				// Item selected
				if len(filteredItems) > 0 {
					screen.Clear()
					screen.Show()
					screen.Fini()
					return &filteredItems[selectedIndex]
				}
			case tcell.KeyPgUp:
				// Page up: move up by pageSize
				if selectedIndex-pageSize >= 0 {
					selectedIndex -= pageSize
				} else {
					selectedIndex = 0
				}
			case tcell.KeyPgDn:
				// Page down: move down by pageSize
				selectedIndex += pageSize
				if selectedIndex >= len(filteredItems) {
					selectedIndex = len(filteredItems) - 1
				}
			case tcell.KeyBackspace, tcell.KeyBackspace2:
				// Remove last character from the query
				if len(query) > 0 {
					query = query[:len(query)-1]
					selectedIndex = 0 // Reset selection after filtering
				}
			case tcell.KeyEscape, tcell.KeyCtrlC:
				// Exit
				screen.Fini()
				return nil
			default:
				// Add typed characters to the query
				if ev.Rune() != 0 {
					query += string(ev.Rune())
					selectedIndex = 0 // Reset selection after filtering
				}
			}
		}
	}
}

// printString is a helper function to display a string on the screen
func printString(s tcell.Screen, x, y int, str string, style tcell.Style) {
	for i, c := range str {
		s.SetContent(x+i, y, c, nil, style)
	}
}

// filterItems filters the items list by checking if each item contains the query as a substring
func filterItems[T any](items []T, labelProvider LabelProvider[T], query string) []T {
	if query == "" {
		return items
	}
	var filtered []T
	for _, item := range items {
		if strings.Contains(strings.ToLower(labelProvider.Label(item)), strings.ToLower(query)) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}
