package generator

import (
	"fmt"
	"strings"

	"dbca_tui/internal/model"
)

// GenerateCommand generates the DBCA silent mode command
func GenerateCommand(config *model.DBConfig) string {
	var args []string

	args = append(args, "dbca", "-silent", "-createDatabase")

	// Template
	if config.TemplateName != model.TemplateCustom {
		args = append(args, fmt.Sprintf("-templateName %s", config.TemplateName))
	}

	// Database identification
	args = append(args, fmt.Sprintf("-gdbname %s", config.GlobalDBName))
	args = append(args, fmt.Sprintf("-sid %s", config.SID))

	// Container database settings
	if config.CreateAsContainerDB {
		args = append(args, "-createAsContainerDatabase true")
		if config.NumberOfPDBs > 0 {
			args = append(args, fmt.Sprintf("-numberOfPDBs %d", config.NumberOfPDBs))
			args = append(args, fmt.Sprintf("-pdbName %s", config.PDBName))
			args = append(args, fmt.Sprintf("-pdbAdminPassword '%s'", maskPassword(config.PDBAdminPassword)))
		}
	} else {
		args = append(args, "-createAsContainerDatabase false")
	}

	// Passwords
	args = append(args, fmt.Sprintf("-sysPassword '%s'", maskPassword(config.SysPassword)))
	args = append(args, fmt.Sprintf("-systemPassword '%s'", maskPassword(config.SystemPassword)))

	// Character set
	args = append(args, fmt.Sprintf("-characterSet %s", config.CharacterSet))
	args = append(args, fmt.Sprintf("-nationalCharacterSet %s", config.NationalCharacterSet))

	// Memory configuration
	if config.MemoryManagement == "AUTO" {
		args = append(args, fmt.Sprintf("-totalMemory %d", config.TotalMemory))
		args = append(args, "-memoryMgmtType AUTO")
	} else if config.MemoryManagement == "AUTO_SGA" {
		args = append(args, fmt.Sprintf("-totalMemory %d", config.TotalMemory))
		args = append(args, "-memoryMgmtType AUTO_SGA")
	} else {
		args = append(args, fmt.Sprintf("-totalMemory %d", config.TotalMemory))
		args = append(args, "-memoryMgmtType CUSTOM")
	}

	// Database type
	args = append(args, fmt.Sprintf("-databaseType %s", config.DatabaseType))

	// Storage configuration
	args = append(args, fmt.Sprintf("-storageType %s", config.StorageType))

	if config.StorageType == model.StorageTypeASM {
		args = append(args, fmt.Sprintf("-diskGroupName %s", config.ASMDiskGroup))
	} else {
		args = append(args, fmt.Sprintf("-datafileDestination '%s'", config.DatafileDestination))
	}

	// Use OMF
	if config.UseOMF {
		args = append(args, "-useOMF true")
	}

	// Fast Recovery Area
	if config.EnableFRA {
		args = append(args, fmt.Sprintf("-recoveryAreaDestination '%s'", config.FRADestination))
		args = append(args, fmt.Sprintf("-recoveryAreaSize %d", config.FRASize))
	}

	// Redo log size
	if config.RedoLogFileSize > 0 {
		args = append(args, fmt.Sprintf("-redoLogFileSize %d", config.RedoLogFileSize))
	}

	// Listener configuration
	if config.ListenerName != "" && config.ListenerName != "LISTENER" {
		args = append(args, fmt.Sprintf("-listeners %s", config.ListenerName))
	}

	// Enterprise Manager configuration
	args = append(args, fmt.Sprintf("-emConfiguration %s", config.EMConfiguration))
	if config.EMConfiguration == model.EMConfigDBExpress {
		args = append(args, fmt.Sprintf("-dbExpressPort %d", config.EMPort))
	}

	// Sample schemas
	if config.EnableSampleSchemas {
		args = append(args, "-sampleSchema true")
	}

	// Archive log mode
	if config.EnableArchiveLog {
		args = append(args, "-archiveLogMode true")
	}

	// Data Vault
	if config.EnableDataVault {
		args = append(args, "-enableDV true")
		args = append(args, fmt.Sprintf("-dvOwnerName %s", config.DataVaultOwner))
		args = append(args, fmt.Sprintf("-dvAccountManagerName %s", config.DataVaultAccountManager))
	}

	// RAC-specific options
	if config.DeploymentType == model.DeploymentRAC {
		args = append(args, "-databaseConfigType RAC")
		if config.NodeList != "" {
			args = append(args, fmt.Sprintf("-nodelist %s", config.NodeList))
		}
	} else if config.DeploymentType == model.DeploymentRACOneNode {
		args = append(args, "-databaseConfigType RACONENODE")
	} else {
		args = append(args, "-databaseConfigType SI")
	}

	// Ignore prerequisites (useful for testing)
	if config.IgnorePreReqs {
		args = append(args, "-ignorePreReqs")
	}

	return strings.Join(args, " \\\n  ")
}

// GenerateCommandWithPasswords generates the command with actual passwords
func GenerateCommandWithPasswords(config *model.DBConfig) string {
	var args []string

	args = append(args, "dbca", "-silent", "-createDatabase")

	// Template
	if config.TemplateName != model.TemplateCustom {
		args = append(args, fmt.Sprintf("-templateName %s", config.TemplateName))
	}

	// Database identification
	args = append(args, fmt.Sprintf("-gdbname %s", config.GlobalDBName))
	args = append(args, fmt.Sprintf("-sid %s", config.SID))

	// Container database settings
	if config.CreateAsContainerDB {
		args = append(args, "-createAsContainerDatabase true")
		if config.NumberOfPDBs > 0 {
			args = append(args, fmt.Sprintf("-numberOfPDBs %d", config.NumberOfPDBs))
			args = append(args, fmt.Sprintf("-pdbName %s", config.PDBName))
			args = append(args, fmt.Sprintf("-pdbAdminPassword '%s'", config.PDBAdminPassword))
		}
	} else {
		args = append(args, "-createAsContainerDatabase false")
	}

	// Passwords (actual values)
	args = append(args, fmt.Sprintf("-sysPassword '%s'", config.SysPassword))
	args = append(args, fmt.Sprintf("-systemPassword '%s'", config.SystemPassword))

	// Character set
	args = append(args, fmt.Sprintf("-characterSet %s", config.CharacterSet))
	args = append(args, fmt.Sprintf("-nationalCharacterSet %s", config.NationalCharacterSet))

	// Memory configuration
	if config.MemoryManagement == "AUTO" {
		args = append(args, fmt.Sprintf("-totalMemory %d", config.TotalMemory))
		args = append(args, "-memoryMgmtType AUTO")
	} else if config.MemoryManagement == "AUTO_SGA" {
		args = append(args, fmt.Sprintf("-totalMemory %d", config.TotalMemory))
		args = append(args, "-memoryMgmtType AUTO_SGA")
	} else {
		args = append(args, fmt.Sprintf("-totalMemory %d", config.TotalMemory))
		args = append(args, "-memoryMgmtType CUSTOM")
	}

	// Database type
	args = append(args, fmt.Sprintf("-databaseType %s", config.DatabaseType))

	// Storage configuration
	args = append(args, fmt.Sprintf("-storageType %s", config.StorageType))

	if config.StorageType == model.StorageTypeASM {
		args = append(args, fmt.Sprintf("-diskGroupName %s", config.ASMDiskGroup))
	} else {
		args = append(args, fmt.Sprintf("-datafileDestination '%s'", config.DatafileDestination))
	}

	// Use OMF
	if config.UseOMF {
		args = append(args, "-useOMF true")
	}

	// Fast Recovery Area
	if config.EnableFRA {
		args = append(args, fmt.Sprintf("-recoveryAreaDestination '%s'", config.FRADestination))
		args = append(args, fmt.Sprintf("-recoveryAreaSize %d", config.FRASize))
	}

	// Redo log size
	if config.RedoLogFileSize > 0 {
		args = append(args, fmt.Sprintf("-redoLogFileSize %d", config.RedoLogFileSize))
	}

	// Listener configuration
	if config.ListenerName != "" && config.ListenerName != "LISTENER" {
		args = append(args, fmt.Sprintf("-listeners %s", config.ListenerName))
	}

	// Enterprise Manager configuration
	args = append(args, fmt.Sprintf("-emConfiguration %s", config.EMConfiguration))
	if config.EMConfiguration == model.EMConfigDBExpress {
		args = append(args, fmt.Sprintf("-dbExpressPort %d", config.EMPort))
	}

	// Sample schemas
	if config.EnableSampleSchemas {
		args = append(args, "-sampleSchema true")
	}

	// Archive log mode
	if config.EnableArchiveLog {
		args = append(args, "-archiveLogMode true")
	}

	// Data Vault
	if config.EnableDataVault {
		args = append(args, "-enableDV true")
		args = append(args, fmt.Sprintf("-dvOwnerName %s", config.DataVaultOwner))
		args = append(args, fmt.Sprintf("-dvAccountManagerName %s", config.DataVaultAccountManager))
	}

	// RAC-specific options
	if config.DeploymentType == model.DeploymentRAC {
		args = append(args, "-databaseConfigType RAC")
		if config.NodeList != "" {
			args = append(args, fmt.Sprintf("-nodelist %s", config.NodeList))
		}
	} else if config.DeploymentType == model.DeploymentRACOneNode {
		args = append(args, "-databaseConfigType RACONENODE")
	} else {
		args = append(args, "-databaseConfigType SI")
	}

	// Ignore prerequisites (useful for testing)
	if config.IgnorePreReqs {
		args = append(args, "-ignorePreReqs")
	}

	return strings.Join(args, " \\\n  ")
}

func maskPassword(pwd string) string {
	if pwd == "" {
		return "<PASSWORD>"
	}
	return "<PASSWORD>"
}

// GenerateSummary generates a human-readable summary of the configuration
func GenerateSummary(config *model.DBConfig) string {
	var b strings.Builder

	b.WriteString("Database Configuration Summary\n")
	b.WriteString("==============================\n\n")

	// Creation mode
	if config.CreationMode == model.CreationModeTypical {
		b.WriteString("Creation Mode: Typical\n")
	} else {
		b.WriteString("Creation Mode: Advanced\n")
	}

	// Deployment type
	switch config.DeploymentType {
	case model.DeploymentSingleInstance:
		b.WriteString("Deployment: Single Instance\n")
	case model.DeploymentRAC:
		b.WriteString("Deployment: RAC\n")
	case model.DeploymentRACOneNode:
		b.WriteString("Deployment: RAC One Node\n")
	}

	b.WriteString("\n")

	// Database identification
	b.WriteString(fmt.Sprintf("Global Database Name: %s\n", config.GlobalDBName))
	b.WriteString(fmt.Sprintf("SID: %s\n", config.SID))

	if config.CreateAsContainerDB {
		b.WriteString("Container Database: Yes\n")
		b.WriteString(fmt.Sprintf("Number of PDBs: %d\n", config.NumberOfPDBs))
		if config.NumberOfPDBs > 0 {
			b.WriteString(fmt.Sprintf("PDB Name: %s\n", config.PDBName))
		}
	} else {
		b.WriteString("Container Database: No\n")
	}

	b.WriteString("\n")

	// Storage
	if config.StorageType == model.StorageTypeASM {
		b.WriteString(fmt.Sprintf("Storage Type: ASM (%s)\n", config.ASMDiskGroup))
	} else {
		b.WriteString(fmt.Sprintf("Storage Type: File System\n"))
		b.WriteString(fmt.Sprintf("Data Files: %s\n", config.DatafileDestination))
	}

	// FRA
	if config.EnableFRA {
		b.WriteString(fmt.Sprintf("Fast Recovery Area: %s (%d MB)\n", config.FRADestination, config.FRASize))
	}

	b.WriteString("\n")

	// Configuration
	b.WriteString(fmt.Sprintf("Total Memory: %d MB\n", config.TotalMemory))
	b.WriteString(fmt.Sprintf("Character Set: %s\n", config.CharacterSet))

	// Enterprise Manager
	if config.EMConfiguration != model.EMConfigNone {
		b.WriteString(fmt.Sprintf("Enterprise Manager: %s\n", config.EMConfiguration))
	}

	return b.String()
}
