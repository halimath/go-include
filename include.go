// Package include provides declarations for functions that include file content into regular go files.
// The functions exported by this package do not actually include content but work as a placeholder to make
// the go files compile in your IDE. Use the generator to generate a version with calls to functions of this
// package being replaced.
package include

// String includes the file named filename by replacing the call with a string literal of the file's content.
func String(filename string) string {
	return ""
}

// Bytes includes the file named filename by replacing the call with a byte slice of the file's content.
func Bytes(filename string) []byte {
	return nil
}
