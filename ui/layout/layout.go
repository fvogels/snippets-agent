package layout

import (
	"code-snippets/util"

	"github.com/charmbracelet/lipgloss"
)

type Layout interface {
	Render() string
	UpdateSize(size util.Size)
}

type VerticalPart struct {
	Height int
	Render func(size util.Size) string
}

type HorizontalPart struct {
	Width  int
	Render func(size util.Size) string
}

func Vertical(totalSize util.Size, parts ...VerticalPart) string {
	renderedParts := []string{}

	for _, part := range parts {
		style := lipgloss.NewStyle().Width(totalSize.Width).Height(part.Height)
		subsize := util.Size{Width: totalSize.Width, Height: part.Height}
		renderedPart := style.Render(part.Render(subsize))
		renderedParts = append(renderedParts, renderedPart)
	}

	return lipgloss.JoinVertical(0, renderedParts...)
}

func Horizontal(totalSize util.Size, parts ...HorizontalPart) string {
	renderedParts := []string{}

	for _, part := range parts {
		style := lipgloss.NewStyle().Width(part.Width).Height(totalSize.Height)
		subsize := util.Size{Width: part.Width, Height: totalSize.Height}
		renderedPart := style.Render(part.Render(subsize))
		renderedParts = append(renderedParts, renderedPart)
	}

	return lipgloss.JoinHorizontal(0, renderedParts...)
}
