---
name: api-generator
description: Specialized skill for developing and refactoring the DST (Don't Starve Together) Admin Go project. Use this skill whenever the user mentions DST, game server management, API refactoring, adding handlers/services, creating CRUD endpoints, or working within the dst-admin-go codebase. This skill ensures all code follows the project's architectural patterns including dependency injection, layered architecture (Handler → Service → Model), and specific naming conventions. Make sure to use this skill when the user asks to refactor existing APIs, add new features, fix bugs, or modify any part of the dst-admin-go project structure.
license: MIT
metadata:
  author: DST Admin Go Team
  version: "2.0.0"
  domain: application
  triggers: DST, dst-admin-go, Don't Starve Together, game server, refactor API, add handler, create service, GORM, Gin, cluster management, backup service, player management, mod management, CRUD generator, generate API, create endpoint, 生成API, 创建接口
  role: specialist
  scope: implementation
  output-format: code
  related-skills: golang-pro, golang-patterns
---

# DST Admin Go API Generator

Specialized skill for developing and refactoring the DST (Don't Starve Together) Admin Go project. Use this skill whenever the user mentions DST, game server management, API refactoring, adding handlers/services, creating CRUD endpoints, or working within the dst-admin-go codebase. This skill ensures all code follows the project's architectural patterns including dependency injection, layered architecture (Handler → Service → Model), and specific naming conventions.

---

## Project Context

DST Admin Go is a web-based management panel for "Don't Starve Together" game servers written in Go. It follows a three-layer architecture:

1. **Handler Layer** (`internal/api/handler/`): HTTP request handling, validation, response formatting
2. **Service Layer** (`internal/service/`): Business logic, orchestration
3. **Model Layer** (`internal/model/`): Database entities (GORM)

### Architecture Patterns

- **Dependency Injection**: All services instantiated in `internal/api/router.go`, injected via constructors
- **No Global Variables**: Pass config and dependencies through constructors
- **Platform Abstraction**: Factory patterns for Linux/Windows differences
- **Naming Conventions**:
  - Files: `snake_case.go` (e.g., `backup_service.go`)
  - Structs: `PascalCase` without suffix (e.g., `BackupInfo`, NOT `BackupInfoVO`)
  - Constructors: `New{Name}` (e.g., `NewBackupService`)
  - Interfaces: Simple names (e.g., `Config`, `Process`)

### Project Structure

```
dst-admin-go/
├── cmd/server/main.go           # Entry point
├── internal/
│   ├── api/
│   │   ├── handler/             # HTTP handlers (controllers)
│   │   │   └── {entity}_handler.go
│   │   └── router.go            # Route registration & DI
│   ├── model/                   # GORM models
│   │   └── {entity}.go
│   ├── service/                 # Business logic
│   │   └── {domain}/
│   │       ├── {domain}_service.go
│   │       ├── factory.go       # (if platform-specific)
│   │       ├── linux_{domain}.go
│   │       └── window_{domain}.go
│   └── pkg/
│       ├── response/            # Standard responses
│       └── utils/               # Utilities
└── config.yml                   # Application config
```

---

## Your Instructions

You are an expert Go developer specializing in the DST Admin Go project. When the user requests creating a CRUD module, refactoring APIs, or adding new functionality:

### Step 1: Gather Requirements

Ask the user for the following information (use a friendly, concise tone):

1. **Entity/Model Name**: What is the entity called? (e.g., "Announcement", "ModConfig", "Player")
2. **Chinese Name**: What is the Chinese name for Swagger docs? (e.g., "公告", "模组配置")
3. **Database Fields**: What fields does this entity have?
   - Field name, type (string, int, bool, time.Time, etc.)
   - GORM constraints (e.g., `unique`, `not null`, `default`)
   - JSON tag name (camelCase for API responses)
   - Example: `Title string, required, JSON: "title"` or `IsActive bool, default true, JSON: "isActive"`
4. **Operations Needed**: Which operations should be available?
   - Standard CRUD: Create, Read (Get by ID), Update, Delete, List (with pagination)
   - Custom operations: Batch delete, toggle status, search, etc.
5. **Business Logic Notes**: Any special validation, processing, or relationships?
   - E.g., "Must validate expiresAt is in the future", "Needs to interact with game process"
6. **Cluster Context**: Does this entity belong to a specific cluster/server?
   - If yes: Endpoints will be under `/api/{clusterName}/{entity}` and need cluster middleware

### Step 2: Analyze Dependencies

Based on the requirements, determine:

- **Always needed**: `*gorm.DB` for database operations
- **Game-related**: If interacting with game state → `game.Process`
- **File operations**: If handling files/archives → `archive.PathResolver`
- **DST config**: If reading/writing DST configs → `dstConfig.Config`
- **Level config**: If working with world configs → `levelConfig.LevelConfigUtils`
- **Platform-specific**: If operations differ on Linux vs Windows → factory pattern

### Step 3: Generate Model File

Create `internal/model/{entity}.go`:

```go
package model

import (
	"gorm.io/gorm"
	"time"
)

// {EntityName} {中文描述}
type {EntityName} struct {
	gorm.Model
	// Fields based on user input
	Title       string     `json:"title" gorm:"type:varchar(255);not null"`
	Content     string     `json:"content" gorm:"type:text"`
	IsActive    bool       `json:"isActive" gorm:"default:true"`
	ExpiresAt   *time.Time `json:"expiresAt"`
}
```

**Guidelines**:
- Always embed `gorm.Model` (provides ID, CreatedAt, UpdatedAt, DeletedAt)
- Use GORM tags for constraints: `type:`, `not null`, `unique`, `default:`
- Use JSON tags in camelCase
- Add struct comment with Chinese description
- Use pointer types for nullable fields (e.g., `*time.Time`, `*string`)

### Step 4: Generate Service File

Create `internal/service/{domain}/{domain}_service.go`:

```go
package {domain}

import (
	"dst-admin-go/internal/model"
	"gorm.io/gorm"
)

type {Entity}Service struct {
	db *gorm.DB
	// Add other dependencies as detected
}

func New{Entity}Service(db *gorm.DB /* other deps */) *{Entity}Service {
	return &{Entity}Service{
		db: db,
	}
}

// List{Entity} 获取{中文名}列表
func (s *{Entity}Service) List{Entity}(page, pageSize int) ([]model.{Entity}, int64, error) {
	var list []model.{Entity}
	var total int64

	offset := (page - 1) * pageSize

	err := s.db.Model(&model.{Entity}{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = s.db.Offset(offset).Limit(pageSize).Find(&list).Error
	return list, total, err
}

// Get{Entity} 获取{中文名}详情
func (s *{Entity}Service) Get{Entity}(id uint) (*model.{Entity}, error) {
	var entity model.{Entity}
	err := s.db.First(&entity, id).Error
	return &entity, err
}

// Create{Entity} 创建{中文名}
func (s *{Entity}Service) Create{Entity}(entity *model.{Entity}) error {
	return s.db.Create(entity).Error
}

// Update{Entity} 更新{中文名}
func (s *{Entity}Service) Update{Entity}(entity *model.{Entity}) error {
	return s.db.Save(entity).Error
}

// Delete{Entity} 删除{中文名}
func (s *{Entity}Service) Delete{Entity}(id uint) error {
	return s.db.Delete(&model.{Entity}{}, id).Error
}
```

**Guidelines**:
- Constructor accepts all dependencies
- Chinese comments for all exported methods
- Use GORM for database operations
- Add pagination support for List (offset/limit)
- Return errors, don't handle them here
- Add custom business logic methods as needed

### Step 5: Generate Handler File

Create `internal/api/handler/{entity}_handler.go`:

```go
package handler

import (
	"dst-admin-go/internal/model"
	"dst-admin-go/internal/pkg/response"
	"dst-admin-go/internal/service/{domain}"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type {Entity}Handler struct {
	service *{domain}.{Entity}Service
}

func New{Entity}Handler(service *{domain}.{Entity}Service) *{Entity}Handler {
	return &{Entity}Handler{service: service}
}

func (h *{Entity}Handler) RegisterRoute(router *gin.RouterGroup) {
	router.GET("/api/{entity}/list", h.List)
	router.GET("/api/{entity}/:id", h.Get)
	router.POST("/api/{entity}", h.Create)
	router.PUT("/api/{entity}/:id", h.Update)
	router.DELETE("/api/{entity}/:id", h.Delete)
}

// List 获取{中文名}列表
// @Summary 获取{中文名}列表
// @Description 分页获取{中文名}列表
// @Tags {entity}
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} response.Response{data=object{list=[]model.{Entity},total=int,page=int,pageSize=int}}
// @Router /api/{entity}/list [get]
func (h *{Entity}Handler) List(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	list, total, err := h.service.List{Entity}(page, pageSize)
	if err != nil {
		response.FailWithMessage("获取列表失败: "+err.Error(), ctx)
		return
	}

	response.OkWithData(gin.H{
		"list":     list,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	}, ctx)
}

// Get 获取{中文名}详情
// @Summary 获取{中文名}详情
// @Description 根据ID获取{中文名}详情
// @Tags {entity}
// @Accept json
// @Produce json
// @Param id path int true "{中文名}ID"
// @Success 200 {object} response.Response{data=model.{Entity}}
// @Router /api/{entity}/{id} [get]
func (h *{Entity}Handler) Get(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.FailWithMessage("无效的ID", ctx)
		return
	}

	entity, err := h.service.Get{Entity}(uint(id))
	if err != nil {
		response.FailWithMessage("获取详情失败: "+err.Error(), ctx)
		return
	}

	response.OkWithData(entity, ctx)
}

// Create 创建{中文名}
// @Summary 创建{中文名}
// @Description 创建新的{中文名}
// @Tags {entity}
// @Accept json
// @Produce json
// @Param data body model.{Entity} true "{中文名}信息"
// @Success 200 {object} response.Response{data=model.{Entity}}
// @Router /api/{entity} [post]
func (h *{Entity}Handler) Create(ctx *gin.Context) {
	var entity model.{Entity}
	if err := ctx.ShouldBindJSON(&entity); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), ctx)
		return
	}

	if err := h.service.Create{Entity}(&entity); err != nil {
		response.FailWithMessage("创建失败: "+err.Error(), ctx)
		return
	}

	response.OkWithData(entity, ctx)
}

// Update 更新{中文名}
// @Summary 更新{中文名}
// @Description 更新{中文名}信息
// @Tags {entity}
// @Accept json
// @Produce json
// @Param id path int true "{中文名}ID"
// @Param data body model.{Entity} true "{中文名}信息"
// @Success 200 {object} response.Response{data=model.{Entity}}
// @Router /api/{entity}/{id} [put]
func (h *{Entity}Handler) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.FailWithMessage("无效的ID", ctx)
		return
	}

	var entity model.{Entity}
	if err := ctx.ShouldBindJSON(&entity); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), ctx)
		return
	}

	entity.ID = uint(id)
	if err := h.service.Update{Entity}(&entity); err != nil {
		response.FailWithMessage("更新失败: "+err.Error(), ctx)
		return
	}

	response.OkWithData(entity, ctx)
}

// Delete 删除{中文名}
// @Summary 删除{中文名}
// @Description 根据ID删除{中文名}
// @Tags {entity}
// @Accept json
// @Produce json
// @Param id path int true "{中文名}ID"
// @Success 200 {object} response.Response
// @Router /api/{entity}/{id} [delete]
func (h *{Entity}Handler) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.FailWithMessage("无效的ID", ctx)
		return
	}

	if err := h.service.Delete{Entity}(uint(id)); err != nil {
		response.FailWithMessage("删除失败: "+err.Error(), ctx)
		return
	}

	response.OkWithMessage("删除成功", ctx)
}
```

**Guidelines**:
- Complete Swagger annotations for every handler
- Use `response.Response` helpers: `OkWithData`, `FailWithMessage`, `OkWithMessage`
- Validate all input parameters
- Chinese error messages
- Parse ID from path parameter as uint
- Bind JSON request body with `ShouldBindJSON`
- Return HTTP 200 even for errors (error code in response body)

### Step 6: Update Router

Read `internal/api/router.go` and make these changes:

1. **Add import** (if not exists):
```go
"{domain}" "dst-admin-go/internal/service/{domain}"
```

2. **In `Register()` function, add service initialization** (after existing services):
```go
// {entity} service
{entity}Service := {domain}.New{Entity}Service(db /* add detected dependencies */)
```

3. **Add handler initialization** (after existing handlers):
```go
{entity}Handler := handler.New{Entity}Handler({entity}Service)
```

4. **Add route registration** (after existing routes):
```go
{entity}Handler.RegisterRoute(router)
```

**Important**: Maintain the existing order and formatting. Add new code at the end of each section.

### Step 7: Handle Platform-Specific Code (if needed)

If the service needs platform-specific behavior:

1. Create `internal/service/{domain}/factory.go`:
```go
package {domain}

import (
	"gorm.io/gorm"
	"runtime"
)

func New{Entity}Service(db *gorm.DB) {Entity}Service {
	if runtime.GOOS == "windows" {
		return &Windows{Entity}Service{db: db}
	}
	return &Linux{Entity}Service{db: db}
}
```

2. Create interface in `{domain}_service.go`:
```go
type {Entity}Service interface {
	List{Entity}(page, pageSize int) ([]model.{Entity}, int64, error)
	// ... other methods
}
```

3. Create platform implementations:
   - `internal/service/{domain}/linux_{domain}.go`
   - `internal/service/{domain}/window_{domain}.go`

### Step 8: Verify and Test

After generating all files:

1. **Run compilation check**:
```bash
go mod tidy
go build cmd/server/main.go
```

2. **Report to user**:
   - List all generated files
   - Show modified files (router.go)
   - Provide curl test commands
   - Mention Swagger UI location

3. **Provide test commands**:
```bash
# List
curl -X GET "http://localhost:8082/api/{entity}/list?page=1&pageSize=10"

# Create
curl -X POST "http://localhost:8082/api/{entity}" \
  -H "Content-Type: application/json" \
  -d '{"field1": "value1", "field2": "value2"}'

# Get
curl -X GET "http://localhost:8082/api/{entity}/1"

# Update
curl -X PUT "http://localhost:8082/api/{entity}/1" \
  -H "Content-Type: application/json" \
  -d '{"field1": "new_value"}'

# Delete
curl -X DELETE "http://localhost:8082/api/{entity}/1"
```

4. **Remind user**: Access Swagger UI at `http://localhost:8082/swagger/index.html` (after running the server)

---

## Common Patterns Reference

### Response Helpers

Located in `internal/pkg/response/response.go`:

```go
// Success with data
response.OkWithData(data, ctx)

// Success with message
response.OkWithMessage("操作成功", ctx)

// Error with message
response.FailWithMessage("操作失败: "+err.Error(), ctx)
```

### Common Dependencies

- `*gorm.DB`: Database access
- `*gin.Context`: HTTP request context
- `dstConfig.Config`: DST configuration interface
- `game.Process`: Game process management interface
- `archive.PathResolver`: Archive path resolution
- `levelConfig.LevelConfigUtils`: Level config parsing

### Cluster-Aware Endpoints

If entity belongs to a cluster, routes should be:
```go
router.GET("/api/:cluster_name/{entity}/list", h.List)
router.GET("/api/:cluster_name/{entity}/:id", h.Get)
// etc.
```

**IMPORTANT**: Handler should use cluster context helper to get clusterName:
```go
import "dst-admin-go/internal/pkg/context"

clusterName := context.GetClusterName(ctx)
// Use clusterName in service calls
```

**DO NOT** use `ctx.Query("clusterName")` or `ctx.Param("cluster_name")` directly. Always use `context.GetClusterName(ctx)` which retrieves the clusterName set by the cluster middleware.

### Pagination Pattern

Always use this pattern for list endpoints:
```go
page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))

if page < 1 {
	page = 1
}
if pageSize < 1 || pageSize > 100 {
	pageSize = 10
}

offset := (page - 1) * pageSize
```

### GORM Common Tags

- `type:varchar(255)` - String with length
- `type:text` - Long text
- `not null` - Required field
- `unique` - Unique constraint
- `default:true` - Default value
- `index` - Create index
- `foreignKey:UserID` - Foreign key

---

## Examples

### Example 1: Simple Announcement System

**User Input**:
- Entity: Announcement
- Chinese: 公告
- Fields: Title (string, required), Content (text), IsActive (bool, default true), ExpiresAt (time, nullable)
- Operations: Full CRUD + List
- Cluster context: No

**Generated Files**:
- `internal/model/announce.go`
- `internal/service/announce/announce_service.go`
- `internal/api/handler/announce_handler.go`
- Modified: `internal/api/router.go`

### Example 2: Mod Management (Game-Related)

**User Input**:
- Entity: ModInfo
- Chinese: 模组信息
- Fields: ModId (string, unique), Name (string), Author (string), Version (string), IsEnabled (bool), ClusterName (string)
- Operations: Full CRUD + List + Toggle status
- Business logic: Needs to interact with game process when toggling
- Cluster context: Yes

**Generated Files**:
- `internal/model/modInfo.go`
- `internal/service/mod/mod_service.go` (with `game.Process` dependency)
- `internal/api/handler/mod_handler.go` (with cluster-aware routes)
- Modified: `internal/api/router.go`

**Additional method in service**:
```go
// ToggleModStatus 切换模组启用状态
func (s *ModService) ToggleModStatus(clusterName, modId string) error {
	// Update database
	// Interact with game process
	return nil
}
```

---

## Critical Checklist

Before completing the task, verify:

- [ ] Model file has `gorm.Model` embedded
- [ ] All fields have both `json` and `gorm` tags
- [ ] Service has constructor accepting all dependencies
- [ ] All service methods have Chinese comments
- [ ] Handler has `RegisterRoute` method
- [ ] All handler methods have complete Swagger annotations
- [ ] Swagger tags use entity name (e.g., `@Tags announcement`)
- [ ] Chinese names used in Swagger summaries
- [ ] router.go has import added
- [ ] router.go has service initialization
- [ ] router.go has handler initialization
- [ ] router.go has route registration
- [ ] Naming follows conventions (snake_case files, PascalCase types, no VO suffix)
- [ ] Code compiles without errors
- [ ] Test commands provided to user

---

## DST-Specific Domain Knowledge

### Cluster Architecture

A **cluster** contains multiple **levels** (worlds):
- **Master** - Overworld (surface world)
- **Caves** - Underground world

Each level runs as a separate process with its own configuration.

### Important Paths

```go
// Use pathResolver service for all path operations
pathResolver.ClusterPath(clusterName)          // e.g., ~/.klei/DoNotStarveTogether/MyCluster/
pathResolver.LevelPath(clusterName, "Master")  // e.g., ~/.klei/DoNotStarveTogether/MyCluster/Master/
pathResolver.SavePath(clusterName, "Master")   // e.g., ~/.klei/DoNotStarveTogether/MyCluster/Master/save/
```

### Configuration Files

- `cluster.ini` - Cluster settings (game mode, max players, passwords)
- `leveldataoverride.lua` - World generation settings (Lua table format)
- `modoverrides.lua` - Enabled mods configuration (Lua table format)
- `server.ini` - Server-specific settings

Use `luaUtils` package to parse/generate Lua configuration files.

### Process Management

**Linux**: Uses `screen` sessions
```go
screenName := fmt.Sprintf("dst_%s_%s", clusterName, levelName)
```

**Windows**: Uses custom CLI wrapper (`windowGameCli.go`)

---

## Common Utilities

### File Operations

```go
import "dst-admin-go/internal/pkg/utils/fileUtils"

fileUtils.PathExists(path)
fileUtils.CreateDir(path)
fileUtils.CopyFile(src, dst)
fileUtils.Unzip(zipPath, destPath)
```

### Shell Commands

```go
import "dst-admin-go/internal/pkg/utils/shellUtils"

output, err := shellUtils.ExecuteCommand("ls", "-la")
```

### Lua Configuration

```go
import "dst-admin-go/internal/pkg/utils/luaUtils"

// Parse Lua table to map
config, err := luaUtils.ParseLuaTable(luaContent)

// Generate Lua table from map
luaContent := luaUtils.GenerateLuaTable(configMap)
```

---

## Tips

- **Be concise**: Don't over-explain. Generate code efficiently.
- **Follow patterns**: Always reference existing code in the project for consistency.
- **Chinese comments**: Use Chinese for inline comments and method descriptions, English for Swagger and exported names.
- **Dependency detection**: Ask about business logic to determine dependencies accurately.
- **Platform awareness**: Ask if operations differ on Windows vs Linux.
- **Cluster context**: Ask if entity belongs to a specific cluster/server.
- **Validation**: Add appropriate validation in handlers (required fields, format checks).
- **Error messages**: Always include the original error in Chinese messages: "操作失败: " + err.Error()

---

## When NOT to Use This Skill

Don't use this skill if:
- User is asking general Go questions (not specific to DST Admin Go)
- User wants to modify frontend code
- User is asking about Docker, deployment, or infrastructure
- Task doesn't involve creating/modifying handlers, services, or models

In those cases, handle the request normally without the skill context.

---

## Summary

This skill automates the creation of complete CRUD modules in DST Admin Go following the project's three-layer architecture. It handles model generation, service creation with dependency injection, handler implementation with Swagger docs, and router integration. Always gather complete requirements first, analyze dependencies, generate code following established patterns, and verify compilation before reporting success to the user.
