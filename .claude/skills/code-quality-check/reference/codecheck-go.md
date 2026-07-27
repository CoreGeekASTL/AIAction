# CodeCheck Go规则参考

CodeCheck Go语言的代码质量检查规则，涵盖函数复杂度、注释规范、命名规范、安全规范等。

## 1. 检查项汇总

| 检查项 | 阈值 | 问题类型 | 严重度 | 规则ID |
| --- | --- | --- | --- | --- |
| **超大函数** | 最大50行 | 可读性问题 | 一般 | 超大函数[GO] |
| **超大函数深度** | 最大4层 | 可读性问题 | 一般 | 超大函数深度[GO] |
| **冗余代码** | 不允许 | 可读性问题 | 一般 | 冗余代码[GO] |
| **文件头版权声明** | 必须包含 | 规范问题 | 提示 | G.CMT.01[GO] |
| **包注释** | 每个包必须有 | 规范问题 | 提示 | G.CMT.02[GO] |
| **注释位置** | 注释置于代码上方或右边 | 规范问题 | 提示 | G.CMT.05[GO] |
| **chan方向限定** | chan参数指定方向 | 规范问题 | 一般 | G.CON.01[GO] |
| **导出错误命名** | ErrXxx格式 | 规范问题 | 提示 | G.NAM.03[GO] |
| **文件命名** | 全小写+下划线 | 规范问题 | 一般 | G.NAM.05[GO] |
| **命名风格统一** | 统一命名格式 | 规范问题 | 提示 | G.NAM.01[GO] |
| **包名与目录名一致** | 包名=目录名 | 规范问题 | 一般 | G.NAM.08[GO] |
| **魔法数字** | 提取为常量 | 可读性问题 | 一般 | G.DCL.01[GO] |
| **资源泄漏** | 资源必须释放 | 安全问题 | 严重 | G.RES.01[GO] |
| **错误返回值检查** | 必须检查 | 安全问题 | 严重 | G.ERR.01[GO] |
| **文件权限指定** | 显式指定权限 | 安全问题 | 严重 | G.FIO.02[GO] |
| **gofmt缩进规范** | 使用gofmt格式化代码 | 规范问题 | 一般 | G.FMT.01[GO] |
| **switch有default** | 必须有default分支 | 安全问题 | 一般 | G.CTL.05[GO] |

---

## 2. 函数行数检查

### 2.1 检查规则

| 检查项 | 阈值 | 问题类型 | 检查方法 | 严重度 |
| --- | --- | --- | --- | --- |
| **超大函数** | 最大50行 | 可读性问题 | 函数体行数超出阈值 | 一般 |

### 2.2 超标示例

| 文件 | 函数 | 行数 | 阈值 | 超出量 |
| --- | --- | --- | --- | --- |
| `auth_controller.go` | ImportIMEIList() | 55 | 50 | +5 |
| `plugin_service.go` | UploadPluginPackage() | 56 | 50 | +6 |
| `plugin_service.go` | readMetaFromZip() | 52 | 50 | +2 |
| `auditlog.go` | AuditsLog() | 52 | 50 | +2 |
| `auth_service_test.go` | TestAuthServiceImpl_ImportIMEIList() | 130 | 50 | +80 |
| `auth_service_test.go` | TestAuthServiceImpl_ExportIMEIList() | 88 | 50 | +38 |
| `auth_service_test.go` | TestAuthServiceImpl_CheckIMEI() | 67 | 50 | +17 |
| `auth_service_test.go` | Test_parseIMEIRange() | 54 | 50 | +4 |
| `base_dao_test.go` | TestBaseDao() | 71 | 50 | +21 |
| `minio_test.go` | TestGetObject() | 163 | 50 | +113 |
| `builder_test.go` | Test_request() | 81 | 50 | +31 |
| `https_server_test.go` | TestGetLocalIP() | 57 | 50 | +7 |
| `file_service_test.go` | Test_CleanFileName() | 55 | 50 | +5 |
| `monitor_service_test.go` | TestMonitorSchedule() | 98 | 50 | +48 |

### 2.3 修正策略

| 策略 | 适用场景 | 修正方式 |
| --- | --- | --- |
| **拆分函数** | 函数承担过多职责 | 按职责拆分为多个子函数 |
| **提取辅助函数** | 存在重复逻辑 | 提取为独立的辅助函数 |
| **按步骤拆分** | 流程型函数 | 每个步骤提取为子函数 |

### 2.4 修正示例

#### 示例1：UploadPluginPackage() - 56行→22行

```go
// ❌ 问题：函数行数56行，超出50行阈值
func (p *PluginServiceImpl) UploadPluginPackage(req *req.UploadPluginPackageReq) error {
    defer func(File multipart.File) {
        err := File.Close()
        // ...
    }(req.File)
    
    pkg, err := p.readPackageMeta(req)
    // ... 检查数据是否存在 (10行)
    // ... 上传软件包 (8行)
    // ... 数据库保存 (20行)
    return p.ppd.DoTxWithCtx(...)
}

// ✅ 修正：拆分为3个子函数
func (p *PluginServiceImpl) UploadPluginPackage(req *req.UploadPluginPackageReq) error {
    defer func(File multipart.File) {
        err := File.Close()
        if err != nil {
            logger.Warnf("file close error: %s", err)
        }
    }(req.File)

    pkg, err := p.readPackageMeta(req)
    if err != nil {
        return err
    }
    if err := p.checkPackageExists(pkg); err != nil {
        return err
    }
    content, err := p.readFileContent(req)
    if err != nil {
        return err
    }
    return p.savePackage(pkg, content, req.Size)
}

func (p *PluginServiceImpl) checkPackageExists(pkg *db.PluginPackage) error {
    oldPP := &db.PluginPackage{Field: pkg.GetField()}
    err := p.ppd.Get(oldPP)
    if err != nil && err != orm.ErrNoRows {
        logger.Errorf("check key %s exist failed, err is %v", pkg.GetKey(), err)
        return err
    }
    if err == nil {
        return fmt.Errorf("key %s is exist, forbidden upload", pkg.GetKey())
    }
    return nil
}

func (p *PluginServiceImpl) readFileContent(req *req.UploadPluginPackageReq) ([]byte, error) {
    _, err := req.File.Seek(0, io.SeekStart)
    if err != nil {
        logger.Errorf("seek file to start error: %s", err)
        return nil, err
    }
    content, err := io.ReadAll(req.File)
    if err != nil {
        logger.Errorf("read request file failed: %v", err)
        return nil, err
    }
    return content, nil
}

func (p *PluginServiceImpl) savePackage(pkg *db.PluginPackage, content []byte, size int64) error {
    return p.ppd.DoTxWithCtx(context.Background(), func(ctx context.Context, txOrm orm.TxOrmer) error {
        pkg.Status = db.NotStart
        pkg.IfActive = false
        f := &db.File{
            Bucket:    pkg.PackageBucket,
            Name:      pkg.PackageName,
            Content:   content,
            Size:      size,
            CreatedAt: time.Now().Format(time.DateTime),
        }
        if err := p.fd.InsertWithOrm(ctx, txOrm, f); err != nil {
            logger.Errorf("put file %s to storage failed, err is %v", pkg.PackageName, err)
            return err
        }
        pkg.Field = pkg.GetField()
        pkg.CreatedAt = time.Now().Format(time.DateTime)
        if err := p.ppd.InsertWithOrm(ctx, txOrm, pkg); err != nil {
            logger.Errorf("Set data to redis failed, data: %v, err %v", pkg, err)
            return err
        }
        return nil
    })
}
```

#### 示例2：AuditsLog() - 52行→10行

```go
// ❌ 问题：审计日志函数52行
func AuditsLog(auditsPara *AuditsPara, requestURL string) {
    headers := make(http.Header)
    headers.Set("Content-Type", "application/json")
    body := make(map[string]interface{})
    // ... 构建body (20行)
    // ... 序列化 (10行)
    // ... 发送请求 (15行)
}

// ✅ 修正：拆分为4个子函数
func AuditsLog(auditsPara *AuditsPara, requestURL string) {
    headers := buildAuditLogHeaders()
    body := buildAuditLogBody(auditsPara, requestURL)
    bs2, err := marshalAuditLogBody(body)
    if err != nil {
        return
    }
    sendAuditLog(requestURL, headers, bs2)
}

func buildAuditLogHeaders() http.Header { ... }
func buildAuditLogBody(auditsPara *AuditsPara, requestURL string) map[string]interface{} { ... }
func marshalAuditLogBody(body map[string]interface{}) ([]byte, error) { ... }
func sendAuditLog(requestURL string, headers http.Header, bs2 []byte) { ... }
```

---

## 3. 函数深度检查

### 3.1 检查规则

| 检查项 | 阈值 | 问题类型 | 检查方法 | 严重度 |
| --- | --- | --- | --- | --- |
| **超大函数深度** | 最大4层 | 可读性问题 | 函数嵌套深度超出阈值 | 一般 |

### 3.2 深度计算规则

| 结构 | 深度贡献 |
| --- | --- |
| 函数定义 | 1层 |
| for循环 | +1层 |
| if条件 | +1层 |
| goroutine | +1层 |
| select语句 | +1层 |
| switch/case | +1层 |
| 匿名函数 | +1层 |
| t.Run() | +1层 |

### 3.3 超标示例

| 文件 | 函数 | 深度 | 阈值 | 超出量 | 嵌套结构 |
| --- | --- | --- | --- | --- | --- |
| `https_server.go` | Run() | 5 | 4 | +1 | func+go+for+select+if |
| `minio_test.go` | TestGetObject() | 6 | 4 | +2 | for+t.Run+if+if+if+if |
| `request_entity_test.go` | TestDeleteCacheRequest_Validate() | 5 | 4 | +1 | for+t.Run+if+if+if |
| `request_entity_test.go` | TestFileUploadRequest_Validate() | 5 | 4 | +1 | for+t.Run+if+if+if |
| `file_service_test.go` | Test_CleanFileName() | 5 | 4 | +1 | for+t.Run+if+if+if |

### 3.4 修正策略

| 策略 | 适用场景 | 修正方式 |
| --- | --- | --- |
| **提取子函数** | goroutine+for+select嵌套 | 提取monitor/handler函数 |
| **Guard Clause** | 多层if嵌套 | 使用提前返回减少嵌套 |
| **提取验证辅助函数** | 测试代码深度超标 | 提取validate辅助函数 |

### 3.5 修正示例

#### 示例1：Run() - 深度5层→1层

```go
// ❌ 问题：深度5层
// 函数定义(1) + goroutine(1) + for循环(1) + select(1) + if(1) = 5层
func (b *BeegoHttpsServer) Run() {
    go func() {
        for {
            select {
            case certInfo := <-b.restartChan:
                if !b.needStartServer() {
                    // ...
                }
            }
        }
    }()
}

// ✅ 修正：提取子函数，深度降至1层
func (b *BeegoHttpsServer) Run() {
    go b.monitorCertificate()
}

func (b *BeegoHttpsServer) monitorCertificate() {
    for {
        certInfo := <-b.restartChan
        b.handleCertificateUpdate(certInfo)
    }
}

func (b *BeegoHttpsServer) handleCertificateUpdate(certInfo CertInfo) {
    b.updateCertInfo(certInfo)
    if !b.needStartServer() {
        return
    }
    b.server.Server.TLSConfig = GetTLS(b.certInfo, ServerType)
    go b.server.Run("")
    b.isServerReady = true
}
```

#### 示例2：TestGetObject() - 深度6层→3层

```go
// ❌ 问题：测试代码深度6层
// for(1) + t.Run(1) + if(1) + if(1) + if(1) + if(1) = 6层
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        object, err := Instance().GetObject(tt.args.ctx, tt.args.options)
        if tt.wantErr {
            if err == nil {
                t.Errorf("expected error but got nil")
            } else if err.Error() != tt.wantErrMsg {
                t.Errorf("expected error %q, got %q", tt.wantErrMsg, err.Error())
            }
            if object != nil {
                if closeErr := object.Close(); closeErr != nil {
                    t.Errorf("failed to close object: %v", closeErr)
                }
            }
        } else {
            // ...
        }
    })
}

// ✅ 修正：提取验证辅助函数，深度降至3层
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        object, err := Instance().GetObject(tt.args.ctx, tt.args.options)
        validateObjectResult(t, object, err, tt.wantErr, tt.wantErrMsg)
    })
}

func validateObjectResult(t *testing.T, object io.ReadCloser, err error, wantErr bool, wantErrMsg string) {
    if wantErr {
        validateError(t, err, wantErrMsg)
        closeObjectIfNeeded(t, object)
        return
    }
    validateSuccess(t, object, err)
}

func validateError(t *testing.T, err error, wantErrMsg string) { ... }
func closeObjectIfNeeded(t *testing.T, object io.ReadCloser) { ... }
func validateSuccess(t *testing.T, object io.ReadCloser, err error) { ... }
```

---

## 4. 冗余代码检查

### 4.1 检查规则

| 检查项 | 阈值 | 问题类型 | 检查方法 | 严重度 | 规则ID |
| --- | --- | --- | --- | --- | --- |
| **冗余代码** | 不允许 | 可读性问题 | 代码中存在注释掉的代码、冗余空行等无用内容 | 一般 | 冗余代码[GO] |

### 4.2 冗余代码类型

| 冗余类型 | 说明 | 典型示例 |
| --- | --- | --- |
| **注释掉的代码** | 被注释掉但保留的代码行，应删除而非注释保留 | `// defer resp.Body.Close()` |
| **冗余空行** | 版权声明与package声明之间应保留一个空行 | `// Copyright...\npackage xxx` 应改为 `// Copyright...\n\npackage xxx` |
| **冗余注释** | 代码逻辑已足够清晰，注释只是重复描述代码 | `endpoints[ipPort] = struct{}{} // 将 IP 和端口加入到 endpoints 中` |

### 4.3 修正策略

| 冗余类型 | 修正方式 |
| --- | --- |
| 注释掉的代码 | 直接删除注释行，需要时可通过Git历史找回 |
| 冗余空行（版权与package间缺少空行） | 添加空行，版权声明与package声明之间保留一行空行 |
| 冗余注释 | 删除对代码逻辑的简单重复描述，保留有价值的业务说明 |

### 4.4 修正示例

#### 示例1：注释掉的代码

```go
// ❌ 问题：注释掉的代码行
if err != nil {
    return resp, err
}
// defer resp.Body.Close()

return resp, nil

// ✅ 修正：删除注释行
if err != nil {
    return resp, err
}

return resp, nil
```

#### 示例2：版权声明与package间冗余空行

```go
// ❌ 问题：版权与package间缺少空行
// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.
package service

// ✅ 修正：版权与package间保留空行
// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.

package service
```

---

## 5. 注释规范检查

### 5.1 G.CMT.01 - 文件头版权声明

| 检查项 | 阈值 | 问题类型 | 严重度 | 规则ID |
| --- | --- | --- | --- | --- |
| **文件头版权声明** | 必须包含 | 规范问题 | 提示 | huawei-copyright |

**规则**：每个Go源文件首行必须包含版权声明 `// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.`

**修正方式**：在文件首行添加版权声明，版权声明与package声明之间保留一行空行。

```go
// ❌ 问题：缺少版权声明
package service

// ✅ 修正：添加版权声明（版权与package间保留空行）
// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.

package service
```

### 5.2 G.CMT.02 - 包注释

| 检查项 | 阈值 | 问题类型 | 严重度 | 规则ID |
| --- | --- | --- | --- | --- |
| **包注释** | 每个包必须有 | 规范问题 | 提示 | huawei-package-comments |

**规则**：每个包必须有包级注释，格式为 `// Package xxx 功能说明`，放在package声明之前。

```go
// ❌ 问题：缺少包注释
// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.

package cert

// ✅ 修正：添加包注释（包注释在空行之后、package之前）
// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved.

// Package cert 证书管理
package cert
```

### 5.3 G.CMT.05 - 注释位置

| 检查项 | 阈值 | 问题类型 | 严重度 | 规则ID |
| --- | --- | --- | --- | --- |
| **注释位置** | 注释置于代码上方或右边 | 规范问题 | 提示 | huawei-comments-location |

**规则**：注释应置于对应代码的上方或右边，不允许注释与代码之间有空行。

```go
// ❌ 问题：注释与代码间有空行
resp, err = c.client.Do(req)

// defer resp.Body.Close()

return resp, nil

// ✅ 修正：删除注释掉的代码（或紧贴代码上方放置注释）
resp, err = c.client.Do(req)
return resp, nil
```

---

## 6. 命名规范检查

### 6.1 G.CON.01 - chan方向限定

| 检查项 | 阈值 | 问题类型 | 严重度 | 规则ID |
| --- | --- | --- | --- | --- |
| **chan方向限定** | 函数参数chan必须指定方向 | 规范问题 | 一般 | huawei-incorrect-chan |

**规则**：chan作为函数参数时，必须限定为只发送(`chan<-`)或只接收(`<-chan`)方向，除非确实需要双向。

```go
// ❌ 问题：chan未指定方向
func gsfStartHandler(done chan bool, podName string) {
    done <- true
}

// ✅ 修正：指定为只发送方向
func gsfStartHandler(done chan<- bool, podName string) {
    done <- true
}
```

### 6.2 G.NAM.03 - 导出错误变量命名

| 检查项 | 阈值 | 问题类型 | 严重度 | 规则ID |
| --- | --- | --- | --- | --- |
| **导出错误命名** | ErrXxx格式 | 规范问题 | 提示 | huawei-wrong-format-error-naming |

**规则**：导出的错误类型应使用`XxxError`格式命名，而非`Err`。

```go
// ❌ 问题：导出错误类型名为Err
type Err string

// ✅ 修正：使用XxxError格式
type AppError string
```

### 6.3 G.NAM.05 - 文件命名

| 检查项 | 阈值 | 问题类型 | 严重度 | 规则ID |
| --- | --- | --- | --- | --- |
| **文件命名** | 全小写+下划线 | 规范问题 | 一般 | huawei-file-name |

**规则**：Go文件名应使用全小写，允许数字和下划线，不允许大写字母。

```go
// ❌ 问题：文件名含大写
AlarmService.go
VideoService.go

// ✅ 修正：全小写+下划线
alarm_service.go
video_service.go
```

### 6.4 G.NAM.01 - 命名风格统一

| 9查项 | 阈值 | 问题类型 | 严重度 | 规则ID |
| --- | --- | --- | --- | --- |
| **命名风格统一** | 统一命名格式 | 规范问题 | 提示 | huawei-name |

**规则**：使用统一的命名格式。Go测试函数应使用`TestXxx`或`TestXxx_Yyy`格式，避免`Test_xxx`风格。

### 6.5 G.NAM.08 - 包名与目录名一致

| 检查项 | 阈值 | 问题类型 | 严重度 | 规则ID |
| --- | --- | --- | --- | --- |
| **包名与目录名一致** | 包名=目录名 | 规范问题 | 一般 | huawei-directory-name-consistency |

**规则**：Go包名应与所在目录名一致。

**注意**：此规则修复涉及全局import路径更新，风险较高，需谨慎处理。当包名与目录名不一致时（如`src/common/error/`使用`package common`），建议在评估影响范围后再决定是否重构。

---

## 7. 声明与魔法数字检查

### 7.1 G.DCL.01 - 魔法数字

| 检查项 | 阈值 | 问题类型 | 严重度 | 规则ID |
| --- | --- | --- | --- | --- |
| **魔法数字** | 提取为常量 | 可读性问题 | 一般 | huawei-magic-number |

**规则**：不要使用难以理解的字面量，应提取为命名常量。

**常见场景与修正**：

| 场景 | 魔法数字 | 修正方式 |
| --- | --- | --- |
| HTTP状态码 | `200` | 使用`http.StatusOK` |
| 端口默认值 | `9997`, `9990` | 提取为`defaultXxxPort`常量 |
| 时间间隔 | `10 * time.Second` | 提取为`xxxInterval`常量 |
| 阈值比率 | `0.75`, `0.85` | 提取为`xxxRatio`常量 |
| 退出码 | `os.Exit(3)` | 提取为`restartExitCode`常量 |
| 数量阈值 | `10*60*1000` | 提取为`alarmSuppressThreshold`常量 |

```go
// ❌ 问题：使用魔法数字
httpsPort := beego.AppConfig.DefaultInt("httpsport", 9997)
if response.StatusCode != 200 { ... }
time.Sleep(10 * time.Second)

// ✅ 修正：提取为常量
const defaultInternalHttpsPort = 9997
const respOK = http.StatusOK
const alarmRetrySleepSeconds = 10

httpsPort := beego.AppConfig.DefaultInt("httpsport", defaultInternalHttpsPort)
if response.StatusCode != respOK { ... }
time.Sleep(alarmRetrySleepSeconds * time.Second)
```

---

## 8. 安全规范检查

### 8.1 G.RES.01 - 资源泄漏

| 检查项 | 阈值 | 问题类型 | 严重度 | 规则ID |
| --- | --- | --- | --- | --- |
| **资源泄漏** | 资源必须释放 | 安全问题 | 严重 | Resource_Leak |

**规则**：所有通过`os.Create`、`os.Open`、`os.OpenFile`等创建的文件资源必须在使用后关闭，使用`defer file.Close()`确保释放。

**常见问题**：
- `os.Open()`/`os.OpenFile()`打开的文件未调用`defer file.Close()`关闭
- goroutine没有退出机制，无法保证正常退出

**推荐修复模式**：在`os.Open`/`os.OpenFile`成功后，紧跟`defer file.Close()`，与存量代码风格一致（参考`fileutil.go`、`zip_util.go`）。

```go
// ❌ 问题：文件打开后未defer关闭
file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
if err != nil {
    return err
}
// 使用file但没有defer file.Close()

// ✅ 修正：紧跟defer file.Close()（标准模式，与存量代码一致）
file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
if err != nil {
    return err
}
defer file.Close()
```

```go
// ❌ 问题：goroutine无退出机制
go func() {
    for {
        certInfo := <-ch
        // 处理...
    }
}()

// ✅ 修正：添加stopChan退出机制
type Server struct {
    stopChan   chan struct{}
    restartChan chan CertInfo
}

func (s *Server) monitor() {
    for {
        select {
        case <-s.stopChan:
            return
        case certInfo := <-s.restartChan:
            s.handle(certInfo)
        }
    }
}

func (s *Server) Stop() { close(s.stopChan) }
```

### 8.2 G.ERR.01 - 错误返回值检查

| 检查项 | 阈值 | 问题类型 | 严重度 | 规则ID |
| --- | --- | --- | --- | --- |
| **错误返回值检查** | 必须检查 | 安全问题 | 严重 | Check_Return_Error |

**规则**：函数的错误返回值必须被检查处理，不允许忽略。

```go
// ❌ 问题：忽略json.Marshal错误
respData, _ := json.Marshal(data)
logger.Infof("response is %s", string(respData))

// ✅ 修正：检查错误
respData, err := json.Marshal(data)
if err != nil {
    logger.Errorf("marshal response data failed: %v", err)
} else {
    logger.Infof("response is %s", string(respData))
}

// ❌ 问题：忽略resp.Body.Close()错误
resp.Body.Close()

// ✅ 修正：检查错误
if closeErr := resp.Body.Close(); closeErr != nil {
    logger.Errorf("close response body failed: %v", closeErr)
}
```

### 8.3 G.FIO.02 - 文件权限指定

| 检查项 | 阈值 | 问题类型 | 严重度 | 规则ID |
| --- | --- | --- | --- | --- |
| **文件权限指定** | 显式指定权限 | 安全问题 | 严重 | File_Permission |

**规则**：创建文件时必须显式指定合适的文件访问权限，使用`os.OpenFile`替代`os.Create`。

```go
// ❌ 问题：os.Create未指定权限
file, err := os.Create(absolutePath)

// ✅ 修正：os.OpenFile指定权限
file, err := os.OpenFile(absolutePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
```

### 8.4 G.CTL.05 - switch必须有default

| 检查项 | 阈值 | 问题类型 | 严重度 | 规则ID |
| --- | --- | --- | --- | --- |
| **switch有default** | 必须有default分支 | 安全问题 | 一般 | Missing_Default_Branch |

**规则**：所有switch语句必须包含default分支。

```go
// ❌ 问题：switch缺少default
switch tlsType {
case ServerType:
    tlsConfig.ClientCAs = caPool
case ClientType:
    tlsConfig.RootCAs = caPool
}

// ✅ 修正：添加default分支
switch tlsType {
case ServerType:
    tlsConfig.ClientCAs = caPool
case ClientType:
    tlsConfig.RootCAs = caPool
default:
    logger.Infof("unknown tls type: %s", tlsType)
}
```

---

## 9. 检查步骤

### 9.1 步骤流程

| 步骤 | 操作 | 方法 |
| --- | --- | --- |
| 1 | 识别超标函数 | 从CodeCheck报告中获取超标函数列表 |
| 2 | 分析超标原因 | 行数超标：职责过多；深度超标：嵌套过多 |
| 3 | 制定拆分策略 | 按职责/层次/复用拆分 |
| 4 | 验证拆分效果 | 行数≤50行、深度≤4层、测试通过 |

### 9.2 拆分策略选择

| 超标类型 | 拆分策略 | 适用场景 |
| --- | --- | --- |
| 行数超标 | 按职责拆分 | 函数承担多个独立职责 |
| 行数超标 | 按步骤拆分 | 流程型函数（解析→验证→处理→响应） |
| 深度超标 | 提取子函数 | goroutine+for+select嵌套 |
| 深度超标 | Guard Clause | 多层if嵌套检查参数 |
| 深度超标 | 提取辅助函数 | 测试代码验证逻辑 |

---

## 10. 注意事项

### 10.1 测试函数特殊性

| 问题 | 建议 |
| --- | --- |
| 表驱动测试行数较长 | 合理的Go测试风格，重点优化深度而非强行压缩行数 |
| 测试验证逻辑重复 | 提取验证辅助函数是推荐做法 |
| 测试数据准备 | 可提取数据准备函数 |

### 10.2 可读性要求

| 要求 | 说明 |
| --- | --- |
| 函数命名清晰 | 子函数命名应表达职责 |
| 避免过度拆分 | 防止逻辑碎片化 |
| 保持调用关系 | 子函数应有合理的调用层次 |

### 10.3 重构成本权衡

| 超标程度 | 处理建议 |
| --- | --- |
| 严重超标（行数≥100，深度≥6） | 必须修复 |
| 中度超标（行数51-80，深度5） | 建议修复 |
| 略超阈值（行数51-55，深度5） | 可酌情处理 |

---

## 11. 补充规范（CodeCheck实际检出规则）

### 11.1 G.CMT.04 - 注释符与注释内容间要有1个空格

| 检查项 | 阈值 | 问题类型 | 严重度 | 规则ID |
| --- | --- | --- | --- | --- |
| **注释空格格式** | 注释符后必须有空格 | 规范问题 | 提示 | huawei-comment-space-format |

**规则**：注释符(`//`)与注释内容之间必须保留至少1个空格。

```go
// ❌ 问题：注释符后无空格
//TODO implement me
//逃生

// ✅ 修正：注释符后加空格
// TODO implement me
// 逃生
```

### 11.2 G.ERR.02 - 错误信息不应大写或以标点符号结尾

| 检查项 | 阈值 | 问题类型 | 严重度 | 规则ID |
| --- | --- | --- | --- | --- |
| **错误信息格式** | 小写且不以标点结尾 | 规范问题 | 提示 | huawei-wrong-format-error-string |

**规则**：`errors.New()`和`fmt.Errorf()`创建的错误信息应以小写字母开头，且不以标点符号结尾。

```go
// ❌ 问题：错误信息大写或以标点结尾
return errors.New("MultiTableRequest is nil")
return fmt.Errorf("BrowserGW returned non-OK status: %d", resp.StatusCode)
return fmt.Errorf("Insert database error.")

// ✅ 修正：小写开头，无标点结尾
return errors.New("multiTableRequest is nil")
return fmt.Errorf("browserGW returned non-OK status: %d", resp.StatusCode)
return fmt.Errorf("insert database error")
```

### 11.3 G.ERR.03 - 包外可见的函数禁止向外抛出panic

| 检查项 | 阈值 | 问题类型 | 严重度 | 规则ID |
| --- | --- | --- | --- | --- |
| **禁止panic** | 导出函数不得panic | 安全问题 | 严重 | No_Panic |

**规则**：包外可见（导出）的函数禁止使用`panic`，应返回error。mock/stub中的`panic("implement me")`应替换为`fmt.Errorf("not implemented")`。

```go
// ❌ 问题：mock中使用panic
func (m *MockCSEService) GetAllMicroServices() []MicroService {
    panic("implement me")
}

// ✅ 修正：返回error
func (m *MockCSEService) GetAllMicroServices() ([]MicroService, error) {
    return nil, fmt.Errorf("not implemented")
}
```

### 11.4 G.FUN.03 - 函数的返回值应避免使用命名返回值

| 检查项 | 阈值 | 问题类型 | 严重度 | 规则ID |
| --- | --- | --- | --- | --- |
| **命名返回值** | 避免使用命名返回值 | 规范问题 | 一般 | huawei-function-results |

**规则**：函数返回值应避免使用命名返回值（named return values），除非为了提高可读性或实现裸返回（naked return）。在大多数情况下应使用匿名返回值。

```go
// ❌ 问题：使用命名返回值
func (f *fakeReader) Read(p []byte) (n int, err error) { ... }
func (p *PluginActive) MarshalBinary() (data []byte, err error) { ... }

// ✅ 修正：使用匿名返回值
func (f *fakeReader) Read(p []byte) (int, error) { ... }
func (p *PluginActive) MarshalBinary() ([]byte, error) { ... }
```

### 11.5 G.SIT.02 - 必须处理类型断言的失败

| 检查项 | 阈值 | 问题类型 | 严重度 | 规则ID |
| --- | --- | --- | --- | --- |
| **类型断言安全** | 必须使用comma-ok模式 | 安全问题 | 严重 | Assert_Error |

**规则**：类型断言必须使用comma-ok模式（`v, ok := x.(T)`），避免断言失败时panic。

```go
// ❌ 问题：类型断言未检查失败
factory := GetFactory().(*eventStorageFactory)
instance := v.(browsergateway.ServiceInstance)
cacheEntry := value.(*IMEICache)

// ✅ 修正：使用comma-ok模式
factory, ok := GetFactory().(*eventStorageFactory)
if !ok {
    t.Fatalf("factory type assertion failed")
}

instance, ok := v.(browsergateway.ServiceInstance)
if !ok {
    logger.Warnf("type assertion failed for ServiceInstance")
    return true
}
```

### 11.6 G.PKG.02 - 禁止使用.来简化导入包

| 检查项 | 阈值 | 问题类型 | 严重度 | 规则ID |
| --- | --- | --- | --- | --- |
| **禁止dot导入** | 禁止使用.简化导入 | 规范问题 | 一般 | huawei-imports-dot |

**规则**：禁止使用`.`来简化导入包（dot import），应使用显式别名。

```go
// ❌ 问题：使用dot导入
. "github.com/smartystreets/goconvey/convey"
. "GIDS/test/util"

// ✅ 修正：使用显式别名
convey "github.com/smartystreets/goconvey/convey"
testutil "GIDS/test/util"
```

### 11.7 G.NAM.07 - 包名采用全小写单词

| 检查项 | 阈值 | 问题类型 | 严重度 | 规则ID |
| --- | --- | --- | --- | --- |
| **包名格式** | 全小写无下划线 | 规范问题 | 一般 | huawei-package-name |

**规则**：Go包名应采用全小写单词，允许包含数字，不允许使用下划线。

```go
// ❌ 问题：包名含下划线或大写
package eventsStorage
package MyPackage

// ✅ 修正：全小写无下划线
package event
package mypackage
```

### 11.8 G.RES.05 - 避免过多的time.After函数调用导致消耗大量的资源

| 检查项 | 阈值 | 问题类型 | 严重度 | 规则ID |
| --- | --- | --- | --- | --- |
| **time.After循环优化** | 循环中使用time.NewTimer替代time.After | 安全问题 | 严重 | TimeAfter_In_Loop |

**规则**：在for循环中使用`time.After`会不断创建新的Timer导致资源浪费，应使用`time.NewTimer`并重置。

```go
// ❌ 问题：循环中使用time.After（每次循环创建新Timer）
for {
    select {
    case <-time.After(sleepDuration):
        // 处理...
    }
}

// ✅ 修正：使用time.NewTimer（复用Timer）
timer := time.NewTimer(sleepDuration)
for {
    select {
    case <-timer.C:
        timer.Reset(sleepDuration)
        // 处理...
    }
}
```

### 11.9 G.EXP.01 - 比较表达式常量放右侧

| 检查项 | 阈值 | 问题类型 | 严重度 | 规则ID |
| --- | --- | --- | --- | --- |
| **比较表达式方向** | 左侧变化右侧不变 | 规范问题 | 一般 | huawei-const-left-value |

**规则**：表达式的比较应遵循左侧倾向于变化、右侧倾向于不变的原则。常量值应放在比较运算符的右侧。

```go
// ❌ 问题：常量在左侧
if "" == item.Obj { ... }

// ✅ 修正：常量在右侧
if item.Obj == "" { ... }
```

### 11.10 G.FMT.01 - 使用gofmt格式化代码

| 检查项 | 阈值 | 问题类型 | 严重度 | 规则ID |
| --- | --- | --- | --- | --- |
| **gofmt缩进规范** | 源码必须经gofmt格式化 | 规范问题 | 一般 | huawei-file-format |

**规则**：源代码必须使用`gofmt`或`goimports`工具进行格式化，确保缩进、对齐、空格等符合Go标准规范。

**常见问题**：
- Tab与空格混用（Go标准缩进为Tab）
- 行内多余空格或空行
- 多行对齐不一致

```bash
# 格式化单个文件
gofmt -w file.go

# 格式化整个项目
gofmt -w src/
```

### 11.11 G.OTH.02 - 禁止代码中包含公网地址

| 检查项 | 阈值 | 问题类型 | 严重度 | 规则ID |
| --- | --- | --- | --- | --- |
| **禁止硬编码公网IP** | 提取为常量或配置 | 安全问题 | 一般 | HardCode_Addr |

**规则**：代码中禁止硬编码公网IP地址，应提取到配置文件中读取。

**推荐修复模式**：将公网IP地址移到项目配置文件（如Beego项目的`app.conf`），代码中通过`beego.AppConfig.DefaultString()`读取，不再出现任何公网IP字符串。

```go
// ❌ 问题：硬编码公网IP常量
const (
    defaultRedisEndpoint        = "135.242.67.210:6379"
    defaultNodeExternalEndpoint = "135.242.67.210:9090"
    defaultNodeHttpsEndpoint    = "41.203.73.2:40051"
    defaultOSSEndpoint          = "135.242.67.210:9000"
)

// ✅ 修正：移到app.conf，代码从配置文件读取
// app.conf:
// [redis]
// endpoint=135.242.67.210:6379
// [node]
// endpoint=135.242.67.210:9090
// httpsendpoint=41.203.73.2:40051
// [oss]
// endpoint=135.242.67.210:9000

// config.go:
Redis: RedisConfig{Endpoint: beego.AppConfig.DefaultString("redis::endpoint", "")},
Node: NodeConfig{
    ExternalEndpoint: beego.AppConfig.DefaultString("node::endpoint", ""),
    HttpsEndpoint:    beego.AppConfig.DefaultString("node::httpsendpoint", ""),
},
OSS: OSSConfig{Endpoint: beego.AppConfig.DefaultString("oss::endpoint", "")},
```

**注意**：测试代码中的公网IP fixture可替换为私有IP（如`10.0.0.x`），生产配置中的公网IP只出现在配置文件中。

---

## 12. CodeCheck规则完整对照表

| 规则ID | 规则名称 | 严重度 | 对应章节 |
| --- | --- | --- | --- |
| G.CMT.01 | 文件头版权声明 | 提示 | 5.1 |
| G.CMT.02 | 包注释 | 提示 | 5.2 |
| G.CMT.04 | 注释符与注释内容间要有1个空格 | 提示 | 11.1 |
| G.CMT.05 | 注释位置 | 提示 | 5.3 |
| G.CON.01 | chan方向限定 | 一般 | 6.1 |
| G.CTL.05 | switch有default | 一般 | 8.4 |
| G.DCL.01 | 魔法数字 | 一般 | 7.1 |
| G.ERR.01 | 错误返回值检查 | 严重 | 8.2 |
| G.ERR.02 | 错误信息不应大写或以标点结尾 | 提示 | 11.2 |
| G.ERR.03 | 包外可见函数禁止panic | 严重 | 11.3 |
| G.EXP.01 | 比较表达式常量放右侧 | 一般 | 11.1 |
| G.FIO.02 | 文件权限指定 | 严重 | 8.3 |
| G.FMT.01 | gofmt缩进规范 | 一般 | 11.10 |
| G.FUN.03 | 避免命名返回值 | 一般 | 11.4 |
| G.NAM.01 | 命名风格统一 | 提示 | 6.4 |
| G.NAM.03 | 导出错误命名 | 提示 | 6.2 |
| G.NAM.05 | 文件命名 | 一般 | 6.3 |
| G.NAM.07 | 包名全小写无下划线 | 一般 | 11.7 |
| G.NAM.08 | 包名与目录名一致 | 一般 | 6.5 |
| G.OTH.02 | 禁止硬编码公网地址 | 一般 | 11.11 |
| G.PKG.02 | 禁止dot导入 | 一般 | 11.6 |
| G.RES.01 | 资源泄漏 | 严重 | 8.1 |
| G.RES.05 | time.After循环优化 | 严重 | 11.8 |
| G.SIT.02 | 类型断言comma-ok | 严重 | 11.5 |

---
