<div style="text-align: center"></div>
  <p align="center">
  <img src="https://user-images.githubusercontent.com/42825450/201523059-ed28e427-e1b6-443a-8326-100460e6dec9.jpg" width="250px" height="220px">
      <br>
      <i>Make the project development and release simpler, easier and more efficient.</i>
  </p>
</div>



## What is DKube
This is a K8s cluster management platform；DKube Provides a wizard-style operation interface for K8s cluster management to help your team manage your cluster environment quickly and easily

## 前端 | Front
- 前端项目仓库地址：[DKube-Web](https://github.com/dqzboy/DKube-Web)

## 功能 | Function

<details>
  <summary><b> K8s集群管理 | K8s cluster management</b></summary>
</details>

<details>
  <summary><b> 平台化界面 | platform interface</b></summary>
</details>

<details>
  <summary><b> 更加便捷的管理K8s | More convenient management of K8s</b></summary>
</details>

<details>
  <summary><b> 支持YAML信息查看\变更 | Support YAML information view\change</b></summary>
</details>

<details>
  <summary><b> 平台化管理控制器 | Platform Management Controller</b></summary>
</details>

<details>
  <summary><b> 实时查看容器日志 | View container logs in real time</b></summary>
</details>

## 截图 | screenshot
<br/>
<table>
    <tr>
      <td width="50%" align="center"><b>登入认证管理 | Login authentication management</b></td>
      <td width="50%" align="center"><b>集群信息状态 | Cluster Information Status</b></td>
    </tr>
    <tr>
        <td width="50%" align="center"><img src="https://user-images.githubusercontent.com/42825450/193593148-4d258b30-b972-4583-b359-32978a8a8637.jpg?raw=true"></td>
        <td width="50%" align="center"><img src="https://user-images.githubusercontent.com/42825450/193593170-3373dabd-8d5d-4a01-a59f-49851f11f433.jpg?raw=true"></td>
    </tr>
    <tr>
      <td width="50%" align="center"><b>节点资源管理 | Node resource management</b></td>
      <td width="50%" align="center"><b>名称空间管理 | namespace management</b></td>
    </tr>
        <td width="50%" align="center"><img src="https://user-images.githubusercontent.com/42825450/193593569-daebc649-f6c4-45a2-88f6-2aa4860c3dea.jpg?raw=true"></td>
        <td width="50%" align="center"><img src="https://user-images.githubusercontent.com/42825450/193593579-e0539ab0-6b22-4060-b254-c6495fb87cbd.jpg?raw=true"></td>
    <tr>
    </tr>
    <tr>
      <td width="50%" align="center"><b>YAML信息管理 | YAML information management</b></td>
      <td width="50%" align="center"><b>Pod副本管理 | Pod copy management</b></td>
    </tr>
        <td width="50%" align="center"><img src="https://user-images.githubusercontent.com/42825450/193593867-4a98bd0f-a910-4b90-92e3-6a3164d0c241.jpg?raw=true"></td>
        <td width="50%" align="center"><img src="https://user-images.githubusercontent.com/42825450/193593871-ee004cb8-42cb-427a-a0cc-fa1e15e7d466.jpg?raw=true"></td>
    <tr>
    </tr>
</table>


## 部署 | deploy
### 后端
- 修改kube-config路径: `config/config.go`

```shell
Kubeconfig = "/root/.kube/config"
```

- 后端
```go
env GOOS=linux GOARCH=amd64 go build -o DKube-Server main.go

 ~]# chmod +x DKube-Server 
 ~]# ./DKube-Server
```

- 数据库|database: `config/config.go`
```shell
//数据库配置|database configuration
DbType = "mysql"
DbHost = "xx.xx.xx.xx"
DbPort = 3306
DbName = "dkube"
DbUser = "root"
DbPwd  = "123456"
```

### 前端
```shell
npm install
npm run build
```

## 推广


<table>
  <thead>
    <tr>
      <th width="50%" align="center">描述信息</th>
      <th width="50%" align="center">图文介绍</th>
    </tr>
  </thead>
  <tbody>
    <!-- 第一个广告：RackNerd -->
    <tr>
      <td width="50%" align="left">
        <a href="https://dqzboy.github.io/proxyui/racknerd" target="_blank">提供高性价比的海外VPS，支持多种操作系统，适合搭建Docker代理服务。</a>
      </td>
      <td width="50%" align="center">
        <a href="https://dqzboy.github.io/proxyui/racknerd" target="_blank">
          <img src="https://cdn.jsdelivr.net/gh/dqzboy/Images/dqzboy-proxy/Image_2025-07-07_16-14-49.png?raw=true" alt="RackNerd" width="200" height="120">
        </a>
      </td>
    </tr>
    <!-- 第二个广告：CloudCone -->
    <tr>
      <td width="50%" align="left">
        <a href="https://dqzboy.github.io/proxyui/CloudCone" target="_blank">CloudCone 提供灵活的云服务器方案，支持按需付费，适合个人和企业用户。</a>
      </td>
      <td width="50%" align="center">
        <a href="https://dqzboy.github.io/proxyui/CloudCone" target="_blank">
          <img src="https://cdn.jsdelivr.net/gh/dqzboy/Images/dqzboy-proxy/111.png?raw=true" alt="CloudCone" width="200" height="120">
        </a>
      </td>
    </tr>
  </tbody>
</table>


[本项目 CDN 加速及安全防护由 Tencent EdgeOne 赞助](https://edgeone.ai/zh?from=github)

<img src="https://github.com/user-attachments/assets/2d7a5081-a0e7-4d4d-9148-c3c424f23cc8" width="300" height="80">

