package webserver

import (
	"net/http"

	"github.com/av-belyakov/application_template/components"
	"github.com/av-belyakov/application_template/internal/memorystatistics"
)

// RouteMemoryStatistics статистика использования памяти
func (is *InformationServer) RouteMemoryStatistics(w http.ResponseWriter, r *http.Request) {
	is.getBasePage(
		components.TemplateMemoryStats(memorystatistics.GetMemoryStats()),
		components.BaseComponentScripts(),
	).Component.Render(r.Context(), w)
}
