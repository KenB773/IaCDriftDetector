# IaC Drift Detector

A lightweight CLI tool that identifies **infrastructure drift** between your live AWS environment and your Terraform state files. Built in Go (and therefore my first project using a non-Python language!), this tool helps DevSecOps engineers, SREs, and cloud security professionals detect unmanaged changes, enforce compliance, and improve IaC governance.

---

## Features

- ✅ **AWS resource comparison** against `.tfstate`
- 🔍 Detects:
  - Missing (deleted) resources
  - Extra (unmanaged) resources
  - Configuration changes
- 🧪 **Dry-run mode** for testing without live cloud access
- 🎯 **Filter by resource type** (e.g., only EC2, S3)
- 📊 **Severity levels** for prioritization (`info`, `warning`, `critical`)
- 🖨️ Outputs in JSON, Markdown, or CLI table
- 📁 Export results to file
- 🧹 Easily extensible for other cloud providers

---

## Getting Started

### 1. Clone the Repo

```bash
git clone https://github.com/KenB773/IaCDriftDetector.git
cd IaCDriftDetector
```

### 2. Build the CLI

```bash
go build -o drift-detector ./cmd
```

### 3. Run Against a Terraform State File

```bash
./drift-detector --state-file ./examples/sample.tfstate --region us-east-1 --output json
```

---

## CLI Flags

| Flag             | Description                                   |
|------------------|-----------------------------------------------|
| `--state-file`   | Path to your Terraform `.tfstate` file        |
| `--region`       | AWS region to scan                            |
| `--dry-run`      | Skip AWS API calls (for testing only)         |
| `--include`      | Comma-separated resource types to scan        |
| `--output`       | Output format: `json`, `markdown`, `table`    |
| `--output-file`  | Save output to a file                         |
| `--config`       | Optional YAML config file                     |
| `--version`      | Print version info                            |

---

## Example Output

```json
{
  "drift": [
    {
      "resource": "aws_instance.web",
      "change": "Tags differ",
      "severity": "warning"
    },
    {
      "resource": "aws_s3_bucket.logs",
      "change": "Missing from state file",
      "severity": "critical"
    }
  ]
}
```

---

## Security Considerations

This tool **does not make any changes** to your infrastructure. It only performs **read-only scans** and comparisons.

---

## Roadmap

- [ ] Add support for Azure & GCP
- [ ] Integrate with CI/CD for automated drift checks
- [ ] Visualize drift via HTML dashboard
- [ ] Add support for Terraform plans (not just `.tfstate`)

---

## Author

**Ken Brigham**  
[GitHub](https://github.com/KenB773) | [Portfolio](https://kenb773.github.io)

---

## License

MIT License — use freely, credit appreciated.
