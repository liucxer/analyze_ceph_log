# Ceph 日志分析工具

一个全面的 Ceph OSD 日志分析工具，支持分析 AIO 操作、OSD repop 操作和 OSD 操作。

## 功能特性

- **多类型分析**：支持分析三种类型的操作：
  - AIO (异步 I/O) 操作
  - OSD repop 操作
  - OSD 操作
- **HTML 输出**：生成结构良好的 HTML 报告，带有标签页界面
- **统计摘要**：为每种操作类型显示全面的统计信息
- **筛选功能**：允许按时间范围和持续时间/延迟进行筛选
- **持续时间分布**：显示操作按持续时间/延迟范围的分布情况
- **遗留脚本**：在 `old/` 目录中包含遗留的分析脚本

## 安装

### 前提条件
- Go 1.16 或更高版本

### 设置
1. 克隆仓库：
   ```bash
   git clone https://github.com/liucxer/analyze_ceph_log.git
   cd analyze_ceph_log
   ```

2. 构建项目：
   ```bash
   go build -o analyze_ceph .
   ```

## 使用方法

### 命令语法
```bash
# 使用构建好的二进制文件
./analyze_ceph <log_file> <analysis_type> [output.html]

# 使用 go run
go run analyze_ceph.go <log_file> <analysis_type> [output.html]
```

### 分析类型
- `aio` - 分析 AIO 操作
- `repop` - 分析 OSD repop 操作
- `op` - 分析 OSD 操作
- `all` - 分析所有操作类型

### 示例

1. 分析所有操作并生成 HTML 报告：
   ```bash
   go run analyze_ceph.go /path/to/ceph-osd.log all analysis.html
   ```

2. 仅分析 AIO 操作：
   ```bash
   go run analyze_ceph.go /path/to/ceph-osd.log aio aio_analysis.html
   ```

3. 分析 OSD 操作：
   ```bash
   go run analyze_ceph.go /path/to/ceph-osd.log op osd_analysis.html
   ```

## 输出

该工具生成的 HTML 报告包含：

- **标签页界面**：每种操作类型都有单独的标签页
- **摘要部分**：每个标签页顶部的关键统计信息
- **筛选表单**：每个表格的时间范围和持续时间/延迟筛选
- **数据表格**：详细的操作列表，包含时间戳和持续时间
- **持续时间分布**：按持续时间/延迟范围细分的操作

## 项目结构

```
analyze_ceph_log/
├── analyze_ceph.go       # 主要分析工具
├── go.mod                # Go 模块文件
├── .gitignore            # Git 忽略文件
├── README.md             # 英文 README
├── README_zh.md          # 中文 README
└── old/                  # 遗留分析脚本
    ├── analyze_aio.go
    ├── analyze_osd_op.go
    └── analyze_osd_repop.go
```

## 工作原理

1. **日志解析**：使用正则表达式从 Ceph 日志中提取相关信息
2. **事件匹配**：匹配开始和结束事件以计算持续时间
3. **数据分析**：计算统计信息和持续时间分布
4. **HTML 生成**：创建结构良好的 HTML 报告，带有筛选功能

## 贡献

欢迎贡献！请随时提交 Pull Request。

## 许可证

本项目是开源的，根据 MIT 许可证提供。
