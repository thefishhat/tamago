package component

import "fmt"

type componentItem struct {
	name   string
	value  interface{}
	input  *toggleInput
	errMsg *errorMsg
}

func (i componentItem) Title() string       { return "Type: " + i.name }
func (i componentItem) FilterValue() string { return i.name + fmt.Sprintf("%v", i.value) }
func (i componentItem) Description() string {
	var renderedItem string
	if i.input.IsEditing() {
		renderedItem = i.input.View()
	} else {
		renderedItem = "Value: " + fmt.Sprintf("%v", i.value)
	}
	if i.errMsg.msg != "" {
		renderedItem += " " + i.errMsg.View()
	}
	return renderedItem
}

func newComponentItem(name string, value interface{}) componentItem {
	return componentItem{
		name:   name,
		value:  value,
		input:  newToggleInput(value),
		errMsg: newErrorMsg(),
	}
}
