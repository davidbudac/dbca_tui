# DBCA TUI - Project Guidance

## Project Overview

DBCA TUI is a Terminal User Interface that mimics Oracle's Database Configuration Assistant (DBCA) wizard. It generates `dbca -silent` commands for Oracle 19c database creation without actually executing anything.

**Target Oracle Version:** 19c (Long-term support release)

## Tech Stack

- **Language:** Go 1.21+
- **TUI Framework:** [bubbletea](https://github.com/charmbracelet/bubbletea) (Elm architecture)
- **Styling:** [lipgloss](https://github.com/charmbracelet/lipgloss)
- **Components:** [bubbles](https://github.com/charmbracelet/bubbles)

## Development

### Build

```bash
go build -o dbca_tui .
```

### Run

```bash
./dbca_tui
```

### Cross-compile for macOS ARM64

```bash
GOOS=darwin GOARCH=arm64 go build -o dbca_tui-darwin-arm64 .
```

## Architecture

### Wizard Flow

The application uses a wizard pattern with the following components:

1. **Wizard** (`internal/wizard/wizard.go`) - Main controller that manages step navigation
2. **Step Interface** (`internal/wizard/steps.go`) - Each step implements this interface
3. **Steps** (`internal/steps/`) - Individual wizard step implementations
4. **Model** (`internal/model/dbconfig.go`) - Holds all configuration state
5. **Generator** (`internal/generator/command.go`) - Builds the DBCA command string
6. **UI** (`internal/ui/`) - Reusable styles and components

### Step Interface

Each step must implement:

```go
type Step interface {
    Init(config *DBConfig) tea.Cmd      // Initialize with current config
    Update(msg tea.Msg) (Step, StepResult, tea.Cmd)  // Handle input
    View() string                        // Render the step
    Title() string                       // Step title for header
    Apply(config *DBConfig)              // Save changes to config
    ShouldSkip(config *DBConfig) bool    // Conditional display
}
```

### Adding a New Step

1. Create a new file in `internal/steps/`
2. Implement the `Step` interface
3. Add the step to the wizard in `main.go`

### Conditional Steps

Steps can be skipped based on configuration. For example:
- Network, Data Vault, and Management steps only show in **Advanced** mode
- PDB-related fields only show when **Container Database** is enabled

## Coding Conventions

- Use the `ui` package for consistent styling
- All user input validation happens in step's `validate()` method
- Use `tea.Cmd` for async operations (though this app doesn't need them)
- Keep steps focused - one concern per step

## DBCA Parameters Reference

Key Oracle 19c DBCA parameters used:

| Parameter | Description |
|-----------|-------------|
| `-templateName` | Database template (.dbt file) |
| `-gdbname` | Global database name |
| `-sid` | Oracle SID (max 12 chars) |
| `-createAsContainerDatabase` | Enable CDB architecture |
| `-numberOfPDBs` | Number of pluggable databases |
| `-pdbName` | PDB name/prefix |
| `-storageType` | FS or ASM |
| `-datafileDestination` | Data file location |
| `-recoveryAreaDestination` | FRA location |
| `-totalMemory` | Total memory in MB |
| `-characterSet` | Database character set |
| `-emConfiguration` | Enterprise Manager setup |

## Testing

Manual testing workflow:
1. Run through Typical mode - verify fewer steps shown
2. Run through Advanced mode - verify all steps shown
3. Test CDB vs non-CDB paths
4. Test FS vs ASM storage paths
5. Verify generated command syntax against Oracle documentation

## Future Enhancements

Potential improvements:
- Add support for Oracle 21c/23ai parameters
- Add response file generation (in addition to command)
- Add database deletion command generation
- Add configuration templates (save/load)
- Add input validation against Oracle naming rules
