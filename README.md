# DBCA TUI

A Terminal User Interface (TUI) that mimics Oracle's Database Configuration Assistant (DBCA) wizard. It walks you through all configuration options and generates the corresponding `dbca -silent` command for Oracle 19c database creation or deletion.

## Features

- **Database Creation & Deletion**: Generate commands for both operations
- **Interactive wizard** following Oracle DBCA's flow
- **Two modes**: Typical (simplified) and Advanced (full control)
- **Supports all deployment types**: Single Instance, RAC, RAC One Node
- **Container database support**: CDB/PDB configuration
- **Storage options**: File System or ASM
- **Archive Log Mode**: Easy toggle for ARCHIVELOG/NOARCHIVELOG
- **Complete command generation**: Ready-to-use `dbca -silent` command
- **Save to file**: Export command as executable shell script

## Requirements

- Go 1.21 or later

## Building

### Build for your current platform

```bash
go build -o dbca_tui .
```

### Cross-compile for different platforms

**macOS (Apple Silicon / ARM64):**
```bash
GOOS=darwin GOARCH=arm64 go build -o dbca_tui-darwin-arm64 .
```

**macOS (Intel / AMD64):**
```bash
GOOS=darwin GOARCH=amd64 go build -o dbca_tui-darwin-amd64 .
```

**Linux (AMD64):**
```bash
GOOS=linux GOARCH=amd64 go build -o dbca_tui-linux-amd64 .
```

**Linux (ARM64):**
```bash
GOOS=linux GOARCH=arm64 go build -o dbca_tui-linux-arm64 .
```

**Windows (AMD64):**
```bash
GOOS=windows GOARCH=amd64 go build -o dbca_tui-windows-amd64.exe .
```

### Build all platforms at once

Use the included build script to cross-compile for all supported platforms:

```bash
./build.sh
```

This creates binaries in the `dist/` directory for:
- macOS (ARM64 + AMD64)
- Linux (AMD64 + ARM64 + 386)
- Windows (AMD64 + 386)
- Solaris (AMD64)
- FreeBSD (AMD64)

**Note:** AIX is not supported due to clipboard library limitations.

### GitHub Actions (Automated Builds)

The repository includes a GitHub Actions workflow that automatically:

1. **On every push to `main`**: Builds all binaries and uploads them as artifacts
2. **On tag push (`v*`)**: Creates a GitHub Release with all binaries and checksums

To create a release:

```bash
git tag v1.0.0
git push origin v1.0.0
```

The workflow will automatically build all platforms and create a release with downloadable binaries.

## Usage

Run the compiled binary:

```bash
./dbca_tui
```

### Navigation

| Key | Action |
|-----|--------|
| `↑` / `↓` | Navigate options |
| `Enter` | Select / Confirm |
| `Tab` | Next field |
| `Shift+Tab` | Previous field |
| `Esc` | Go back |
| `q` | Quit |

### Step-specific Keys

| Key | Step | Action |
|-----|------|--------|
| `a` | Recovery & Archive Log | Toggle Archive Log Mode |
| `f` | Recovery & Archive Log | Toggle Fast Recovery Area |
| `f` | Delete Database | Toggle Force Delete |
| `c` | Credentials | Toggle common password mode |
| `p` | Summary | Toggle password visibility |
| `s` | Summary | Save to file |
| `g` | Summary | Generate command and exit |

### Wizard Steps

#### Create Database Flow

1. **Operation** - Create or Delete database
2. **Creation Mode** - Typical (fewer steps) or Advanced (full control)
3. **Deployment Type** - Single Instance, RAC, or RAC One Node
4. **Database Template** - General Purpose, Data Warehouse, or Custom
5. **Database Identification** - Global name, SID, CDB/PDB settings
6. **Storage Configuration** - File System or ASM
7. **Recovery & Archive Log** - FRA settings and Archive Log Mode (ARCHIVELOG/NOARCHIVELOG)
8. **Network Configuration** - Listener settings (Advanced mode)
9. **Data Vault** - Security configuration (Advanced mode)
10. **Configuration Options** - Memory, character set, connection mode
11. **Management Options** - Enterprise Manager (Advanced mode)
12. **Credentials** - Database passwords
13. **Summary** - Review and generate command

#### Delete Database Flow

1. **Operation** - Select "Delete a Database"
2. **Delete Configuration** - Database SID, SYS password, force delete option
3. **Summary** - Review and generate delete command

### Output

At the end of the wizard, you'll see a preview of the generated command. You can:

- Press `g` or `Enter` to **generate the command and exit** - the command will be printed to your terminal
- Press `p` to toggle password visibility in the preview
- Press `s` to save the command to a shell script file (`dbca_<SID>.sh` or `dbca_delete_<SID>.sh`)
- Press `q` to exit without printing

When you select "Generate command and exit", the wizard closes and prints the complete `dbca -silent` command to your terminal, making it easy to copy or pipe to other commands.

## Example Output

### Create Database Command

```bash
dbca -silent -createDatabase \
  -templateName General_Purpose.dbt \
  -gdbname orcl.example.com \
  -sid orcl \
  -createAsContainerDatabase true \
  -numberOfPDBs 1 \
  -pdbName orclpdb \
  -pdbAdminPassword '<PASSWORD>' \
  -sysPassword '<PASSWORD>' \
  -systemPassword '<PASSWORD>' \
  -characterSet AL32UTF8 \
  -nationalCharacterSet AL16UTF16 \
  -totalMemory 2048 \
  -memoryMgmtType AUTO \
  -databaseType MULTIPURPOSE \
  -storageType FS \
  -datafileDestination '/u01/app/oracle/oradata' \
  -useOMF true \
  -recoveryAreaDestination '/u01/app/oracle/fast_recovery_area' \
  -recoveryAreaSize 10240 \
  -redoLogFileSize 50 \
  -emConfiguration NONE \
  -archiveLogMode true \
  -databaseConfigType SI
```

### Delete Database Command

```bash
dbca -silent -deleteDatabase \
  -sourceDB orcl \
  -sysDBAUserName SYS \
  -sysDBAPassword '<PASSWORD>' \
  -forceArchiveLogDeletion
```

## Project Structure

```
dbca_tui/
├── main.go                     # Entry point
├── go.mod                      # Go module definition
├── build.sh                    # Cross-platform build script
├── internal/
│   ├── wizard/
│   │   ├── wizard.go           # Wizard controller
│   │   └── steps.go            # Step interface
│   ├── steps/                  # Individual wizard steps
│   │   ├── operation.go        # Create/Delete selection
│   │   ├── creation_mode.go
│   │   ├── deployment.go
│   │   ├── template.go
│   │   ├── identification.go
│   │   ├── storage.go
│   │   ├── recovery.go         # FRA & Archive Log settings
│   │   ├── network.go
│   │   ├── datavault.go
│   │   ├── config.go
│   │   ├── management.go
│   │   ├── credentials.go
│   │   ├── delete.go           # Delete database configuration
│   │   └── summary.go
│   ├── model/
│   │   └── dbconfig.go         # Configuration struct
│   ├── generator/
│   │   └── command.go          # DBCA command generator (create & delete)
│   └── ui/
│       ├── styles.go           # Terminal styles
│       └── components.go       # UI components
└── README.md
```

## Dependencies

- [bubbletea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- [bubbles](https://github.com/charmbracelet/bubbles) - TUI components

## License

MIT
