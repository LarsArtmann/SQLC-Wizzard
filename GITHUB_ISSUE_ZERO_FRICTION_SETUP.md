# 🚀 Zero-Friction Setup: One-Command SQLC Project Creation

**Priority:** HIGH (#1)  
**Complexity:** HIGH  
**Estimated Time:** 3-4 days  
**Impact:** CRITICAL - Core user experience

---

## 🎯 **Problem Statement**

Currently, creating a SQLC project requires:
- Manual sqlc.yaml configuration
- Complex database setup
- Schema file creation
- Migration initialization
- Project structure setup

**This creates significant friction and makes SQLC adoption difficult for new users.**

---

## 🎪 **Solution Vision**

Implement magical one-command project setup:

```bash
# Create perfect SQLC project instantly
sqlc-wizard create my-service --type microservice --database postgresql

# Creates complete project with:
✅ Optimized sqlc.yaml configuration
✅ Database schema templates  
✅ Migration files structure
✅ Go package structure
✅ Development scripts
✅ Docker configuration
✅ Makefile with common tasks
✅ README with usage instructions
```

---

## 🏗️ **Technical Implementation**

### **Phase 1: Core Template System**
```go
type ProjectTemplate struct {
    Name        string
    Type        ProjectType
    Database    DatabaseType
    Config      SqlcConfig
    Structure   DirectoryStructure
}

type ProjectCreator struct {
    templates map[ProjectType]*ProjectTemplate
    fs        FileSystemAdapter
    config    ConfigAdapter
}
```

### **Phase 2: Intelligent Defaults**
```go
// Generate optimized sqlc.yaml based on project type
func (pc *ProjectCreator) GenerateConfig(ptype ProjectType, db DatabaseType) (*SqlcConfig, error)

// Create project-specific directory structure
func (pc *ProjectCreator) CreateStructure(basePath string, template *ProjectTemplate) error

// Initialize migration system
func (pc *ProjectCreator) SetupMigrations(basePath string, db DatabaseType) error
```

### **Phase 3: Template Library**
```bash
sqlc-wizard create my-api --template restful
sqlc-wizard create my-lib --template library  
sqlc-wizard create my-service --template microservice
sqlc-wizard create my-webapp --template fullstack
```

---

## 📁 **Directory Structure Generated**

```
my-service/
├── sqlc.yaml                 # Optimized configuration
├── db/
│   ├── schema/               # Database schemas
│   │   ├── 001_users.sql
│   │   ├── 002_posts.sql
│   │   └── 003_comments.sql
│   └── migrations/          # Migration files
│       ├── 001_create_users.up.sql
│       ├── 001_create_users.down.sql
│       └── 002_create_posts.up.sql
├── internal/
│   └── db/                 # Generated Go code
│       ├── db.go
│       ├── models.go
│       └── queries/
├── cmd/
│   └── server/              # Application entry point
├── pkg/
│   └── config/             # Application config
├── scripts/                # Development scripts
│   ├── dev.sh
│   ├── migrate.sh
│   └── seed.sh
├── docker-compose.yml      # Development environment
├── Dockerfile            # Production container
├── Makefile             # Common tasks
└── README.md           # Documentation
```

---

## 🧪 **Example Usage Scenarios**

### **Scenario 1: Microservice API**
```bash
sqlc-wizard create user-service --type microservice --database postgresql --include auth

# Creates:
✅ Microservice-optimized sqlc.yaml
✅ PostgreSQL-specific schemas
✅ JWT authentication setup
✅ API endpoint templates
✅ Docker environment
```

### **Scenario 2: Library/Package**
```bash
sqlc-wizard create my-db-lib --type library --database sqlite

# Creates:
✅ Library-focused sqlc.yaml
✅ SQLite database with test data
✅ Public API structure
✅ Testing setup
✅ Documentation templates
```

### **Scenario 3: Full-Stack Application**
```bash
sqlc-wizard create my-webapp --type fullstack --database postgresql --include frontend

# Creates:
✅ Full-stack sqlc.yaml
✅ Backend API structure
✅ Frontend integration setup
✅ Authentication system
✅ Deployment configuration
```

---

## 🎯 **Acceptance Criteria**

### **Core Functionality**
- [ ] `sqlc-wizard create <name>` works without additional parameters
- [ ] `--type` flag supports: microservice, library, fullstack, api
- [ ] `--database` flag supports: postgresql, mysql, sqlite
- [ ] Generated `sqlc.yaml` is valid and optimized
- [ ] All commands in generated Makefile work

### **Template Quality**
- [ ] Generated schemas follow SQLC best practices
- [ ] Directory structure follows Go conventions
- [ ] All generated files compile and work
- [ ] Docker configuration runs out-of-the-box
- [ ] README has clear setup and usage instructions

### **User Experience**
- [ ] Setup takes < 30 seconds from command to working project
- [ ] No manual configuration required
- [ ] Generated project passes `sqlc validate`
- [ ] All generated tests pass
- [ ] User can run project with single command

---

## 🛠️ **Implementation Plan**

### **Day 1: Template Engine**
- [ ] Create `ProjectCreator` core structure
- [ ] Implement template loading system
- [ ] Add basic microservice template
- [ ] Create directory structure generator

### **Day 2: Smart Configuration**
- [ ] Implement sqlc.yaml generation logic
- [ ] Add database-specific optimizations
- [ ] Create validation system
- [ ] Add project type customizations

### **Day 3: Ecosystem Integration**
- [ ] Add Docker configuration templates
- [ ] Create Makefile generation
- [ ] Implement development script generation
- [ ] Add README template system

### **Day 4: Polish & Testing**
- [ ] Add comprehensive templates for all project types
- [ ] Create integration tests for all scenarios
- [ ] Add help system and documentation
- [ ] Performance optimization and error handling

---

## 🔧 **Technical Architecture**

### **New Components**
```go
// internal/commands/create.go
func NewCreateCommand() *cobra.Command

// internal/templates/project/
type ProjectTemplate interface {
    Name() string
    Generate(basePath string, config CreateConfig) error
}

// internal/creators/
type ProjectCreator struct {
    templateEngine *TemplateEngine
    fsAdapter     FileSystemAdapter
}
```

### **Integration Points**
- Extend existing `internal/templates/` system
- Use current `internal/adapters/` for file operations
- Leverage `generated/` types for validation
- Add to main CLI command registry

---

## 📊 **Success Metrics**

### **User Experience**
- [ ] Time to working project < 30 seconds
- [ ] Zero manual configuration required
- [ ] Generated project success rate > 95%

### **Technical Quality**
- [ ] All generated code compiles and works
- [ ] sqlc validate passes for all templates
- [ ] Comprehensive test coverage (> 90%)

### **Adoption**
- [ ] Users can create and run projects without SQLC knowledge
- [ ] Generated projects are production-ready
- [ ] Community feedback rate > 4.5/5

---

## 🎯 **Why HIGH PRIORITY (#1)**

This feature is **CRITICAL** for SQLC adoption because:

1. **Reduces Learning Curve:** New users can start without SQLC knowledge
2. **Eliminates Setup Friction:** No manual configuration required
3. **Ensures Best Practices:** Generated configs follow SQLC conventions
4. **Accelerates Development:** Users start productive immediately
5. **Improves User Retention:** Better first-time experience

**This is the most impactful feature we can implement for SQLC-Wizard adoption.**

---

## 📋 **Definition of Done**

- [ ] `sqlc-wizard create` command implemented and working
- [ ] All project types (microservice, library, fullstack, api) supported
- [ ] All databases (postgresql, mysql, sqlite) supported
- [ ] Generated projects compile and run successfully
- [ ] Comprehensive test coverage (> 90%)
- [ ] Documentation with examples and screenshots
- [ ] User guide with troubleshooting section

---

**This issue will transform SQLC-Wizard into a magical project setup tool that makes SQLC accessible to everyone!** 🚀✨

---

*Created: 2025-11-15*  
*Priority: HIGH (#1)*  
*Ready for implementation* 🎯