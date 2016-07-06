# Aliyun builder (for packer.io)

阿里云镜像构建工具能够创建在阿里云上使用的弹性计算ECS镜像。 

该镜像构建工具通过调用阿里云的API，以如下步骤构建镜像：

* 创建并启动阿里云弹性计算ECS实例；
* 在ECS实例上安装和运行必要的程序(provisioning)；
* 基于ECS实例创建对应的快照(snapshot)；
* 基于快照创建镜像(image)。

创建的镜像后续可以作为在阿里云上启动新ECS服务器的基础，支持所谓镜像部署。

镜像管理工具仅负责创建镜像，镜像创建以后存放在阿里云上，镜像管理工具并不负责后续的镜像管理任务。

## 安装


## 基本的例子

下面的模板是经过测试可以工作的，它创建一个基本的ubuntu 14镜像:

```JSON

{
  "variables": {
    "ali_access_key_id": "{{env  `ALIYUN_ACCESS_KEY_ID`}}",
    "ali_access_key_secret": "{{env `ALIYUN_ACCESS_KEY_SECRET`}}",
    "ali_security_group_id": "{{env `ALIYUN_SECURITY_GROUP_ID`}}"
  },
  "provisioners": [
    {
      "type": "shell",
      "inline": [
        "echo ubuntu"
      ]
    }
  ],
  "builders": [{
    "type": "aliyun",
    "access_key_id": "{{user `ali_access_key_id`}}",
    "access_key_secret": "{{user `ali_access_key_secret`}}",
    "region_id": "cn-shanghai",
    "base_image_id": "ubuntu1404_64_40G_cloudinit_20160427.raw",
    "instance_type": "ecs.s2.large",
    "security_group_id": "{{user `ali_security_group_id`}}",
    "image_name": "packer-ubuntu-{{timestamp}}",
    "image_description": "ubuntu image created by packer at {{timestamp}}"
  }]
}

```

关于模板含义和如何运行packer，请参考[packer.io的文档](https://www.packer.io/docs/)。

## 模板配置参考

下面列出支持的必选和可选参数。

### 必选参数

* `access_key_id` (string) - 访问阿里云API的**Access Key Id**。 如果未指定，则从**ALIYUN_ACCESS_KEY_ID**环境变量中获取值。
* `access_key_secret` (string) - 访问阿里云API的**Access Key Secret**。如果未指定，则从**ALIYUN_ACCESS_KEY_SECRET**环境变量中获取值。
* `security_group_id` (string) - 阿里云安全组id，需要在阿里云上预创建安全组。如果未指定，则从**ALIYUN_SECURITY_GROUP_ID**环境变量中获取。
* `region_id` (string) - 阿里云区域id，例如*cn-shanghai*，可通过阿里云API *DescribeRegions*查询。
* `base_image_id` (string) - 基础镜像文件ID，表示启动实例时选择的镜像资源。可通过阿里云API *DescribeImages*查询，[这里](examples/ali_base_image.txt)有一份备份的调用返回供参考。
* `instance_type` (string) - 实例的资源规则。取值参见[实例资源规格对照表](https://help.aliyun.com/document_detail/25685.html?spm=5176.doc25499.2.5.HgFKqE)，也可通过阿里云API *DescribeInstanceTypes*查询。

### 可选参数

* `image_name` (string) - 目标镜像的名称，缺省值*packer-{{timestamp}}*, 命名规则参考阿里云API *CreateImage*。
* `image_description` (string) - 目标镜像的描述，缺省值*Created by Packer for Aliyun*, 命名规则参考阿里云API *CreateImage*。

## 贡献

欢迎任何对Packer builder for Aliyun的意见、建议和贡献！

如有疑问可以随时提问，你只需通过[Create an issue](https://github.com/archcentric/packer-builder-aliyun/issues)进行询问。

## TODO
* 目前ssh仅通过阿里云内网IP访问来完成软件安装过程(provisioning)，也就是说运行packer的机器和目标ECS实例都要在阿里云内部才能访问到，后续考虑增加支持ssh通过阿里云公网IP访问的方式。
* 创建快照过程中的进度判断，目前仅判断progress完成度百分比，后续考虑增加状态status判断，需要完善第三方的[aliyun golang sdk](https://github.com/denverdino/aliyungo)。


