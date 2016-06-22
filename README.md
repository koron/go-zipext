# zip.File#Extra parser/accessor

Sample code:

```go
func lsext(name string) error {
	zr, err := zip.OpenReader(name)
	if err != nil {
		return err
	}
	defer zr.Close()
	for _, zf := range zr.File {
		ex := zipext.Parse(zf)
		fmt.Printf("%s\n", zf.Name)
		fmt.Printf("  ModTime: %s\n", zf.ModTime())
		fmt.Printf("  (EX)ModTime: %s\n", ex.ModTime())
		fmt.Printf("  (EX)AcTime:  %s\n", ex.AcTime())
		fmt.Printf("  (EX)CrTime:  %s\n", ex.CrTime())
	}
	return nil
}
```

Example:

```
$ go run ./cmd/lsext/lsext.go netupvim-v1.1.zip
netupvim.exe
  ModTime: 2016-06-20 01:25:12 +0000 UTC
  (EX)ModTime: 2016-06-20 01:25:12 +0900 JST
  (EX)AcTime:  0001-01-01 00:00:00 +0000 UTC
  (EX)CrTime:  0001-01-01 00:00:00 +0000 UTC
UPDATE.bat
  ModTime: 2016-05-04 10:40:36 +0000 UTC
  (EX)ModTime: 2016-05-04 10:40:36 +0900 JST
  (EX)AcTime:  0001-01-01 00:00:00 +0000 UTC
  (EX)CrTime:  0001-01-01 00:00:00 +0000 UTC
RESTORE.bat
  ModTime: 2016-05-04 10:40:36 +0000 UTC
  (EX)ModTime: 2016-05-04 10:40:36 +0900 JST
  (EX)AcTime:  0001-01-01 00:00:00 +0000 UTC
  (EX)CrTime:  0001-01-01 00:00:00 +0000 UTC
```
