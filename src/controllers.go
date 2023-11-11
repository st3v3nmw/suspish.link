package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

var SEPARATORS = []string{
	".", "-", "+", "~",
}

var ESCAPABLE_SEPARATORS = []string{
	" ", "<", ">", "#", "%", "*", "/", "{", "}", "|", "$",
	"\\", "^", "[", "]", "'", ";", "?", ":", "@", "=", "&",
}

var VERBS = []string{
	"install", "phish", "exploit", "download", "bypass", "poison",
	"hack", "mine", "rip", "crack", "spy", "wire", "pwn", "access",
	"execute", "implant", "fingerprint", "obfuscate", "panic", "root",
}

var NOUNS = []string{
	"bot", "virus", "wannacry", "zeroday", "c2", "bitcoin", "device",
	"worm", "auth", "malware", "crypto", "trojan", "psyop", "server",
	"bitcoin", "camera", "ransom", "shell", "net", "darknet", "phone",
	"logger", "backdoor", "network", "kernel", "env", "ripgrep",
}

var FILE_FORMATS_AND_FRIENDS = []string{
	".sh", ".apk", ".exe", ".dmg", ".elf", ".cmd", ".bat", ".apk", ".msi",
	".dll", ".jar", ".bin", ".deb", ".rpm", ".appimage", ".tar.gz", ".run",
	".rar", ".zip", ".xlsx", ".vbs", ".pkg", ".so", ".script", ".ssh", ".ftp",
	".torrent", ".telnet", ".whois", ".netbios", ".x86", ".x32", ".arm",
	".x86_64", ".mips",
}

var TARGETS = []string{
	"target_id", "prey_id", "bait_id", "tracker_id",
}

func ShortenURL(c *gin.Context) {
	var RequestBody struct {
		LongURL string `json:"long_url"`
	}

	if err := c.ShouldBindJSON(&RequestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide the long_url"})
		return
	}

	longURL := RequestBody.LongURL
	if !govalidator.IsURL(longURL) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide a valid URL"})
		return
	}

	// Check if the URL is already shortened
	var link Link
	scheme := GetHttpScheme(c)
	hosts := strings.Split(os.Getenv("HOSTS"), ",")
	host := hosts[rand.Int()%len(hosts)]

	if err := CachedLink(FindLinkByLongURL)(&link, longURL); err == nil {
		susURL := fmt.Sprintf("%s://%s/r/%s", scheme, host, url.QueryEscape(link.SusURI))
		c.JSON(http.StatusOK, gin.H{"sus_url": susURL})
		return
	}

	// Shorten URL
	var susURIBuilder strings.Builder
	nWords := 6
	for i := 0; i < nWords; i++ {
		if i%2 == 0 {
			susURIBuilder.WriteString(VERBS[rand.Int()%len(VERBS)])
		} else {
			susURIBuilder.WriteString(NOUNS[rand.Int()%len(NOUNS)])
		}

		if i < nWords-1 {
			if rand.Float64() > 0.2 {
				susURIBuilder.WriteString(SEPARATORS[rand.Int()%len(SEPARATORS)])
			} else {
				susURIBuilder.WriteString(ESCAPABLE_SEPARATORS[rand.Int()%len(ESCAPABLE_SEPARATORS)])
			}
		}
	}

	if rand.Float64() <= 0.4 {
		susURIBuilder.WriteString(FILE_FORMATS_AND_FRIENDS[rand.Int()%len(FILE_FORMATS_AND_FRIENDS)])
	}

	susURI := susURIBuilder.String()
	target := fmt.Sprintf("&%s=%s", TARGETS[rand.Int()%len(TARGETS)], GenerateRandomString(8))
	escapedSusURI := url.QueryEscape(susURI) + target

	susURL := fmt.Sprintf("%s://%s/r/%s", scheme, host, escapedSusURI)

	link = Link{LongURL: longURL, SusURI: susURI + target}
	CreateLink(&link)

	c.JSON(http.StatusOK, gin.H{"sus_url": susURL})
}

func ResolveURL(c *gin.Context) {
	susURI := c.Param("susURI")[1:]

	// Check if the URI exists
	var link Link
	if err := CachedLink(FindLinkBySusURI)(&link, susURI); err != nil {
		c.String(http.StatusNotFound, "404: Suspish link not found")
		return
	}

	// Redirect
	c.Redirect(http.StatusPermanentRedirect, link.LongURL)
}
