# FileChecker

Utility to check whether files (e.g. uploaded) are of the authorised types before saving on disk

## Usage

```go
var uploadedFile *multipart.FileHeader

fc := GetFileChecker(uploadedFile)
if !fc.IsAuthorised() {
    fmt.Println("File not authorised")
}
```
