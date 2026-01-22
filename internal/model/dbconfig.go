package model

// Operation represents the DBCA operation type
type Operation string

const (
	OperationCreate Operation = "create"
	OperationDelete Operation = "delete"
)

// CreationMode represents the database creation mode
type CreationMode string

const (
	CreationModeTypical  CreationMode = "typical"
	CreationModeAdvanced CreationMode = "advanced"
)

// DeploymentType represents the database deployment type
type DeploymentType string

const (
	DeploymentSingleInstance DeploymentType = "SI"
	DeploymentRAC            DeploymentType = "RAC"
	DeploymentRACOneNode     DeploymentType = "RACONENODE"
)

// DatabaseTemplate represents the database template
type DatabaseTemplate string

const (
	TemplateGeneralPurpose DatabaseTemplate = "General_Purpose.dbt"
	TemplateDataWarehouse  DatabaseTemplate = "Data_Warehouse.dbt"
	TemplateCustom         DatabaseTemplate = "Custom"
)

// StorageType represents the storage type
type StorageType string

const (
	StorageTypeFS  StorageType = "FS"
	StorageTypeASM StorageType = "ASM"
)

// EMConfiguration represents Enterprise Manager configuration
type EMConfiguration string

const (
	EMConfigNone     EMConfiguration = "NONE"
	EMConfigDBExpress EMConfiguration = "DBEXPRESS"
	EMConfigCentral  EMConfiguration = "CENTRAL"
)

// DatabaseType represents the database workload type
type DatabaseType string

const (
	DatabaseTypeMultipurpose DatabaseType = "MULTIPURPOSE"
	DatabaseTypeDataWarehouse DatabaseType = "DATA_WAREHOUSING"
	DatabaseTypeOLTP         DatabaseType = "OLTP"
)

// DBConfig holds all database configuration options
type DBConfig struct {
	// Operation type
	Operation Operation

	// Step 1: Creation Mode (for create operation)
	CreationMode CreationMode

	// Step 2: Deployment Type
	DeploymentType DeploymentType
	NodeList       string // Comma-separated list for RAC

	// Step 3: Template
	TemplateName DatabaseTemplate
	DatabaseType DatabaseType

	// Step 4: Database Identification
	GlobalDBName          string
	SID                   string
	CreateAsContainerDB   bool
	NumberOfPDBs          int
	PDBName               string
	PDBPrefix             string

	// Step 5: Storage
	StorageType              StorageType
	DatafileDestination      string
	RedoLogDestination       string
	ASMDiskGroup             string
	UseOMF                   bool // Oracle Managed Files

	// Step 6: Fast Recovery Area
	EnableFRA                bool
	FRADestination           string
	FRASize                  int // In MB
	EnableArchiveLog         bool

	// Step 7: Network
	ListenerName             string
	ListenerPort             int
	CreateNewListener        bool

	// Step 8: Data Vault (Advanced only)
	EnableDataVault          bool
	DataVaultOwner           string
	DataVaultAccountManager  string

	// Step 9: Configuration Options
	MemoryManagement         string // AUTO_SGA, MANUAL, AUTO
	TotalMemory              int    // In MB
	SGASize                  int    // In MB
	PGASize                  int    // In MB
	CharacterSet             string
	NationalCharacterSet     string
	ConnectionMode           string // DEDICATED, SHARED
	EnableSampleSchemas      bool

	// Step 10: Management Options
	EMConfiguration          EMConfiguration
	EMPort                   int
	CloudControlAgent        string

	// Step 11: Credentials
	UseCommonPassword        bool
	CommonPassword           string
	SysPassword              string
	SystemPassword           string
	PDBAdminPassword         string

	// Additional Options
	RedoLogFileSize          int  // In MB
	IgnorePreReqs            bool
	InitParams               map[string]string

	// Delete Operation Options
	DeleteSID                string
	DeleteForce              bool  // Force delete even if database is running
	DeleteExpressMode        bool  // Express mode (no prompts)
}

// NewDBConfig creates a new DBConfig with sensible defaults
func NewDBConfig() *DBConfig {
	return &DBConfig{
		Operation:            OperationCreate,
		CreationMode:         CreationModeTypical,
		DeploymentType:       DeploymentSingleInstance,
		TemplateName:         TemplateGeneralPurpose,
		DatabaseType:         DatabaseTypeMultipurpose,
		GlobalDBName:         "orcl",
		SID:                  "orcl",
		CreateAsContainerDB:  true,
		NumberOfPDBs:         1,
		PDBName:              "orclpdb",
		StorageType:          StorageTypeFS,
		DatafileDestination:  "/u01/app/oracle/oradata",
		UseOMF:               true,
		EnableFRA:            true,
		FRADestination:       "/u01/app/oracle/fast_recovery_area",
		FRASize:              10240,
		EnableArchiveLog:     false,
		ListenerName:         "LISTENER",
		ListenerPort:         1521,
		CreateNewListener:    false,
		EnableDataVault:      false,
		MemoryManagement:     "AUTO",
		TotalMemory:          2048,
		CharacterSet:         "AL32UTF8",
		NationalCharacterSet: "AL16UTF16",
		ConnectionMode:       "DEDICATED",
		EnableSampleSchemas:  false,
		EMConfiguration:      EMConfigNone,
		EMPort:               5500,
		UseCommonPassword:    true,
		RedoLogFileSize:      50,
		IgnorePreReqs:        false,
		InitParams:           make(map[string]string),
	}
}
