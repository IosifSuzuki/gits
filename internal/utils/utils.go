package utils

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gomarkdown/markdown/ast"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type Number interface {
	int | int8 | int16 | int32 |
	uint | uint8 | uint16 | uint32 |
	float32 | float64
}

func HashPassword(password string) (*string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return nil, err
	}
	hashStr := string(hash)
	return &hashStr, nil
}

func CompareHashAndPassword(hash string, password string) (bool, error) {
	lhsBytes := []byte(hash)
	rhsBytes := []byte(password)
	err := bcrypt.CompareHashAndPassword(lhsBytes, rhsBytes)
	if err != nil {
		return false, err
	}
	return true, err
}

func DateFormat(time time.Time, format string) string {
	return time.Format(format)
}

func Add(lsh, rhs int) int {
	return lsh + rhs
}

func Mul(lhs, rhs int) int {
	return lhs * rhs
}

func PrefixString(text string, words int) string {
	parts := strings.Split(text, " ")
	var actualWords int
	if len(parts)-1 > words {
		actualWords = words
	} else {
		actualWords = len(parts) - 1
	}
	return strings.Join(parts[:actualWords], " ")
}

func Max[T Number](x, y T) T {
	if x > y {
		return x
	}

	return y
}

func Min[T Number](x, y T) T {
	if x < y {
		return x
	}

	return y
}

func MarkdownHasTextNode(node ast.Node) bool {
	if leaf := node.AsLeaf(); leaf != nil {
		if textNode, ok := node.(*ast.Text); ok && len(textNode.Literal) > 0 {
			return true
		}
	}
	for _, child := range node.GetChildren() {
		if MarkdownHasTextNode(child) {
			return true
		}
	}

	return false
}

func MarkdownTextContent(node ast.Node) string {
	if leaf := node.AsLeaf(); leaf != nil {
		switch node := node.(type) {
		case *ast.Text:
			return string(node.Literal)
		default:
			return ""
		}
	}
	var result string
	for _, child := range node.GetChildren() {
		childText := MarkdownTextContent(child)
		if len(childText) > 0 {
			result += childText
		}
	}

	return result
}

func String(str string) *string {
	return &str
}

func MarshalQueryParam(value any) (*string, error) {
	valueJSON, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	base64ValueText := base64.StdEncoding.EncodeToString(valueJSON)
	return &base64ValueText, err
}

func Unmarshal(base64QueryParam string, value any) error {
	jsonQueryParam, err := base64.StdEncoding.DecodeString(base64QueryParam)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(jsonQueryParam, &value); err != nil {
		return err
	}

	return err
}
