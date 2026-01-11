# Subman - Subscription Manager

A privacy-focused desktop application for managing your online subscriptions. Built with Go and the Fyne UI library.

## Features

- **Track Subscriptions**: Manage all your online subscriptions in one place
- **Payment History**: Automatic payment tracking with Year-to-Date spending calculations
- **Pause Subscriptions**: Temporarily pause subscriptions without losing data
- **Custom Images**: Add logos/images to subscriptions with category-based defaults
- **Cost Analysis**: View monthly, yearly, and YTD cost summaries at a glance
- **Search & Filter**: Find subscriptions by name, category, or billing cycle
- **Sort Options**: Sort by name, cost, or next payment date
- **Export Data**: Export your subscription data to CSV or JSON
- **Theme Selection**: Choose between light, dark, or system default themes
- **Privacy First**: All data stored locally on your machine - no cloud, no third parties
- **Cross-Platform**: Works on macOS, Linux, and Windows

## Data Tracked

For each subscription, you can track:
- Name
- Cost
- Billing cycle (monthly or yearly)
- Next payment date
- Start date
- Category (Streaming, Software, Utilities, Gaming, News, Education, Creator, Other)
- Custom image/logo
- Pause status
- Notes
- Payment history (automatically tracked)

## Installation

### Prerequisites

- Go 1.21 or later
- C compiler (for Fyne dependencies)
  - macOS: Xcode Command Line Tools
  - Linux: gcc
  - Windows: TDM-GCC or MinGW

### Build from Source

```bash
# Clone the repository
git clone https://github.com/douglasbarnum-cmyk/subman.git
cd subman

# Download dependencies
go mod tidy

# Build the application
go build -o subman

# Run the application
./subman
```

### Quick Start with Sample Data

To test the application with sample data, copy the test data to your config directory:

**macOS:**
```bash
mkdir -p ~/Library/Application\ Support/subman
cp testdata/sample_subscriptions.json ~/Library/Application\ Support/subman/subscriptions.json
./subman
```

**Linux:**
```bash
mkdir -p ~/.config/subman
cp testdata/sample_subscriptions.json ~/.config/subman/subscriptions.json
./subman
```

**Windows:**
```bash
mkdir %APPDATA%\subman
copy testdata\sample_subscriptions.json %APPDATA%\subman\subscriptions.json
subman.exe
```

## Data Storage Location

Your subscription data is stored locally in a JSON file:

- **macOS**: `~/Library/Application Support/subman/subscriptions.json`
- **Linux**: `~/.config/subman/subscriptions.json`
- **Windows**: `%APPDATA%\subman\subscriptions.json`

You can back up this file to preserve your data or transfer it to another machine.

## Usage

### Adding a Subscription

1. Click the "Add Subscription" button
2. Fill in the subscription details:
   - Name (e.g., "Netflix")
   - Cost (e.g., 15.99)
   - Billing Cycle (monthly or yearly)
   - Category
   - Next Payment Date (YYYY-MM-DD format)
   - Start Date (YYYY-MM-DD format)
   - Notes (optional)
3. Click "Submit"

### Editing a Subscription

1. Click the "Edit" button on any subscription card
2. Modify the details
3. Click "Submit" to save changes

### Deleting a Subscription

1. Click the "Delete" button on any subscription card
2. Confirm the deletion

### Filtering Subscriptions

Use the filter panel at the top to:
- Search by name or notes
- Filter by category
- Filter by billing cycle (monthly/yearly)
- Click "Clear Filters" to reset

### Sorting Subscriptions

1. Click the "Sort" button
2. Choose sort criteria:
   - Sort by Name
   - Sort by Cost
   - Sort by Next Payment
3. Toggle between ascending and descending order

### Exporting Data

1. Click the "Export" button
2. Select format (CSV or JSON)
3. Choose where to save the file
4. Click "Export"

## Dashboard

The dashboard at the top displays:
- **Monthly Total**: Total monthly cost (yearly subscriptions converted to monthly equivalent)
- **Yearly Total**: Total yearly cost
- **Year to Date**: Actual amount spent from January 1st to today (based on payment history)
- **Active Subscriptions**: Number of subscriptions being tracked

## Architecture

```
subman/
├── internal/
│   ├── models/         # Data models and types
│   ├── storage/        # JSON storage implementation
│   ├── service/        # Business logic
│   └── ui/             # Fyne UI components
├── pkg/
│   ├── calculator/     # Cost calculation utilities
│   └── export/         # CSV and JSON exporters
└── main.go             # Application entry point
```

## Development

### Project Structure

- **Data Layer**: JSON file storage with thread-safe operations
- **Business Logic**: CRUD operations, filtering, sorting, cost calculations
- **Presentation Layer**: Fyne-based desktop UI

### Running Tests

```bash
go test ./...
```

### Building for Different Platforms

**macOS:**
```bash
go build -o subman
```

**Linux:**
```bash
go build -o subman
```

**Windows:**
```bash
go build -o subman.exe
```

### Cross-Platform Builds with fyne-cross

Build binaries for multiple platforms and architectures:

```bash
# Install fyne-cross
go install github.com/fyne-io/fyne-cross@latest

# Build for specific platforms
fyne-cross darwin -arch=amd64,arm64 -app-id=com.subman.app  # macOS Intel + Apple Silicon
fyne-cross linux -arch=amd64,arm64 -app-id=com.subman.app   # Linux x86-64 + ARM64
fyne-cross windows -arch=amd64,arm64 -app-id=com.subman.app # Windows x86-64 + ARM64
```

Packaged binaries will be in `fyne-cross/dist/[platform]-[arch]/`

### Creating a Release

Releases are built and published manually. Docker must be running for fyne-cross.

**Prerequisites:**
```bash
# Install fyne-cross
go install github.com/fyne-io/fyne-cross@latest

# Ensure Docker is running
docker ps
```

**Release Process:**

1. **Tag and push the release:**
   ```bash
   git tag v1.0.3
   git push origin v1.0.3
   ```

2. **Build binaries for all platforms:**
   ```bash
   # macOS (Intel + Apple Silicon)
   fyne-cross darwin -arch=amd64,arm64 -app-id=com.subman.app

   # Linux (amd64 + arm64)
   fyne-cross linux -arch=amd64,arm64 -app-id=com.subman.app

   # Windows (amd64 + arm64)
   fyne-cross windows -arch=amd64,arm64 -app-id=com.subman.app
   ```

3. **Package macOS .app bundles:**
   ```bash
   cd fyne-cross/dist/darwin-amd64 && zip -r ../subman-macos-amd64.zip subman.app && cd ../../..
   cd fyne-cross/dist/darwin-arm64 && zip -r ../subman-macos-arm64.zip subman.app && cd ../../..
   ```

4. **Create the release on GitHub:**
   ```bash
   gh release create v1.0.3 --title "v1.0.3" --generate-notes
   ```

5. **Upload all binaries:**
   ```bash
   gh release upload v1.0.3 \
     fyne-cross/dist/subman-macos-amd64.zip \
     fyne-cross/dist/subman-macos-arm64.zip \
     fyne-cross/dist/linux-amd64/subman.tar.xz \
     fyne-cross/dist/linux-arm64/subman.tar.xz \
     fyne-cross/dist/windows-amd64/subman.exe.zip \
     fyne-cross/dist/windows-arm64/subman.exe.zip
   ```

**Download releases from:**
`https://github.com/douglasbarnum-cmyk/subman/releases`

## Future Enhancements

Potential features for future versions:
- Payment reminders and notifications
- Charts and visualizations for spending trends
- Multi-currency support
- Cloud sync (optional)
- Free trial expiration tracking
- Recurring payment calendar view

## Privacy

Subman is designed with privacy as the top priority:
- **100% Local**: All data stored on your machine
- **No Network Calls**: Application never connects to the internet
- **No Telemetry**: No usage tracking or analytics
- **No Third-Party APIs**: You manually enter and manage all data
- **Open Source**: Code is transparent and auditable

## License

MIT License - feel free to use and modify as needed.

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

## Support

For issues or questions, please open an issue on the project repository.
