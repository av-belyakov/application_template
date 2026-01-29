package webserver

import (
	"net/http"

	"github.com/av-belyakov/application_template/components"
)

func (is *InformationServer) RouteLogs(w http.ResponseWriter, r *http.Request) {
	is.getBasePage(
		components.TemplateLogs(is.storage.GetLogs(), is.storage.LogMaxSize()),
		components.BaseComponentScripts(),
	).Component.Render(r.Context(), w)
}
