<p align="center">
  <img src="assets/logo.png" alt="MetaGo Logo" width="200"/>
</p>

# MetaGo · 🧰📁

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Build](https://img.shields.io/badge/build-passing-brightgreen)]()

A minimal, concurrent CLI tool to extract file metadata and hashes (SHA256/MD5) from files or directories.

---

## ✨ Features

- Extract filename, size, modified time.
- Compute SHA256 or MD5 hashes.
- Directory traversal with goroutines.
- Output in plaintext or JSON (WIP).
- Simple, clean structure with extensibility in mind.

---

## 📦 Installation

### 🔧 Option 1: Build from source

```bash
git clone https://github.com/0xM3H51N/MetaGo.git
cd MetaGo
go build -o metago
```

### 📥 Option 2: Install via Go

```bash
go install github.com/0xM3H51N/MetaGo@latest
```

Make sure `GOBIN` is in your `$PATH`.

---

## 🚀 Usage

```bash
metago --file /path/to/file.txt
metago --dir /path/to/folder --hash md5
metago --file /path/to/file --json
```

### Flags

| Flag       | Description                         | Default  |
|------------|-------------------------------------|----------|
| `--file`   | Path to a single file               |          |
| `--dir`    | Path to a directory                 |          |
| `--json`   | Output in JSON format               | `false`  |
| `--hash`   | Hashing algorithm (md5 or sha256)   | `sha256` |

> Only one of `--file` or `--dir` can be used at a time.

---

## 🧪 Output (Example)

```
Name: report.pdf
Size: 1.2 MB
Hash: 98af71c...
ModTime: 2025-07-20 13:12:01 +0100
```

---

## 📁 Project Structure

<pre>
MetaGo/
├── cmd/         # CLI commands and orchestration
├── internal/    # Internal packages (hashing)
├── core/        # Shared types/interfaces
├── main.go      # Entry point
├── go.mod
└── README.md
</pre>

---

## 🛠️ TODO

- ✅ Metadata extraction
- ✅ File hashing
- ✅ JSON output
- [ ] Recursive directory walk
- [ ] File type detection
- [ ] Unit tests

---

## 🧑‍💻 Author

**[@0xM3H51N](https://github.com/0xM3H51N)**

---

Licensed under the MIT License.
