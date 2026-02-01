# Release Notes - v1.2.0

## 🎉 New Features

### 1. Optional Email Validation
- Users can now toggle email syntax checking on/off
- When enabled, performs basic validation (e[.e]@e.c pattern)
- When disabled, processes all lines that contain data before the separator
- Helps with flexibility when processing various types of data files

### 2. Configurable Output Directory
- Users can now specify where output files should be saved
- Default: `output/` directory in the selected folder
- Supports both relative and absolute paths
- Better organization for large batch processing jobs

### 3. Optional Invalid Files Output
- Users can choose whether to generate [INVALID] files
- Helps reduce clutter when you only care about valid output
- Filename format: "original file name [INVALID].txt"

### 4. Delete After Separator Mode
- New option to remove separator and everything after it
- Useful for cleaning data files where you only need the first part
- Example: `user@example.com:password123` becomes `user@example.com`

### 5. Enhanced Error Handling
- Comprehensive error messages for all file operations
- Clear explanations when paths don't exist or lack permissions
- User-friendly error displays in the frontend
- Better validation of input parameters

### 6. Windows Compatibility Improvements
- Full support for Windows paths (both forward and backward slashes)
- Updated documentation with Windows-specific examples
- Verified compatibility with WINDOWS-INSTALL.md guidelines

## 🔧 Improvements

- Updated UI to show all new configuration options
- Added version display (v1.2.0) in the frontend
- Improved progress messages with user-friendly formatting
- Better confirmation dialogs showing all selected options

## 📝 Documentation Updates

- Updated README.md with all new features
- Added detailed usage examples for each new option
- Clarified Windows path usage recommendations
- Added examples of different processing modes

## 🧪 Testing

All features have been thoroughly tested:
- ✅ Email validation toggle (enabled/disabled)
- ✅ Custom output directories
- ✅ Invalid file output toggle
- ✅ Delete-after-separator mode
- ✅ Error handling for invalid paths
- ✅ Cross-platform path handling

## 🚀 How to Upgrade

### For Binary Users:
1. Download the new v1.2.0 release for your platform
2. Extract and replace your existing installation
3. Run `start.bat` (Windows) or `start.sh` (Linux/Mac)

### For Source Users:
1. Pull the latest changes: `git pull origin main`
2. Rebuild: `cd backend && go build -o email-processor main.go processor.go`
3. Restart the application

## 💡 Usage Examples

### Example 1: Process without email validation
- Uncheck "Validate email syntax"
- Useful for processing non-email data or when you want to accept all formats

### Example 2: Clean output only
- Uncheck "Output invalid emails to [INVALID] files"
- Get only the valid/processed output without clutter

### Example 3: Extract emails only
- Check "Delete separator and everything after it"
- Perfect for extracting just email addresses from combo lists

### Example 4: Custom organization
- Set output directory to "processed-2024-02"
- Keep different batch jobs organized

## 🙏 Feedback

If you encounter any issues or have suggestions, please open an issue on GitHub:
https://github.com/1-kabir/db-parser/issues

---

**Full Changelog**: https://github.com/1-kabir/db-parser/compare/v1.1.0...v1.2.0
