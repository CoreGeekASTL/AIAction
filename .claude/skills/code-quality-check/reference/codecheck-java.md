# CodeCheck Java规则参考

CodeCheck Java语言的代码质量检查规则，涵盖格式规范、类型安全、异常处理、日志规范等。

## 1 规则汇总

| 规则ID | 规则名称 | 工具规则名 | 严重度 | 语言 | 描述 |
| --- | --- | --- | --- | --- | --- |
| **G.FMT.05** | NeedBrace | NeedBrace | 提示 | JAVA | 在条件语句和循环块中应该使用大括号 |
| **G.FMT.10** | LineLength | LineLength | 提示 | JAVA | 行宽不超过120个窄字符 |
| **G.FMT.20** | UpperEllRule | UpperEll | 提示 | JAVA | 数字字面量应该设置合适的后缀，long类型应该使用L作为后缀 |
| **G.TYP.08** | Locale | Locale | 一般 | JAVA | 字符串大小写转换、数字格式化为西方数字时，必须加上Locale.ROOT或Locale.ENGLISH |
| **G.ERR.01** | EmptyCatch | EmptyCatch | 一般 | JAVA | 不要通过一个空的catch块忽略异常 |
| **G.ERR.05** | ThrowRawException | ThrowRawException | 一般 | JAVA | 方法抛出的异常，应该与本身的抽象层次相对应 |
| **G.LOG.02** | LogModifier | LogModifier | 一般 | JAVA | 日志工具Logger类的实例必须声明为private static final或者private final |
| **G.LOG.04** | LogWithoutChinese | LogWithoutChinese | 一般 | JAVA | 非仅限于中文区销售产品禁止用中文打印日志 |
| **G.CMT.07** | TodoComment | TodoComment | 提示 | JAVA | 正式交付给客户的代码不应包含TODO/FIXME注释 |
| **G.PRM.07** | ResourceRelease | ResourceRelease | 严重 | JAVA | 进行IO类操作时，必须在try-with-resource或finally里关闭资源 |
| **G.OTH.03** | CodeInComment | CodeInComment | 一般 | JAVA | 不用代码段（含import）应直接删除，不要注释掉 |
| **G.OTH.03** | UnusedImport | UnusedImport | 一般 | JAVA | 不要import未使用的类型 |
| **G.LOG.01** | DoNotLogWithSystemPrint | DoNotLogWithSystemPrint | 一般 | JAVA | 不要使用System.out/err打印日志，应使用日志框架 |
| **G.TYP.09** | AssignCharset | AssignCharset | 一般 | JAVA | 字符与字节的互相转换操作，要指明正确的编码方式 |

---

## 2 G.FMT.05 - 条件语句和循环块必须使用大括号 (NeedBrace)

**规则**：if、else、for、while、do等控制语句的body必须用大括号`{}`包围，即使只有一行语句。

```java
// ❌ 问题：条件语句缺少大括号
if (e.isDirectory()) continue;
if (!name.endsWith(".class")) continue;
if (name.equals("module-info.class")) continue;
if (read == -1) break;
if (read == 0) return null;
if (read < 512) throw new EOFException("invalid tar header");
if (isZeroBlock(header)) return null;
if (b != 0) return false;
if (!isNotExpired) continue;
if (StrUtil.isBlank(value) || value.length() <= 1) continue;
if (!cookie.containsKey("expires") || cookie.isNull("expires")) return true;
if (local == null || local.isEmpty()) return new JSONArray();

// ❌ 问题：测试代码中for循环缺少大括号（线程并发测试常见）
for (Thread t : threads) t.start();
for (Thread t : threads) t.join();

// ✅ 修正：添加大括号
for (Thread t : threads) { t.start(); }
for (Thread t : threads) { t.join(); }

// ✅ 修正：或换行规范写法
for (Thread t : threads) {
    t.start();
}
for (Thread t : threads) {
    t.join();
}

// ✅ 修正：添加大括号
if (e.isDirectory()) { continue; }
if (!name.endsWith(".class")) { continue; }
if (name.equals("module-info.class")) { continue; }
if (read == -1) { break; }
if (read == 0) { return null; }
if (read < 512) { throw new EOFException("invalid tar header"); }
if (isZeroBlock(header)) { return null; }
if (b != 0) { return false; }
if (!isNotExpired) { continue; }
if (StrUtil.isBlank(value) || value.length() <= 1) { continue; }
if (!cookie.containsKey("expires") || cookie.isNull("expires")) { return true; }
if (local == null || local.isEmpty()) { return new JSONArray(); }
```

---

## 3 G.FMT.10 - 行宽不超过120个窄字符 (LineLength)

**规则**：每行代码不应超过120个窄字符（ASCII字符计为1，宽字符计为2）。

**修正策略**：

| 策略 | 适用场景 | 修正方式 |
| --- | --- | --- |
| 方法签名换行 | 构造函数/方法参数过多 | 在参数列表逗号后换行 |
| 长字符串拆分 | JSON模板/URL拼接 | 使用多行字符串或StringBuilder |
| 长链式调用换行 | Stream/API调用 | 在操作符前换行 |
| 注解值换行 | @Value/@Schema注解值过长 | 注解值内换行 |
| 三元表达式换行 | 条件表达式过长 | 在`?`和`:`后换行 |
| catch声明换行 | 多异常类型catch | 在`|`前换行 |

```java
// ❌ 问题：行宽超过120字符
public MuenConfig(List<ChromeConfig> chromeConfigList, List<RouteAppConfig> routeAppConfigList, List<UrlConfig> urlConfigList) {

// ✅ 修正：在逗号后换行
public MuenConfig(List<ChromeConfig> chromeConfigList,
                  List<RouteAppConfig> routeAppConfigList,
                  List<UrlConfig> urlConfigList) {

// ❌ 问题：@Value注解行过长
@Value("#{'${browsergw.report.control-tls-endpoint}'.empty ? '41.203.73.4:30011' : '${browsergw.report.control-tls-endpoint}'}")

// ✅ 修正：注解值内换行
@Value("#{'${browsergw.report.control-tls-endpoint}'.empty ? "
        + "'41.203.73.4:30011' : '${browsergw.report.control-tls-endpoint}'}")

// ❌ 问题：catch声明行过长
} catch (InstantiationException | IllegalAccessException | InvocationTargetException | NoSuchMethodException e) {

// ✅ 修正：在|前换行
} catch (InstantiationException
        | IllegalAccessException
        | InvocationTargetException
        | NoSuchMethodException e) {

// ❌ 问题：三元表达式行过长
return this.privateKeyPassword == null ? new byte[0] : Arrays.copyOf(this.privateKeyPassword, this.privateKeyPassword.length);

// ✅ 修正：拆分三元表达式
if (this.privateKeyPassword == null) {
    return new byte[0];
}
return Arrays.copyOf(this.privateKeyPassword, this.privateKeyPassword.length);
```

---

## 4 G.FMT.20 - long类型数字字面量使用L后缀 (UpperEll)
**规则**：浮点型变量值应该以`f`结尾，Long类型变量值应该以`L`结尾（而非小写`l`，避免与数字1混淆）。

```java
// ❌ 问题：long字面量缺少L后缀
long waitTime = 30000; // 30s间隔

// ✅ 修正：添加L后缀
long waitTime = 30000L; // 30s间隔

// ✅ 更佳：提取为常量并使用下划线分隔提高可读性
private static final long RETRY_INTERVAL_MS = 30_000L;
```

---

## 5 G.TYP.08 - 字符串大小写转换和数字格式化必须指定Locale (Locale)

**规则**：
- 使用`String.toLowerCase()`或`String.toUpperCase()`时，必须传入`Locale.ROOT`或`Locale.ENGLISH`
- 使用`String.format()`格式化数字（%d、%f等）时，必须传入`Locale.ROOT`作为第一个参数

**原因**：不指定Locale时，大小写转换和数字格式化结果依赖于系统默认Locale。例如在土耳其Locale下，`"i".toLowerCase()`会返回`"ı"`（无点i），而非`"i"`；在德语Locale下，`String.format("%d", 1000)`可能输出`"1.000"`而非`"1000"`。

```java
// ❌ 问题：toLowerCase未指定Locale
switch (metricType.toLowerCase()) { ... }
if (!sourceJsonPath.toLowerCase().endsWith(".json")) { ... }

// ✅ 修正：添加Locale.ROOT
switch (metricType.toLowerCase(Locale.ROOT)) { ... }
if (!sourceJsonPath.toLowerCase(Locale.ROOT).endsWith(".json")) { ... }

// ❌ 问题：String.format格式化数字未指定Locale
getLogger().info(String.format("tcp started on port(s): %d", port));
throw new IOException(String.format("data length is: %d, but read size: %d", buffer.length, totalRead));

// ✅ 修正：添加Locale.ROOT作为第一个参数
getLogger().info(String.format(Locale.ROOT, "tcp started on port(s): %d", port));
throw new IOException(String.format(Locale.ROOT, "expected %d bytes but read %d", buffer.length, totalRead));
```

**import声明**：使用Locale时需要添加`import java.util.Locale;`

---

## 6 G.ERR.01 - 不要通过空的catch块忽略异常 (EmptyCatch)

**规则**：catch块不应为空，catch块可以只有非空的单行注释或多行注释。忽略异常可能导致问题难以排查。

```java
// ❌ 问题：catch块为空（使用ignored变量名掩盖忽略行为）
} catch (ParseException ignored) {
}

// ✅ 修正：至少添加日志记录
} catch (ParseException e) {
    log.warn("date parse failed, using raw value: {}", rawValue, e);
}

// ✅ 也可接受：如果确实不需要处理，添加说明注释
} catch (ParseException e) {
    // date parsing is optional, fall through to default handling
}
```

---

## 7 G.ERR.05 - 方法抛出的异常应与抽象层次相对应 (ThrowRawException)

**规则**：不要直接抛出`RuntimeException`这种底层异常，应使用与业务层次匹配的自定义异常类。

项目已提供以下自定义异常：

| 异常类 | 适用场景 | 包路径 |
| --- | --- | --- |
| `BrowserGatewayException` | 通用网关业务异常（浏览器创建失败、锁超时、文件操作失败等） | `com.huawei.browsergateway.exception.BrowserGatewayException` |
| `CodecException` | 音视频编解码处理异常（编码/解码失败、上下文分配失败等） | `com.huawei.browsergateway.exception.CodecException` |
| `IllegalArgumentException` | 输入参数校验异常（无效参数、无效用户绑定信息等） | JDK内置`java.lang.IllegalArgumentException` |
| `IllegalStateException` | 状态校验异常（驱动为null、资源不可用等） | JDK内置`java.lang.IllegalStateException` |

```java
// ❌ 问题：直接抛出RuntimeException
throw new RuntimeException("failed to create browsers");
throw new RuntimeException("failed to create browsers" + e.getMessage());
throw new RuntimeException(e);

// ✅ 修正：使用匹配业务层次的自定义异常
throw new BrowserGatewayException("failed to create browser");
throw new BrowserGatewayException("failed to create browser", e);  // 保留异常链
throw new BrowserGatewayException(e);                               // 仅包装cause

// ✅ 编解码场景使用CodecException
throw new CodecException("Cannot find audio stream");
throw new CodecException("Failed to allocate AVIO context");

// ✅ 参数校验场景使用IllegalArgumentException
throw new IllegalArgumentException("lcd width or height is invalid");
throw new IllegalArgumentException("user token is invalid");

// ✅ 状态校验场景使用IllegalStateException
throw new IllegalStateException("muenDriver is null");
```

**自定义异常类结构**：

```java
// BrowserGatewayException
public class BrowserGatewayException extends RuntimeException {
    public BrowserGatewayException(String message) { super(message); }
    public BrowserGatewayException(String message, Throwable cause) { super(message, cause); }
    public BrowserGatewayException(Throwable cause) { super(cause); }
}

// CodecException
public class CodecException extends RuntimeException {
    public CodecException(String message) { super(message); }
    public CodecException(String message, Throwable cause) { super(message, cause); }
    public CodecException(Throwable cause) { super(cause); }
}
```

**import声明**：使用自定义异常时需要添加对应import：
- `import com.huawei.browsergateway.exception.BrowserGatewayException;`
- `import com.huawei.browsergateway.exception.CodecException;`

---

## 8 G.LOG.02 - Logger实例必须声明为private static final (LogModifier)

**规则**：日志工具Logger类的实例必须声明为`private static final`，确保Logger在类生命周期内不变且不可被外部修改。

```java
// ❌ 问题：Logger缺少final修饰符
private static Logger log = LoggerFactory.getLogger(MyClass.class);

// ✅ 修正：添加final修饰符
private static final Logger log = LoggerFactory.getLogger(MyClass.class);

// ✅ Log4j2风格
private static final Logger log = LogManager.getLogger(MyClass.class);
```

---

## 9 G.LOG.04 - 日志中禁止使用中文字符 (LogWithoutChinese)

**规则**：非仅限于中文区销售的产品，日志输出中不应包含中文字符（包括中文标点如全角冒号`：`、全角逗号`，`等）。

```java
// ❌ 问题：日志中使用中文标点（全角冒号：）
log.info("BrowserOptions full config：\n{}", JSON.toJSONString(options));

// ✅ 修正：使用英文标点（半角冒号:）
log.info("BrowserOptions full config:\n{}", JSON.toJSONString(options));
```

---

## 10 G.CMT.07 - 正式交付代码不应包含TODO/FIXME注释 (TodoComment)

**规则**：正式交付给客户的代码不应包含`TODO`或`FIXME`注释。如需保留待办事项，应在内部项目管理工具中记录。

```java
// ❌ 问题：代码中包含TODO注释
// TODO 可以优化掉，不需要获取userbind
// todo：待增加trunck和fabric平面
// todo:定义常量

// ✅ 修正：删除TODO注释，如需跟踪改进项使用项目管理工具
// （删除整行TODO注释，或替换为不包含TODO的说明性注释）
```

---

## 11 G.PRM.07 - IO资源必须在try-with-resources或finally中关闭 (ResourceRelease)

**规则**：进行IO类操作时，必须在try-with-resources语句或finally块中关闭资源，防止资源泄漏。

```java
// ❌ 问题：BufferedReader未在try-with-resources中关闭
BufferedReader reader = new BufferedReader(new InputStreamReader(process.getInputStream()));
String line;
while ((line = reader.readLine()) != null) {
    log.info("Script output: {}", line);
}
// reader未关闭！

// ✅ 修正：使用try-with-resources自动关闭
try (BufferedReader reader = new BufferedReader(
        new InputStreamReader(process.getInputStream()))) {
    String line;
    while ((line = reader.readLine()) != null) {
        log.info("Script output: {}", line);
    }
}
```

---

## 12 G.OTH.03 - 不用代码段（含import）应直接删除，不要注释掉 (CodeInComment / UnusedImport)

**规则**：不用代码段（含import）应直接删除，不要注释掉。未使用的import也应删除。

```java
// ❌ 问题：注释掉的代码段
// import cn.hutool.json.JSONUtil;
// extensionPath = "/opt/host/extension/record";
// return 0;
// callback.onCertificateUpdate(certUpdateDTO);

// ❌ 问题：未使用的import
import org.springframework.beans.factory.annotation.Autowired;
import java.util.HashSet;
import com.huawei.browsergateway.BrowserGatewayApplication;

// ❌ 问题：测试文件中未使用的import（测试代码同样受CodeCheck约束）
import java.util.Optional;                                                      // DevToolsProxyTest.java
import org.springframework.beans.factory.config.BeanExpressionContext;          // ServerEndpointExporterTest.java

// ✅ 修正：直接删除注释代码和未用import（不做任何保留）
```

**注意**：注释掉的代码和未用import应彻底删除，不保留任何痕迹。如需标记后续改进，应在项目管理工具中记录。**测试代码（src/test/java）同样受CodeCheck约束**，测试文件中的未使用import也必须删除。

---

## 13 G.LOG.01 - 不要使用System.out/err打印日志 (DoNotLogWithSystemPrint)

**规则**：记录日志应该使用Facade模式的日志框架（如SLF4J/Log4j2），不要使用System.out.println或System.err.println打印日志。

```java
// ❌ 问题：使用System.out/err打印日志
System.out.println("Framework.initCsp initialization started");
System.out.println("application.yaml path is " + ap);
System.err.println("JSON file compression failed: " + e.getMessage());

// ✅ 修正：使用日志框架
private static final Logger log = LogManager.getLogger(MyClass.class);
log.info("Framework.initCsp initialization started");
log.info("application.yaml path is {}", ap);
log.error("JSON file compression failed: {}", e.getMessage());
```

**注意**：
- Logger必须声明为`private static final`（参见G.LOG.02）
- 日志中禁止使用中文标点（参见G.LOG.04）
- 使用`{}`占位符而非字符串拼接，避免不必要的字符串构造开销

---

## 14 G.TYP.09 - 字符与字节互相转换必须指定编码 (AssignCharset)

**规则**：字符与字节的互相转换操作，必须指明正确的编码方式。包括`String.getBytes()`、`new String(byte[])`、`new InputStreamReader(InputStream)`等方法。

```java
// ❌ 问题：未指定charset（使用平台默认编码，行为不确定）
caContent.getBytes()
new String(certEntity.getPrivateKeyPassword())
new InputStreamReader(process.getInputStream())
((String) value).getBytes()

// ❌ 问题：测试代码中同样违规（测试setup和helper方法）
certEntity.setPrivateKeyPassword("".getBytes());                      // 测试setup
initialCert.setPrivateKeyPassword("".getBytes());                     // 测试setup
new BufferedReader(new InputStreamReader(inputStream))                // 测试helper

// ✅ 修正：显式指定StandardCharsets.UTF_8
caContent.getBytes(StandardCharsets.UTF_8)
new String(certEntity.getPrivateKeyPassword(), StandardCharsets.UTF_8)
new InputStreamReader(process.getInputStream(), StandardCharsets.UTF_8)
((String) value).getBytes(StandardCharsets.UTF_8)
```

**import声明**：使用`StandardCharsets.UTF_8`需添加：
- `import java.nio.charset.StandardCharsets;`

**注意**：优先使用`StandardCharsets.UTF_8`而非字符串`"UTF-8"`，前者是编译期安全常量，后者可能在运行时抛出UnsupportedEncodingException。

---

## 15 已修复案例汇总（Java）

| # | 文件 | 方法 | 原规则 | 修正方式 |
| --- | --- | --- | --- | --- |
| 1 | `MuenPluginClassLoader.java` | listAllClassName() | G.FMT.05 NeedBrace | 3处if+continue添加大括号 |
| 2 | `GzipUtil.java` | unGzip()/readNextEntry()/isZeroBlock() | G.FMT.05 NeedBrace | 5处单行if添加大括号 |
| 3 | `UserdataSlimmer.java` | slimCookies()/slimLocalStorage() | G.FMT.05 NeedBrace | 3处if+continue/return添加大括号 |
| 4 | `ServiceReporter.java` | reportChainInfoWithRetry() | G.FMT.20 UpperEll | 30000 → 提取为常量RETRY_INTERVAL_MS=30_000L |
| 5 | `AbstractTcpServer.java` | start() | G.TYP.08 Locale | String.format → String.format(Locale.ROOT, ...) |
| 6 | `ZstdUtil.java` | compressJson() | G.TYP.08 Locale | toLowerCase() → toLowerCase(Locale.ROOT) |
| 7 | `Tlv.java` | readAndCheckSize() | G.TYP.08 Locale | String.format → String.format(Locale.ROOT, ...) |
| 8 | `ClientImpl.java` | BrowserImpl.create()/ContextImpl.create() | G.ERR.05 ThrowRawException | RuntimeException → BrowserGatewayException (4处) |
| 9 | `WebElementImpl.java` | convertDate() | G.ERR.05+G.ERR.01 | RuntimeException→IllegalArgumentException; 空catch添加log.warn |
| 10 | `ExtensionManageService.java` | decompress()/findJarPath() | G.ERR.05 ThrowRawException | RuntimeException → BrowserGatewayException (3处) |
| 11 | `ChromeSetImpl.java` | create() | G.ERR.05 ThrowRawException | RuntimeException → BrowserGatewayException |
| 12 | `HWCallbackImpl.java` | getFile() | G.ERR.05 ThrowRawException | RuntimeException → BrowserGatewayException(保留cause链) |
| 13 | `RemoteImpl.java` | createChrome() | G.ERR.05 ThrowRawException | RuntimeException → BrowserGatewayException |
| 14 | `ControlTcpServerHandler.java` | processLogin() | G.ERR.05 ThrowRawException | RuntimeException → IllegalArgumentException (3处) |
| 15 | `MediaTcpServerHandle.java` | processLogin() | G.ERR.05 ThrowRawException | RuntimeException → IllegalArgumentException (2处) |
| 16 | `HttpUtil.java` | request() | G.ERR.05 ThrowRawException | RuntimeException → BrowserGatewayException(保留cause链) |
| 17 | `UserdataSlimmer.java` | slimInplace() | G.ERR.05 ThrowRawException | RuntimeException → BrowserGatewayException |
| 18 | `AudioCodecProcessor.java` | 多个方法 | G.ERR.05 ThrowRawException | RuntimeException → CodecException (15处) |
| 19 | `FfmpegCodecService.java` | initInputContext() | G.ERR.05 ThrowRawException | RuntimeException → CodecException (5处) |
| 20 | `VideoCodecProcessor.java` | 多个方法 | G.ERR.05 ThrowRawException | RuntimeException → CodecException (11处) |
| 21 | `BrowserProxyLogDump.java` | executeShellScript() | G.PRM.07+G.TYP.09 | BufferedReader使用try-with-resources; InputStreamReader指定UTF-8 |
| 22 | `NetWorkInterfaceCheck.java` | 字段注释 | G.CMT.07 TodoComment | 删除TODO注释 |
| 23 | `AudioCodecProcessor.java` | setEncodeContext() | G.CMT.07 TodoComment | 删除todo注释 |
| 24 | `RemoteImpl.java` | createInstance() | G.OTH.03 CodeInComment | 删除注释掉的extensionPath代码 |
| 25 | `CertInfo.java` | SetDeviceContent() | G.OTH.03 CodeInComment | 删除注释掉的SetDeviceContent代码块(6行) |
| 26 | `FfmpegStreamProcessor.java` | readHook() | G.OTH.03 CodeInComment | 删除注释掉的return 0 |
| 27 | `ApplicationConfig.java` | propertySourcesPlaceholderConfigurer() | G.LOG.01 DoNotLogWithSystemPrint | System.out.println → log.info |
| 28 | `ZstdUtil.java` | compressJson()/decompressJson() | G.LOG.01+G.LOG.04 | System.err.println → log.error; 中文→英文 |
| 29 | `CertInfo.java` | Ca()/Device()/convertPkcs1ToPkcs8Stream() | G.TYP.09 AssignCharset | getBytes() → getBytes(StandardCharsets.UTF_8) (3处) |
| 30 | `TlvCodec.java` | setFieldValue() | G.TYP.09 AssignCharset | new String(data) → new String(data, StandardCharsets.UTF_8); getBytes() → getBytes(StandardCharsets.UTF_8) |
| 31 | 约21个文件 | 多个方法/字段 | G.FMT.10 LineLength | 长行拆分(方法签名/注解/字符串/表达式等) |
| 32 | `BrowserProxyLogDump.java` | executeShellScript() | G.FMT.10 LineLength | try-with-resources长行拆分为2行(135→2行) |
| 33 | `Tlv.java` | readAndCheckSize() | G.FMT.10 LineLength | String.format长行拆分(129→2行) |
| 34 | `TlvCodec.java` | setFieldValue() | G.FMT.10 LineLength | IllegalArgumentException含中文长行拆分(121窄字符→2行) |
| 35 | `Config.java`等5个文件 | - | G.OTH.03 UnusedImport | 删除5个未使用import(Autowired/JsonProperty/Objects/BrowserGatewayApplication/ConcurrentHashMap+Map+Session) |
| 36 | `BrowserCloserTask.java` | - | G.CMT.07 TodoComment | 删除TODO注释 |
| 37 | `DateTimeUtilTest.java` | millisToDate() | G.LOG.01 SystemPrint | System.out.println替换为log/断言(4处) |
| 38 | `UserChrome.java` | createBrowser() | G.LOG.04 LogWithoutChinese | 日志中中文冒号：替换为英文冒号: |

**待修复（测试代码CodeCheck问题）**：

| # | 文件 | 方法/位置 | 行号 | 规则 | 问题描述 | 修正方式 |
| --- | --- | --- | --- | --- | --- | --- |
| T1 | `DevToolsProxyTest.java` | import区 | 14 | G.OTH.03 UnusedImport | `import java.util.Optional;`未使用 | 删除该import |
| T2 | `ServerEndpointExporterTest.java` | import区 | 8 | G.OTH.03 UnusedImport | `import org.springframework.beans.factory.config.BeanExpressionContext;`未使用 | 删除该import |
| T3 | `CertInfoTest.java` | testSetDeviceContent() | 69 | G.TYP.09 AssignCharset | `"".getBytes()`未指定编码 | `"".getBytes(StandardCharsets.UTF_8)` |
| T4 | `CertInfoTest.java` | testSetDeviceContentWithNull() | 86 | G.TYP.09 AssignCharset | `"".getBytes()`未指定编码 | `"".getBytes(StandardCharsets.UTF_8)` |
| T5 | `CertInfoTest.java` | testIsCertReadyWhenReady() | 114 | G.TYP.09 AssignCharset | `"".getBytes()`未指定编码 | `"".getBytes(StandardCharsets.UTF_8)` |
| T6 | `CertInfoTest.java` | testDeviceInputStream() | 140 | G.TYP.09 AssignCharset | `"".getBytes()`未指定编码 | `"".getBytes(StandardCharsets.UTF_8)` |
| T7 | `CertInfoTest.java` | readInputStream() | 152 | G.TYP.09 AssignCharset | `new InputStreamReader(inputStream)`未指定编码 | `new InputStreamReader(inputStream, StandardCharsets.UTF_8)` |
| T8 | `DateTimeUtilTest.java` | testMillisToDate_threadSafety() | 42 | G.FMT.05 NeedBrace | `for (Thread t : threads) t.start();`缺少大括号 | `for (Thread t : threads) { t.start(); }` |
| T9 | `DateTimeUtilTest.java` | testMillisToDate_threadSafety() | 43 | G.FMT.05 NeedBrace | `for (Thread t : threads) t.join();`缺少大括号 | `for (Thread t : threads) { t.join(); }` |

**未修复（FORBIDDEN文件）**：

| # | 文件 | 规则 | 原因 |
| --- | --- | --- | --- |
| 1 | `BrowserGatewayApplication.java` | G.LOG.02 LogModifier | 文件禁止修改 |
| 2 | `BrowserGatewayApplication.java` | G.LOG.01 SystemPrint | 文件禁止修改 |
| 3 | `BrowserGatewayApplication.java` | G.CMT.07 TodoComment | 文件禁止修改 |
| 4 | `CseImpl.java` | G.OTH.03 UnusedImport | 文件禁止修改 |
| 5 | `CustomResourceMonitorAdapter.java` | G.TYP.08 Locale | adapter目录禁止修改 |
| 6 | `CertEntity.java` | G.FMT.10+G.TYP.09 | adapter目录禁止修改 |
| 7 | `CspAlarmAdapter.java` | G.FMT.10 LineLength | adapter目录禁止修改 |
| 8 | `CspCertificateAdapter.java` | G.FMT.10+G.TYP.09 | adapter目录禁止修改 |
| 9 | `CustomCertificateAdapter.java` | G.TYP.09 AssignCharset | adapter目录禁止修改 |

---

## 16 测试代码常见CodeCheck问题（src/test/java）

**重要**：测试代码（`src/test/java`目录下的文件）同样受CodeCheck约束。以下为测试代码中高频出现的CodeCheck违规模式，检查时必须覆盖测试文件。

### 16.1 测试文件CodeCheck检查清单

| 检查项 | 规则 | 常见场景 | 检查方法 |
| --- | --- | --- | --- |
| **测试文件未使用import** | G.OTH.03 UnusedImport | 测试编写过程中import了类型但最终未使用（如Optional、Spring内部类等） | 逐个检查import是否在测试方法中被引用 |
| **测试setup中getBytes()无编码** | G.TYP.09 AssignCharset | `certEntity.setXxx("".getBytes())`在测试setup中设置空密码/空内容 | 搜索`getBytes()`替换为`getBytes(StandardCharsets.UTF_8)` |
| **测试helper中InputStreamReader无编码** | G.TYP.09 AssignCharset | `new BufferedReader(new InputStreamReader(inputStream))`在readInputStream等helper方法中 | 搜索`new InputStreamReader(`检查是否有charset参数 |
| **测试并发代码for循环无大括号** | G.FMT.05 NeedBrace | 线程并发测试中`for (Thread t : threads) t.start();`单行写法 | 搜索`for.*).*;`单行for循环模式 |

### 16.2 测试代码违规模式与修正

```java
// === G.OTH.03 UnusedImport - 测试文件未使用import ===
// ❌ 问题：DevToolsProxyTest.java导入了Optional但未使用
import java.util.Optional;
// ❌ 问题：ServerEndpointExporterTest.java导入了BeanExpressionContext但未使用
import org.springframework.beans.factory.config.BeanExpressionContext;
// ✅ 修正：删除未使用的import行

// === G.TYP.09 AssignCharset - 测试setup中getBytes()无编码 ===
// ❌ 问题：CertInfoTest.java多个测试方法中设置空密码未指定编码
certEntity.setPrivateKeyPassword("".getBytes());
initialCert.setPrivateKeyPassword("".getBytes());
// ✅ 修正：指定StandardCharsets.UTF_8
certEntity.setPrivateKeyPassword("".getBytes(StandardCharsets.UTF_8));
initialCert.setPrivateKeyPassword("".getBytes(StandardCharsets.UTF_8));

// === G.TYP.09 AssignCharset - 测试helper中InputStreamReader无编码 ===
// ❌ 问题：CertInfoTest.java readInputStream() helper方法
try (BufferedReader reader = new BufferedReader(new InputStreamReader(inputStream))) {
// ✅ 修正：指定StandardCharsets.UTF_8
try (BufferedReader reader = new BufferedReader(
        new InputStreamReader(inputStream, StandardCharsets.UTF_8))) {

// === G.FMT.05 NeedBrace - 测试并发代码for循环无大括号 ===
// ❌ 问题：DateTimeUtilTest.java线程安全测试
for (Thread t : threads) t.start();
for (Thread t : threads) t.join();
// ✅ 修正：添加大括号
for (Thread t : threads) {
    t.start();
}
for (Thread t : threads) {
    t.join();
}
```

### 16.3 测试代码检查步骤

1. **扫描测试文件目录**：`src/test/java/**/*.java`下所有文件
2. **检查import区**：逐个验证import是否在测试方法体中被引用
3. **搜索getBytes()**：`grep "getBytes()" --include="*.java" src/test/`
4. **搜索InputStreamReader**：`grep "new InputStreamReader" --include="*.java" src/test/`
5. **搜索单行for循环**：`grep "for.*).*;" --include="*.java" src/test/`
6. **搜索单行if语句**：`grep "if.*).*;" --include="*.java" src/test/`

---

## 17 Java CodeCheck检查步骤

| 步骤 | 操作 | 方法 |
| --- | --- | --- |
| 1 | 从CodeCheck报告获取问题列表 | 按规则ID分类整理 |
| 2 | 按规则优先级排序 | 严重→一般→提示 |
| 3 | 逐文件逐行修复 | 参照修正示例和自定义异常表 |
| 4 | 编译验证 | `mvn compile` 确认无编译错误 |
| 5 | 确认import声明 | 检查新增Locale/自定义异常import是否完整 |
