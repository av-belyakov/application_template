package datamodels

// TemplBasePage тип для templ, базовая страница
type TemplBasePage struct {
	Title        string
	AppName      string
	AppVersion   string
	AppShortInfo string
	MenuLinks    []struct {
		Name string
		Link string
		Icon string
	}
}
