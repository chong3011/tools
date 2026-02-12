# GDBus 通讯：会话总线 vs 系统总线

## 概述

GDBus 是 GLib 中用于 D-Bus（桌面总线）通讯的高级 API。D-Bus 提供了两种主要的进程间通信（IPC）总线类型：**会话总线（Session Bus）** 和 **系统总线（System Bus）**。理解这两者之间的差别对于正确的应用程序设计至关重要。

## 会话总线 (Session Bus)

### 定义
**会话总线** 是每个用户登录会话专属的消息总线。每个用户会话都有自己独立的会话总线实例。

### 特征
- **作用域**: 单个用户会话
- **生命周期**: 仅在用户登录会话期间存在
- **访问权限**: 仅可被属于同一用户会话的进程访问
- **安全性**: 不同用户会话之间相互隔离

### 使用场景
会话总线适用于：
1. **桌面应用程序** - 用户级 GUI 应用程序之间的通信
2. **用户专属服务** - 媒体播放器、通知守护进程、剪贴板管理器
3. **应用程序集成** - 文件管理器与桌面环境的通信
4. **用户设置** - 用户偏好设置的配置服务

### 会话总线上的示例服务
- 音乐播放器（MPRIS 接口）
- 通知服务
- 屏幕保护程序
- 剪贴板管理器
- 用户级应用程序启动器

### 代码示例（C 语言使用 GDBus）
```c
// 连接到会话总线
GDBusConnection *connection = g_bus_get_sync(
    G_BUS_TYPE_SESSION,  // 会话总线类型
    NULL,
    &error
);

// 示例：通过会话总线发送通知
g_dbus_connection_call_sync(
    connection,
    "org.freedesktop.Notifications",  // 知名名称
    "/org/freedesktop/Notifications", // 对象路径
    "org.freedesktop.Notifications",  // 接口
    "Notify",                          // 方法
    parameters,
    G_VARIANT_TYPE("(u)"),
    G_DBUS_CALL_FLAGS_NONE,
    -1,
    NULL,
    &error
);
```

## 系统总线 (System Bus)

### 定义
**系统总线** 是系统范围的消息总线，被所有用户会话和系统服务共享。

### 特征
- **作用域**: 系统范围（所有用户）
- **生命周期**: 从系统启动到关闭
- **访问权限**: 由安全策略控制（通常需要 root/管理员权限）
- **安全性**: 通过策略文件进行严格的权限控制

### 使用场景
系统总线适用于：
1. **系统服务** - 硬件管理、网络配置
2. **设备监控** - USB 设备检测、磁盘挂载
3. **系统范围配置** - 网络设置、电源管理
4. **硬件抽象** - 打印机服务、音频系统
5. **系统管理** - 包管理器、系统更新

### 系统总线上的示例服务
- NetworkManager（网络配置）
- UDisks2（磁盘管理）
- systemd（服务管理）
- BlueZ（蓝牙协议栈）
- PulseAudio 系统实例

### 代码示例（C 语言使用 GDBus）
```c
// 连接到系统总线
GDBusConnection *connection = g_bus_get_sync(
    G_BUS_TYPE_SYSTEM,  // 系统总线类型
    NULL,
    &error
);

// 示例：通过系统总线查询 NetworkManager
GVariant *result = g_dbus_connection_call_sync(
    connection,
    "org.freedesktop.NetworkManager",     // 知名名称
    "/org/freedesktop/NetworkManager",    // 对象路径
    "org.freedesktop.DBus.Properties",    // 接口
    "Get",                                 // 方法
    g_variant_new("(ss)", 
                  "org.freedesktop.NetworkManager", 
                  "State"),
    G_VARIANT_TYPE("(v)"),
    G_DBUS_CALL_FLAGS_NONE,
    -1,
    NULL,
    &error
);
```

## 主要差别对比

| 特性 | 会话总线 | 系统总线 |
|------|---------|---------|
| **作用域** | 单个用户会话 | 系统范围 |
| **生命周期** | 登录到登出 | 启动到关闭 |
| **可访问性** | 单个用户会话 | 所有用户（需权限）|
| **使用场景** | 用户应用程序 | 系统服务 |
| **安全性** | 用户级隔离 | 系统级策略 |
| **权限要求** | 会话内通常宽松 | 严格，常需特权 |
| **示例服务** | 媒体播放器、通知 | NetworkManager、磁盘管理 |
| **配置位置** | `$XDG_RUNTIME_DIR/bus` 或 `$DBUS_SESSION_BUS_ADDRESS` | `/var/run/dbus/system_bus_socket` |

## 安全考虑

### 会话总线
- **限制较少** - 用户会话内的任何应用程序通常都可以通信
- **用户边界** - 与其他用户会话隔离
- **信任模型** - 同一会话内的应用程序相互信任

### 系统总线
- **高度限制** - 访问由策略文件控制（`/etc/dbus-1/system.d/`）
- **权限要求** - 许多操作需要 root 或特定权限
- **策略强制** - 每个方法调用都可以被安全策略控制

### 策略示例
```xml
<!-- 系统总线策略文件示例 -->
<busconfig>
  <policy user="root">
    <allow own="org.example.SystemService"/>
    <allow send_destination="org.example.SystemService"/>
  </policy>
  
  <policy context="default">
    <deny own="org.example.SystemService"/>
    <deny send_destination="org.example.SystemService"/>
  </policy>
</busconfig>
```

## 如何选择正确的总线

### 使用会话总线的情况：
- 构建面向用户的应用程序
- 管理用户特定的数据或偏好设置
- 不需要系统范围的作用域
- 不需要在用户会话之间持久化

### 使用系统总线的情况：
- 构建系统服务
- 管理硬件或系统资源
- 需要跨用户会话通信
- 需要系统范围的可见性
- 需要服务在用户登出后继续运行

## 常见陷阱

1. **不必要地使用系统总线** - 不要为用户应用程序使用系统总线；这会增加复杂性和安全限制
2. **权限不足** - 系统总线操作可能因缺乏适当策略而失败
3. **服务发现错误** - 一个总线上的服务无法从另一个总线访问
4. **会话总线假设** - 会话总线可能在无头系统或系统服务中不存在

## 最佳实践

1. **适当选择** - 根据服务的作用域选择总线类型
2. **处理连接失败** - 为总线连接实现适当的错误处理
3. **遵守安全策略** - 不要试图绕过系统总线的安全限制
4. **清理资源** - 正确释放总线连接和资源
5. **在目标环境中测试** - 验证部署场景中的总线可用性

## 环境变量

### 会话总线
```bash
# 会话总线地址（在用户会话中自动设置）
echo $DBUS_SESSION_BUS_ADDRESS
# 示例：unix:path=/run/user/1000/bus
```

### 系统总线
```bash
# 系统总线套接字（固定位置）
# /var/run/dbus/system_bus_socket
# 或
# unix:path=/var/run/dbus/system_bus_socket
```

## 测试和调试

### 列出总线上的服务
```bash
# 列出会话总线服务
dbus-send --session --dest=org.freedesktop.DBus \
    --type=method_call --print-reply \
    /org/freedesktop/DBus org.freedesktop.DBus.ListNames

# 列出系统总线服务
dbus-send --system --dest=org.freedesktop.DBus \
    --type=method_call --print-reply \
    /org/freedesktop/DBus org.freedesktop.DBus.ListNames
```

### 监控总线流量
```bash
# 监控会话总线
dbus-monitor --session

# 监控系统总线（需要权限）
sudo dbus-monitor --system
```

## 结论

在 GDBus 通讯中选择会话总线还是系统总线取决于应用程序的作用域和需求：

- **会话总线**：用户级应用程序、桌面集成、用户专属服务
- **系统总线**：系统服务、硬件管理、跨用户功能

理解这些差别可以确保基于 D-Bus 的应用程序具有正确的架构和安全性。

## 参考资料

- [D-Bus 规范](https://dbus.freedesktop.org/doc/dbus-specification.html)
- [GDBus API 参考](https://docs.gtk.org/gio/class.DBusConnection.html)
- [D-Bus 教程](https://dbus.freedesktop.org/doc/dbus-tutorial.html)

## 快速参考表

### 什么时候用会话总线？
✓ 桌面应用程序  
✓ 用户级服务  
✓ 会话内通信  
✓ 不需要 root 权限  
✓ 用户登出后无需运行  

### 什么时候用系统总线？
✓ 系统服务  
✓ 硬件访问  
✓ 跨用户通信  
✓ 需要系统级持久性  
✓ 需要管理员权限  

## 总结

**会话总线**和**系统总线**各有其特定的用途：

- **会话总线**专为用户应用程序和桌面集成设计，提供简单的会话内进程通信
- **系统总线**为系统级服务设计，提供受控的跨用户进程通信和硬件访问

正确理解和使用这两种总线类型是开发稳定、安全的 D-Bus 应用程序的关键。
