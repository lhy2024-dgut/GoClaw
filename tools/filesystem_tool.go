package tools

import (
	"context"
	"os"
)

type FileSystemTool struct{}

func NewFileSystemTool() *FileSystemTool {
	return &FileSystemTool{}
}

func (t *FileSystemTool) Name() string {
	return "filesystem"
}

func (t *FileSystemTool) Description() string {
	return "用于文件系统操作的工具，包括创建文件夹、读写文件等"
}

func (t *FileSystemTool) Parameters() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"action": map[string]interface{}{
				"type":        "string",
				"description": "操作类型：create_dir, list_dir, read_file, write_file, delete",
			},
			"path": map[string]interface{}{
				"type":        "string",
				"description": "文件或文件夹路径",
			},
			"content": map[string]interface{}{
				"type":        "string",
				"description": "文件内容（仅用于 write_file 操作）",
			},
		},
		"required": []string{"action", "path"},
	}
}

func (t *FileSystemTool) Execute(ctx context.Context, params map[string]interface{}) (string, error) {
	action, _ := params["action"].(string)
	path, _ := params["path"].(string)
	content, _ := params["content"].(string)

	switch action {
	case "create_dir":
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return "", err
		}
		return "成功创建文件夹: " + path, nil

	case "list_dir":
		entries, err := os.ReadDir(path)
		if err != nil {
			return "", err
		}
		result := "目录 " + path + " 的内容:\n"
		for _, entry := range entries {
			result += "- " + entry.Name()
			if entry.IsDir() {
				result += " (目录)"
			}
			result += "\n"
		}
		return result, nil

	case "read_file":
		data, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}
		return "文件 " + path + " 的内容:\n```\n" + string(data) + "\n```", nil

	case "write_file":
		err := os.WriteFile(path, []byte(content), 0644)
		if err != nil {
			return "", err
		}
		return "成功写入文件: " + path, nil

	case "delete":
		err := os.Remove(path)
		if err != nil {
			return "", err
		}
		return "成功删除: " + path, nil

	default:
		return "未知操作: " + action, nil
	}
}
