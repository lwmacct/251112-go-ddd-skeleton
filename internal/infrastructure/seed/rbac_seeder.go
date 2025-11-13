package seed

import (
	"context"
	"embed"
	"fmt"
	"log"

	"gopkg.in/yaml.v3"
	"gorm.io/gorm"

	"github.com/lwmacct/251112-go-ddd-skeleton/internal/domain/rbac"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/infrastructure/persistence/mapper"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/infrastructure/persistence/model"
)

//go:embed data/*.yaml
var seedDataFS embed.FS

// RoleData 角色数据结构
type RoleData struct {
	Code        string `yaml:"code"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	IsActive    bool   `yaml:"is_active"`
}

// PermissionData 权限数据结构
type PermissionData struct {
	Code        string `yaml:"code"`
	Name        string `yaml:"name"`
	Resource    string `yaml:"resource"`
	Action      string `yaml:"action"`
	Description string `yaml:"description"`
}

// MenuData 菜单数据结构
type MenuData struct {
	FixedID     string  `yaml:"fixed_id"`
	Name        string  `yaml:"name"`
	Path        string  `yaml:"path"`
	Icon        string  `yaml:"icon"`
	Type        string  `yaml:"type"`
	ParentID    *string `yaml:"parent_id"`
	SortOrder   int     `yaml:"sort_order"`
	IsVisible   bool    `yaml:"is_visible"`
	Component   string  `yaml:"component"`
	Permission  string  `yaml:"permission"`
	Description string  `yaml:"description"`
}

// RBACSeeder RBAC seed 实现
type RBACSeeder struct{}

// NewRBACSeeder 创建 RBAC seeder
func NewRBACSeeder() *RBACSeeder {
	return &RBACSeeder{}
}

// Name 返回 seeder 名称
func (s *RBACSeeder) Name() string {
	return "RBAC_SEED"
}

// ShouldRun 检查是否应该执行
func (s *RBACSeeder) ShouldRun(ctx context.Context, db *gorm.DB) (bool, error) {
	// 检查 seed_history 表
	var count int64
	err := db.Model(&SeedHistory{}).
		Where("name = ? AND status = ?", s.Name(), "success").
		Count(&count).Error
	return count == 0, err
}

// Run 执行 seed
func (s *RBACSeeder) Run(ctx context.Context, db *gorm.DB) error {
	log.Println("  📦 Loading YAML data...")

	// 1. 加载 YAML 数据
	rolesData, err := s.loadRoles()
	if err != nil {
		return fmt.Errorf("failed to load roles: %w", err)
	}

	permissionsData, err := s.loadPermissions()
	if err != nil {
		return fmt.Errorf("failed to load permissions: %w", err)
	}

	menusData, err := s.loadMenus()
	if err != nil {
		return fmt.Errorf("failed to load menus: %w", err)
	}

	rolePermData, err := s.loadRolePermissions()
	if err != nil {
		return fmt.Errorf("failed to load role_permissions: %w", err)
	}

	roleMenuData, err := s.loadRoleMenus()
	if err != nil {
		return fmt.Errorf("failed to load role_menus: %w", err)
	}

	// 2. 插入角色
	log.Println("  ✓ Creating roles...")
	roleIDMap, err := s.insertRoles(db, rolesData)
	if err != nil {
		return err
	}

	// 3. 插入权限
	log.Println("  ✓ Creating permissions...")
	permIDMap, err := s.insertPermissions(db, permissionsData)
	if err != nil {
		return err
	}

	// 4. 插入菜单（两层结构）
	log.Println("  ✓ Creating menus...")
	menuIDMap, err := s.insertMenus(db, menusData)
	if err != nil {
		return err
	}

	// 5. 创建角色-权限关联
	log.Println("  ✓ Assigning permissions to roles...")
	if err := s.insertRolePermissions(db, roleIDMap, permIDMap, rolePermData); err != nil {
		return err
	}

	// 6. 创建角色-菜单关联
	log.Println("  ✓ Assigning menus to roles...")
	if err := s.insertRoleMenus(db, roleIDMap, menuIDMap, roleMenuData); err != nil {
		return err
	}

	log.Printf("  ✅ Created %d roles, %d permissions, %d menus",
		len(rolesData), len(permissionsData), len(menusData))

	return nil
}

// insertRoles 插入角色
func (s *RBACSeeder) insertRoles(db *gorm.DB, rolesData []RoleData) (map[string]string, error) {
	roleIDMap := make(map[string]string) // code -> ulid

	for _, roleData := range rolesData {
		// 检查是否已存在
		var existing model.Role
		err := db.Where("code = ?", roleData.Code).First(&existing).Error
		if err == nil {
			roleIDMap[roleData.Code] = existing.ID
			continue // 已存在，跳过
		}

		// 创建新角色
		role, err := rbac.NewRole(roleData.Name, roleData.Code, roleData.Description)
		if err != nil {
			return nil, fmt.Errorf("failed to create role %s: %w", roleData.Code, err)
		}
		role.ID = generateULID()
		role.IsActive = roleData.IsActive
		roleIDMap[roleData.Code] = role.ID

		roleModel := mapper.RoleToModel(role)
		if err := db.Create(roleModel).Error; err != nil {
			return nil, fmt.Errorf("failed to insert role %s: %w", roleData.Code, err)
		}
	}

	return roleIDMap, nil
}

// insertPermissions 插入权限
func (s *RBACSeeder) insertPermissions(db *gorm.DB, permissionsData []PermissionData) (map[string]string, error) {
	permIDMap := make(map[string]string) // code -> ulid

	for _, permData := range permissionsData {
		// 检查是否已存在
		var existing model.Permission
		err := db.Where("code = ?", permData.Code).First(&existing).Error
		if err == nil {
			permIDMap[permData.Code] = existing.ID
			continue
		}

		// 创建新权限
		perm, err := rbac.NewPermission(permData.Name, permData.Resource, permData.Action, permData.Description)
		if err != nil {
			return nil, fmt.Errorf("failed to create permission %s: %w", permData.Code, err)
		}
		perm.ID = generateULID()
		perm.Code = permData.Code
		permIDMap[permData.Code] = perm.ID

		permModel := mapper.PermissionToModel(perm)
		if err := db.Create(permModel).Error; err != nil {
			return nil, fmt.Errorf("failed to insert permission %s: %w", permData.Code, err)
		}
	}

	return permIDMap, nil
}

// insertMenus 插入菜单（两层结构）
func (s *RBACSeeder) insertMenus(db *gorm.DB, menusData []MenuData) (map[string]string, error) {
	menuIDMap := make(map[string]string) // fixed_id -> ulid

	// 第一遍：插入父菜单
	for _, menuData := range menusData {
		if menuData.ParentID != nil {
			continue // 跳过子菜单
		}

		// 检查是否已存在（通过 path 判断）
		var existing model.Menu
		err := db.Where("path = ?", menuData.Path).First(&existing).Error
		if err == nil {
			menuIDMap[menuData.FixedID] = existing.ID
			continue
		}

		// 创建新菜单
		menu, err := rbac.NewMenu(menuData.Name, menuData.Path, menuData.Icon, rbac.MenuType(menuData.Type), nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create menu %s: %w", menuData.Name, err)
		}
		menu.ID = generateULID()
		menu.SortOrder = menuData.SortOrder
		menu.IsVisible = menuData.IsVisible
		menu.Component = menuData.Component
		menu.Permission = menuData.Permission
		menu.Description = menuData.Description
		menuIDMap[menuData.FixedID] = menu.ID

		menuModel := mapper.MenuToModel(menu)
		if err := db.Create(menuModel).Error; err != nil {
			return nil, fmt.Errorf("failed to insert menu %s: %w", menuData.Name, err)
		}
	}

	// 第二遍：插入子菜单
	for _, menuData := range menusData {
		if menuData.ParentID == nil {
			continue // 跳过父菜单
		}

		// 检查是否已存在
		var existing model.Menu
		err := db.Where("path = ?", menuData.Path).First(&existing).Error
		if err == nil {
			menuIDMap[menuData.FixedID] = existing.ID
			continue
		}

		// 获取父菜单 ULID
		parentULID, ok := menuIDMap[*menuData.ParentID]
		if !ok {
			return nil, fmt.Errorf("parent menu not found: %s", *menuData.ParentID)
		}

		// 创建新菜单
		menu, err := rbac.NewMenu(menuData.Name, menuData.Path, menuData.Icon, rbac.MenuType(menuData.Type), &parentULID)
		if err != nil {
			return nil, fmt.Errorf("failed to create menu %s: %w", menuData.Name, err)
		}
		menu.ID = generateULID()
		menu.SortOrder = menuData.SortOrder
		menu.IsVisible = menuData.IsVisible
		menu.Component = menuData.Component
		menu.Permission = menuData.Permission
		menu.Description = menuData.Description
		menuIDMap[menuData.FixedID] = menu.ID

		menuModel := mapper.MenuToModel(menu)
		if err := db.Create(menuModel).Error; err != nil {
			return nil, fmt.Errorf("failed to insert menu %s: %w", menuData.Name, err)
		}
	}

	return menuIDMap, nil
}

// insertRolePermissions 插入角色-权限关联
func (s *RBACSeeder) insertRolePermissions(db *gorm.DB, roleIDMap, permIDMap map[string]string, data map[string][]string) error {
	for roleCode, permCodes := range data {
		roleID, ok := roleIDMap[roleCode]
		if !ok {
			return fmt.Errorf("role not found: %s", roleCode)
		}

		for _, permCode := range permCodes {
			permID, ok := permIDMap[permCode]
			if !ok {
				return fmt.Errorf("permission not found: %s", permCode)
			}

			// 检查是否已存在
			var count int64
			db.Model(&model.RolePermission{}).
				Where("role_id = ? AND permission_id = ?", roleID, permID).
				Count(&count)
			if count > 0 {
				continue
			}

			// 创建关联
			rp := model.RolePermission{
				RoleID:       roleID,
				PermissionID: permID,
			}
			if err := db.Create(&rp).Error; err != nil {
				return fmt.Errorf("failed to assign permission %s to role %s: %w", permCode, roleCode, err)
			}
		}
	}

	return nil
}

// insertRoleMenus 插入角色-菜单关联
func (s *RBACSeeder) insertRoleMenus(db *gorm.DB, roleIDMap, menuIDMap map[string]string, data map[string][]string) error {
	for roleCode, menuFixedIDs := range data {
		roleID, ok := roleIDMap[roleCode]
		if !ok {
			return fmt.Errorf("role not found: %s", roleCode)
		}

		for _, fixedID := range menuFixedIDs {
			menuID, ok := menuIDMap[fixedID]
			if !ok {
				return fmt.Errorf("menu not found: %s", fixedID)
			}

			// 检查是否已存在
			var count int64
			db.Model(&model.RoleMenu{}).
				Where("role_id = ? AND menu_id = ?", roleID, menuID).
				Count(&count)
			if count > 0 {
				continue
			}

			// 创建关联
			rm := model.RoleMenu{
				RoleID: roleID,
				MenuID: menuID,
			}
			if err := db.Create(&rm).Error; err != nil {
				return fmt.Errorf("failed to assign menu %s to role %s: %w", fixedID, roleCode, err)
			}
		}
	}

	return nil
}

// loadRoles 从 YAML 加载角色数据
func (s *RBACSeeder) loadRoles() ([]RoleData, error) {
	data, err := seedDataFS.ReadFile("data/roles.yaml")
	if err != nil {
		return nil, err
	}

	var result struct {
		Roles []RoleData `yaml:"roles"`
	}

	err = yaml.Unmarshal(data, &result)
	return result.Roles, err
}

// loadPermissions 从 YAML 加载权限数据
func (s *RBACSeeder) loadPermissions() ([]PermissionData, error) {
	data, err := seedDataFS.ReadFile("data/permissions.yaml")
	if err != nil {
		return nil, err
	}

	var result struct {
		Permissions []PermissionData `yaml:"permissions"`
	}

	err = yaml.Unmarshal(data, &result)
	return result.Permissions, err
}

// loadMenus 从 YAML 加载菜单数据
func (s *RBACSeeder) loadMenus() ([]MenuData, error) {
	data, err := seedDataFS.ReadFile("data/menus.yaml")
	if err != nil {
		return nil, err
	}

	var result struct {
		Menus []MenuData `yaml:"menus"`
	}

	err = yaml.Unmarshal(data, &result)
	return result.Menus, err
}

// loadRolePermissions 从 YAML 加载角色-权限关联数据
func (s *RBACSeeder) loadRolePermissions() (map[string][]string, error) {
	data, err := seedDataFS.ReadFile("data/role_permissions.yaml")
	if err != nil {
		return nil, err
	}

	var result struct {
		RolePermissions map[string][]string `yaml:"role_permissions"`
	}

	err = yaml.Unmarshal(data, &result)
	return result.RolePermissions, err
}

// loadRoleMenus 从 YAML 加载角色-菜单关联数据
func (s *RBACSeeder) loadRoleMenus() (map[string][]string, error) {
	data, err := seedDataFS.ReadFile("data/role_menus.yaml")
	if err != nil {
		return nil, err
	}

	var result struct {
		RoleMenus map[string][]string `yaml:"role_menus"`
	}

	err = yaml.Unmarshal(data, &result)
	return result.RoleMenus, err
}
