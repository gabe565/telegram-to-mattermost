// Code generated by "enumer -type TextEntityType -trimprefix Type -transform snake -text -output text_entity_type_string.go"; DO NOT EDIT.

package telegram

import (
	"fmt"
	"strings"
)

const _TextEntityTypeName = "plainlinktext_linkboldhashtagitalicmentionmention_nameemailphonecodeprestrikethroughbank_cardcashtag"

var _TextEntityTypeIndex = [...]uint8{0, 5, 9, 18, 22, 29, 35, 42, 54, 59, 64, 68, 71, 84, 93, 100}

const _TextEntityTypeLowerName = "plainlinktext_linkboldhashtagitalicmentionmention_nameemailphonecodeprestrikethroughbank_cardcashtag"

func (i TextEntityType) String() string {
	if i >= TextEntityType(len(_TextEntityTypeIndex)-1) {
		return fmt.Sprintf("TextEntityType(%d)", i)
	}
	return _TextEntityTypeName[_TextEntityTypeIndex[i]:_TextEntityTypeIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _TextEntityTypeNoOp() {
	var x [1]struct{}
	_ = x[TypePlain-(0)]
	_ = x[TypeLink-(1)]
	_ = x[TypeTextLink-(2)]
	_ = x[TypeBold-(3)]
	_ = x[TypeHashtag-(4)]
	_ = x[TypeItalic-(5)]
	_ = x[TypeMention-(6)]
	_ = x[TypeMentionName-(7)]
	_ = x[TypeEmail-(8)]
	_ = x[TypePhone-(9)]
	_ = x[TypeCode-(10)]
	_ = x[TypePre-(11)]
	_ = x[TypeStrikethrough-(12)]
	_ = x[TypeBankCard-(13)]
	_ = x[TypeCashtag-(14)]
}

var _TextEntityTypeValues = []TextEntityType{TypePlain, TypeLink, TypeTextLink, TypeBold, TypeHashtag, TypeItalic, TypeMention, TypeMentionName, TypeEmail, TypePhone, TypeCode, TypePre, TypeStrikethrough, TypeBankCard, TypeCashtag}

var _TextEntityTypeNameToValueMap = map[string]TextEntityType{
	_TextEntityTypeName[0:5]:         TypePlain,
	_TextEntityTypeLowerName[0:5]:    TypePlain,
	_TextEntityTypeName[5:9]:         TypeLink,
	_TextEntityTypeLowerName[5:9]:    TypeLink,
	_TextEntityTypeName[9:18]:        TypeTextLink,
	_TextEntityTypeLowerName[9:18]:   TypeTextLink,
	_TextEntityTypeName[18:22]:       TypeBold,
	_TextEntityTypeLowerName[18:22]:  TypeBold,
	_TextEntityTypeName[22:29]:       TypeHashtag,
	_TextEntityTypeLowerName[22:29]:  TypeHashtag,
	_TextEntityTypeName[29:35]:       TypeItalic,
	_TextEntityTypeLowerName[29:35]:  TypeItalic,
	_TextEntityTypeName[35:42]:       TypeMention,
	_TextEntityTypeLowerName[35:42]:  TypeMention,
	_TextEntityTypeName[42:54]:       TypeMentionName,
	_TextEntityTypeLowerName[42:54]:  TypeMentionName,
	_TextEntityTypeName[54:59]:       TypeEmail,
	_TextEntityTypeLowerName[54:59]:  TypeEmail,
	_TextEntityTypeName[59:64]:       TypePhone,
	_TextEntityTypeLowerName[59:64]:  TypePhone,
	_TextEntityTypeName[64:68]:       TypeCode,
	_TextEntityTypeLowerName[64:68]:  TypeCode,
	_TextEntityTypeName[68:71]:       TypePre,
	_TextEntityTypeLowerName[68:71]:  TypePre,
	_TextEntityTypeName[71:84]:       TypeStrikethrough,
	_TextEntityTypeLowerName[71:84]:  TypeStrikethrough,
	_TextEntityTypeName[84:93]:       TypeBankCard,
	_TextEntityTypeLowerName[84:93]:  TypeBankCard,
	_TextEntityTypeName[93:100]:      TypeCashtag,
	_TextEntityTypeLowerName[93:100]: TypeCashtag,
}

var _TextEntityTypeNames = []string{
	_TextEntityTypeName[0:5],
	_TextEntityTypeName[5:9],
	_TextEntityTypeName[9:18],
	_TextEntityTypeName[18:22],
	_TextEntityTypeName[22:29],
	_TextEntityTypeName[29:35],
	_TextEntityTypeName[35:42],
	_TextEntityTypeName[42:54],
	_TextEntityTypeName[54:59],
	_TextEntityTypeName[59:64],
	_TextEntityTypeName[64:68],
	_TextEntityTypeName[68:71],
	_TextEntityTypeName[71:84],
	_TextEntityTypeName[84:93],
	_TextEntityTypeName[93:100],
}

// TextEntityTypeString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func TextEntityTypeString(s string) (TextEntityType, error) {
	if val, ok := _TextEntityTypeNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _TextEntityTypeNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to TextEntityType values", s)
}

// TextEntityTypeValues returns all values of the enum
func TextEntityTypeValues() []TextEntityType {
	return _TextEntityTypeValues
}

// TextEntityTypeStrings returns a slice of all String values of the enum
func TextEntityTypeStrings() []string {
	strs := make([]string, len(_TextEntityTypeNames))
	copy(strs, _TextEntityTypeNames)
	return strs
}

// IsATextEntityType returns "true" if the value is listed in the enum definition. "false" otherwise
func (i TextEntityType) IsATextEntityType() bool {
	for _, v := range _TextEntityTypeValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalText implements the encoding.TextMarshaler interface for TextEntityType
func (i TextEntityType) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for TextEntityType
func (i *TextEntityType) UnmarshalText(text []byte) error {
	var err error
	*i, err = TextEntityTypeString(string(text))
	return err
}
