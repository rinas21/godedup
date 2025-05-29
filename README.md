# godedup

> 🔍 Find and list duplicate files on your system using SHA-256 hashes.  
> 💡 Written in Go. Lightweight, fast, and easy to use.

---

## 🛠 Features

- Recursively scans a directory for duplicate files
- Compares files using SHA-256 hashes
- Displays:
  - File paths
  - File size
  - Last modified time
  - Change time (ctime)
- No automatic deletion — safe for manual cleanup

---

## 🚀 Installation

```bash
git clone https://github.com/rinas21/godedup
cd godedup
```

## Run the project

### Build the project

```bash
go build -o godedup dedup.go
```

---

## 📝 Usage

To find duplicate files in specific directories, run:

```bash
go run dedup.go /dir /dir2
```

Or, if you have built the binary:

```bash
./godedup /dir /dir2
```
