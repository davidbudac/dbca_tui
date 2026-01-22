package generator

import (
	"fmt"
	"strings"

	"dbca_tui/internal/model"
)

// GenerateCommand generates the DBCA silent mode command (with masked passwords)
func GenerateCommand(config *model.DBConfig) string {
	if config.Operation == model.OperationDelete {
		return generateDeleteCommand(config, true)
	}
	return generateCreateCommand(config, true)
}

// GenerateCommandWithPasswords generates the command with actual passwords
func GenerateCommandWithPasswords(config *model.DBConfig) string {
	if config.Operation == model.OperationDelete {
		return generateDeleteCommand(config, false)
	}
	return generateCreateCommand(config, false)
}

// generateDeleteCommand generates the DBCA delete command
func generateDeleteCommand(config *model.DBConfig, maskPwd bool) string {
	var args []string

	args = append(args, "dbca", "-silent", "-deleteDatabase")

	// Database SID
	args = append(args, fmt.Sprintf("-sourceDB %s", config.DeleteSID))

	// SYS password
	pwd := config.SysPassword
	if maskPwd {
		pwd = "<PASSWORD>"
	}
	args = append(args, "-sysDBAUserName SYS")
	args = append(args, fmt.Sprintf("-sysDBAPassword '%s'", pwd))

	// Force delete option
	if config.DeleteForce {
		args = append(args, "-forceArchiveLogDeletion")
	}

	return strings.Join(args, " \\\n  ")
}

// generateCreateCommand generates the DBCA create command
func generateCreateCommand(config *model.DBConfig, maskPwd bool) string {
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
			pdbPwd := config.PDBAdminPassword
			if maskPwd {
				pdbPwd = "<PASSWORD>"
			}
			args = append(args, fmt.Sprintf("-pdbAdminPassword '%s'", pdbPwd))
		}
	} else {
		args = append(args, "-createAsContainerDatabase false")
	}

	// Passwords
	sysPwd := config.SysPassword
	systemPwd := config.SystemPassword
	if maskPwd {
		sysPwd = "<PASSWORD>"
		systemPwd = "<PASSWORD>"
	}
	args = append(args, fmt.Sprintf("-sysPassword '%s'", sysPwd))
	args = append(args, fmt.Sprintf("-systemPassword '%s'", systemPwd))

	// Character set
	args = append(args, fmt.Sprintf("-characterSet %s", config.CharacterSet))
	args = append(args, fmt.Sprintf("-nationalCharacterSet %s", config.NationalCharacterSet))

	// Memory configuration
	args = append(args, fmt.Sprintf("-totalMemory %d", config.TotalMemory))
	switch config.MemoryManagement {
	case "AUTO":
		args = append(args, "-memoryMgmtType AUTO")
	case "AUTO_SGA":
		args = append(args, "-memoryMgmtType AUTO_SGA")
	default:
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
	switch config.DeploymentType {
	case model.DeploymentRAC:
		args = append(args, "-databaseConfigType RAC")
		if config.NodeList != "" {
			args = append(args, fmt.Sprintf("-nodelist %s", config.NodeList))
		}
	case model.DeploymentRACOneNode:
		args = append(args, "-databaseConfigType RACONENODE")
	default:
		args = append(args, "-databaseConfigType SI")
	}

	// Ignore prerequisites
	if config.IgnorePreReqs {
		args = append(args, "-ignorePreReqs")
	}

	return strings.Join(args, " \\\n  ")
}
