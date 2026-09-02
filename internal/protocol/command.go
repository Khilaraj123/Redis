// Parsed command struct
package protocol

// Command is one parsed client command, e.g. GET with args ["mykey"].
type Command struct {
	Name string
	Args []string
}
