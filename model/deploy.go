package model

type DeployReq struct {
	RepoURL       string `json:"repo_url"`       // Git 仓库地址
	Branch        string `json:"branch"`         // 分支名称
	LocalDir      string `json:"local_dir"`      // 本地目录（可选，视业务需求）
	ImageName     string `json:"image_name"`     // Docker 镜像名称
	ContainerName string `json:"container_name"` // 容器名称
	PortMapping   string `json:"port_mapping"`   // 端口映射，例如 "8080:80"
}
