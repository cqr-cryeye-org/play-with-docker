package config

import (
	"flag"
	"github.com/play-with-docker/play-with-docker/pwd/types"
	"os"
	"regexp"
	"time"

	"github.com/gorilla/securecookie"

	"golang.org/x/oauth2"
)

const (
	PWDHostnameRegex      = "[0-9]{1,3}-[0-9]{1,3}-[0-9]{1,3}-[0-9]{1,3}"
	PortRegex             = "[0-9]{1,5}"
	AliasnameRegex        = "[0-9|a-z|A-Z|-]*"
	AliasSessionRegex     = "[0-9|a-z|A-Z]{8}"
	AliasGroupRegex       = "(" + AliasnameRegex + ")-(" + AliasSessionRegex + ")"
	PWDHostPortGroupRegex = "^.*ip(" + PWDHostnameRegex + ")(?:-?(" + PortRegex + "))?(?:\\..*)?$"
	AliasPortGroupRegex   = "^.*pwd" + AliasGroupRegex + "(?:-?(" + PortRegex + "))?\\..*$"
)

var NameFilter = regexp.MustCompile(PWDHostPortGroupRegex)
var AliasFilter = regexp.MustCompile(AliasPortGroupRegex)

var (
	PWDContainerName         string
	L2ContainerName          string
	SessionsFile             string
	L2Subdomain              string
	HashKey                  string
	SSHKeyPath               string
	L2RouterIP               string
	CookieHashKey            string
	CookieBlockKey           string
	PortNumber               string
	LetsEncryptCertsDir      string
	AdminToken               string
	PlaygroundDomain         string
	SegmentId                string
	DefaultDinDInstanceImage string
	CryeyeClientID           string
	CryeyeClientSecret       string
)
var (
	ExternalDindVolume bool
	NoWindows          bool
	UseLetsEncrypt     bool
	ForceTLS           bool
	// Unsafe enables a number of unsafe features when set. It is principally
	// intended to be used in development. For example, it allows the caller to
	// specify the Docker networks to join.
	Unsafe bool
)

var MaxLoadAvg float64
var SecureCookie *securecookie.SecureCookie

// Providers TODO move this to a sync map so it can be updated on demand when the configuration for a playground changes
var Providers = map[string]map[string]*oauth2.Config{}

var LocalPlayground types.Playground

func ParseFlags() {
	flag.StringVar(&LetsEncryptCertsDir, "letsencrypt-certs-dir", "/certs", "Path where let's encrypt certs will be stored")
	flag.BoolVar(&UseLetsEncrypt, "letsencrypt-enable", false, "Enabled let's encrypt tls certificates")
	flag.BoolVar(&ForceTLS, "tls", false, "Use TLS to connect to docker daemons")
	flag.StringVar(&PortNumber, "port", "3000", "Port number")
	flag.StringVar(&SessionsFile, "save", "./pwd/sessions", "Tell where to store sessions file")
	flag.StringVar(&PWDContainerName, "name", "pwd", "Container name used to run PWD (used to be able to connect it to the networks it creates)")
	flag.StringVar(&L2ContainerName, "l2", "l2", "Container name used to run L2 Router")
	flag.StringVar(&L2RouterIP, "l2-ip", "", "Host IP address for L2 router ping response")
	flag.StringVar(&L2Subdomain, "l2-subdomain", "direct", "Subdomain to the L2 Router")
	flag.StringVar(&HashKey, "hash_key", "salmonrosado", "Hash key to use for cookies")
	flag.BoolVar(&NoWindows, "win-disable", false, "Disable windows instances")
	flag.BoolVar(&ExternalDindVolume, "dind-external-volume", false, "Use external dind volume though XFS volume driver")
	flag.Float64Var(&MaxLoadAvg, "max-load", 100, "Maximum allowed load average before failing ping requests")
	flag.StringVar(&SSHKeyPath, "ssh_key_path", "", "SSH Private Key to use")
	flag.StringVar(&CookieHashKey, "cookie-hash-key", os.Getenv("COOKIE_HASH_KEY"), "Hash key to use to validate cookies")
	flag.StringVar(&CookieBlockKey, "cookie-block-key", os.Getenv("COOKIE_BLOCK_KEY"), "Block key to use to encrypt cookies")
	flag.StringVar(&PlaygroundDomain, "playground-domain", os.Getenv("PLAYGROUND_DOMAIN"), "Domain to use for the playground")
	flag.StringVar(&AdminToken, "admin-token", "", "Token to validate admin user for admin endpoints")
	flag.StringVar(&SegmentId, "segment-id", "", "Segment id to post metrics")
	flag.StringVar(&CryeyeClientID, "cryeye-client-id", os.Getenv("CRYEYE_CLIENT_ID"), "Client ID for oauth authorization")
	flag.StringVar(&CryeyeClientSecret, "cryeye-client-secret", os.Getenv("CRYEYE_CLIENT_SECRET"), "Client secret for oauth authorization")
	flag.BoolVar(&Unsafe, "unsafe", os.Getenv("PWD_UNSAFE") == "true", "Operate in unsafe mode")
	flag.Parse()

	SecureCookie = securecookie.New([]byte(CookieHashKey), []byte(CookieBlockKey))
	DefaultDinDInstanceImage = "cryeye/dind"

	LocalPlayground = types.Playground{
		Domain:                      "localhost",
		DefaultDinDInstanceImage:    DefaultDinDInstanceImage,
		AvailableDinDInstanceImages: []string{DefaultDinDInstanceImage},
		AllowWindowsInstances:       NoWindows,
		DefaultSessionDuration:      time.Hour * 24,
		Tasks:                       []string{".*"},
		Extras:                      map[string]interface{}{"LoginRedirect": "http://localhost/"},
		Privileged:                  true,
	}
}
