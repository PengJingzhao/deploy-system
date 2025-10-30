package deployer

import (
	"fmt"
	"os"
	"os/exec"
)

// GitHubDeployer 类似一个“类”，负责拉取github项目、构建镜像、运行容器
type GitHubDeployer struct {
	RepoURL       string // GitHub仓库地址，如：https://github.com/username/repo
	Branch        string // 分支，默认为main或master
	LocalDir      string // 克隆到本地的目录
	ImageName     string // 构建的Docker镜像名称
	ContainerName string // 运行的容器名称
}

// NewGitHubDeployer 构造函数，初始化部署器
func NewGitHubDeployer(repoURL, imageName, containerName string) *GitHubDeployer {
	return &GitHubDeployer{
		RepoURL:       repoURL,
		Branch:        "main", // 默认分支
		LocalDir:      "",     // 默认为空，自动设置
		ImageName:     imageName,
		ContainerName: containerName,
	}
}

// SetBranch 可选：设置分支
func (d *GitHubDeployer) SetBranch(branch string) {
	d.Branch = branch
}

// Deploy 执行完整流程：拉取代码 -> 构建镜像 -> 运行容器
func (d *GitHubDeployer) Deploy() error {
	var err error

	// 1. 拉取 GitHub 项目
	if err = d.CloneRepo(); err != nil {
		return fmt.Errorf("克隆仓库失败: %v", err)
	}

	// 2. 构建 Docker 镜像
	if err = d.BuildDockerImage(); err != nil {
		return fmt.Errorf("构建Docker镜像失败: %v", err)
	}

	// 3. 运行 Docker 容器
	if err = d.RunDockerContainer(); err != nil {
		return fmt.Errorf("运行Docker容器失败: %v", err)
	}

	fmt.Println("✅ 部署成功！")
	return nil
}

// CloneRepo 使用 git clone 拉取代码
func (d *GitHubDeployer) CloneRepo() error {
	// 如果没有设置本地目录，使用仓库名作为目录名
	if d.LocalDir == "" {
		// 从 URL 中提取仓库名，例如 https://github.com/user/repo -> repo
		// 简单处理，仅适用于标准 HTTPS GitHub URL
		parts := splitGitHubURL(d.RepoURL)
		if len(parts) < 2 {
			return fmt.Errorf("无法解析 GitHub 仓库地址: %s", d.RepoURL)
		}
		repoName := parts[1]
		if d.Branch != "main" {
			// 可选：根据分支命名目录
			d.LocalDir = repoName + "-" + d.Branch
		} else {
			d.LocalDir = repoName
		}
	}

	// 如果目录已存在，可选择删除或跳过。这里简单起见，报错提示。
	if _, err := os.Stat(d.LocalDir); err == nil {
		return fmt.Errorf("本地目录 %s 已存在，请删除或选择其他目录", d.LocalDir)
	}

	// 执行 git clone
	cloneCmd := exec.Command("git", "clone", "--branch", d.Branch, d.RepoURL, d.LocalDir)
	cloneCmd.Stdout = os.Stdout
	cloneCmd.Stderr = os.Stderr
	if err := cloneCmd.Run(); err != nil {
		return err
	}

	fmt.Printf("📥 已克隆仓库到目录: %s\n", d.LocalDir)
	return nil
}

// BuildDockerImage 构建 Docker 镜像
func (d *GitHubDeployer) BuildDockerImage() error {
	if d.ImageName == "" {
		return fmt.Errorf("镜像名称不能为空")
	}

	// 切换到项目目录
	if err := os.Chdir(d.LocalDir); err != nil {
		return fmt.Errorf("无法进入目录 %s: %v", d.LocalDir, err)
	}

	// 执行 docker build
	buildCmd := exec.Command("docker", "build", "-t", d.ImageName, ".")
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	if err := buildCmd.Run(); err != nil {
		return err
	}

	fmt.Printf("🐳 已构建 Docker 镜像: %s\n", d.ImageName)
	return nil
}

// RunDockerContainer 运行容器
func (d *GitHubDeployer) RunDockerContainer() error {
	if d.ContainerName == "" {
		d.ContainerName = d.ImageName + "-container"
	}

	// 可以根据需要添加端口映射、环境变量等参数
	// 例如：docker run -d --name mycontainer -p 8080:80 myimage
	runCmd := exec.Command("docker", "run", "--name", d.ContainerName, "-d", d.ImageName)
	runCmd.Stdout = os.Stdout
	runCmd.Stderr = os.Stderr
	if err := runCmd.Run(); err != nil {
		return err
	}

	fmt.Printf("🚀 已运行容器: %s (基于镜像: %s)\n", d.ContainerName, d.ImageName)
	return nil
}

// 工具函数：从 GitHub URL 中提取用户名和仓库名
func splitGitHubURL(url string) []string {
	// 简单实现，仅支持 https://github.com/username/repo.git 或 https://github.com/username/repo
	url = trimGitSuffix(url)
	parts := splitPath(url)
	if len(parts) >= 3 && parts[0] == "https:" && parts[1] == "github.com" {
		return []string{parts[2], parts[3]}
	}
	return []string{}
}

func trimGitSuffix(s string) string {
	if len(s) > 4 && s[len(s)-4:] == ".git" {
		return s[:len(s)-4]
	}
	return s
}

func splitPath(path string) []string {
	// 假设 path 是类似 "https://github.com/user/repo" 的字符串
	// 简单按 "/" 分割
	var parts []string
	start := 0
	for i := 0; i < len(path); i++ {
		if path[i] == '/' {
			if i > start {
				parts = append(parts, path[start:i])
			}
			start = i + 1
		}
	}
	if start < len(path) {
		parts = append(parts, path[start:])
	}
	return parts
}
