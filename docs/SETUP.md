# 🚀 Setup Guide - Shopee Livestream Bot

Complete step-by-step guide to set up and run the Shopee Livestream Auto-Purchase Bot.

## 📋 Prerequisites Checklist

Before you begin, ensure you have:

- [ ] **Go 1.21+** installed
- [ ] **Chrome/Chromium** browser installed
- [ ] **Git** installed
- [ ] **Shopee Thailand account** with:
  - [ ] Verified email or phone
  - [ ] Payment method configured
  - [ ] Shipping address set
  - [ ] Sufficient balance (if using e-wallet)
- [ ] **Fast internet connection** (fiber recommended)
- [ ] **Text editor** (VS Code, Sublime, nano, etc.)

## 🖥️ System Requirements

### Minimum:
- CPU: 2 cores
- RAM: 2GB
- Storage: 500MB free space
- Internet: 10 Mbps

### Recommended:
- CPU: 4+ cores
- RAM: 4GB+
- Storage: 1GB+ free space (SSD preferred)
- Internet: 100+ Mbps (fiber optic)
- Location: Thailand or nearby (low latency)

---

## 📥 Installation Steps

### Step 1: Install Go

#### Ubuntu/Debian:
```bash
# Remove old version (if any)
sudo rm -rf /usr/local/go

# Download and install Go 1.21
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz

# Add to PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Verify installation
go version
```

#### macOS:
```bash
# Using Homebrew
brew install go@1.21

# Verify installation
go version
```

#### Windows:
1. Download installer from https://go.dev/dl/
2. Run the MSI installer
3. Open new terminal and verify:
```cmd
go version
```

### Step 2: Install Chrome/Chromium

#### Ubuntu/Debian:
```bash
# Option 1: Chromium
sudo apt update
sudo apt install chromium-browser

# Option 2: Google Chrome
wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb
sudo dpkg -i google-chrome-stable_current_amd64.deb
sudo apt-get install -f
```

#### macOS:
```bash
brew install --cask google-chrome
```

#### Windows:
Download and install from https://www.google.com/chrome/

### Step 3: Clone the Repository

```bash
# Clone the repository
git clone https://github.com/LLionNg/shopee-livestream-bot.git

# Navigate to project directory
cd shopee-livestream-bot

# Verify structure
ls -la
```

You should see:
```
cmd/
internal/
pkg/
configs/
data/
go.mod
README.md
Makefile
```

---

## ⚙️ Configuration

### Step 4: Initial Setup

Run the setup command:
```bash
make setup
```

This will:
- ✅ Create `.env` file from template
- ✅ Create necessary directories (data/cookies, data/logs, etc.)
- ✅ Display next steps

### Step 5: Configure Environment Variables

Edit the `.env` file:
```bash
nano .env
```

Fill in your Shopee credentials:
```bash
# Required: Your Shopee login credentials
SHOPEE_USERNAME=your_email@example.com
SHOPEE_PASSWORD=your_secure_password_here
SHOPEE_PHONE=+66812345678

# Optional: Proxy configuration (if you have one)
PROXY_URL=

# Optional: Notifications
WEBHOOK_URL=

# Application settings
ENV=development
LOG_LEVEL=info
```

**Security Tips:**
- ⚠️ Never share your `.env` file
- ⚠️ Never commit `.env` to git (already in .gitignore)
- ✅ Use a strong password
- ✅ Keep this file secure

### Step 6: Configure Livestreams

Edit `configs/config.yaml`:
```bash
nano configs/config.yaml
```

**Important:** Add the livestream URLs you want to monitor:

```yaml
shopee:
  base_url: "https://shopee.co.th"
  api_url: "https://shopee.co.th/api/v4"
  livestream_urls:
    # Replace these with actual livestream URLs
    - "https://live.shopee.co.th/pc/123456789"
    - "https://live.shopee.co.th/pc/987654321"
```

**How to find livestream URLs:**
1. Go to Shopee website/app
2. Navigate to "Shopee Live"
3. Open a livestream
4. Copy the URL from browser address bar
5. Paste it in `config.yaml`

### Step 7: Configure Purchase Settings

In `configs/config.yaml`, adjust these settings:

```yaml
purchase:
  max_retries: 3              # How many times to retry on failure
  retry_delay: 1              # Seconds between retries
  checkout_timeout: 5         # Maximum time for checkout (seconds)
  auto_checkout: true         # false = only add to cart, don't checkout
  pre_fill_cart: false        # true = keep items in cart before livestream
  payment_method: "ShopeePay" # Your preferred payment method
```

**Testing Mode:**
For testing without actual purchases, set:
```yaml
purchase:
  auto_checkout: false  # Only adds to cart
```

### Step 8: Configure Browser Settings

```yaml
browser:
  headless: false  # false = show browser (good for debugging)
                   # true = hide browser (good for production)
  timeout: 30
  viewport:
    width: 1920
    height: 1080
```

**Development:** Use `headless: false` to see what's happening  
**Production:** Use `headless: true` for background operation

---

## 📦 Install Dependencies

Install all required Go packages:
```bash
make install
```

Or manually:
```bash
go mod download
go mod tidy
```

This will download:
- chromedp (browser automation)
- viper (configuration)
- logrus (logging)
- And other dependencies

---

## 🏃 Running the Bot

### Method 1: Run Directly (Development)

```bash
# Using Makefile
make run

# Or using go run
go run cmd/bot/main.go
```

### Method 2: Build and Run (Production)

```bash
# Build binary
make build

# Run binary
./bin/bot
```

### Method 3: Build for Specific Platform

```bash
# Build for Linux (VPS deployment)
make build-linux

# Build for Windows
make build-windows

# Build for all platforms
make build-all
```

---

## ✅ Verification Steps

### 1. Check if Bot Starts

When you run the bot, you should see:
```
╔═══════════════════════════════════════════════════════════╗
║                                                           ║
║     🛒  SHOPEE LIVESTREAM AUTO-PURCHASE BOT  🛒          ║
║                                                           ║
║              Version: 1.0.0                               ║
║              Made with ❤️  in Go                          ║
║                                                           ║
║     ⚠️  Use responsibly & at your own risk ⚠️            ║
║                                                           ║
╚═══════════════════════════════════════════════════════════╝

INFO[...] Starting Shopee Livestream Bot...
INFO[...] Configuration loaded successfully
INFO[...] Initializing browser...
```

### 2. Check if Login Works

You should see:
```
INFO[...] Authenticating with Shopee...
INFO[...] Authentication successful!
```

**If login fails:**
- Check credentials in `.env`
- Try logging in manually first
- Check for CAPTCHA (may need manual solving)

### 3. Check if Monitoring Starts

You should see:
```
INFO[...] Starting livestream monitor...
🚀 Bot is now running! Monitoring livestreams...
🔍 Starting livestream monitoring...
📺 Monitoring 2 livestream(s)
🎥 [Stream 1] Starting monitor: https://live.shopee.co.th/...
✅ [Stream 1] Successfully loaded livestream
```

### 4. Test Purchase Flow (Optional)

For testing, set `auto_checkout: false` and manually trigger by:
1. Bot will add items to cart when detected
2. You can manually complete checkout
3. Verify everything works correctly

---

## 🔧 Troubleshooting Common Issues

### Issue 1: "go: command not found"

**Solution:**
```bash
# Verify Go installation
which go

# If not found, check PATH
echo $PATH

# Add Go to PATH
export PATH=$PATH:/usr/local/go/bin
```

### Issue 2: "chrome not found"

**Solution:**
```bash
# Ubuntu/Debian
sudo apt install chromium-browser

# macOS
brew install --cask google-chrome

# Or specify Chrome location in code if needed
```

### Issue 3: "permission denied" when running binary

**Solution:**
```bash
# Make binary executable
chmod +x bin/bot

# Then run
./bin/bot
```

### Issue 4: "cannot find package"

**Solution:**
```bash
# Clean and reinstall dependencies
rm -rf vendor/
go clean -modcache
go mod download
```

### Issue 5: Login fails / CAPTCHA appears

**Solution:**
1. Try logging in manually in Chrome first
2. May need to solve CAPTCHA manually
3. Once logged in, bot will save session
4. Subsequent runs won't need login

### Issue 6: Livestream URLs not working

**Solution:**
1. Verify URL format is correct
2. Make sure livestream is actually live
3. Check if URL has changed (Shopee updates these)
4. Try accessing URL manually in browser first

---

## 📊 Monitoring and Logs

### View Logs in Real-Time

```bash
# Follow log file
tail -f data/logs/app.log

# Last 100 lines
tail -n 100 data/logs/app.log

# Search logs
grep "ERROR" data/logs/app.log
```

### Understanding Log Levels

```
DEBUG - Detailed information for debugging
INFO  - General informational messages
WARN  - Warning messages (not critical)
ERROR - Error messages (something went wrong)
FATAL - Critical errors (bot will stop)
```

---

## 🎯 Optimization Tips

### 1. Network Optimization

```bash
# Check ping to Shopee
ping shopee.co.th

# Should be < 50ms for best results
# If higher, consider:
# - Using a VPS in Thailand
# - Upgrading internet connection
# - Using a wired connection instead of WiFi
```

### 2. System Optimization

```bash
# Close unnecessary programs
# Allocate more RAM if possible
# Use SSD for faster I/O
# Keep CPU temperature low
```

### 3. Bot Configuration

```yaml
# Faster checking (more aggressive)
monitoring:
  check_interval: 0.5  # Check every 500ms (faster but more CPU)

# More concurrent streams
monitoring:
  max_concurrent_streams: 10  # Monitor up to 10 streams

# Faster purchase
purchase:
  checkout_timeout: 3  # Reduce timeout (faster but less reliable)
```

---

## 🚀 Advanced Setup (Optional)

### Running as System Service (Linux)

Create systemd service file:
```bash
sudo nano /etc/systemd/system/shopee-bot.service
```

Content:
```ini
[Unit]
Description=Shopee Livestream Bot
After=network.target

[Service]
Type=simple
User=youruser
WorkingDirectory=/path/to/shopee-livestream-bot
ExecStart=/path/to/shopee-livestream-bot/bin/bot
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

Enable and start:
```bash
sudo systemctl enable shopee-bot
sudo systemctl start shopee-bot
sudo systemctl status shopee-bot
```

### Using Screen (Alternative for VPS)

```bash
# Start a screen session
screen -S shopee-bot

# Run bot
./bin/bot

# Detach: Ctrl+A then D
# Reattach: screen -r shopee-bot
```

---

## 📝 Next Steps

After successful setup:

1. ✅ **Test with low-value items first**
   - Don't start with expensive products
   - Verify the bot works correctly

2. ✅ **Monitor logs carefully**
   - Watch for errors
   - Check success rate
   - Adjust configuration as needed

3. ✅ **Prepare your account**
   - Ensure payment method works
   - Have sufficient balance
   - Verify shipping address

4. ✅ **Time your runs**
   - Know when livestreams start
   - Start bot 5-10 minutes early
   - Be ready to intervene if needed

5. ✅ **Stay ethical**
   - Don't abuse the system
   - Respect other buyers
   - Don't use for scalping

---

## 🆘 Getting Help

If you encounter issues:

1. **Check logs:** `tail -f data/logs/app.log`
2. **Review troubleshooting section** above
3. **Check README.md** for common issues
4. **Open GitHub issue** with:
   - Error messages
   - Log excerpts (remove sensitive data!)
   - Configuration (remove credentials!)
   - Steps to reproduce

---

## ✨ You're Ready!

Your bot is now set up and ready to use. Good luck! 🍀

**Remember:**
- Use responsibly
- Start with testing
- Monitor carefully
- Be ethical

Happy shopping! 🛍️