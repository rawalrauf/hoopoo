// Generic KDL config file editor for niri
// Supports reading, modifying, and writing KDL config files while preserving comments

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	kdl "github.com/calico32/kdl-go"
)

// SetProperty sets a property value in a KDL config file
// file: relative path to config file in ~/.config/niri/ (e.g., "monitors.kdl")
// blockType: node type to find (e.g., "output", "window-rule")
// blockName: identifier for the block (e.g., "eDP-1" for output "eDP-1")
// property: property name to set (e.g., "scale", "mode")
// value: property value (e.g., "1.25", "off")
func SetProperty(file, blockType, blockName, property, value string) error {
	if blockName == "" || property == "" {
		return fmt.Errorf("blockName and property cannot be empty")
	}

	configPath, err := getConfigPath(file)
	if err != nil {
		return err
	}

	// Read and parse config file
	doc, err := parseConfigFile(configPath, file)
	if err != nil {
		return err
	}

	// Find the target block or create it if it doesn't exist
	targetNode := findOrCreateBlock(doc, blockType, blockName)

	// Normalize all block quotes
	normalizeBlockQuotes(doc, blockType)

	// Get or create children document
	children := targetNode.Children()
	if children == nil {
		targetNode.AddChildrenFunc(func(doc *kdl.Document) {})
		children = targetNode.Children()
	}

	// Find existing property or add new one
	setOrUpdateProperty(targetNode, property, value)

	// Write back to file
	return writeConfigFile(configPath, doc, file)
}

// RemoveProperty removes a property from a KDL config file
func RemoveProperty(file, blockType, blockName, property string) error {
	if blockName == "" || property == "" {
		return fmt.Errorf("blockName and property cannot be empty")
	}

	configPath, err := getConfigPath(file)
	if err != nil {
		return err
	}

	// Read and parse config file
	doc, err := parseConfigFile(configPath, file)
	if err != nil {
		return err
	}

	// Find target block
	targetNode := findBlock(doc, blockType, blockName)
	if targetNode == nil {
		return nil // Nothing to remove
	}

	// Normalize all block quotes
	normalizeBlockQuotes(doc, blockType)

	children := targetNode.Children()
	if children == nil {
		return nil // No children, nothing to remove
	}

	// Find and remove property (keep the block even if empty)
	for i, child := range children.Nodes {
		if child.Name() == property {
			children.Nodes = append(children.Nodes[:i], children.Nodes[i+1:]...)
			break
		}
	}

	// Write back to file
	return writeConfigFile(configPath, doc, file)
}

// GetProperty reads a property value from a KDL config file
func GetProperty(file, blockType, blockName, property string) (string, error) {
	if blockName == "" || property == "" {
		return "", fmt.Errorf("blockName and property cannot be empty")
	}

	configPath, err := getConfigPath(file)
	if err != nil {
		return "", err
	}

	// Read and parse config file
	doc, err := parseConfigFile(configPath, file)
	if err != nil {
		return "", err
	}

	// Find target block
	targetNode := findBlock(doc, blockType, blockName)
	if targetNode == nil {
		return "", fmt.Errorf("block not found: %s \"%s\"", blockType, blockName)
	}

	children := targetNode.Children()
	if children == nil {
		return "", fmt.Errorf("property not found")
	}

	for _, child := range children.Nodes {
		if child.Name() == property {
			childArgs := child.Arguments()
			if len(childArgs) > 0 {
				return childArgs[0].String(), nil
			}
			// Property exists but has no value (like "off")
			return property, nil
		}
	}

	return "", fmt.Errorf("property not found")
}

// Helper functions

// parseConfigFile reads and parses a KDL config file
func parseConfigFile(configPath, file string) (*kdl.Document, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", file, err)
	}

	doc, err := kdl.Parse(strings.NewReader(string(data)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", file, err)
	}

	return doc, nil
}

// writeConfigFile formats and writes a KDL document to file
func writeConfigFile(configPath string, doc *kdl.Document, file string) error {
	output, err := kdl.FormatToString(doc,
		kdl.WithFormatIndentStr("\t"),
		kdl.WithFormatMaxLineLen(20), // Force multi-line formatting for blocks
	)
	if err != nil {
		return fmt.Errorf("failed to format %s: %w", file, err)
	}

	// Fix slashdashed blocks: Format removes quotes from identifiers in slashdashed nodes
	output = fixSlashdashQuotes(output)

	if err := os.WriteFile(configPath, []byte(output), 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", file, err)
	}

	return nil
}

// findBlock finds a block node by type and name
func findBlock(doc *kdl.Document, blockType, blockName string) *kdl.Node {
	for _, node := range doc.Nodes {
		if node.Name() == blockType {
			args := node.Arguments()
			if len(args) > 0 && args[0].String() == blockName {
				return node
			}
		}
	}
	return nil
}

// findOrCreateBlock finds or creates a block node
func findOrCreateBlock(doc *kdl.Document, blockType, blockName string) *kdl.Node {
	targetNode := findBlock(doc, blockType, blockName)
	if targetNode == nil {
		targetNode = kdl.NewNode(blockType, kdl.NewString(blockName).WithLiteral(`"`+blockName+`"`))
		doc.Nodes = append(doc.Nodes, targetNode)
	}
	return targetNode
}

// normalizeBlockQuotes ensures all blocks of a given type have quoted arguments
func normalizeBlockQuotes(doc *kdl.Document, blockType string) {
	for _, node := range doc.Nodes {
		if node.Name() == blockType {
			args := node.Arguments()
			if len(args) > 0 {
				name := args[0].String()
				node.SetArg(0, kdl.NewString(name).WithLiteral(`"`+name+`"`))
			}
		}
	}
}

// setOrUpdateProperty sets or updates a property in a node
func setOrUpdateProperty(targetNode *kdl.Node, property, value string) {
	children := targetNode.Children()
	if children == nil {
		return
	}

	// Find existing property
	var propertyNode *kdl.Node
	for _, child := range children.Nodes {
		if child.Name() == property {
			propertyNode = child
			break
		}
	}

	if propertyNode == nil {
		// Create new property node
		if value == "off" || value == "on" {
			propertyNode = kdl.NewNode(property)
		} else {
			val := parseValue(value)
			propertyNode = kdl.NewNode(property, val)
		}
		targetNode.AddChild(propertyNode)
	} else {
		// Update existing property
		if value == "off" || value == "on" {
			// Remove all arguments for boolean-like keywords
			for len(propertyNode.Arguments()) > 0 {
				propertyNode.RemoveArgument(0)
			}
		} else {
			// Update the value
			val := parseValue(value)
			if len(propertyNode.Arguments()) > 0 {
				propertyNode.SetArg(0, val)
			} else {
				propertyNode.AddArgument(val)
			}
		}
	}
}

func getConfigPath(file string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".config", "niri", file), nil
}

func parseValue(value string) kdl.Value {
	// Try to parse as float
	if strings.Contains(value, ".") {
		var f float64
		if _, err := fmt.Sscanf(value, "%f", &f); err == nil {
			return kdl.NewFloat(f)
		}
	}
	
	// Try to parse as int  
	if !strings.Contains(value, ".") && !strings.Contains(value, "\"") {
		var i int
		if _, err := fmt.Sscanf(value, "%d", &i); err == nil {
			return kdl.NewInt(i)
		}
	}
	
	// Default to string
	return kdl.NewString(value)
}

// fixSlashdashQuotes re-quotes identifiers in slashdashed blocks
// Format removes quotes from identifiers like /- output eDP-1, but niri requires /- output "eDP-1"
func fixSlashdashQuotes(content string) string {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		// Match slashdashed output lines: /- output IDENTIFIER {
		if strings.HasPrefix(strings.TrimSpace(line), "/- output ") {
			// Find the identifier (word after "output")
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				identifier := parts[2]
				// Remove trailing { if present
				identifier = strings.TrimSuffix(identifier, "{")
				identifier = strings.TrimSpace(identifier)
				
				// If not already quoted, add quotes
				if !strings.HasPrefix(identifier, `"`) {
					// Reconstruct line with quoted identifier
					indent := strings.Split(line, "/-")[0]
					lines[i] = indent + "/- output \"" + identifier + "\" {"
				}
			}
		}
	}
	return strings.Join(lines, "\n")
}
