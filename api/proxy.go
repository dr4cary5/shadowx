package handler

import (
    "compress/gzip"
    "fmt"
    "io"
    "net/http"
    "strings"
    "time"
)

func Handler(w http.ResponseWriter, r *http.Request) {
    url := r.URL.Query().Get("url")

    // صفحه اصلی با کادر ورود URL
    if url == "" {
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        w.Write([]byte(homePage))
        return
    }

    // اضافه کردن پروتکل
    if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
        url = "https://" + url
    }

    // ساخت کلاینت با تنظیمات پیشرفته
    client := &http.Client{
        Timeout: 30 * time.Second,
        CheckRedirect: func(req *http.Request, via []*http.Request) error {
            if len(via) >= 5 {
                return fmt.Errorf("too many redirects")
            }
            return nil
        },
    }

    // ساخت درخواست
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // تنظیم هدرهای واقعی
    req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36")
    req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
    req.Header.Set("Accept-Language", "en-US,en;q=0.5")
    req.Header.Set("Accept-Encoding", "gzip, deflate")
    req.Header.Set("Cache-Control", "no-cache")
    req.Header.Set("Pragma", "no-cache")

    // ارسال درخواست
    resp, err := client.Do(req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadGateway)
        return
    }
    defer resp.Body.Close()

    // کپی هدرهای پاسخ
    skipHeaders := map[string]bool{
        "Content-Encoding":    true,
        "Content-Length":      true,
        "Transfer-Encoding":   true,
        "Strict-Transport-Security": true,
        "Content-Security-Policy":   true,
    }

    for key, values := range resp.Header {
        if skipHeaders[key] {
            continue
        }
        for _, value := range values {
            w.Header().Add(key, value)
        }
    }

    w.WriteHeader(resp.StatusCode)

    // پشتیبانی از gzip
    var reader io.Reader = resp.Body
    if resp.Header.Get("Content-Encoding") == "gzip" {
        gz, err := gzip.NewReader(resp.Body)
        if err == nil {
            defer gz.Close()
            reader = gz
        }
    }

    io.Copy(w, reader)
}

const homePage = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>ShadowX Proxy</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: 'Segoe UI', system-ui, sans-serif;
            background: #0a0a0a;
            color: #e0e0e0;
            display: flex;
            justify-content: center;
            align-items: center;
            min-height: 100vh;
            flex-direction: column;
            padding: 20px;
        }
        .container {
            background: #111;
            padding: 3rem;
            border-radius: 20px;
            box-shadow: 0 0 60px rgba(110, 0, 255, 0.2);
            width: 100%;
            max-width: 700px;
            border: 1px solid #222;
        }
        .logo {
            text-align: center;
            font-size: 4rem;
            margin-bottom: 0.5rem;
        }
        h1 {
            text-align: center;
            color: #8a2be2;
            margin-bottom: 0.3rem;
            font-size: 2.5rem;
            letter-spacing: 3px;
        }
        .tagline {
            text-align: center;
            color: #555;
            font-size: 0.9rem;
            margin-bottom: 2rem;
            font-style: italic;
        }
        .url-form {
            display: flex;
            gap: 10px;
            margin-bottom: 1.5rem;
        }
        input {
            flex: 1;
            padding: 16px;
            font-size: 16px;
            border: 2px solid #333;
            border-radius: 12px;
            background: #0d0d0d;
            color: #fff;
            outline: none;
            transition: border-color 0.3s;
        }
        input:focus { border-color: #6e00ff; }
        button {
            padding: 16px 28px;
            font-size: 16px;
            background: #6e00ff;
            color: #fff;
            border: none;
            border-radius: 12px;
            cursor: pointer;
            font-weight: bold;
            transition: all 0.3s;
        }
        button:hover { background: #8a2be2; transform: scale(1.02); }
        .features {
            display: flex;
            gap: 12px;
            justify-content: center;
            flex-wrap: wrap;
            margin-top: 1.5rem;
        }
        .feature {
            background: #1a1a1a;
            padding: 8px 16px;
            border-radius: 20px;
            font-size: 0.8rem;
            color: #999;
            border: 1px solid #2a2a2a;
        }
        .footer {
            text-align: center;
            margin-top: 2rem;
            color: #444;
            font-size: 0.75rem;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="logo">🕶️</div>
        <h1>ShadowX Proxy</h1>
        <p class="tagline">In shadow we trust, in freedom we surf</p>

        <form onsubmit="openProxy(event)">
            <div class="url-form">
                <input type="text" id="url" placeholder="Enter URL (e.g. google.com)" autofocus>
                <button type="submit">Go</button>
            </div>
        </form>

        <div class="features">
            <span class="feature">⚡ Serverless</span>
            <span class="feature">🔒 Encrypted</span>
            <span class="feature">🌍 Global CDN</span>
            <span class="feature">🎭 Stealth Mode</span>
        </div>
    </div>

    <div class="footer">ShadowX Proxy v2.0 · Powered by Vercel Edge Network</div>

    <script>
        function openProxy(e) {
            e.preventDefault();
            const url = document.getElementById('url').value.trim();
            if (url) {
                window.open('/api/proxy?url=' + encodeURIComponent(url), '_blank');
            }
        }
    </script>
</body>
</html>`
