package main

import (
	"fmt"
	"image/color"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

 
type darkTheme struct{}

func (darkTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.NRGBA{R: 8, G: 12, B: 20, A: 255}
	case theme.ColorNameButton:
		return color.NRGBA{R: 0, G: 150, B: 0, A: 255}
	case theme.ColorNameForeground:
		return color.NRGBA{R: 0, G: 255, B: 0, A: 255}
	case theme.ColorNamePrimary:
		return color.NRGBA{R: 0, G: 255, B: 0, A: 255}
	case theme.ColorNameInputBackground:
		return color.NRGBA{R: 20, G: 30, B: 40, A: 255}
	case theme.ColorNamePlaceHolder:
		return color.NRGBA{R: 100, G: 150, B: 100, A: 255}
	case theme.ColorNameDisabled:
		return color.NRGBA{R: 50, G: 80, B: 50, A: 255}
	default:
		return color.NRGBA{R: 0, G: 255, B: 0, A: 255}
	}
}

func (darkTheme) Font(s fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(s)
}

func (darkTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (darkTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

func main() {
	myApp := app.New()
	myApp.Settings().SetTheme(&darkTheme{})

	myWindow := myApp.NewWindow("⚡ PySecTools - WordPress Login Checker ⚡")
	myWindow.Resize(fyne.NewSize(600, 500)) 
	myWindow.CenterOnScreen()

 
	titleLabel := widget.NewLabelWithStyle(
		"┌─[ PySecTools ]─[ WordPress Login Checker ]─[ v1.0 ]─┐\n└─────────────────────────────────────────────────────┘",
		fyne.TextAlignCenter,
		fyne.TextStyle{Monospace: true, Bold: true},
	)

 
	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("https://target-site.com")
	urlEntry.TextStyle.Monospace = true

	userEntry := widget.NewEntry()
	userEntry.SetPlaceHolder("admin / root / user@example.com")
	userEntry.TextStyle.Monospace = true

	passEntry := widget.NewPasswordEntry()
	passEntry.SetPlaceHolder("**********")
	passEntry.TextStyle.Monospace = true

 
	resultLabel := widget.NewLabel("[ SYSTEM ] ➜ Ready for action...")
	resultLabel.Alignment = fyne.TextAlignCenter
	resultLabel.TextStyle.Monospace = true
	resultLabel.Wrapping = fyne.TextWrapWord

	// Progress bar
	progressBar := widget.NewProgressBarInfinite()
	progressBar.Hide()

	 
	logText := widget.NewMultiLineEntry()
	logText.SetText("[>] PySecTools initialized\n[>] Dark mode activated\n[>] Ready for login check\n")
	logText.Disable()
	logText.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}
	logText.SetMinRowsVisible(6)  

 
	telegramLink := widget.NewHyperlink("https://t.me/PySecTools", &url.URL{
		Scheme: "https",
		Host:   "t.me",
		Path:   "/PySecTools",
	})
	telegramLink.Alignment = fyne.TextAlignCenter
	telegramLink.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}

	var testBtn, checkBtn, clearBtn *widget.Button

	 
	addLog := func(msg string) {
		currentLog := logText.Text
		logText.SetText(currentLog + "\n" + time.Now().Format("15:04:05") + " " + msg)
		logText.CursorRow = len(logText.Text)
	}

 
	clearBtn = widget.NewButtonWithIcon("Clear Log", theme.DeleteIcon(), func() {
		logText.SetText("[>] PySecTools - Log cleared\n")
		addLog("[INFO] Log history deleted")
	})

	 
	testBtn = widget.NewButtonWithIcon("Test Connection", theme.ComputerIcon(), func() {
		site := strings.TrimSpace(urlEntry.Text)

		if site == "" {
			resultLabel.SetText("[ ERROR ] ➜ No target specified!")
			addLog("[ERROR] No URL provided")
			return
		}

		if !strings.HasPrefix(site, "http://") && !strings.HasPrefix(site, "https://") {
			site = "https://" + site
		}

		testBtn.Disable()
		checkBtn.Disable()
		clearBtn.Disable()
		progressBar.Show()
		resultLabel.SetText("[ STATUS ] ➜ Testing connection...")
		addLog("[INFO] Testing connection to: " + site)

		go func() {
			client := &http.Client{Timeout: 10 * time.Second}
			startTime := time.Now()
			resp, err := client.Get(site)
			responseTime := time.Since(startTime)

			fyne.Do(func() {
				if err != nil {
					resultLabel.SetText(fmt.Sprintf("[ FAILED ] ➜ Cannot reach %s", site))
					addLog(fmt.Sprintf("[ERROR] Connection failed: %v", err))
				} else {
					defer resp.Body.Close()
					resultLabel.SetText(fmt.Sprintf("[ SUCCESS ] ➜ Online | Status: %d | Time: %v", resp.StatusCode, responseTime))
					addLog(fmt.Sprintf("[SUCCESS] Target reachable | Status: %d | Response: %v", resp.StatusCode, responseTime))
				}
				progressBar.Hide()
				testBtn.Enable()
				checkBtn.Enable()
				clearBtn.Enable()
			})
		}()
	})

	 
	checkBtn = widget.NewButtonWithIcon("Login Check", theme.ConfirmIcon(), func() {
		siteUrl := strings.TrimSpace(urlEntry.Text)
		username := strings.TrimSpace(userEntry.Text)
		password := passEntry.Text

		if siteUrl == "" || username == "" || password == "" {
			resultLabel.SetText("[ ERROR ] ➜ Missing parameters!")
			addLog("[ERROR] Required fields empty")
			return
		}

		if !strings.HasPrefix(siteUrl, "http://") && !strings.HasPrefix(siteUrl, "https://") {
			siteUrl = "https://" + siteUrl
		}

		checkBtn.Disable()
		testBtn.Disable()
		clearBtn.Disable()
		progressBar.Show()
		resultLabel.SetText("[ EXECUTING ] ➜ Launching login check...")
		addLog(fmt.Sprintf("[ATTACK] Targeting: %s | User: %s", siteUrl, username))

		go func() {
			startTime := time.Now()
			msg, success := checkWordPressLogin(siteUrl, username, password)
			responseTime := time.Since(startTime)

			fyne.Do(func() {
				if success {
					resultLabel.SetText(fmt.Sprintf("[ ACCESS GRANTED ] ➜ %s | Time: %v", msg, responseTime))
					addLog(fmt.Sprintf("[SUCCESS] ✅ Login valid! Time: %v", responseTime))
				} else {
					resultLabel.SetText(fmt.Sprintf("[ ACCESS DENIED ] ➜ %s | Time: %v", msg, responseTime))
					addLog(fmt.Sprintf("[FAILED] ❌ Login invalid | Time: %v", responseTime))
				}
				progressBar.Hide()
				checkBtn.Enable()
				testBtn.Enable()
				clearBtn.Enable()
			})
		}()
	})

 
	sysInfo := widget.NewLabelWithStyle(
		"┌─[ PySecTools System Info ]─────────────────┐\n│ Tool: WordPress Login Checker v1.0        │\n│ Author: PySecTools                         │\n│ Method: Cookie-based detection             │\n│ Telegram: @PySecTools                      │\n└─────────────────────────────────────────────┘",
		fyne.TextAlignCenter,
		fyne.TextStyle{Monospace: true},
	)

	 
	helpContent := container.NewVBox(
		widget.NewLabelWithStyle(
			"\n[ HOW TO USE ]\n"+
			"━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"+
			"1. Enter target WordPress site URL\n"+
			"2. Input username/email\n"+
			"3. Input password\n"+
			"4. Click 'Login Check'\n\n"+
			"[ FEATURES ]\n"+
			"━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"+
			"• Dark/Hacker theme with neon green text\n"+
			"• Cookie-based login detection\n"+
			"• Real-time console logging\n"+
			"• Response time measurement\n"+
			"• Multi-tab interface\n\n"+
			"[ CONTACT ]\n"+
			"━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"+
			"Telegram: @PySecTools\n"+
			"⚠️ Use only on authorized targets!\n",
			fyne.TextAlignCenter,
			fyne.TextStyle{Monospace: true},
		),
		telegramLink,
	)

	 
	content := container.NewBorder(
		container.NewVBox(
			titleLabel,
			widget.NewSeparator(),
		),
		container.NewVBox(
			widget.NewSeparator(),
			sysInfo,
		),
		nil,
		nil,
		container.NewAppTabs(
			container.NewTabItem("🎯 Target", container.NewVBox(
				widget.NewLabel("└─ URL:"),
				urlEntry,
				widget.NewLabel("└─ Username:"),
				userEntry,
				widget.NewLabel("└─ Password:"),
				passEntry,
				widget.NewSeparator(),
				container.NewHBox(
					checkBtn,
					testBtn,
					clearBtn,
				),
				progressBar,
				widget.NewSeparator(),
				resultLabel,
			)),
			container.NewTabItem("📜 Console Log", container.NewVBox(
				logText,
			)),
			container.NewTabItem("ℹ️ Help", helpContent),
		),
	)

	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}

func checkWordPressLogin(siteUrl, username, password string) (string, bool) {
	siteUrl = strings.TrimRight(siteUrl, "/")
	loginUrl := siteUrl + "/wp-login.php"

	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Timeout: 20 * time.Second,
		Jar:     jar,
	}

	respGet, err := client.Get(loginUrl)
	if err != nil {
		return "Cannot reach login page", false
	}
	respGet.Body.Close()

	data := url.Values{}
	data.Set("log", username)
	data.Set("pwd", password)
	data.Set("wp-submit", "Log In")
	data.Set("redirect_to", siteUrl+"/wp-admin/")
	data.Set("testcookie", "1")

	respPost, err := client.PostForm(loginUrl, data)
	if err != nil {
		return "Connection error", false
	}
	defer respPost.Body.Close()

	parsedUrl, _ := url.Parse(siteUrl)
	cookies := jar.Cookies(parsedUrl)

	for _, cookie := range cookies {
		if strings.Contains(cookie.Name, "wordpress_logged_in") {
			return "Access granted ✓", true
		}
	}

	body, _ := io.ReadAll(respPost.Body)
	bodyStr := string(body)

	if strings.Contains(bodyStr, "incorrect") ||
		strings.Contains(bodyStr, "Invalid username") ||
		strings.Contains(bodyStr, "wrong password") ||
		strings.Contains(bodyStr, "ERROR") {
		return "Access denied ✗", false
	}

	return fmt.Sprintf("Unknown response (HTTP %d)", respPost.StatusCode), false
}
