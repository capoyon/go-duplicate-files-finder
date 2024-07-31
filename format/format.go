package format

import "fmt"

// formatBytes converts a byte size into a human-readable string.
func FileSize(bytes int64) string {
    const (
        _       = iota
        KB = 1 << (10 * iota)
        MB
        GB
        TB
        PB
        EB
        ZB
        YB
    )

    var value float64 = float64(bytes)
    var unit string

    switch {
    case value >= PB:
        value /= PB
        unit = "PB"
    case value >= TB:
        value /= TB
        unit = "TB"
    case value >= GB:
        value /= GB
        unit = "GB"
    case value >= MB:
        value /= MB
        unit = "MB"
    case value >= KB:
        value /= KB
        unit = "KB"
    default:
        unit = "B"
    }

    return fmt.Sprintf("%.2f %s", value, unit)
}
