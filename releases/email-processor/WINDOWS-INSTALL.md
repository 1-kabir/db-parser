# 🪟 Windows Installation Guide

Super simple setup for non-technical users!

## 📥 Download & Install

### **Option 1: Pre-built Release (Easiest!)**

1. **Download the latest release:**
   - Go to: [Releases](https://github.com/1-kabir/db-parser/releases)
   - Download: `email-processor-windows-amd64.zip`

2. **Extract the ZIP file:**
   - Right-click → "Extract All"
   - Choose a folder (e.g., `C:\email-processor`)

3. **Run the application:**
   - Double-click `start.bat`
   - Your browser will open automatically

4. **Done!** 🎉
   - The app is now running at: `http://localhost:8080`
   - Press `Ctrl+C` in the terminal window to stop

---

### **Option 2: Build from Source**

If there's no pre-built release yet, follow these steps:

#### **Step 1: Install Go**
1. Download Go from: https://golang.org/dl/
2. Download the Windows installer (`.msi` file)
3. Run the installer (keep all default settings)
4. Restart your computer

#### **Step 2: Download the Code**
1. Download this repository as ZIP
2. Extract to a folder (e.g., `C:\email-processor`)

#### **Step 3: Build the Application**
1. Open Command Prompt (search "cmd" in Start menu)
2. Navigate to the project folder:
   ```cmd
   cd C:\path\to\email-processor\backend
   ```
3. Build the executable:
   ```cmd
   go build -o email-processor.exe main.go processor.go
   ```

#### **Step 4: Run the Application**
1. Go back to the main folder:
   ```cmd
   cd ..
   ```
2. Double-click `start-windows.bat` (or run it from cmd)

---

## 🎯 How to Use

1. **Start the app** (double-click `start.bat` or `start-windows.bat`)
2. **Open browser** to `http://localhost:8080`
3. **Enter your folder path** (e.g., `C:\Users\YourName\Documents\emails`)
4. **Click "Load Files"** to see all `.txt` files
5. **Select files** you want to process
6. **Click "Process Files"** and watch the magic happen!
7. **Find your results** in the `outputs` folder inside your original folder

---

## 🛠️ Troubleshooting

**Port 8080 already in use?**
- Create a file named `.env` in the main folder
- Add this line: `PORT=9000`
- Restart the app
- Now open: `http://localhost:9000`

**Can't find my files?**
- Make sure you use the full path (e.g., `C:\Users\YourName\Documents\emails`)
- Use forward slashes `/` or double backslashes `\\` in the path

**Antivirus blocking the app?**
- The `.exe` file is safe (it's compiled from open-source Go code)
- Add an exception in your antivirus for `email-processor.exe`

---

## 📞 Need Help?

Open an issue on GitHub: [github.com/1-kabir/db-parser/issues](https://github.com/1-kabir/db-parser/issues)
