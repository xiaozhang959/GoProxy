package webui

const loginHTML = `<!DOCTYPE html>
<html lang="zh">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>GoProxy — 身份验证</title>
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;600;700&family=Share+Tech+Mono&display=swap" rel="stylesheet">
<script>
(function(){
  try {
    document.documentElement.setAttribute('data-theme', localStorage.getItem('goproxy-theme') || 'warm');
  } catch (e) {
    document.documentElement.setAttribute('data-theme', 'warm');
  }
})();
</script>
<style>
*{box-sizing:border-box;margin:0;padding:0}
:root{color-scheme:light;--login-bg:#f7efe4;--login-card:#fffaf4;--login-input:#fff8ef;--login-border:#ead2b8;--login-accent:#b5652d;--login-muted:#8f6b50;--login-text:#3c2415;--login-tip:#7f6b5a;--login-button-text:#fffaf4;--login-scanline:rgba(125,69,36,0.025);--login-ambient:rgba(204,126,56,0.13);--login-shadow:rgba(126,74,37,0.24);--login-button-shadow:rgba(181,101,45,0.24);--login-button-hover:rgba(181,101,45,0.38)}
html[data-theme="paper"]{color-scheme:light;--login-bg:#f4f5f3;--login-card:#ffffff;--login-input:#fbfbfa;--login-border:#d8dde3;--login-accent:#475569;--login-muted:#64748b;--login-text:#111827;--login-tip:#6b7280;--login-button-text:#ffffff;--login-scanline:rgba(55,65,81,0.018);--login-ambient:rgba(100,116,139,0.08);--login-shadow:rgba(15,23,42,0.16);--login-button-shadow:rgba(15,23,42,0.14);--login-button-hover:rgba(15,23,42,0.22)}
html[data-theme="night"]{color-scheme:dark;--login-bg:#0f1724;--login-card:#111b2c;--login-input:#151f31;--login-border:#26364d;--login-accent:#6aa7ff;--login-muted:#93a4bb;--login-text:#edf5ff;--login-tip:#93a4bb;--login-button-text:#07111f;--login-scanline:rgba(106,167,255,0.026);--login-ambient:rgba(75,126,207,0.12);--login-shadow:rgba(57,102,168,0.30);--login-button-shadow:rgba(57,102,168,0.30);--login-button-hover:rgba(57,102,168,0.42)}
body{background:var(--login-bg);color:var(--login-text);font-family:JetBrains Mono,monospace;display:flex;align-items:center;justify-content:center;min-height:100vh;position:relative}
body::before{content:'';position:fixed;top:0;left:0;width:100%;height:100%;background:repeating-linear-gradient(0deg,var(--login-scanline) 0px,transparent 1px,transparent 2px,var(--login-scanline) 3px);pointer-events:none;z-index:9999}
body::after{content:'';position:fixed;top:0;left:0;width:100%;height:100%;background:radial-gradient(ellipse at center,var(--login-ambient) 0%,transparent 70%);pointer-events:none;z-index:9998}
.card{border:1px solid var(--login-accent);padding:64px;width:440px;background:var(--login-card);box-shadow:0 0 40px var(--login-shadow);position:relative;z-index:1}
h1{font-size:32px;font-weight:700;margin-bottom:8px;letter-spacing:0.15em;text-transform:uppercase;color:var(--login-accent);text-shadow:0 0 15px var(--login-accent)}
.sub{color:var(--login-muted);font-size:12px;margin-bottom:48px;font-family:JetBrains Mono,monospace;letter-spacing:0.08em;text-transform:uppercase}
label{display:block;font-size:10px;text-transform:uppercase;letter-spacing:0.1em;color:var(--login-muted);margin-bottom:8px;font-weight:600}
input[type=password]{width:100%;padding:16px;background:var(--login-input);border:1px solid var(--login-border);color:var(--login-accent);font-size:16px;font-family:JetBrains Mono,monospace;outline:none;transition:all 0.2s}
input[type=password]:focus{border-color:var(--login-accent);background:var(--login-card);box-shadow:0 0 10px var(--login-shadow)}
button{width:100%;margin-top:24px;padding:16px;background:var(--login-accent);color:var(--login-button-text);border:1px solid var(--login-accent);font-size:12px;font-weight:700;cursor:pointer;transition:all 0.2s;text-transform:uppercase;letter-spacing:0.1em;font-family:JetBrains Mono,monospace;box-shadow:0 0 15px var(--login-button-shadow)}
button:hover{box-shadow:0 0 25px var(--login-button-hover);transform:translateY(-1px)}
.logo{font-size:64px;margin-bottom:24px;line-height:1;font-weight:700;letter-spacing:0.1em;color:var(--login-accent);text-shadow:0 0 20px var(--login-accent)}
.tip{color:var(--login-tip);font-size:10px;margin-top:24px;line-height:1.6;letter-spacing:0.05em;text-align:center}
.tip a{color:var(--login-accent);text-decoration:none;border-bottom:1px solid transparent;transition:all 0.2s}
.tip a:hover{border-bottom-color:var(--login-accent);text-shadow:0 0 8px var(--login-accent)}
.github{position:absolute;top:20px;right:20px;color:var(--login-accent);opacity:0.6;transition:all 0.3s}
.github:hover{opacity:1;transform:scale(1.1);filter:drop-shadow(0 0 8px var(--login-accent))}
</style>
</head>
<body>
<a href="https://github.com/isboyjc/ProxyGo" target="_blank" class="github" title="GitHub">
  <svg width="32" height="32" viewBox="0 0 16 16" fill="currentColor">
    <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z"/>
  </svg>
</a>
<div class="card">
  <div class="logo">[GP]</div>
  <h1>GoProxy</h1>
  <p class="sub">// Intelligent Proxy Pool</p>
  <form method="POST" action="/login">
    <label>> Password</label>
    <input type="password" name="password" placeholder="****************" autofocus>
    <button type="submit">[ AUTHENTICATE ]</button>
  </form>
  <p class="tip">访客模式可<a href="/">查看数据</a>，管理员登录后可完全控制</p>
</div>
</body>
</html>`

const loginHTMLWithError = `<!DOCTYPE html>
<html lang="zh">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>GoProxy — 身份验证</title>
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;600;700&family=Share+Tech+Mono&display=swap" rel="stylesheet">
<script>
(function(){
  try {
    document.documentElement.setAttribute('data-theme', localStorage.getItem('goproxy-theme') || 'warm');
  } catch (e) {
    document.documentElement.setAttribute('data-theme', 'warm');
  }
})();
</script>
<style>
*{box-sizing:border-box;margin:0;padding:0}
:root{color-scheme:light;--login-bg:#f7efe4;--login-card:#fffaf4;--login-input:#fff8ef;--login-border:#ead2b8;--login-accent:#b5652d;--login-muted:#8f6b50;--login-text:#3c2415;--login-tip:#7f6b5a;--login-button-text:#fffaf4;--login-scanline:rgba(125,69,36,0.025);--login-ambient:rgba(204,126,56,0.13);--login-shadow:rgba(126,74,37,0.24);--login-button-shadow:rgba(181,101,45,0.24);--login-button-hover:rgba(181,101,45,0.38)}
html[data-theme="paper"]{color-scheme:light;--login-bg:#f4f5f3;--login-card:#ffffff;--login-input:#fbfbfa;--login-border:#d8dde3;--login-accent:#475569;--login-muted:#64748b;--login-text:#111827;--login-tip:#6b7280;--login-button-text:#ffffff;--login-scanline:rgba(55,65,81,0.018);--login-ambient:rgba(100,116,139,0.08);--login-shadow:rgba(15,23,42,0.16);--login-button-shadow:rgba(15,23,42,0.14);--login-button-hover:rgba(15,23,42,0.22)}
html[data-theme="night"]{color-scheme:dark;--login-bg:#0f1724;--login-card:#111b2c;--login-input:#151f31;--login-border:#26364d;--login-accent:#6aa7ff;--login-muted:#93a4bb;--login-text:#edf5ff;--login-tip:#93a4bb;--login-button-text:#07111f;--login-scanline:rgba(106,167,255,0.026);--login-ambient:rgba(75,126,207,0.12);--login-shadow:rgba(57,102,168,0.30);--login-button-shadow:rgba(57,102,168,0.30);--login-button-hover:rgba(57,102,168,0.42)}
body{background:var(--login-bg);color:var(--login-text);font-family:JetBrains Mono,monospace;display:flex;align-items:center;justify-content:center;min-height:100vh;position:relative}
body::before{content:'';position:fixed;top:0;left:0;width:100%;height:100%;background:repeating-linear-gradient(0deg,var(--login-scanline) 0px,transparent 1px,transparent 2px,var(--login-scanline) 3px);pointer-events:none;z-index:9999}
body::after{content:'';position:fixed;top:0;left:0;width:100%;height:100%;background:radial-gradient(ellipse at center,var(--login-ambient) 0%,transparent 70%);pointer-events:none;z-index:9998}
.card{border:1px solid var(--login-accent);padding:64px;width:440px;background:var(--login-card);box-shadow:0 0 40px var(--login-shadow);position:relative;z-index:1}
h1{font-size:32px;font-weight:700;margin-bottom:8px;letter-spacing:0.15em;text-transform:uppercase;color:var(--login-accent);text-shadow:0 0 15px var(--login-accent)}
.sub{color:var(--login-muted);font-size:12px;margin-bottom:48px;font-family:JetBrains Mono,monospace;letter-spacing:0.08em;text-transform:uppercase}
label{display:block;font-size:10px;text-transform:uppercase;letter-spacing:0.1em;color:var(--login-muted);margin-bottom:8px;font-weight:600}
input[type=password]{width:100%;padding:16px;background:var(--login-input);border:1px solid var(--login-border);color:var(--login-accent);font-size:16px;font-family:JetBrains Mono,monospace;outline:none;transition:all 0.2s}
input[type=password]:focus{border-color:var(--login-accent);background:var(--login-card);box-shadow:0 0 10px var(--login-shadow)}
button{width:100%;margin-top:24px;padding:16px;background:var(--login-accent);color:var(--login-button-text);border:1px solid var(--login-accent);font-size:12px;font-weight:700;cursor:pointer;transition:all 0.2s;text-transform:uppercase;letter-spacing:0.1em;font-family:JetBrains Mono,monospace;box-shadow:0 0 15px var(--login-button-shadow)}
button:hover{box-shadow:0 0 25px var(--login-button-hover);transform:translateY(-1px)}
.logo{font-size:64px;margin-bottom:24px;line-height:1;font-weight:700;letter-spacing:0.1em;color:var(--login-accent);text-shadow:0 0 20px var(--login-accent)}
.error{background:#ff0033;color:#fff;padding:16px;font-size:11px;margin-bottom:24px;font-family:JetBrains Mono,monospace;font-weight:600;border:1px solid #ff0033;box-shadow:0 0 15px rgba(255,0,51,0.5);text-transform:uppercase;letter-spacing:0.05em}
.tip{color:var(--login-tip);font-size:10px;margin-top:24px;line-height:1.6;letter-spacing:0.05em;text-align:center}
.tip a{color:var(--login-accent);text-decoration:none;border-bottom:1px solid transparent;transition:all 0.2s}
.tip a:hover{border-bottom-color:var(--login-accent);text-shadow:0 0 8px var(--login-accent)}
.github{position:absolute;top:20px;right:20px;color:var(--login-accent);opacity:0.6;transition:all 0.3s}
.github:hover{opacity:1;transform:scale(1.1);filter:drop-shadow(0 0 8px var(--login-accent))}
</style>
</head>
<body>
<a href="https://github.com/isboyjc/ProxyGo" target="_blank" class="github" title="GitHub">
  <svg width="32" height="32" viewBox="0 0 16 16" fill="currentColor">
    <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z"/>
  </svg>
</a>
<div class="card">
  <div class="logo">[GP]</div>
  <h1>GoProxy</h1>
  <p class="sub">// Intelligent Proxy Pool</p>
  <div class="error">[!] ACCESS DENIED - INVALID PASSWORD</div>
  <form method="POST" action="/login">
    <label>> Password</label>
    <input type="password" name="password" placeholder="****************" autofocus>
    <button type="submit">[ AUTHENTICATE ]</button>
  </form>
  <p class="tip">访客模式可<a href="/">查看数据</a>，管理员登录后可完全控制</p>
</div>
</body>
</html>`

// dashboardHTML 已移至 dashboard.go
