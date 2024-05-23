package DialogBuilder

import (
	"fmt"
	"strings"

	"github.com/eikarna/gotps/handler/utils"
)

type Vec2i struct {
	X int
	Y int
}

type DialogBuilder struct {
	dialog string
}

func NewDialogBuilder(color string) *DialogBuilder {
	if color == "" {
		color = "`0"
	}
	return &DialogBuilder{dialog: "set_default_color|" + color + "\n"}
}

func (db *DialogBuilder) String() string {
	return db.dialog
}

func (db *DialogBuilder) Raw(dialogText string) *DialogBuilder {
	db.dialog += dialogText
	return db
}

func (db *DialogBuilder) SetCustomSpacing(x, y int) *DialogBuilder {
	db.dialog += fmt.Sprintf("\nset_custom_spacing|x:%d;y:%d|", x, y)
	return db
}

func (db *DialogBuilder) AddBreak() *DialogBuilder {
	db.dialog += "\nadd_custom_break|"
	return db
}

func (db *DialogBuilder) TextScaling(label string) *DialogBuilder {
	db.dialog += fmt.Sprintf("\ntext_scaling_string|%s|", label)
	return db
}

func (db *DialogBuilder) AddStatic(name string, id int, message string) *DialogBuilder {
	db.dialog += fmt.Sprintf("\nadd_button_with_icon|%s|%s|frame|%d||", name, message, id)
	return db
}

func (db *DialogBuilder) AddFriendImageLabelButton(name, label, texturePath string, size float64, texture Vec2i) *DialogBuilder {
	db.dialog += fmt.Sprintf("\nadd_friend_image_label_button|%s|%s|%s|%.2f|%d|%d|32|false|false|",
		name, label, texturePath, size, texture.X, texture.Y)
	return db
}

func (db *DialogBuilder) AddImageButton(name, file, extra string) *DialogBuilder {
	if extra == "" {
		extra = "bannerlayout"
	}
	db.dialog += fmt.Sprintf("\nadd_image_button|%s|interface/large/%s.rttex|%s|||", name, file, extra)
	return db
}

func (db *DialogBuilder) AddItemPicker(name, message string) *DialogBuilder {
	db.dialog += fmt.Sprintf("\nadd_item_picker|%s|%s|%s|", name, message, "Choose an item from your inventory")
	return db
}

func (db *DialogBuilder) AddPlayerInfo(name string, currentLevel, currentExp, expRequired int) *DialogBuilder {
	db.dialog += fmt.Sprintf("\nadd_player_info|%s|%d|%d|%d|", name, currentLevel, currentExp, expRequired)
	return db
}

func (db *DialogBuilder) AddCheckbox(checked bool, name, message string) *DialogBuilder {
	db.dialog += fmt.Sprintf("\nadd_checkbox|%s|%s|%d|", name, message, utils.BoolToInt(checked))
	return db
}

func (db *DialogBuilder) AddSelectorCheckbox(id int, name string, messages []string, indexChecked int) *DialogBuilder {
	for i, msg := range messages {
		db.dialog += fmt.Sprintf("\nadd_checkbox|%s_%d_%d|%s|%d|", name, i, id, msg, utils.BoolToInt(i == indexChecked))
	}
	return db
}

func (db *DialogBuilder) AddSmallText(message string) *DialogBuilder {
	db.dialog += fmt.Sprintf("\nadd_smalltext|%s|", message)
	return db
}

func (db *DialogBuilder) EndList() *DialogBuilder {
	db.dialog += "\nadd_button_with_icon||END_LIST|noflags|0|0|"
	return db
}

func (db *DialogBuilder) AddDualLayer(big, iconLeft bool, foreground, background int, size float64, message string) *DialogBuilder {
	db.dialog += fmt.Sprintf("\nadd_dual_layer_icon_label|%s|%s|left|%d|%d|%.2f|%d|", utils.BoolToStr(big, "big", "small"), message, background, foreground, size, utils.BoolToInt(!iconLeft))
	return db
}

func (db *DialogBuilder) AddTextInput(length int, name, message, defaultInput string) *DialogBuilder {
	db.dialog += fmt.Sprintf("\nadd_text_input|%s|%s|%s|%d|", name, message, defaultInput, length)
	return db
}

func (db *DialogBuilder) AddSeedIcon(itemId string) *DialogBuilder {
	db.dialog += fmt.Sprintf("\nadd_seed_color_icons|%s|", itemId)
	return db
}

func (db *DialogBuilder) AddStaticIconButton(name string, id int, message, hoverNumber string) *DialogBuilder {
	db.dialog += fmt.Sprintf("\nadd_button_with_icon|%s|%s|staticBlueFrame|%d|%s|", name, message, id, hoverNumber)
	return db
}

func (db *DialogBuilder) AddLabelIcon(big bool, id int, message string) *DialogBuilder {
	db.dialog += fmt.Sprintf("\nadd_label_with_icon|%s|%s|left|%d|", utils.BoolToStr(big, "big", "small"), message, id)
	return db
}

func (db *DialogBuilder) AddIconButton(btnName, text, option string, itemID int, unkVal string) *DialogBuilder {
	db.dialog += fmt.Sprintf("\nadd_button_with_icon|%s|%s|%s|%d|%s|", btnName, text, option, itemID, unkVal)
	return db
}

func (db *DialogBuilder) AddKitDisabledButton(btnName, progress string, itemID int) *DialogBuilder {
	db.dialog += fmt.Sprintf("\nadd_button_with_icon|%s|`4%s|staticGreyFrame,no_padding_x,is_count_label,disabled|%d||", btnName, progress, itemID)
	return db
}

func (db *DialogBuilder) AddDisabledButton(btnName, progress string, itemID int) *DialogBuilder {
	db.dialog += fmt.Sprintf("\nadd_button_with_icon|%s|`w%s`|staticGreyFrame,no_padding_x,is_count_label,disabled|%d||", btnName, progress, itemID)
	return db
}

func (db *DialogBuilder) AddKitClaimButton(btnName, underText string, itemID int) *DialogBuilder {
	db.dialog += fmt.Sprintf("\nadd_button_with_icon|%s|`2%s`|staticYellowFrame,no_padding_x,is_count_label|%d||", btnName, underText, itemID)
	return db
}

func (db *DialogBuilder) AddKitClaimedButton(btnName string, itemID int) *DialogBuilder {
	db.dialog += fmt.Sprintf("\nadd_button_with_icon|%s|`5CLAIMED`|staticBlueFrame,no_padding_x,is_count_label|%d||", btnName, itemID)
	return db
}

func (db *DialogBuilder) AddCenterButton(btnName, label string, itemID int) *DialogBuilder {
	db.dialog += fmt.Sprintf("\nadd_button_with_icon|%s|%s|staticBlueFrame,no_padding_x,is_count_label|%d||", btnName, label, itemID)
	return db
}

func (db *DialogBuilder) AddLabelIconButton(big bool, message string, id int, name string) *DialogBuilder {
	db.dialog += fmt.Sprintf("\nadd_label_with_icon_button|%s|%s|left|%d|%s|", utils.BoolToStr(big, "big", "small"), message, id, name)
	return db
}

func (db *DialogBuilder) AddSpacer(big bool) *DialogBuilder {
	db.dialog += fmt.Sprintf("\nadd_spacer|%s|", utils.BoolToStr(big, "big", "small"))
	return db
}

func (db *DialogBuilder) AddTextbox(message string) *DialogBuilder {
	db.dialog += fmt.Sprintf("\nadd_textbox|%s|", message)
	return db
}

func (db *DialogBuilder) AddQuickExit() *DialogBuilder {
	db.dialog += "\nadd_quick_exit|"
	return db
}

func (db *DialogBuilder) StartCustomTabs() *DialogBuilder {
	db.dialog += "\nstart_custom_tabs|"
	return db
}

func (db *DialogBuilder) ResetPlacementX() *DialogBuilder {
	db.dialog += "\nreset_placement_x|"
	return db
}

func (db *DialogBuilder) ResetPlacementY() *DialogBuilder {
	db.dialog += "\nreset_placement_y|"
	return db
}

func (db *DialogBuilder) AddCustomMargin(x, y int) *DialogBuilder {
	db.dialog += fmt.Sprintf("\nadd_custom_margin|x:%d;y:%d|", x, y)
	return db
}

func (db *DialogBuilder) AddPlayerPicker(name, button string) *DialogBuilder {
	db.dialog += fmt.Sprintf("\nadd_player_picker|%s|%s|", name, button)
	return db
}

func (db *DialogBuilder) AddInput(length int, name, message, defaultInput string) *DialogBuilder {
	db.dialog += fmt.Sprintf("\nadd_text_input|%s|%s|%s|%d|", name, message, defaultInput, length)
	return db
}

func (db *DialogBuilder) EndDialog(name, cancel, accept string) *DialogBuilder {
	db.dialog += fmt.Sprintf("\nend_dialog|%s|%s|%s|", name, cancel, accept)
	return db
}

func (db *DialogBuilder) AddLabel(big bool, message string) *DialogBuilder {
	db.dialog += fmt.Sprintf("\nadd_label|%s|%s|left|", utils.BoolToStr(big, "big", "small"), message)
	return db
}

func (db *DialogBuilder) AddButton(name, button string) *DialogBuilder {
	db.dialog += fmt.Sprintf("\nadd_button|%s|%s|noflags|0|0|", name, button)
	return db
}

func (db *DialogBuilder) AddSmallFontButton(name, button string) *DialogBuilder {
	db.dialog += fmt.Sprintf("\nadd_small_font_button|%s|%s|noflags|0|0|", name, button)
	return db
}

func (db *DialogBuilder) AddDisabledButtonAlt(name, button string) *DialogBuilder {
	db.dialog += fmt.Sprintf("\nadd_button|%s|%s|off|0|0|", name, button)
	return db
}

func (db *DialogBuilder) AddSmallFontDisabledButton(name, button string) *DialogBuilder {
	db.dialog += fmt.Sprintf("\nadd_small_font_button|%s|%s|off|0|0|", name, button)
	return db
}

func (db *DialogBuilder) AddCustomButton(name, option string) *DialogBuilder {
	db.dialog += fmt.Sprintf("\nadd_custom_button|%s|%s|", name, option)
	return db
}

func (db *DialogBuilder) AddCustomLabel(option1, option2 string) *DialogBuilder {
	db.dialog += fmt.Sprintf("\nadd_custom_label|%s|%s|", option1, option2)
	return db
}

func (db *DialogBuilder) AddCustomSpacer(x float64) *DialogBuilder {
	db.dialog += fmt.Sprintf("\nadd_custom_spacer|x:%f|", x)
	return db
}

func (db *DialogBuilder) AddCustomTextbox(text, size string) *DialogBuilder {
	db.dialog += fmt.Sprintf("\nadd_custom_textbox|%s|size:%s|", text, strings.ToLower(size))
	return db
}

func (db *DialogBuilder) EmbedData(pushFront bool, embed, data string) *DialogBuilder {
	if pushFront {
		db.dialog = fmt.Sprintf("\nembed_data|%s|%s\n", embed, data) + db.dialog
	} else {
		db.dialog += fmt.Sprintf("\nembed_data|%s|%s", embed, data)
	}
	return db
}

func (db *DialogBuilder) AddAchieveButton(achName, achToGet string, achID int, unk string) *DialogBuilder {
	db.dialog += fmt.Sprintf("\nadd_achieve_button|%s|%s|left|%d|%s|", achName, achToGet, achID, unk)
	return db
}

/*
func main() {
	db := NewDialogBuilder("`0")
	db.SetCustomSpacing(10, 20).
		AddBreak().
		TextScaling("Example").
		AddStatic("btn1", 1, "Click me").
		AddFriendImageLabelButton("friendBtn", "Friend", "texture/path", 1.5, Vec2i{10, 20}).
		AddImageButton("imgBtn", "file", "extra")

	fmt.Println(db.String())
}*/
