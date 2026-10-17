<!-- SPDX-License-Identifier: MIT; Copyright (c) 2026 Phoobobo -->

# 在哪儿 · Zainaer

别问我在哪儿，问“在哪儿”在哪儿。

Zainaer 是面向家庭的自托管物品定位工具：记录物品是什么、放在哪里，以及如何快速找到。它基于 [HomeBox](https://github.com/sysadminsmedia/homebox) 改造，现以本仓库的 **`main`** 为独立产品基线维护，不再与上游合并或同步。

[仓库](https://github.com/Phoobobo/zainaer) · [问题反馈](https://github.com/Phoobobo/zainaer/issues) · [许可证](LICENSE)

## 功能

- **家庭物品记录**：名称、标签、照片、描述和存放位置，收起购买、保修、销售和维护等库存管理信息。
- **多级位置**：按“客厅 → 餐边柜 → 左下第三个抽屉”建立位置层级，无需使用收纳箱。
- **文字助手**：多轮对话查找物品，显示真实记录和完整位置路径；准备新增物品、位置和标签。
- **确认后入库**：模型只生成草稿。用户点击确认后，服务端重新校验并保存；草稿绑定用户与集合，10 分钟过期，防止重复提交。
- **服务端 BYOK**：集合所有者管理提供商、主模型和两个备选模型。支持 OpenRouter、阿里云百炼（北京）和 DeepSeek。
- **故障切换**：主模型请求失败、超时或返回空响应时，依次尝试备选模型；共用所选提供商与 API Key，不因回答质量触发重试。
- **中英双语、移动端布局**，以及创建物品、位置、标签的悬浮入口。
- **家庭集合与邀请**：数据按集合隔离；模型查询使用当前登录用户的集合权限。

当前版本以 [VERSION](VERSION) 为准，独立版本号从 `0.0.1` 开始。语音、AI 照片识别、助手修改/删除现有记录尚未实现；修改和删除仍使用常规界面。

## 模型配置

登录后进入 **集合 → 模型配置**，选择提供商，填写完整模型 ID 与 API Key，保存后测试连接。OpenRouter 的模型 ID 需要带厂商前缀；备选模型可留空。模型必须支持工具调用，费用及免费额度以提供商规则为准。

密钥使用 AES-GCM 加密，保存在服务端存储目录的 `ai/` 子目录，文件权限为 `0600`；页面只返回是否已配置，不回显密钥。**备份必须同时保留 `ai/*.json` 和 `ai/master.key`**。主密钥与密文保存在同一服务端存储目录，加密不能替代服务器和备份的访问控制。

助手会将对话和本次工具查询返回的记录发送给配置的模型提供商。模型不持有用户登录凭据，不执行 shell；新增记录必须经过独立确认接口。对话历史目前仅保留在当前页面内存中，刷新后清空；待确认草稿在服务重启后失效。

## 构建

后端为 Go，前端为 Vue / Nuxt。生产程序嵌入前端静态资源，可作为单个可执行文件部署；默认使用 SQLite。

需要 Go 1.26+、Node.js 22+ 和 pnpm 10.28。

```sh
git clone https://github.com/Phoobobo/zainaer.git
cd zainaer
pnpm --dir frontend install --frozen-lockfile

# 构建 Linux amd64 程序
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  HOUSEHOLD_BUILD_OUTPUT="$PWD/zainaer" \
  sh scripts/build-household.sh
```

构建脚本读取 `VERSION` 和当前 Git 提交，生成前端资源并嵌入后端。若需其他平台，请调整 `GOOS`、`GOARCH`。

### 服务器运行示例

将程序复制到服务器，在独立数据目录中运行。以下环境变量用于说明配置方式，不包含可直接使用的密钥：

```sh
export HBOX_WEB_HOST=127.0.0.1
export HBOX_WEB_PORT=7745
export HBOX_WEB_WRITE_TIMEOUT=180s
export HBOX_AUTH_API_KEY_PEPPER='<自行生成并持久保存的随机密钥>'
export HBOX_OPTIONS_ALLOW_ANALYTICS=false
export HBOX_OPTIONS_GITHUB_RELEASE_CHECK=false
./zainaer
```

工作目录下的 `.data/` 存放默认数据库、附件与模型配置。生产环境使用进程管理器保持运行，通过 HTTPS 反向代理或受控隧道访问。模型可能需要多轮请求，响应超时应覆盖助手的 150 秒总执行窗口。

出于与现有部署的兼容性，内部 Go module/import 路径与 `HBOX_*` 环境变量暂保留 HomeBox 命名；这不代表本项目继续跟踪上游版本。

## 验证

```sh
cd backend
go test ./app/api/handlers/v1 -run TestAI -count=1
```

这些测试覆盖密钥加密、配置校验、集合隔离、草稿确认/取消、重复提交和模型故障切换。网络模拟测试需要临时回环端口，不调用真实模型。真实提供商检查是显式启用的可选测试。

前端构建及 ESLint 使用 `frontend/package.json` 中的脚本。全量类型检查目前存在继承代码中的 Ref / 实体类型问题，不能将前端构建成功描述为全量类型检查通过。

## 部署与备份

维护者当前仅以北京 ECS 为生产服务；原 Mac 应用及 Companion 已退役，本机 SSH 隧道仅用于访问云端。具体主机信息和凭据不写入仓库。

更新前备份当前程序及 SQLite 数据库，上传后核对校验和，再替换程序并重启。不要使用旧的本地数据库覆盖云端数据。备份数据库时应使用 SQLite backup API，或停服后取得一致快照；同时保存附件、模型加密配置及解密主密钥。

## 来源与许可证

感谢 [HomeBox / sysadminsmedia](https://github.com/sysadminsmedia/homebox) 及其贡献者。本项目保留相关历史与归属，整体及继承代码遵循 [GNU AGPL-3.0](LICENSE)。明确标注的 Zainaer 原创文件另按宽松的 [MIT](LICENSE-MIT) 许可提供，授权边界见 [LICENSING.md](LICENSING.md)；这不是将整个项目改为 MIT。既有上游文档和自动化可能仍包含 HomeBox 的部署方式或链接，不应视为 Zainaer 的官方发布配置。

---

**English:** Zainaer is a self-hosted household item locator, independently maintained from HomeBox on `main`. It offers nested locations, photos and tags, Chinese/English UI, server-side BYOK with two fallback models, and a text assistant whose creation drafts require explicit user confirmation. The existing Go module paths and `HBOX_*` settings remain for compatibility. The combined application remains AGPL-3.0; explicitly marked original files are available under MIT. See [license scope](LICENSING.md).
