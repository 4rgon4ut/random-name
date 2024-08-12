package types

type CommandType string

const (
    AddItem     CommandType = "addItem"
    DeleteItem  CommandType = "deleteItem"
    GetItem     CommandType = "getItem"
    GetAllItems CommandType = "getAllItems"
)

type Command struct {
    Type  CommandType `json:"type"`
    Key   string      `json:"key,omitempty"`
    Value string      `json:"value,omitempty"`
}

// TODO prob add logging?
func (cmd *Command) IsValid() bool {
    switch cmd.Type {
    case AddItem:
        if cmd.Key == "" || cmd.Value == "" {
            return false
        }
    case DeleteItem, GetItem:
        if cmd.Key == "" {
            return false
        }
    case GetAllItems:
		// do nothing
	default:
		// unknown cmd
        return false
    }
    return true
}

