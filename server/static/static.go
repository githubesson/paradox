package static

import (
	"path/filepath"

	"paradox_server/models"
)

var SupportedBrowsers = []models.BrowserConfig{
	{Name: "Chrome", PathString: filepath.Join("Browsers", "Chrome"), PrintName: "b'Chrome Safe Storage'"},
	{Name: "Edge", PathString: filepath.Join("Browsers", "Edge"), PrintName: "b'Edge Safe Storage'"},
	{Name: "Brave", PathString: filepath.Join("Browsers", "Brave"), PrintName: "b'Brave Safe Storage'"},
	{Name: "Vivaldi", PathString: filepath.Join("Browsers", "Vivaldi"), PrintName: "b'Vivaldi Safe Storage'"},
	{Name: "Opera", PathString: filepath.Join("Browsers", "Opera"), PrintName: "b'Opera Safe Storage'"},
}
