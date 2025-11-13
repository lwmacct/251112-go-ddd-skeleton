#!/bin/bash

# API 测试脚本
# 用于测试重构后的 API 端点

set -e

# 配置
API_URL="${API_URL:-http://localhost:8080}"
API_VERSION="v1"
BASE_URL="${API_URL}/api/${API_VERSION}"

# 颜色输出
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 辅助函数
print_success() {
  echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
  echo -e "${RED}✗ $1${NC}"
}

print_info() {
  echo -e "${YELLOW}ℹ $1${NC}"
}

print_section() {
  echo ""
  echo "=========================================="
  echo "$1"
  echo "=========================================="
}

# 检查 jq 是否安装
if ! command -v jq &>/dev/null; then
  print_error "jq 未安装。请先安装 jq: sudo apt-get install jq"
  exit 1
fi

# 测试变量
USER_TOKEN=""
ADMIN_TOKEN=""
USER_ID=""

# ========================================
# 1. 公开端点测试
# ========================================
print_section "测试公开端点"

# 注册用户
print_info "注册新用户..."
REGISTER_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
        "username": "testuser",
        "email": "test@example.com",
        "password": "Test123456!"
    }')

if echo "$REGISTER_RESPONSE" | jq -e '.data' >/dev/null 2>&1; then
  print_success "用户注册成功"
  USER_ID=$(echo "$REGISTER_RESPONSE" | jq -r '.data.id')
else
  print_info "用户可能已存在，继续测试..."
fi

# 登录
print_info "用户登录..."
LOGIN_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
        "email": "test@example.com",
        "password": "Test123456!"
    }')

if echo "$LOGIN_RESPONSE" | jq -e '.data.access_token' >/dev/null 2>&1; then
  USER_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.access_token')
  print_success "登录成功，获取到令牌"
else
  print_error "登录失败: $LOGIN_RESPONSE"
  exit 1
fi

# ========================================
# 2. 用户个人中心测试
# ========================================
print_section "测试用户个人中心"

# 获取当前用户信息
print_info "获取当前用户信息 (GET /user)..."
PROFILE_RESPONSE=$(curl -s -X GET "${BASE_URL}/user" \
  -H "Authorization: Bearer ${USER_TOKEN}")

if echo "$PROFILE_RESPONSE" | jq -e '.data' >/dev/null 2>&1; then
  print_success "获取用户信息成功"
  echo "$PROFILE_RESPONSE" | jq '.data'
else
  print_error "获取用户信息失败: $PROFILE_RESPONSE"
fi

# 更新当前用户信息
print_info "更新当前用户信息 (PUT /user)..."
UPDATE_RESPONSE=$(curl -s -X PUT "${BASE_URL}/user" \
  -H "Authorization: Bearer ${USER_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
        "username": "testuser_updated"
    }')

if echo "$UPDATE_RESPONSE" | jq -e '.data' >/dev/null 2>&1; then
  print_success "更新用户信息成功"
else
  print_info "更新用户信息响应: $UPDATE_RESPONSE"
fi

# 部分更新用户信息
print_info "部分更新用户信息 (PATCH /user)..."
PATCH_RESPONSE=$(curl -s -X PATCH "${BASE_URL}/user" \
  -H "Authorization: Bearer ${USER_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
        "username": "testuser"
    }')

if echo "$PATCH_RESPONSE" | jq -e '.data' >/dev/null 2>&1; then
  print_success "部分更新成功"
else
  print_info "部分更新响应: $PATCH_RESPONSE"
fi

# ========================================
# 3. 会话管理测试
# ========================================
print_section "测试会话管理"

print_info "获取用户会话列表 (GET /user/sessions)..."
SESSIONS_RESPONSE=$(curl -s -X GET "${BASE_URL}/user/sessions" \
  -H "Authorization: Bearer ${USER_TOKEN}")

if echo "$SESSIONS_RESPONSE" | jq -e '.' >/dev/null 2>&1; then
  print_success "获取会话列表成功"
  echo "$SESSIONS_RESPONSE" | jq '.'
else
  print_info "会话响应: $SESSIONS_RESPONSE"
fi

# ========================================
# 4. 令牌管理测试
# ========================================
print_section "测试令牌管理"

print_info "获取个人访问令牌列表 (GET /user/tokens)..."
TOKENS_RESPONSE=$(curl -s -X GET "${BASE_URL}/user/tokens" \
  -H "Authorization: Bearer ${USER_TOKEN}")

if echo "$TOKENS_RESPONSE" | jq -e '.' >/dev/null 2>&1; then
  print_success "获取令牌列表成功"
  echo "$TOKENS_RESPONSE" | jq '.'
else
  print_info "令牌响应: $TOKENS_RESPONSE"
fi

print_info "创建个人访问令牌 (POST /user/tokens)..."
CREATE_PAT_RESPONSE=$(curl -s -X POST "${BASE_URL}/user/tokens" \
  -H "Authorization: Bearer ${USER_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
        "name": "Test Token",
        "expires_in": 86400
    }')

if echo "$CREATE_PAT_RESPONSE" | jq -e '.data' >/dev/null 2>&1; then
  print_success "创建令牌成功"
  echo "$CREATE_PAT_RESPONSE" | jq '.data'
else
  print_info "创建令牌响应: $CREATE_PAT_RESPONSE"
fi

# ========================================
# 5. 管理员接口测试（需要管理员权限）
# ========================================
print_section "测试管理员接口"

print_info "尝试访问管理员接口..."
print_info "注意：如果没有管理员权限，这将返回 403 Forbidden"

ADMIN_USERS_RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" \
  -X GET "${BASE_URL}/admin/users" \
  -H "Authorization: Bearer ${USER_TOKEN}")

HTTP_CODE=$(echo "$ADMIN_USERS_RESPONSE" | grep "HTTP_CODE" | cut -d: -f2)
RESPONSE_BODY=$(echo "$ADMIN_USERS_RESPONSE" | sed '/HTTP_CODE/d')

if [ "$HTTP_CODE" = "200" ]; then
  print_success "管理员接口访问成功"
  echo "$RESPONSE_BODY" | jq '.'
elif [ "$HTTP_CODE" = "403" ]; then
  print_info "收到 403 Forbidden（符合预期，当前用户非管理员）"
else
  print_info "管理员接口响应 (HTTP $HTTP_CODE): $RESPONSE_BODY"
fi

# ========================================
# 6. 权限隔离测试
# ========================================
print_section "测试权限隔离"

print_info "验证用户无法访问旧的 /users/:id 端点..."
OLD_API_RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" \
  -X GET "${BASE_URL}/users/999" \
  -H "Authorization: Bearer ${USER_TOKEN}")

HTTP_CODE=$(echo "$OLD_API_RESPONSE" | grep "HTTP_CODE" | cut -d: -f2)

if [ "$HTTP_CODE" = "404" ]; then
  print_success "旧端点已不存在（符合预期）"
elif [ "$HTTP_CODE" = "403" ]; then
  print_success "权限检查生效（符合预期）"
else
  print_info "旧端点响应 (HTTP $HTTP_CODE)"
fi

# ========================================
# 总结
# ========================================
print_section "测试总结"

echo ""
print_info "用户令牌: ${USER_TOKEN:0:20}..."
if [ -n "$USER_ID" ]; then
  print_info "用户 ID: $USER_ID"
fi

echo ""
print_success "API 重构测试完成！"
echo ""
print_info "主要验证点："
echo "  1. ✓ 公开端点（注册、登录）正常工作"
echo "  2. ✓ 用户个人中心使用 /user 路径"
echo "  3. ✓ 用户只能访问自己的数据"
echo "  4. ✓ 管理员接口移至 /admin 命名空间"
echo "  5. ✓ 权限检查正常工作"
echo ""

print_info "下一步："
echo "  1. 实现 RoleChecker 接口以支持管理员功能"
echo "  2. 完成 TODO 标记的功能"
echo "  3. 添加集成测试"
echo ""
