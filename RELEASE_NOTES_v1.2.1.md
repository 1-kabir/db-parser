# Release Notes - v1.2.1

## 🎯 Key Improvement

This is a minor release that improves output file naming for better clarity and consistency.

### Enhanced Filename Handling

**What Changed:**
- Output files now have ALL curly bracket content removed from the original filename
- The line count is added at the end in curly brackets with proper comma formatting
- Multiple spaces are cleaned up automatically

**Examples:**

| Original Filename | Lines | Old Behavior (v1.2.0) | New Behavior (v1.2.1) |
|-------------------|-------|-----------------------|------------------------|
| `sample {100}.txt` | 1,500 | `sample {1500}.txt` | `sample {1,500}.txt` |
| `data {old} file {test}.txt` | 2,500 | `data {old} file {2500}.txt` | `data file {2,500}.txt` |
| `emails {v1} {backup}.txt` | 100,000 | `emails {v1} {100000}.txt` | `emails {100,000}.txt` |
| `simple.txt` | 999 | `simple {999}.txt` | `simple {999}.txt` |

**Benefits:**
- ✅ Cleaner, more professional output filenames
- ✅ All old metadata in brackets is removed
- ✅ Numbers are easier to read with comma separators (international format)
- ✅ Consistent naming across all processed files
- ✅ Automatic cleanup of extra spaces

## 🔧 Technical Details

### Number Formatting
Numbers are formatted using the international comma separator system:
- 999 → `999` (no commas for numbers under 1,000)
- 1,000 → `1,000`
- 1,500 → `1,500`
- 12,345 → `12,345`
- 100,000 → `100,000`
- 1,000,000 → `1,000,000`

### Space Handling
Multiple consecutive spaces resulting from bracket removal are automatically cleaned:
- `sample  {a}  {b}  test.txt` → `sample test {count}.txt` (single spaces)

## 🚀 Upgrade Notes

This is a drop-in replacement for v1.2.0:
- No configuration changes needed
- No breaking changes
- All v1.2.0 features are preserved

### For Binary Users:
1. Download the new v1.2.1 release for your platform
2. Replace your existing installation
3. Run `start.bat` (Windows) or `start.sh` (Linux/Mac)

### For Source Users:
1. Pull the latest changes
2. Rebuild: `cd backend && go build -o email-processor main.go processor.go`
3. Restart the application

## 📦 What's Included from v1.2.0

All features from v1.2.0 are included:
- ✅ Optional email validation
- ✅ Configurable output directory
- ✅ Optional invalid file output
- ✅ Delete-after-separator mode
- ✅ Enhanced error handling
- ✅ Windows path support

## 🧪 Testing

Thoroughly tested with:
- ✅ Files with single brackets
- ✅ Files with multiple brackets
- ✅ Files with no brackets
- ✅ Various line counts (5, 999, 1,500, 100,000)
- ✅ Files with multiple spaces
- ✅ All formatting scenarios

## 🙏 Feedback

If you encounter any issues or have suggestions, please open an issue on GitHub:
https://github.com/1-kabir/db-parser/issues

---

**Full Changelog**: https://github.com/1-kabir/db-parser/compare/v1.2.0...v1.2.1
